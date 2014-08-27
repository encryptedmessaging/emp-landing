<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/objects/localmsg.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/objects/localmsg.go</h1>




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
<span id="L12" class="ln">    12</span>	package objects
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;bytes&#34;
<span id="L16" class="ln">    16</span>		&#34;encoding/binary&#34;
<span id="L17" class="ln">    17</span>		&#34;errors&#34;
<span id="L18" class="ln">    18</span>		&#34;emp/encryption&#34;
<span id="L19" class="ln">    19</span>		&#34;time&#34;
<span id="L20" class="ln">    20</span>	)
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	type MetaMessage struct {
<span id="L23" class="ln">    23</span>		TxidHash  Hash      `json:&#34;txid_hash&#34;` <span class="comment">// Hash of random identifier</span>
<span id="L24" class="ln">    24</span>		Timestamp time.Time `json:&#34;sent&#34;`      <span class="comment">// Time message was sent</span>
<span id="L25" class="ln">    25</span>		Purged    bool      `json:&#34;read&#34;`      <span class="comment">// Whether purge token has been received</span>
<span id="L26" class="ln">    26</span>		Sender    string    `json:&#34;sender&#34;`    <span class="comment">// String representation of sender&#39;s address, if available.</span>
<span id="L27" class="ln">    27</span>		Recipient string    `json:&#34;recipient&#34;` <span class="comment">// String representation of recipient&#39;s address, if available.</span>
<span id="L28" class="ln">    28</span>	}
<span id="L29" class="ln">    29</span>	
<span id="L30" class="ln">    30</span>	type FullMessage struct {
<span id="L31" class="ln">    31</span>		MetaMessage MetaMessage                  `json:&#34;info&#34;`
<span id="L32" class="ln">    32</span>		Decrypted   *DecryptedMessage            `json:&#34;decrypted&#34;`
<span id="L33" class="ln">    33</span>		Encrypted   *encryption.EncryptedMessage `json:&#34;encrypted&#34;`
<span id="L34" class="ln">    34</span>	}
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>	type DecryptedMessage struct {
<span id="L37" class="ln">    37</span>		Txid      [16]byte <span class="comment">// Randomly generated identifier and purge token.</span>
<span id="L38" class="ln">    38</span>		Pubkey    [65]byte <span class="comment">// Sender&#39;s 65-byte Public Key</span>
<span id="L39" class="ln">    39</span>		Subject   string   <span class="comment">// Human-readable subject of this message</span>
<span id="L40" class="ln">    40</span>		MimeType  string   <span class="comment">// Mime-Type of Content</span>
<span id="L41" class="ln">    41</span>		Length    uint32   <span class="comment">// Length of Content in bytes</span>
<span id="L42" class="ln">    42</span>		Content   string   <span class="comment">// Content of message, could be any data.</span>
<span id="L43" class="ln">    43</span>		Signature [65]byte <span class="comment">// Sender&#39;s Signature of entire message.</span>
<span id="L44" class="ln">    44</span>	}
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>	func (d *DecryptedMessage) GetBytes() []byte {
<span id="L47" class="ln">    47</span>		if d == nil {
<span id="L48" class="ln">    48</span>			return nil
<span id="L49" class="ln">    49</span>		}
<span id="L50" class="ln">    50</span>	
<span id="L51" class="ln">    51</span>		ret := append(d.Txid[:], d.Pubkey[:]...)
<span id="L52" class="ln">    52</span>		ret = append(ret, d.Subject...)
<span id="L53" class="ln">    53</span>		ret = append(ret, 0)
<span id="L54" class="ln">    54</span>		ret = append(ret, d.MimeType...)
<span id="L55" class="ln">    55</span>		ret = append(ret, 0)
<span id="L56" class="ln">    56</span>	
<span id="L57" class="ln">    57</span>		leng := make([]byte, 4, 4)
<span id="L58" class="ln">    58</span>	
<span id="L59" class="ln">    59</span>		binary.BigEndian.PutUint32(leng, d.Length)
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>		ret = append(ret, leng...)
<span id="L62" class="ln">    62</span>		ret = append(ret, d.Content...)
<span id="L63" class="ln">    63</span>		ret = append(ret, d.Signature[:]...)
<span id="L64" class="ln">    64</span>	
<span id="L65" class="ln">    65</span>		return ret
<span id="L66" class="ln">    66</span>	}
<span id="L67" class="ln">    67</span>	
<span id="L68" class="ln">    68</span>	func (d *DecryptedMessage) FromBytes(data []byte) error {
<span id="L69" class="ln">    69</span>		if d == nil {
<span id="L70" class="ln">    70</span>			return errors.New(&#34;Can&#39;t fill nil object!&#34;)
<span id="L71" class="ln">    71</span>		}
<span id="L72" class="ln">    72</span>	
<span id="L73" class="ln">    73</span>		var err error
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		buf := bytes.NewBuffer(data)
<span id="L76" class="ln">    76</span>		copy(d.Txid[:], buf.Next(16))
<span id="L77" class="ln">    77</span>		copy(d.Pubkey[:], buf.Next(65))
<span id="L78" class="ln">    78</span>		d.Subject, err = buf.ReadString(0)
<span id="L79" class="ln">    79</span>		if err != nil {
<span id="L80" class="ln">    80</span>			return err
<span id="L81" class="ln">    81</span>		}
<span id="L82" class="ln">    82</span>		d.Subject = d.Subject[:len(d.Subject)-1]
<span id="L83" class="ln">    83</span>		d.MimeType, err = buf.ReadString(0)
<span id="L84" class="ln">    84</span>		if err != nil {
<span id="L85" class="ln">    85</span>			return err
<span id="L86" class="ln">    86</span>		}
<span id="L87" class="ln">    87</span>		d.MimeType = d.MimeType[:len(d.MimeType)-1]
<span id="L88" class="ln">    88</span>	
<span id="L89" class="ln">    89</span>		d.Length = binary.BigEndian.Uint32(buf.Next(4))
<span id="L90" class="ln">    90</span>	
<span id="L91" class="ln">    91</span>		d.Content = string(buf.Next(int(d.Length)))
<span id="L92" class="ln">    92</span>		copy(d.Signature[:], buf.Next(65))
<span id="L93" class="ln">    93</span>	
<span id="L94" class="ln">    94</span>		return nil
<span id="L95" class="ln">    95</span>	}
</pre><p><a href="/src/pkg/emp/objects/localmsg.go?m=text">View as plain text</a></p>

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

