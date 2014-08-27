<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/encryption/address.go - The Go Programming Language</title>

<link type="text/css" rel="stylesheet" href="/lib/godoc/style.css">

<script type="text/javascript">window.initFuncs = [];</script>
</head>
<body>

<div id="topbar" class="wide"><div class="container">

<form method="GET" action="/search">
<div id="menu">
<a href="/doc/">Documents</a>
<a href="/pkg/">Packages</a>
<a href="/project/">The Project</a>
<a href="/help/">Help</a>
<a href="/blog/">Blog</a>

<input type="text" id="search" name="q" class="inactive" value="Search" placeholder="Search">
</div>
<div id="heading"><a href="/">The Go Programming Language</a></div>
</form>

</div></div>



<div id="page" class="wide">
<div class="container">


  <div id="plusone"><g:plusone size="small" annotation="none"></g:plusone></div>
  <h1>Source file src/pkg/emp/encryption/address.go</h1>




<div id="nav"></div>


<pre><span id="L1" class="ln">     1</span>	<span class="comment">/**
</span><span id="L2" class="ln">     2</span>	<span class="comment">    Copyright 2014 JARST, LLC.
</span><span id="L3" class="ln">     3</span>	<span class="comment">    
</span><span id="L4" class="ln">     4</span>	<span class="comment">    This file is part of EMP.
</span><span id="L5" class="ln">     5</span>	<span class="comment">
</span><span id="L6" class="ln">     6</span>	<span class="comment">    EMP is distributed in the hope that it will be useful,
</span><span id="L7" class="ln">     7</span>	<span class="comment">    but WITHOUT ANY WARRANTY; without even the implied warranty of
</span><span id="L8" class="ln">     8</span>	<span class="comment">    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the included
</span><span id="L9" class="ln">     9</span>	<span class="comment">    LICENSE file for more details.
</span><span id="L10" class="ln">    10</span>	<span class="comment">**/</span>
<span id="L11" class="ln">    11</span>	
<span id="L12" class="ln">    12</span>	package encryption
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;code.google.com/p/go.crypto/ripemd160&#34;
<span id="L16" class="ln">    16</span>		&#34;crypto/elliptic&#34;
<span id="L17" class="ln">    17</span>		&#34;crypto/rand&#34;
<span id="L18" class="ln">    18</span>		&#34;crypto/sha512&#34;
<span id="L19" class="ln">    19</span>		&#34;encoding/base64&#34;
<span id="L20" class="ln">    20</span>		&#34;math/big&#34;
<span id="L21" class="ln">    21</span>		&#34;strconv&#34;
<span id="L22" class="ln">    22</span>	)
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	<span class="comment">// Create a new Public-Private ECC-256 Keypair.</span>
<span id="L25" class="ln">    25</span>	func CreateKey(log chan string) ([]byte, *big.Int, *big.Int) {
<span id="L26" class="ln">    26</span>		priv, x, y, err := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
<span id="L27" class="ln">    27</span>		if err != nil {
<span id="L28" class="ln">    28</span>			log &lt;- &#34;Key Generation Error&#34;
<span id="L29" class="ln">    29</span>			return nil, nil, nil
<span id="L30" class="ln">    30</span>		}
<span id="L31" class="ln">    31</span>		return priv, x, y
<span id="L32" class="ln">    32</span>	}
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>	<span class="comment">// Convert public key to uncompressed 65-byte slice (2 32-bit integers with prefix 0x04)</span>
<span id="L35" class="ln">    35</span>	func MarshalPubkey(x, y *big.Int) []byte {
<span id="L36" class="ln">    36</span>		return elliptic.Marshal(elliptic.P256(), x, y)
<span id="L37" class="ln">    37</span>	}
<span id="L38" class="ln">    38</span>	
<span id="L39" class="ln">    39</span>	<span class="comment">// Convert 65-byte slice as created by MarshalPubkey() into an ECC-256 Public Key.</span>
<span id="L40" class="ln">    40</span>	func UnmarshalPubkey(data []byte) (x, y *big.Int) {
<span id="L41" class="ln">    41</span>		return elliptic.Unmarshal(elliptic.P256(), data)
<span id="L42" class="ln">    42</span>	}
<span id="L43" class="ln">    43</span>	
<span id="L44" class="ln">    44</span>	func GetCurve() elliptic.Curve {
<span id="L45" class="ln">    45</span>		return elliptic.P256()
<span id="L46" class="ln">    46</span>	}
<span id="L47" class="ln">    47</span>	
<span id="L48" class="ln">    48</span>	<span class="comment">// Convert ECC-256 Public Key to an EMP address (raw 25 bytes).</span>
<span id="L49" class="ln">    49</span>	func GetAddress(log chan string, x, y *big.Int) []byte {
<span id="L50" class="ln">    50</span>		pubKey := elliptic.Marshal(elliptic.P256(), x, y)
<span id="L51" class="ln">    51</span>		ripemd := ripemd160.New()
<span id="L52" class="ln">    52</span>	
<span id="L53" class="ln">    53</span>		sum := sha512.Sum384(pubKey)
<span id="L54" class="ln">    54</span>		sumslice := make([]byte, sha512.Size384, sha512.Size384)
<span id="L55" class="ln">    55</span>		for i := 0; i &lt; sha512.Size384; i++ {
<span id="L56" class="ln">    56</span>			sumslice[i] = sum[i]
<span id="L57" class="ln">    57</span>		}
<span id="L58" class="ln">    58</span>	
<span id="L59" class="ln">    59</span>		ripemd.Write(sumslice)
<span id="L60" class="ln">    60</span>		appender := ripemd.Sum(nil)
<span id="L61" class="ln">    61</span>		appender = appender[len(appender)-20:]
<span id="L62" class="ln">    62</span>		address := make([]byte, 1, 1)
<span id="L63" class="ln">    63</span>	
<span id="L64" class="ln">    64</span>		<span class="comment">// Version 0x01</span>
<span id="L65" class="ln">    65</span>		address[0] = 0x01
<span id="L66" class="ln">    66</span>		address = append(address, appender...)
<span id="L67" class="ln">    67</span>	
<span id="L68" class="ln">    68</span>		sum = sha512.Sum384(address)
<span id="L69" class="ln">    69</span>		sum = sha512.Sum384(sum[:])
<span id="L70" class="ln">    70</span>	
<span id="L71" class="ln">    71</span>		for i := 0; i &lt; 4; i++ {
<span id="L72" class="ln">    72</span>			address = append(address, sum[i])
<span id="L73" class="ln">    73</span>		}
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		return address
<span id="L76" class="ln">    76</span>	}
<span id="L77" class="ln">    77</span>	
<span id="L78" class="ln">    78</span>	<span class="comment">// Determine if address is valid (checksum is correct, correct length, and it starts with a 1-byte number).</span>
<span id="L79" class="ln">    79</span>	func ValidateAddress(addr []byte) bool {
<span id="L80" class="ln">    80</span>		if len(addr) != 25 {
<span id="L81" class="ln">    81</span>			return false
<span id="L82" class="ln">    82</span>		}
<span id="L83" class="ln">    83</span>		ripe := addr[:21]
<span id="L84" class="ln">    84</span>		sum := sha512.Sum384(ripe)
<span id="L85" class="ln">    85</span>		sum = sha512.Sum384(sum[:])
<span id="L86" class="ln">    86</span>	
<span id="L87" class="ln">    87</span>		for i := 0; i &lt; 4; i++ {
<span id="L88" class="ln">    88</span>			if sum[i] != addr[i+21] {
<span id="L89" class="ln">    89</span>				return false
<span id="L90" class="ln">    90</span>			}
<span id="L91" class="ln">    91</span>		}
<span id="L92" class="ln">    92</span>	
<span id="L93" class="ln">    93</span>		return true
<span id="L94" class="ln">    94</span>	}
<span id="L95" class="ln">    95</span>	
<span id="L96" class="ln">    96</span>	<span class="comment">// Converts 25-byte address to String representation.</span>
<span id="L97" class="ln">    97</span>	func AddressToString(addr []byte) string {
<span id="L98" class="ln">    98</span>		if !ValidateAddress(addr) {
<span id="L99" class="ln">    99</span>			return &#34;&#34;
<span id="L100" class="ln">   100</span>		}
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		return strconv.Itoa(int(addr[0])) + base64.StdEncoding.EncodeToString(addr[1:])
<span id="L103" class="ln">   103</span>	}
<span id="L104" class="ln">   104</span>	
<span id="L105" class="ln">   105</span>	<span class="comment">// Converts String representation to 25-byte address.</span>
<span id="L106" class="ln">   106</span>	func StringToAddress(addr string) []byte {
<span id="L107" class="ln">   107</span>		if len(addr) &lt; 2 {
<span id="L108" class="ln">   108</span>			return nil
<span id="L109" class="ln">   109</span>		}
<span id="L110" class="ln">   110</span>		data, err := base64.StdEncoding.DecodeString(addr[1:])
<span id="L111" class="ln">   111</span>		if err != nil {
<span id="L112" class="ln">   112</span>			return nil
<span id="L113" class="ln">   113</span>		}
<span id="L114" class="ln">   114</span>		version := make([]byte, 1, 1)
<span id="L115" class="ln">   115</span>		version[0] = byte(addr[0] - 48)
<span id="L116" class="ln">   116</span>		address := append(version, data...)
<span id="L117" class="ln">   117</span>		if !ValidateAddress(address) {
<span id="L118" class="ln">   118</span>			return nil
<span id="L119" class="ln">   119</span>		}
<span id="L120" class="ln">   120</span>		return address
<span id="L121" class="ln">   121</span>	}
</pre><p><a href="/src/pkg/emp/encryption/address.go?m=text">View as plain text</a></p>

<div id="footer">
Build version go1.2.<br>
Except as <a href="http://code.google.com/policies.html#restrictions">noted</a>,
the content of this page is licensed under the
Creative Commons Attribution 3.0 License,
and code is licensed under a <a href="/LICENSE">BSD license</a>.<br>
<a href="/doc/tos.html">Terms of Service</a> | 
<a href="http://www.google.com/intl/en/policies/privacy/">Privacy Policy</a>
</div>

</div><!-- .container -->
</div><!-- #page -->

<script type="text/javascript" src="/lib/godoc/jquery.js"></script>

<script type="text/javascript" src="/lib/godoc/godocs.js"></script>

</body>
</html>

