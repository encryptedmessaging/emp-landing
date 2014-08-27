<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/db/objects.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/db/objects.go</h1>




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
<span id="L12" class="ln">    12</span>	package db
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;crypto/sha512&#34;
<span id="L16" class="ln">    16</span>		&#34;fmt&#34;
<span id="L17" class="ln">    17</span>		&#34;emp/objects&#34;
<span id="L18" class="ln">    18</span>		&#34;time&#34;
<span id="L19" class="ln">    19</span>	)
<span id="L20" class="ln">    20</span>	
<span id="L21" class="ln">    21</span>	<span class="comment">// Remove all basic and published messages older than (Current Time - duration).</span>
<span id="L22" class="ln">    22</span>	func SweepMessages(duration time.Duration) error {
<span id="L23" class="ln">    23</span>		mutex.Lock()
<span id="L24" class="ln">    24</span>		defer mutex.Unlock()
<span id="L25" class="ln">    25</span>	
<span id="L26" class="ln">    26</span>		deadline := time.Now().Add(-duration).Unix()
<span id="L27" class="ln">    27</span>	
<span id="L28" class="ln">    28</span>		return dbConn.Exec(&#34;DELETE FROM msg WHERE timestamp &lt;= ?&#34;, deadline)
<span id="L29" class="ln">    29</span>	}
<span id="L30" class="ln">    30</span>	
<span id="L31" class="ln">    31</span>	<span class="comment">// Add Encrypted public key to database and hash list.</span>
<span id="L32" class="ln">    32</span>	func AddPubkey(log chan string, pubkey objects.EncryptedPubkey) error {
<span id="L33" class="ln">    33</span>		mutex.Lock()
<span id="L34" class="ln">    34</span>		defer mutex.Unlock()
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>		hash := pubkey.AddrHash.GetBytes()
<span id="L37" class="ln">    37</span>		payload := append(pubkey.IV[:], pubkey.Payload...)
<span id="L38" class="ln">    38</span>	
<span id="L39" class="ln">    39</span>		if hashList == nil || dbConn == nil {
<span id="L40" class="ln">    40</span>			return DBError(EUNINIT)
<span id="L41" class="ln">    41</span>		}
<span id="L42" class="ln">    42</span>		if Contains(pubkey.AddrHash) == PUBKEY {
<span id="L43" class="ln">    43</span>			return nil
<span id="L44" class="ln">    44</span>		}
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>		err := dbConn.Exec(&#34;INSERT INTO pubkey VALUES (?, ?)&#34;, hash, payload)
<span id="L47" class="ln">    47</span>		if err != nil {
<span id="L48" class="ln">    48</span>			log &lt;- fmt.Sprintf(&#34;Error inserting pubkey into db... %s&#34;, err)
<span id="L49" class="ln">    49</span>			return err
<span id="L50" class="ln">    50</span>		}
<span id="L51" class="ln">    51</span>	
<span id="L52" class="ln">    52</span>		Add(pubkey.AddrHash, PUBKEY)
<span id="L53" class="ln">    53</span>		return nil
<span id="L54" class="ln">    54</span>	}
<span id="L55" class="ln">    55</span>	
<span id="L56" class="ln">    56</span>	<span class="comment">// Get Encrypted Public Key from database.</span>
<span id="L57" class="ln">    57</span>	func GetPubkey(log chan string, addrHash objects.Hash) *objects.EncryptedPubkey {
<span id="L58" class="ln">    58</span>		mutex.Lock()
<span id="L59" class="ln">    59</span>		defer mutex.Unlock()
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>		hash := addrHash.GetBytes()
<span id="L62" class="ln">    62</span>	
<span id="L63" class="ln">    63</span>		if hashList == nil || dbConn == nil {
<span id="L64" class="ln">    64</span>			return nil
<span id="L65" class="ln">    65</span>		}
<span id="L66" class="ln">    66</span>		if hashList[string(hash)] != PUBKEY {
<span id="L67" class="ln">    67</span>			return nil
<span id="L68" class="ln">    68</span>		}
<span id="L69" class="ln">    69</span>	
<span id="L70" class="ln">    70</span>		for s, err := dbConn.Query(&#34;SELECT payload FROM pubkey WHERE hash=?&#34;, hash); err == nil; err = s.Next() {
<span id="L71" class="ln">    71</span>			var payload []byte
<span id="L72" class="ln">    72</span>			s.Scan(&amp;payload) <span class="comment">// Assigns 1st column to rowid, the rest to row</span>
<span id="L73" class="ln">    73</span>			pub := new(objects.EncryptedPubkey)
<span id="L74" class="ln">    74</span>			pub.AddrHash = addrHash
<span id="L75" class="ln">    75</span>			copy(pub.IV[:], payload[:16])
<span id="L76" class="ln">    76</span>			pub.Payload = payload[16:]
<span id="L77" class="ln">    77</span>			return pub
<span id="L78" class="ln">    78</span>		}
<span id="L79" class="ln">    79</span>		<span class="comment">// Not Found</span>
<span id="L80" class="ln">    80</span>		return nil
<span id="L81" class="ln">    81</span>	}
<span id="L82" class="ln">    82</span>	
<span id="L83" class="ln">    83</span>	<span class="comment">// Add Purge Token to database, and remove corresponding message if necessary.</span>
<span id="L84" class="ln">    84</span>	func AddPurge(log chan string, p objects.Purge) error {
<span id="L85" class="ln">    85</span>		mutex.Lock()
<span id="L86" class="ln">    86</span>		defer mutex.Unlock()
<span id="L87" class="ln">    87</span>	
<span id="L88" class="ln">    88</span>		txid := p.GetBytes()
<span id="L89" class="ln">    89</span>		hashArr := sha512.Sum384(txid)
<span id="L90" class="ln">    90</span>		hash := hashArr[:]
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>		if hashList == nil || dbConn == nil {
<span id="L93" class="ln">    93</span>			return DBError(EUNINIT)
<span id="L94" class="ln">    94</span>		}
<span id="L95" class="ln">    95</span>		hashObj := new(objects.Hash)
<span id="L96" class="ln">    96</span>		hashObj.FromBytes(hash)
<span id="L97" class="ln">    97</span>	
<span id="L98" class="ln">    98</span>		if Contains(*hashObj) == PURGE {
<span id="L99" class="ln">    99</span>			return nil
<span id="L100" class="ln">   100</span>		}
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		err := dbConn.Exec(&#34;INSERT INTO purge VALUES (?, ?)&#34;, hash, txid)
<span id="L103" class="ln">   103</span>		if err != nil {
<span id="L104" class="ln">   104</span>			log &lt;- fmt.Sprintf(&#34;Error inserting purge into db... %s&#34;, err)
<span id="L105" class="ln">   105</span>			return err
<span id="L106" class="ln">   106</span>		}
<span id="L107" class="ln">   107</span>	
<span id="L108" class="ln">   108</span>		Add(*hashObj, PURGE)
<span id="L109" class="ln">   109</span>		return nil
<span id="L110" class="ln">   110</span>	}
<span id="L111" class="ln">   111</span>	
<span id="L112" class="ln">   112</span>	<span class="comment">// Get purge token from the database.</span>
<span id="L113" class="ln">   113</span>	func GetPurge(log chan string, txidHash objects.Hash) *objects.Purge {
<span id="L114" class="ln">   114</span>		mutex.Lock()
<span id="L115" class="ln">   115</span>		defer mutex.Unlock()
<span id="L116" class="ln">   116</span>	
<span id="L117" class="ln">   117</span>		hash := txidHash.GetBytes()
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>		if hashList == nil || dbConn == nil {
<span id="L120" class="ln">   120</span>			return nil
<span id="L121" class="ln">   121</span>		}
<span id="L122" class="ln">   122</span>		if hashList[string(hash)] != PURGE {
<span id="L123" class="ln">   123</span>			return nil
<span id="L124" class="ln">   124</span>		}
<span id="L125" class="ln">   125</span>	
<span id="L126" class="ln">   126</span>		for s, err := dbConn.Query(&#34;SELECT txid FROM purge WHERE hash=?&#34;, hash); err == nil; err = s.Next() {
<span id="L127" class="ln">   127</span>			var txid []byte
<span id="L128" class="ln">   128</span>			s.Scan(&amp;txid) <span class="comment">// Assigns 1st column to rowid, the rest to row</span>
<span id="L129" class="ln">   129</span>			p := new(objects.Purge)
<span id="L130" class="ln">   130</span>			p.FromBytes(txid)
<span id="L131" class="ln">   131</span>			return p
<span id="L132" class="ln">   132</span>		}
<span id="L133" class="ln">   133</span>		<span class="comment">// Not Found</span>
<span id="L134" class="ln">   134</span>		return nil
<span id="L135" class="ln">   135</span>	}
<span id="L136" class="ln">   136</span>	
<span id="L137" class="ln">   137</span>	<span class="comment">// Add Published Message to database and hash list.</span>
<span id="L138" class="ln">   138</span>	func AddPub(log chan string, msg *objects.Message) error {
<span id="L139" class="ln">   139</span>		mutex.Lock()
<span id="L140" class="ln">   140</span>		defer mutex.Unlock()
<span id="L141" class="ln">   141</span>	
<span id="L142" class="ln">   142</span>		if hashList == nil || dbConn == nil {
<span id="L143" class="ln">   143</span>			return DBError(EUNINIT)
<span id="L144" class="ln">   144</span>		}
<span id="L145" class="ln">   145</span>		if Contains(msg.TxidHash) == MSG {
<span id="L146" class="ln">   146</span>			return nil
<span id="L147" class="ln">   147</span>		}
<span id="L148" class="ln">   148</span>	
<span id="L149" class="ln">   149</span>		err := dbConn.Exec(&#34;INSERT INTO pub VALUES (?, ?, ?, ?)&#34;, msg.TxidHash.GetBytes(), msg.AddrHash.GetBytes(), msg.Timestamp.Unix(), msg.Content.GetBytes())
<span id="L150" class="ln">   150</span>		if err != nil {
<span id="L151" class="ln">   151</span>			log &lt;- fmt.Sprintf(&#34;Error inserting message into db... %s&#34;, err)
<span id="L152" class="ln">   152</span>			return err
<span id="L153" class="ln">   153</span>		}
<span id="L154" class="ln">   154</span>	
<span id="L155" class="ln">   155</span>		Add(msg.TxidHash, PUB)
<span id="L156" class="ln">   156</span>		return nil
<span id="L157" class="ln">   157</span>	}
<span id="L158" class="ln">   158</span>	
<span id="L159" class="ln">   159</span>	<span class="comment">// Add basic message to database and hash list.</span>
<span id="L160" class="ln">   160</span>	func AddMessage(log chan string, msg *objects.Message) error {
<span id="L161" class="ln">   161</span>		mutex.Lock()
<span id="L162" class="ln">   162</span>		defer mutex.Unlock()
<span id="L163" class="ln">   163</span>	
<span id="L164" class="ln">   164</span>		if hashList == nil || dbConn == nil {
<span id="L165" class="ln">   165</span>			return DBError(EUNINIT)
<span id="L166" class="ln">   166</span>		}
<span id="L167" class="ln">   167</span>		if Contains(msg.TxidHash) == MSG {
<span id="L168" class="ln">   168</span>			return nil
<span id="L169" class="ln">   169</span>		}
<span id="L170" class="ln">   170</span>	
<span id="L171" class="ln">   171</span>		err := dbConn.Exec(&#34;INSERT INTO msg VALUES (?, ?, ?, ?)&#34;, msg.TxidHash.GetBytes(), msg.AddrHash.GetBytes(), msg.Timestamp.Unix(), msg.Content.GetBytes())
<span id="L172" class="ln">   172</span>		if err != nil {
<span id="L173" class="ln">   173</span>			log &lt;- fmt.Sprintf(&#34;Error inserting message into db... %s&#34;, err)
<span id="L174" class="ln">   174</span>			return err
<span id="L175" class="ln">   175</span>		}
<span id="L176" class="ln">   176</span>	
<span id="L177" class="ln">   177</span>		Add(msg.TxidHash, MSG)
<span id="L178" class="ln">   178</span>		return nil
<span id="L179" class="ln">   179</span>	
<span id="L180" class="ln">   180</span>	}
<span id="L181" class="ln">   181</span>	
<span id="L182" class="ln">   182</span>	<span class="comment">// Get basic message from database.</span>
<span id="L183" class="ln">   183</span>	func GetMessage(log chan string, txidHash objects.Hash) *objects.Message {
<span id="L184" class="ln">   184</span>		mutex.Lock()
<span id="L185" class="ln">   185</span>		defer mutex.Unlock()
<span id="L186" class="ln">   186</span>	
<span id="L187" class="ln">   187</span>		hash := txidHash.GetBytes()
<span id="L188" class="ln">   188</span>	
<span id="L189" class="ln">   189</span>		if hashList == nil || dbConn == nil {
<span id="L190" class="ln">   190</span>			return nil
<span id="L191" class="ln">   191</span>		}
<span id="L192" class="ln">   192</span>		if hashList[string(hash)] != MSG &amp;&amp; hashList[string(hash)] != PUB {
<span id="L193" class="ln">   193</span>			return nil
<span id="L194" class="ln">   194</span>		}
<span id="L195" class="ln">   195</span>	
<span id="L196" class="ln">   196</span>		msg := new(objects.Message)
<span id="L197" class="ln">   197</span>	
<span id="L198" class="ln">   198</span>		for s, err := dbConn.Query(&#34;SELECT * FROM msg WHERE hash=?&#34;, hash); err == nil; err = s.Next() {
<span id="L199" class="ln">   199</span>			var timestamp int64
<span id="L200" class="ln">   200</span>			encrypted := make([]byte, 0, 0)
<span id="L201" class="ln">   201</span>			txidhash := make([]byte, 0, 0)
<span id="L202" class="ln">   202</span>			addrhash := make([]byte, 0, 0)
<span id="L203" class="ln">   203</span>			s.Scan(&amp;txidhash, &amp;addrhash, &amp;timestamp, &amp;encrypted)
<span id="L204" class="ln">   204</span>	
<span id="L205" class="ln">   205</span>			msg.TxidHash.FromBytes(txidhash)
<span id="L206" class="ln">   206</span>			msg.AddrHash.FromBytes(addrhash)
<span id="L207" class="ln">   207</span>			msg.Timestamp = time.Unix(timestamp, 0)
<span id="L208" class="ln">   208</span>			msg.Content.FromBytes(encrypted)
<span id="L209" class="ln">   209</span>	
<span id="L210" class="ln">   210</span>			return msg
<span id="L211" class="ln">   211</span>		}
<span id="L212" class="ln">   212</span>		<span class="comment">// Not Found</span>
<span id="L213" class="ln">   213</span>		return nil
<span id="L214" class="ln">   214</span>	}
<span id="L215" class="ln">   215</span>	
<span id="L216" class="ln">   216</span>	<span class="comment">// Remove any object from the database and hash list.</span>
<span id="L217" class="ln">   217</span>	func RemoveHash(log chan string, hashObj objects.Hash) error {
<span id="L218" class="ln">   218</span>		mutex.Lock()
<span id="L219" class="ln">   219</span>		defer mutex.Unlock()
<span id="L220" class="ln">   220</span>	
<span id="L221" class="ln">   221</span>		hash := hashObj.GetBytes()
<span id="L222" class="ln">   222</span>	
<span id="L223" class="ln">   223</span>		if hashList == nil || dbConn == nil {
<span id="L224" class="ln">   224</span>			return DBError(EUNINIT)
<span id="L225" class="ln">   225</span>		}
<span id="L226" class="ln">   226</span>	
<span id="L227" class="ln">   227</span>		var sql string
<span id="L228" class="ln">   228</span>	
<span id="L229" class="ln">   229</span>		switch Contains(hashObj) {
<span id="L230" class="ln">   230</span>		case PUBKEY:
<span id="L231" class="ln">   231</span>			sql = &#34;DELETE FROM pubkey WHERE hash=?&#34;
<span id="L232" class="ln">   232</span>		case MSG:
<span id="L233" class="ln">   233</span>			sql = &#34;DELETE FROM msg WHERE hash=?&#34;
<span id="L234" class="ln">   234</span>		case PURGE:
<span id="L235" class="ln">   235</span>			sql = &#34;DELETE FROM purge WHERE hash=?&#34;
<span id="L236" class="ln">   236</span>		default:
<span id="L237" class="ln">   237</span>			return nil
<span id="L238" class="ln">   238</span>		}
<span id="L239" class="ln">   239</span>	
<span id="L240" class="ln">   240</span>		err := dbConn.Exec(sql, hash)
<span id="L241" class="ln">   241</span>		if err != nil {
<span id="L242" class="ln">   242</span>			log &lt;- fmt.Sprintf(&#34;Error deleting hash from db... %s&#34;, err)
<span id="L243" class="ln">   243</span>			return nil
<span id="L244" class="ln">   244</span>		}
<span id="L245" class="ln">   245</span>	
<span id="L246" class="ln">   246</span>		Delete(hashObj)
<span id="L247" class="ln">   247</span>		return nil
<span id="L248" class="ln">   248</span>	}
</pre><p><a href="/src/pkg/emp/db/objects.go?m=text">View as plain text</a></p>

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

