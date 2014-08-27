<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/encryption/message.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/encryption/message.go</h1>




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
<span id="L12" class="ln">    12</span>	<span class="comment">// Package encryption wraps around Go&#39;s native crypto library to provide </span>
<span id="L13" class="ln">    13</span>	<span class="comment">// ECIES and AES-256 encryption for EMP Basic and Published Messages.</span>
<span id="L14" class="ln">    14</span>	package encryption
<span id="L15" class="ln">    15</span>	
<span id="L16" class="ln">    16</span>	import (
<span id="L17" class="ln">    17</span>		&#34;crypto/elliptic&#34;
<span id="L18" class="ln">    18</span>		&#34;crypto/hmac&#34;
<span id="L19" class="ln">    19</span>		&#34;crypto/sha256&#34;
<span id="L20" class="ln">    20</span>		&#34;crypto/sha512&#34;
<span id="L21" class="ln">    21</span>	)
<span id="L22" class="ln">    22</span>	
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	<span class="comment">// Encrypt plainText into an Encrypted Message using the given public key.</span>
<span id="L25" class="ln">    25</span>	func Encrypt(log chan string, dest_pubkey []byte, plainText string) *EncryptedMessage {
<span id="L26" class="ln">    26</span>		<span class="comment">// Generate New Public/Private Key Pair</span>
<span id="L27" class="ln">    27</span>		D1, X1, Y1 := CreateKey(log)
<span id="L28" class="ln">    28</span>		<span class="comment">// Unmarshal the Destination&#39;s Pubkey</span>
<span id="L29" class="ln">    29</span>		X2, Y2 := elliptic.Unmarshal(elliptic.P256(), dest_pubkey)
<span id="L30" class="ln">    30</span>	
<span id="L31" class="ln">    31</span>		<span class="comment">// Point Multiply to get new Pubkey</span>
<span id="L32" class="ln">    32</span>		PubX, PubY := elliptic.P256().ScalarMult(X2, Y2, D1)
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>		<span class="comment">// Generate Pubkey hashes</span>
<span id="L35" class="ln">    35</span>		PubHash := sha512.Sum512(elliptic.Marshal(elliptic.P256(), PubX, PubY))
<span id="L36" class="ln">    36</span>		PubHash_E := PubHash[:32]
<span id="L37" class="ln">    37</span>		PubHash_M := PubHash[32:64]
<span id="L38" class="ln">    38</span>	
<span id="L39" class="ln">    39</span>		IV, cipherText, _ := SymmetricEncrypt(PubHash_E, plainText)
<span id="L40" class="ln">    40</span>	
<span id="L41" class="ln">    41</span>		<span class="comment">// Generate HMAC</span>
<span id="L42" class="ln">    42</span>		mac := hmac.New(sha256.New, PubHash_M)
<span id="L43" class="ln">    43</span>		mac.Write(cipherText)
<span id="L44" class="ln">    44</span>		HMAC := mac.Sum(nil)
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>		ret := new(EncryptedMessage)
<span id="L47" class="ln">    47</span>		copy(ret.IV[:], IV[:])
<span id="L48" class="ln">    48</span>		copy(ret.PublicKey[:], elliptic.Marshal(elliptic.P256(), X1, Y1))
<span id="L49" class="ln">    49</span>		ret.CipherText = cipherText
<span id="L50" class="ln">    50</span>		copy(ret.HMAC[:], HMAC)
<span id="L51" class="ln">    51</span>	
<span id="L52" class="ln">    52</span>		return ret
<span id="L53" class="ln">    53</span>	}
<span id="L54" class="ln">    54</span>	
<span id="L55" class="ln">    55</span>	<span class="comment">// Encrypt plainText into an Encrypted Published Message using the given private key.</span>
<span id="L56" class="ln">    56</span>	func EncryptPub(log chan string, src_privkey []byte, plainText string) *EncryptedMessage {
<span id="L57" class="ln">    57</span>		<span class="comment">// Generate New Public/Private Key Pair</span>
<span id="L58" class="ln">    58</span>		D1, X1, Y1 := CreateKey(log)
<span id="L59" class="ln">    59</span>	
<span id="L60" class="ln">    60</span>		<span class="comment">// Point Multiply to get new Pubkey</span>
<span id="L61" class="ln">    61</span>		PubX, PubY := elliptic.P256().ScalarMult(X1, Y1, src_privkey)
<span id="L62" class="ln">    62</span>	
<span id="L63" class="ln">    63</span>		<span class="comment">// Generate Pubkey hashes</span>
<span id="L64" class="ln">    64</span>		PubHash := sha512.Sum512(elliptic.Marshal(elliptic.P256(), PubX, PubY))
<span id="L65" class="ln">    65</span>		PubHash_E := PubHash[:32]
<span id="L66" class="ln">    66</span>		PubHash_M := PubHash[32:64]
<span id="L67" class="ln">    67</span>	
<span id="L68" class="ln">    68</span>		IV, cipherText, _ := SymmetricEncrypt(PubHash_E, plainText)
<span id="L69" class="ln">    69</span>	
<span id="L70" class="ln">    70</span>		<span class="comment">// Generate HMAC</span>
<span id="L71" class="ln">    71</span>		mac := hmac.New(sha256.New, PubHash_M)
<span id="L72" class="ln">    72</span>		mac.Write(cipherText)
<span id="L73" class="ln">    73</span>		HMAC := mac.Sum(nil)
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		ret := new(EncryptedMessage)
<span id="L76" class="ln">    76</span>		copy(ret.IV[:], IV[:])
<span id="L77" class="ln">    77</span>		copy(ret.PublicKey[:32], D1)
<span id="L78" class="ln">    78</span>		ret.CipherText = cipherText
<span id="L79" class="ln">    79</span>		copy(ret.HMAC[:], HMAC)
<span id="L80" class="ln">    80</span>	
<span id="L81" class="ln">    81</span>		return ret
<span id="L82" class="ln">    82</span>	}
<span id="L83" class="ln">    83</span>	
<span id="L84" class="ln">    84</span>	<span class="comment">// checkMAC returns true if messageMAC is a valid HMAC tag for message.</span>
<span id="L85" class="ln">    85</span>	func checkMAC(message, messageMAC, key []byte) bool {
<span id="L86" class="ln">    86</span>		mac := hmac.New(sha256.New, key)
<span id="L87" class="ln">    87</span>		mac.Write(message)
<span id="L88" class="ln">    88</span>		expectedMAC := mac.Sum(nil)
<span id="L89" class="ln">    89</span>		return hmac.Equal(messageMAC, expectedMAC)
<span id="L90" class="ln">    90</span>	}
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>	<span class="comment">// Decrypt a given Encrypted Message using the given private key. </span>
<span id="L93" class="ln">    93</span>	<span class="comment">// &lt;Nil&gt; is returned if the key fails the HMAC-SHA256 test.</span>
<span id="L94" class="ln">    94</span>	func Decrypt(log chan string, privKey []byte, encrypted *EncryptedMessage) []byte {
<span id="L95" class="ln">    95</span>		if encrypted == nil || privKey == nil || log == nil {
<span id="L96" class="ln">    96</span>			return nil
<span id="L97" class="ln">    97</span>		}
<span id="L98" class="ln">    98</span>	
<span id="L99" class="ln">    99</span>		<span class="comment">// Unmarshal the Sender&#39;s Pubkey</span>
<span id="L100" class="ln">   100</span>		X2, Y2 := elliptic.Unmarshal(elliptic.P256(), encrypted.PublicKey[:])
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		<span class="comment">// Point Multiply to get the new Pubkey</span>
<span id="L103" class="ln">   103</span>		PubX, PubY := elliptic.P256().ScalarMult(X2, Y2, privKey)
<span id="L104" class="ln">   104</span>	
<span id="L105" class="ln">   105</span>		<span class="comment">// Generate Pubkey hashes</span>
<span id="L106" class="ln">   106</span>		PubHash := sha512.Sum512(elliptic.Marshal(elliptic.P256(), PubX, PubY))
<span id="L107" class="ln">   107</span>		PubHash_E := PubHash[:32]
<span id="L108" class="ln">   108</span>		PubHash_M := PubHash[32:64]
<span id="L109" class="ln">   109</span>	
<span id="L110" class="ln">   110</span>		<span class="comment">// Check HMAC</span>
<span id="L111" class="ln">   111</span>		if !checkMAC(encrypted.CipherText[:], encrypted.HMAC[:], PubHash_M) {
<span id="L112" class="ln">   112</span>			log &lt;- &#34;Invalid HMAC Message&#34;
<span id="L113" class="ln">   113</span>			return nil
<span id="L114" class="ln">   114</span>		}
<span id="L115" class="ln">   115</span>	
<span id="L116" class="ln">   116</span>		return SymmetricDecrypt(encrypted.IV, PubHash_E, encrypted.CipherText)
<span id="L117" class="ln">   117</span>	}
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>	<span class="comment">// Decrypt the given Published Message using the given Pubkey. </span>
<span id="L120" class="ln">   120</span>	<span class="comment">// &lt;Nil&gt; is returned if the HMAC-SHA256 test fails.</span>
<span id="L121" class="ln">   121</span>	func DecryptPub(log chan string, pubkey []byte, encrypted *EncryptedMessage) []byte {
<span id="L122" class="ln">   122</span>		if encrypted == nil || pubkey == nil || log == nil {
<span id="L123" class="ln">   123</span>			return nil
<span id="L124" class="ln">   124</span>		}
<span id="L125" class="ln">   125</span>	
<span id="L126" class="ln">   126</span>		<span class="comment">// Unmarshal the Sender&#39;s Pubkey</span>
<span id="L127" class="ln">   127</span>		X2, Y2 := elliptic.Unmarshal(elliptic.P256(), pubkey)
<span id="L128" class="ln">   128</span>	
<span id="L129" class="ln">   129</span>		<span class="comment">// Point Multiply to get the new Pubkey</span>
<span id="L130" class="ln">   130</span>		PubX, PubY := elliptic.P256().ScalarMult(X2, Y2, encrypted.PublicKey[:32])
<span id="L131" class="ln">   131</span>	
<span id="L132" class="ln">   132</span>		<span class="comment">// Generate Pubkey hashes</span>
<span id="L133" class="ln">   133</span>		PubHash := sha512.Sum512(elliptic.Marshal(elliptic.P256(), PubX, PubY))
<span id="L134" class="ln">   134</span>		PubHash_E := PubHash[:32]
<span id="L135" class="ln">   135</span>		PubHash_M := PubHash[32:64]
<span id="L136" class="ln">   136</span>	
<span id="L137" class="ln">   137</span>		<span class="comment">// Check HMAC</span>
<span id="L138" class="ln">   138</span>		if !checkMAC(encrypted.CipherText[:], encrypted.HMAC[:], PubHash_M) {
<span id="L139" class="ln">   139</span>			log &lt;- &#34;Invalid HMAC Message&#34;
<span id="L140" class="ln">   140</span>			return nil
<span id="L141" class="ln">   141</span>		}
<span id="L142" class="ln">   142</span>	
<span id="L143" class="ln">   143</span>		return SymmetricDecrypt(encrypted.IV, PubHash_E, encrypted.CipherText)
<span id="L144" class="ln">   144</span>	}</pre><p><a href="/src/emp/encryption/message.go?m=text">View as plain text</a></p>

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

