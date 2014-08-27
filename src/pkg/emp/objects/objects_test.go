<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/objects/objects_test.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/objects/objects_test.go</h1>




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
<span id="L15" class="ln">    15</span>		&#34;crypto/elliptic&#34;
<span id="L16" class="ln">    16</span>		&#34;fmt&#34;
<span id="L17" class="ln">    17</span>		&#34;net&#34;
<span id="L18" class="ln">    18</span>		&#34;emp/encryption&#34;
<span id="L19" class="ln">    19</span>		&#34;testing&#34;
<span id="L20" class="ln">    20</span>		&#34;time&#34;
<span id="L21" class="ln">    21</span>	)
<span id="L22" class="ln">    22</span>	
<span id="L23" class="ln">    23</span>	func TestVersion(t *testing.T) {
<span id="L24" class="ln">    24</span>		v := new(Version)
<span id="L25" class="ln">    25</span>		v.Version = uint16(1)
<span id="L26" class="ln">    26</span>		v.Timestamp = time.Unix(0, 0)
<span id="L27" class="ln">    27</span>		v.IpAddress = net.ParseIP(&#34;1.2.3.4&#34;)
<span id="L28" class="ln">    28</span>		v.Port = uint16(4444)
<span id="L29" class="ln">    29</span>		v.UserAgent = &#34;Hello World!&#34;
<span id="L30" class="ln">    30</span>	
<span id="L31" class="ln">    31</span>		verBytes := v.GetBytes()
<span id="L32" class="ln">    32</span>		if len(verBytes) != verLen+12 {
<span id="L33" class="ln">    33</span>			fmt.Println(&#34;Incorrect Byte Length: &#34;, verBytes)
<span id="L34" class="ln">    34</span>			t.Fail()
<span id="L35" class="ln">    35</span>		}
<span id="L36" class="ln">    36</span>	
<span id="L37" class="ln">    37</span>		v2 := new(Version)
<span id="L38" class="ln">    38</span>		err := v2.FromBytes(verBytes)
<span id="L39" class="ln">    39</span>		if err != nil {
<span id="L40" class="ln">    40</span>			fmt.Println(&#34;Error Decoding: &#34;, err)
<span id="L41" class="ln">    41</span>			t.FailNow()
<span id="L42" class="ln">    42</span>		}
<span id="L43" class="ln">    43</span>	
<span id="L44" class="ln">    44</span>		if v2.Version != 1 || v2.Timestamp != time.Unix(0, 0) || v2.IpAddress.String() != &#34;1.2.3.4&#34; || v2.Port != 4444 || v2.UserAgent != &#34;Hello World!&#34; {
<span id="L45" class="ln">    45</span>			fmt.Println(&#34;Incorrect decoded version: &#34;, v2)
<span id="L46" class="ln">    46</span>			t.Fail()
<span id="L47" class="ln">    47</span>		}
<span id="L48" class="ln">    48</span>	}
<span id="L49" class="ln">    49</span>	
<span id="L50" class="ln">    50</span>	func TestNodes(t *testing.T) {
<span id="L51" class="ln">    51</span>		n := new(NodeList)
<span id="L52" class="ln">    52</span>		n.Nodes = make(map[string]Node)
<span id="L53" class="ln">    53</span>		n2 := new(NodeList)
<span id="L54" class="ln">    54</span>		n2.Nodes = make(map[string]Node)
<span id="L55" class="ln">    55</span>		node1 := new(Node)
<span id="L56" class="ln">    56</span>		node2 := new(Node)
<span id="L57" class="ln">    57</span>	
<span id="L58" class="ln">    58</span>		node1.IP = net.ParseIP(&#34;1.2.3.4&#34;)
<span id="L59" class="ln">    59</span>		node1.Port = uint16(4444)
<span id="L60" class="ln">    60</span>		node1.LastSeen = time.Now().Round(time.Second)
<span id="L61" class="ln">    61</span>	
<span id="L62" class="ln">    62</span>		node2.IP = net.ParseIP(&#34;5.6.7.8&#34;)
<span id="L63" class="ln">    63</span>		node2.Port = uint16(5555)
<span id="L64" class="ln">    64</span>		node2.LastSeen = time.Now().Round(time.Second)
<span id="L65" class="ln">    65</span>	
<span id="L66" class="ln">    66</span>		n.Nodes[node1.String()] = *node1
<span id="L67" class="ln">    67</span>		n.Nodes[node2.String()] = *node2
<span id="L68" class="ln">    68</span>	
<span id="L69" class="ln">    69</span>		nBytes := n.GetBytes()
<span id="L70" class="ln">    70</span>		if len(nBytes) != 2*nodeLen {
<span id="L71" class="ln">    71</span>			fmt.Println(&#34;Byte Lengh Incorrect: &#34;, nBytes)
<span id="L72" class="ln">    72</span>			t.FailNow()
<span id="L73" class="ln">    73</span>		}
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		err := n2.FromBytes(nBytes)
<span id="L76" class="ln">    76</span>		if err != nil {
<span id="L77" class="ln">    77</span>			fmt.Println(&#34;Error Decoding: &#34;, err)
<span id="L78" class="ln">    78</span>		}
<span id="L79" class="ln">    79</span>	
<span id="L80" class="ln">    80</span>		for key, _ := range n.Nodes {
<span id="L81" class="ln">    81</span>			if n2.Nodes[key].IP.String() != n.Nodes[key].IP.String() || n2.Nodes[key].Port != n.Nodes[key].Port {
<span id="L82" class="ln">    82</span>				fmt.Println(&#34;Nodes don&#39;t match!&#34;, n2.Nodes)
<span id="L83" class="ln">    83</span>				t.FailNow()
<span id="L84" class="ln">    84</span>			}
<span id="L85" class="ln">    85</span>		}
<span id="L86" class="ln">    86</span>	}
<span id="L87" class="ln">    87</span>	
<span id="L88" class="ln">    88</span>	func TestObj(t *testing.T) {
<span id="L89" class="ln">    89</span>		o := new(Obj)
<span id="L90" class="ln">    90</span>		o.HashList = make([]Hash, 2, 2)
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>		o.HashList[0] = MakeHash([]byte{&#39;a&#39;, &#39;b&#39;, &#39;c&#39;, &#39;d&#39;})
<span id="L93" class="ln">    93</span>		o.HashList[1] = MakeHash([]byte{&#39;e&#39;, &#39;f&#39;, &#39;g&#39;, &#39;h&#39;})
<span id="L94" class="ln">    94</span>	
<span id="L95" class="ln">    95</span>		o2 := new(Obj)
<span id="L96" class="ln">    96</span>	
<span id="L97" class="ln">    97</span>		oBytes := o.GetBytes()
<span id="L98" class="ln">    98</span>		if len(oBytes) != 2*hashLen {
<span id="L99" class="ln">    99</span>			fmt.Println(&#34;Incorrect Obj Length! &#34;, oBytes)
<span id="L100" class="ln">   100</span>			t.FailNow()
<span id="L101" class="ln">   101</span>		}
<span id="L102" class="ln">   102</span>	
<span id="L103" class="ln">   103</span>		err := o2.FromBytes(oBytes)
<span id="L104" class="ln">   104</span>		if err != nil {
<span id="L105" class="ln">   105</span>			fmt.Println(&#34;Error while decoding obj: &#34;, err)
<span id="L106" class="ln">   106</span>			t.FailNow()
<span id="L107" class="ln">   107</span>		}
<span id="L108" class="ln">   108</span>	
<span id="L109" class="ln">   109</span>		if string(o2.HashList[0].GetBytes()) != string(o.HashList[0].GetBytes()) || string(o2.HashList[1].GetBytes()) != string(o.HashList[1].GetBytes()) {
<span id="L110" class="ln">   110</span>			fmt.Println(&#34;Incorrect decoding of obj: &#34;, o2)
<span id="L111" class="ln">   111</span>			t.Fail()
<span id="L112" class="ln">   112</span>		}
<span id="L113" class="ln">   113</span>	}
<span id="L114" class="ln">   114</span>	
<span id="L115" class="ln">   115</span>	func TestPubkey(t *testing.T) {
<span id="L116" class="ln">   116</span>		p := new(EncryptedPubkey)
<span id="L117" class="ln">   117</span>		var err error
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>		address := make([]byte, 25, 25)
<span id="L120" class="ln">   120</span>		pubkey := [65]byte{&#39;a&#39;}
<span id="L121" class="ln">   121</span>		p.AddrHash = MakeHash(address)
<span id="L122" class="ln">   122</span>		p.IV, p.Payload, err = encryption.SymmetricEncrypt(address, string(pubkey[:]))
<span id="L123" class="ln">   123</span>		if err != nil {
<span id="L124" class="ln">   124</span>			fmt.Println(&#34;Could not encrypt pubkey: &#34;, err)
<span id="L125" class="ln">   125</span>			t.FailNow()
<span id="L126" class="ln">   126</span>		}
<span id="L127" class="ln">   127</span>	
<span id="L128" class="ln">   128</span>		pBytes := p.GetBytes()
<span id="L129" class="ln">   129</span>		if len(pBytes) != 144 {
<span id="L130" class="ln">   130</span>			fmt.Println(&#34;Incorrect length for pubkey: &#34;, pBytes)
<span id="L131" class="ln">   131</span>			t.FailNow()
<span id="L132" class="ln">   132</span>		}
<span id="L133" class="ln">   133</span>	
<span id="L134" class="ln">   134</span>		pubkey2 := new(EncryptedPubkey)
<span id="L135" class="ln">   135</span>		err = pubkey2.FromBytes(pBytes)
<span id="L136" class="ln">   136</span>		if err != nil {
<span id="L137" class="ln">   137</span>			fmt.Println(&#34;Error decoding pubkey: &#34;, err)
<span id="L138" class="ln">   138</span>			t.Fail()
<span id="L139" class="ln">   139</span>		}
<span id="L140" class="ln">   140</span>		if string(pubkey2.AddrHash.GetBytes()) != string(p.AddrHash.GetBytes()) {
<span id="L141" class="ln">   141</span>			fmt.Println(&#34;Incorrect Address Hash: &#34;, pubkey2.AddrHash)
<span id="L142" class="ln">   142</span>			t.FailNow()
<span id="L143" class="ln">   143</span>		}
<span id="L144" class="ln">   144</span>	
<span id="L145" class="ln">   145</span>		pubkeyTest := encryption.SymmetricDecrypt(pubkey2.IV, address, pubkey2.Payload)
<span id="L146" class="ln">   146</span>		if string(pubkeyTest[:65]) != string(pubkey[:]) {
<span id="L147" class="ln">   147</span>			fmt.Println(&#34;Incorrect public key decryption: &#34;, pubkeyTest)
<span id="L148" class="ln">   148</span>			t.Fail()
<span id="L149" class="ln">   149</span>		}
<span id="L150" class="ln">   150</span>	}
<span id="L151" class="ln">   151</span>	
<span id="L152" class="ln">   152</span>	func TestPurge(t *testing.T) {
<span id="L153" class="ln">   153</span>		p := new(Purge)
<span id="L154" class="ln">   154</span>		p.Txid = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
<span id="L155" class="ln">   155</span>		pBytes := p.GetBytes()
<span id="L156" class="ln">   156</span>		if len(pBytes) != 16 {
<span id="L157" class="ln">   157</span>			fmt.Println(&#34;Error encoding purge: &#34;, pBytes)
<span id="L158" class="ln">   158</span>			t.FailNow()
<span id="L159" class="ln">   159</span>		}
<span id="L160" class="ln">   160</span>	
<span id="L161" class="ln">   161</span>		p2 := new(Purge)
<span id="L162" class="ln">   162</span>		p2.FromBytes(pBytes)
<span id="L163" class="ln">   163</span>		if string(p2.Txid[:]) != string(p.Txid[:]) {
<span id="L164" class="ln">   164</span>			fmt.Println(&#34;Incorrect decoding: &#34;, p2.Txid)
<span id="L165" class="ln">   165</span>			t.Fail()
<span id="L166" class="ln">   166</span>		}
<span id="L167" class="ln">   167</span>	}
<span id="L168" class="ln">   168</span>	
<span id="L169" class="ln">   169</span>	func TestMessage(t *testing.T) {
<span id="L170" class="ln">   170</span>		log := make(chan string, 100)
<span id="L171" class="ln">   171</span>		priv, x, y := encryption.CreateKey(log)
<span id="L172" class="ln">   172</span>		pub := elliptic.Marshal(elliptic.P256(), x, y)
<span id="L173" class="ln">   173</span>		address := encryption.GetAddress(log, x, y)
<span id="L174" class="ln">   174</span>	
<span id="L175" class="ln">   175</span>		msg := new(Message)
<span id="L176" class="ln">   176</span>		msg.AddrHash = MakeHash(address)
<span id="L177" class="ln">   177</span>		msg.TxidHash = MakeHash([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
<span id="L178" class="ln">   178</span>		msg.Timestamp = time.Now().Round(time.Second)
<span id="L179" class="ln">   179</span>		msg.Content = *encryption.Encrypt(log, pub, &#34;Hello World!&#34;)
<span id="L180" class="ln">   180</span>	
<span id="L181" class="ln">   181</span>		mBytes := msg.GetBytes()
<span id="L182" class="ln">   182</span>		if mBytes == nil {
<span id="L183" class="ln">   183</span>			fmt.Println(&#34;Error Encoding Message!&#34;)
<span id="L184" class="ln">   184</span>			t.FailNow()
<span id="L185" class="ln">   185</span>		}
<span id="L186" class="ln">   186</span>	
<span id="L187" class="ln">   187</span>		msg2 := new(Message)
<span id="L188" class="ln">   188</span>		msg2.FromBytes(mBytes)
<span id="L189" class="ln">   189</span>		if string(msg2.AddrHash.GetBytes()) != string(msg.AddrHash.GetBytes()) || string(msg2.TxidHash.GetBytes()) != string(msg.TxidHash.GetBytes()) || msg2.Timestamp.Unix() != msg.Timestamp.Unix() {
<span id="L190" class="ln">   190</span>			fmt.Println(&#34;Message Header incorrect: &#34;, msg2)
<span id="L191" class="ln">   191</span>			t.FailNow()
<span id="L192" class="ln">   192</span>		}
<span id="L193" class="ln">   193</span>	
<span id="L194" class="ln">   194</span>		if string(encryption.Decrypt(log, priv, &amp;msg.Content)[:12]) != &#34;Hello World!&#34; {
<span id="L195" class="ln">   195</span>			fmt.Println(&#34;Message content incorrect: &#34;, string(encryption.Decrypt(log, priv, &amp;msg.Content)[:12]))
<span id="L196" class="ln">   196</span>			t.Fail()
<span id="L197" class="ln">   197</span>		}
<span id="L198" class="ln">   198</span>	}
</pre><p><a href="/src/pkg/emp/objects/objects_test.go?m=text">View as plain text</a></p>

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

