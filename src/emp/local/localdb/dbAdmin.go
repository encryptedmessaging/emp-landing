<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/local/localdb/dbAdmin.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/local/localdb/dbAdmin.go</h1>




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
<span id="L12" class="ln">    12</span>	package localdb
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;emp/encryption&#34;
<span id="L16" class="ln">    16</span>		&#34;emp/objects&#34;
<span id="L17" class="ln">    17</span>		&#34;errors&#34;
<span id="L18" class="ln">    18</span>		&#34;time&#34;
<span id="L19" class="ln">    19</span>	)
<span id="L20" class="ln">    20</span>	
<span id="L21" class="ln">    21</span>	func AddUpdateAddress(address *objects.AddressDetail) error {
<span id="L22" class="ln">    22</span>		localMutex.Lock()
<span id="L23" class="ln">    23</span>		defer localMutex.Unlock()
<span id="L24" class="ln">    24</span>	
<span id="L25" class="ln">    25</span>		var err error
<span id="L26" class="ln">    26</span>	
<span id="L27" class="ln">    27</span>		if address.Address == nil {
<span id="L28" class="ln">    28</span>			address.Address = encryption.StringToAddress(address.String)
<span id="L29" class="ln">    29</span>		}
<span id="L30" class="ln">    30</span>	
<span id="L31" class="ln">    31</span>		if address.Address == nil {
<span id="L32" class="ln">    32</span>			return errors.New(&#34;Invalid Address!&#34;)
<span id="L33" class="ln">    33</span>		}
<span id="L34" class="ln">    34</span>	
<span id="L35" class="ln">    35</span>		addrHash := objects.MakeHash(address.Address)
<span id="L36" class="ln">    36</span>	
<span id="L37" class="ln">    37</span>		if Contains(addrHash) == ADDRESS { <span class="comment">// Exists in message database, update pubkey, privkey, and registration</span>
<span id="L38" class="ln">    38</span>			err = LocalDB.Exec(&#34;UPDATE addressbook SET registered=?, subscribed=?, label=? WHERE hash=?&#34;, address.IsRegistered, address.IsSubscribed, address.Label, addrHash.GetBytes())
<span id="L39" class="ln">    39</span>			if err != nil {
<span id="L40" class="ln">    40</span>				return err
<span id="L41" class="ln">    41</span>			}
<span id="L42" class="ln">    42</span>	
<span id="L43" class="ln">    43</span>			if address.Pubkey != nil {
<span id="L44" class="ln">    44</span>				err = LocalDB.Exec(&#34;UPDATE addressbook SET pubkey=? WHERE hash=?&#34;, address.Pubkey, addrHash.GetBytes())
<span id="L45" class="ln">    45</span>				if err != nil {
<span id="L46" class="ln">    46</span>					return err
<span id="L47" class="ln">    47</span>				}
<span id="L48" class="ln">    48</span>			}
<span id="L49" class="ln">    49</span>	
<span id="L50" class="ln">    50</span>			if address.Privkey != nil {
<span id="L51" class="ln">    51</span>				err = LocalDB.Exec(&#34;UPDATE addressbook SET privkey=? WHERE hash=?&#34;, address.Privkey, addrHash.GetBytes())
<span id="L52" class="ln">    52</span>				if err != nil {
<span id="L53" class="ln">    53</span>					return err
<span id="L54" class="ln">    54</span>				}
<span id="L55" class="ln">    55</span>			}
<span id="L56" class="ln">    56</span>	
<span id="L57" class="ln">    57</span>			if address.EncPrivkey != nil {
<span id="L58" class="ln">    58</span>				err = LocalDB.Exec(&#34;UPDATE addressbook SET encprivkey=? WHERE hash=?&#34;, address.EncPrivkey, addrHash.GetBytes())
<span id="L59" class="ln">    59</span>				if err != nil {
<span id="L60" class="ln">    60</span>					return err
<span id="L61" class="ln">    61</span>				}
<span id="L62" class="ln">    62</span>			}
<span id="L63" class="ln">    63</span>	
<span id="L64" class="ln">    64</span>		} else { <span class="comment">// Doesn&#39;t exist yet, insert it!</span>
<span id="L65" class="ln">    65</span>			err = LocalDB.Exec(&#34;INSERT INTO addressbook VALUES (?, ?, ?, ?, ?, ?, ?, ?)&#34;, addrHash.GetBytes(), address.Address, address.IsRegistered, address.Pubkey, address.Privkey, address.Label, address.IsSubscribed, address.EncPrivkey)
<span id="L66" class="ln">    66</span>			if err != nil {
<span id="L67" class="ln">    67</span>				return err
<span id="L68" class="ln">    68</span>			}
<span id="L69" class="ln">    69</span>			Add(addrHash, ADDRESS)
<span id="L70" class="ln">    70</span>		}
<span id="L71" class="ln">    71</span>	
<span id="L72" class="ln">    72</span>		return nil
<span id="L73" class="ln">    73</span>	}
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>	func GetAddressDetail(addrHash objects.Hash) (*objects.AddressDetail, error) {
<span id="L76" class="ln">    76</span>		localMutex.Lock()
<span id="L77" class="ln">    77</span>		defer localMutex.Unlock()
<span id="L78" class="ln">    78</span>	
<span id="L79" class="ln">    79</span>		if Contains(addrHash) != ADDRESS {
<span id="L80" class="ln">    80</span>			return nil, errors.New(&#34;Address not found!&#34;)
<span id="L81" class="ln">    81</span>		}
<span id="L82" class="ln">    82</span>	
<span id="L83" class="ln">    83</span>		ret := new(objects.AddressDetail)
<span id="L84" class="ln">    84</span>	
<span id="L85" class="ln">    85</span>		s, err := LocalDB.Query(&#34;SELECT address, registered, pubkey, privkey, label, subscribed, encprivkey FROM addressbook WHERE hash=?&#34;, addrHash.GetBytes())
<span id="L86" class="ln">    86</span>		if err == nil {
<span id="L87" class="ln">    87</span>			s.Scan(&amp;ret.Address, &amp;ret.IsRegistered, &amp;ret.Pubkey, &amp;ret.Privkey, &amp;ret.Label, &amp;ret.IsSubscribed, &amp;ret.EncPrivkey)
<span id="L88" class="ln">    88</span>			ret.String = encryption.AddressToString(ret.Address)
<span id="L89" class="ln">    89</span>			return ret, nil
<span id="L90" class="ln">    90</span>		}
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>		return nil, err
<span id="L93" class="ln">    93</span>	}
<span id="L94" class="ln">    94</span>	
<span id="L95" class="ln">    95</span>	func ListAddresses(registered bool) [][2]string {
<span id="L96" class="ln">    96</span>		ret := make([][2]string, 0, 0)
<span id="L97" class="ln">    97</span>	
<span id="L98" class="ln">    98</span>		for s, err := LocalDB.Query(&#34;SELECT address, label FROM addressbook WHERE registered=?&#34;, registered); err == nil; err = s.Next() {
<span id="L99" class="ln">    99</span>			var addr []byte
<span id="L100" class="ln">   100</span>			var label string
<span id="L101" class="ln">   101</span>			s.Scan(&amp;addr, &amp;label)
<span id="L102" class="ln">   102</span>			ret = append(ret, [2]string{encryption.AddressToString(addr), label})
<span id="L103" class="ln">   103</span>		}
<span id="L104" class="ln">   104</span>	
<span id="L105" class="ln">   105</span>		return ret
<span id="L106" class="ln">   106</span>	}
<span id="L107" class="ln">   107</span>	
<span id="L108" class="ln">   108</span>	func GetMessageDetail(txidHash objects.Hash) (*objects.FullMessage, error) {
<span id="L109" class="ln">   109</span>		localMutex.Lock()
<span id="L110" class="ln">   110</span>		defer localMutex.Unlock()
<span id="L111" class="ln">   111</span>	
<span id="L112" class="ln">   112</span>		if Contains(txidHash) &gt; SENDBOX {
<span id="L113" class="ln">   113</span>			return nil, errors.New(&#34;Message not found!&#34;)
<span id="L114" class="ln">   114</span>		}
<span id="L115" class="ln">   115</span>	
<span id="L116" class="ln">   116</span>		ret := new(objects.FullMessage)
<span id="L117" class="ln">   117</span>		ret.Encrypted = new(encryption.EncryptedMessage)
<span id="L118" class="ln">   118</span>		ret.Decrypted = new(objects.DecryptedMessage)
<span id="L119" class="ln">   119</span>	
<span id="L120" class="ln">   120</span>		s, err := LocalDB.Query(&#34;SELECT * FROM msg WHERE txid_hash=?&#34;, txidHash.GetBytes())
<span id="L121" class="ln">   121</span>		if err == nil {
<span id="L122" class="ln">   122</span>			recipient := make([]byte, 0, 0)
<span id="L123" class="ln">   123</span>			sender := make([]byte, 0, 0)
<span id="L124" class="ln">   124</span>			encrypted := make([]byte, 0, 0)
<span id="L125" class="ln">   125</span>			decrypted := make([]byte, 0, 0)
<span id="L126" class="ln">   126</span>			txidHash := make([]byte, 0, 0)
<span id="L127" class="ln">   127</span>			var timestamp int64
<span id="L128" class="ln">   128</span>			var purged bool
<span id="L129" class="ln">   129</span>			var box int
<span id="L130" class="ln">   130</span>	
<span id="L131" class="ln">   131</span>			s.Scan(&amp;txidHash, &amp;recipient, &amp;timestamp, &amp;box, &amp;encrypted, &amp;decrypted, &amp;purged, &amp;sender)
<span id="L132" class="ln">   132</span>			ret.MetaMessage.TxidHash.FromBytes(txidHash)
<span id="L133" class="ln">   133</span>			ret.MetaMessage.Recipient = encryption.AddressToString(recipient)
<span id="L134" class="ln">   134</span>			ret.MetaMessage.Sender = encryption.AddressToString(sender)
<span id="L135" class="ln">   135</span>			ret.MetaMessage.Timestamp = time.Unix(timestamp, 0)
<span id="L136" class="ln">   136</span>			ret.MetaMessage.Purged = purged
<span id="L137" class="ln">   137</span>			ret.Encrypted.FromBytes(encrypted)
<span id="L138" class="ln">   138</span>			if len(decrypted) &gt; 0 {
<span id="L139" class="ln">   139</span>				ret.Decrypted.FromBytes(decrypted)
<span id="L140" class="ln">   140</span>			} else {
<span id="L141" class="ln">   141</span>				ret.Decrypted = nil
<span id="L142" class="ln">   142</span>			}
<span id="L143" class="ln">   143</span>			return ret, nil
<span id="L144" class="ln">   144</span>		}
<span id="L145" class="ln">   145</span>	
<span id="L146" class="ln">   146</span>		return nil, err
<span id="L147" class="ln">   147</span>	
<span id="L148" class="ln">   148</span>	}
<span id="L149" class="ln">   149</span>	
<span id="L150" class="ln">   150</span>	func DeleteMessage(txidHash *objects.Hash) error {
<span id="L151" class="ln">   151</span>		localMutex.Lock()
<span id="L152" class="ln">   152</span>		defer localMutex.Unlock()
<span id="L153" class="ln">   153</span>	
<span id="L154" class="ln">   154</span>		if Contains(*txidHash) &gt; SENDBOX {
<span id="L155" class="ln">   155</span>			return errors.New(&#34;Error Deleting Message: Not Found!&#34;)
<span id="L156" class="ln">   156</span>		}
<span id="L157" class="ln">   157</span>	
<span id="L158" class="ln">   158</span>		return LocalDB.Exec(&#34;DELETE FROM msg WHERE txid_hash=?&#34;, txidHash.GetBytes());
<span id="L159" class="ln">   159</span>	}
<span id="L160" class="ln">   160</span>	
<span id="L161" class="ln">   161</span>	func DeleteAddress(addrHash *objects.Hash) error {
<span id="L162" class="ln">   162</span>		localMutex.Lock()
<span id="L163" class="ln">   163</span>		defer localMutex.Unlock()
<span id="L164" class="ln">   164</span>	
<span id="L165" class="ln">   165</span>		if Contains(*addrHash) &gt; ADDRESS {
<span id="L166" class="ln">   166</span>			return errors.New(&#34;Error Deleting Message: Not Found!&#34;)
<span id="L167" class="ln">   167</span>		}
<span id="L168" class="ln">   168</span>	
<span id="L169" class="ln">   169</span>		return LocalDB.Exec(&#34;DELETE FROM addressbook WHERE hash=?&#34;, addrHash.GetBytes());
<span id="L170" class="ln">   170</span>	}
<span id="L171" class="ln">   171</span>	
<span id="L172" class="ln">   172</span>	func AddUpdateMessage(msg *objects.FullMessage, box int) error {
<span id="L173" class="ln">   173</span>		localMutex.Lock()
<span id="L174" class="ln">   174</span>		defer localMutex.Unlock()
<span id="L175" class="ln">   175</span>	
<span id="L176" class="ln">   176</span>		var err error
<span id="L177" class="ln">   177</span>	
<span id="L178" class="ln">   178</span>		if Contains(msg.MetaMessage.TxidHash) &gt; SENDBOX { <span class="comment">// Insert Message Into Database!</span>
<span id="L179" class="ln">   179</span>	
<span id="L180" class="ln">   180</span>			err = LocalDB.Exec(&#34;INSERT INTO msg VALUES (?, ?, ?, ?, ?, ?, ?, ?)&#34;, msg.MetaMessage.TxidHash.GetBytes(), encryption.StringToAddress(msg.MetaMessage.Recipient),
<span id="L181" class="ln">   181</span>				msg.MetaMessage.Timestamp.Unix(), box, msg.Encrypted.GetBytes(), msg.Decrypted.GetBytes(), msg.MetaMessage.Purged, encryption.StringToAddress(msg.MetaMessage.Sender))
<span id="L182" class="ln">   182</span>			if err != nil {
<span id="L183" class="ln">   183</span>				return err
<span id="L184" class="ln">   184</span>			}
<span id="L185" class="ln">   185</span>	
<span id="L186" class="ln">   186</span>		} else { <span class="comment">// Update recipient, sender, purged, encrypted, decrypted, box</span>
<span id="L187" class="ln">   187</span>			if box &lt; 0 {
<span id="L188" class="ln">   188</span>				err = LocalDB.Exec(&#34;UPDATE msg SET purged=? WHERE txid_hash=?&#34;, msg.MetaMessage.Purged, msg.MetaMessage.TxidHash.GetBytes())
<span id="L189" class="ln">   189</span>				if err != nil {
<span id="L190" class="ln">   190</span>					return err
<span id="L191" class="ln">   191</span>				}
<span id="L192" class="ln">   192</span>			} else {
<span id="L193" class="ln">   193</span>				err = LocalDB.Exec(&#34;UPDATE msg SET box=?, purged=? WHERE txid_hash=?&#34;, box, msg.MetaMessage.Purged, msg.MetaMessage.TxidHash.GetBytes())
<span id="L194" class="ln">   194</span>				if err != nil {
<span id="L195" class="ln">   195</span>					return err
<span id="L196" class="ln">   196</span>				}
<span id="L197" class="ln">   197</span>			}
<span id="L198" class="ln">   198</span>	
<span id="L199" class="ln">   199</span>			if len(msg.MetaMessage.Sender) &gt; 0 {
<span id="L200" class="ln">   200</span>				err = LocalDB.Exec(&#34;UPDATE msg SET sender=? WHERE txid_hash=?&#34;, encryption.StringToAddress(msg.MetaMessage.Sender), msg.MetaMessage.TxidHash.GetBytes())
<span id="L201" class="ln">   201</span>				if err != nil {
<span id="L202" class="ln">   202</span>					return err
<span id="L203" class="ln">   203</span>				}
<span id="L204" class="ln">   204</span>			}
<span id="L205" class="ln">   205</span>	
<span id="L206" class="ln">   206</span>			if len(msg.MetaMessage.Recipient) &gt; 0 {
<span id="L207" class="ln">   207</span>				err = LocalDB.Exec(&#34;UPDATE msg SET recipient=? WHERE txid_hash=?&#34;, encryption.StringToAddress(msg.MetaMessage.Recipient), msg.MetaMessage.TxidHash.GetBytes())
<span id="L208" class="ln">   208</span>				if err != nil {
<span id="L209" class="ln">   209</span>					return err
<span id="L210" class="ln">   210</span>				}
<span id="L211" class="ln">   211</span>			}
<span id="L212" class="ln">   212</span>	
<span id="L213" class="ln">   213</span>			if msg.Encrypted != nil {
<span id="L214" class="ln">   214</span>				err = LocalDB.Exec(&#34;UPDATE msg SET encrypted=? WHERE txid_hash=?&#34;, msg.Encrypted.GetBytes(), msg.MetaMessage.TxidHash.GetBytes())
<span id="L215" class="ln">   215</span>				if err != nil {
<span id="L216" class="ln">   216</span>					return err
<span id="L217" class="ln">   217</span>				}
<span id="L218" class="ln">   218</span>			}
<span id="L219" class="ln">   219</span>	
<span id="L220" class="ln">   220</span>			if msg.Decrypted != nil {
<span id="L221" class="ln">   221</span>				err = LocalDB.Exec(&#34;UPDATE msg SET decrypted=? WHERE txid_hash=?&#34;, msg.Decrypted.GetBytes(), msg.MetaMessage.TxidHash.GetBytes())
<span id="L222" class="ln">   222</span>				if err != nil {
<span id="L223" class="ln">   223</span>					return err
<span id="L224" class="ln">   224</span>				}
<span id="L225" class="ln">   225</span>			}
<span id="L226" class="ln">   226</span>	
<span id="L227" class="ln">   227</span>		}
<span id="L228" class="ln">   228</span>	
<span id="L229" class="ln">   229</span>		Add(msg.MetaMessage.TxidHash, box)
<span id="L230" class="ln">   230</span>		return nil
<span id="L231" class="ln">   231</span>	}
<span id="L232" class="ln">   232</span>	
<span id="L233" class="ln">   233</span>	func GetBox(box int) []objects.MetaMessage {
<span id="L234" class="ln">   234</span>		if box &gt; SENDBOX || box &lt; INBOX {
<span id="L235" class="ln">   235</span>			return nil
<span id="L236" class="ln">   236</span>		}
<span id="L237" class="ln">   237</span>	
<span id="L238" class="ln">   238</span>		localMutex.Lock()
<span id="L239" class="ln">   239</span>		defer localMutex.Unlock()
<span id="L240" class="ln">   240</span>	
<span id="L241" class="ln">   241</span>		ret := make([]objects.MetaMessage, 0, 0)
<span id="L242" class="ln">   242</span>	
<span id="L243" class="ln">   243</span>		for s, err := LocalDB.Query(&#34;SELECT txid_hash, timestamp, purged, sender, recipient FROM msg WHERE box=?&#34;, box); err == nil; err = s.Next() {
<span id="L244" class="ln">   244</span>			mm := new(objects.MetaMessage)
<span id="L245" class="ln">   245</span>			sendBytes := make([]byte, 0, 0)
<span id="L246" class="ln">   246</span>			recvBytes := make([]byte, 0, 0)
<span id="L247" class="ln">   247</span>			txidHash := make([]byte, 0, 0)
<span id="L248" class="ln">   248</span>			var timestamp int64
<span id="L249" class="ln">   249</span>	
<span id="L250" class="ln">   250</span>			s.Scan(&amp;txidHash, &amp;timestamp, &amp;mm.Purged, &amp;sendBytes, &amp;recvBytes)
<span id="L251" class="ln">   251</span>			mm.Sender = encryption.AddressToString(sendBytes)
<span id="L252" class="ln">   252</span>			mm.Recipient = encryption.AddressToString(recvBytes)
<span id="L253" class="ln">   253</span>	
<span id="L254" class="ln">   254</span>			mm.TxidHash.FromBytes(txidHash)
<span id="L255" class="ln">   255</span>			mm.Timestamp = time.Unix(timestamp, 0)
<span id="L256" class="ln">   256</span>	
<span id="L257" class="ln">   257</span>			ret = append(ret, *mm)
<span id="L258" class="ln">   258</span>		}
<span id="L259" class="ln">   259</span>	
<span id="L260" class="ln">   260</span>		return ret
<span id="L261" class="ln">   261</span>	}
<span id="L262" class="ln">   262</span>	
<span id="L263" class="ln">   263</span>	func GetBySender(sender string) []objects.MetaMessage {
<span id="L264" class="ln">   264</span>		localMutex.Lock()
<span id="L265" class="ln">   265</span>		defer localMutex.Unlock()
<span id="L266" class="ln">   266</span>	
<span id="L267" class="ln">   267</span>		ret := make([]objects.MetaMessage, 0, 0)
<span id="L268" class="ln">   268</span>	
<span id="L269" class="ln">   269</span>		for s, err := LocalDB.Query(&#34;SELECT txid_hash, timestamp, purged, sender, recipient FROM msg WHERE sender=?&#34;, sender); err == nil; err = s.Next() {
<span id="L270" class="ln">   270</span>			mm := new(objects.MetaMessage)
<span id="L271" class="ln">   271</span>			sendBytes := make([]byte, 0, 0)
<span id="L272" class="ln">   272</span>			recvBytes := make([]byte, 0, 0)
<span id="L273" class="ln">   273</span>			txidHash := make([]byte, 0, 0)
<span id="L274" class="ln">   274</span>			var timestamp int64
<span id="L275" class="ln">   275</span>	
<span id="L276" class="ln">   276</span>			s.Scan(&amp;txidHash, &amp;timestamp, &amp;mm.Purged, &amp;sendBytes, &amp;recvBytes)
<span id="L277" class="ln">   277</span>			mm.Sender = encryption.AddressToString(sendBytes)
<span id="L278" class="ln">   278</span>			mm.Recipient = encryption.AddressToString(recvBytes)
<span id="L279" class="ln">   279</span>	
<span id="L280" class="ln">   280</span>			mm.TxidHash.FromBytes(txidHash)
<span id="L281" class="ln">   281</span>			mm.Timestamp = time.Unix(timestamp, 0)
<span id="L282" class="ln">   282</span>	
<span id="L283" class="ln">   283</span>			ret = append(ret, *mm)
<span id="L284" class="ln">   284</span>		}
<span id="L285" class="ln">   285</span>	
<span id="L286" class="ln">   286</span>		return ret
<span id="L287" class="ln">   287</span>	}
<span id="L288" class="ln">   288</span>	
<span id="L289" class="ln">   289</span>	func GetByRecipient(recipient string) []objects.MetaMessage {
<span id="L290" class="ln">   290</span>		localMutex.Lock()
<span id="L291" class="ln">   291</span>		defer localMutex.Unlock()
<span id="L292" class="ln">   292</span>	
<span id="L293" class="ln">   293</span>		ret := make([]objects.MetaMessage, 0, 0)
<span id="L294" class="ln">   294</span>	
<span id="L295" class="ln">   295</span>		for s, err := LocalDB.Query(&#34;SELECT txid_hash, timestamp, purged, sender, recipient FROM msg WHERE recipient=?&#34;, recipient); err == nil; err = s.Next() {
<span id="L296" class="ln">   296</span>			mm := new(objects.MetaMessage)
<span id="L297" class="ln">   297</span>			sendBytes := make([]byte, 0, 0)
<span id="L298" class="ln">   298</span>			recvBytes := make([]byte, 0, 0)
<span id="L299" class="ln">   299</span>			txidHash := make([]byte, 0, 0)
<span id="L300" class="ln">   300</span>			var timestamp int64
<span id="L301" class="ln">   301</span>	
<span id="L302" class="ln">   302</span>			s.Scan(&amp;txidHash, &amp;timestamp, &amp;mm.Purged, &amp;sendBytes, &amp;recvBytes)
<span id="L303" class="ln">   303</span>			mm.Sender = encryption.AddressToString(sendBytes)
<span id="L304" class="ln">   304</span>			mm.Recipient = encryption.AddressToString(recvBytes)
<span id="L305" class="ln">   305</span>	
<span id="L306" class="ln">   306</span>			mm.TxidHash.FromBytes(txidHash)
<span id="L307" class="ln">   307</span>			mm.Timestamp = time.Unix(timestamp, 0)
<span id="L308" class="ln">   308</span>	
<span id="L309" class="ln">   309</span>			ret = append(ret, *mm)
<span id="L310" class="ln">   310</span>		}
<span id="L311" class="ln">   311</span>	
<span id="L312" class="ln">   312</span>		return ret
<span id="L313" class="ln">   313</span>	}
<span id="L314" class="ln">   314</span>	
<span id="L315" class="ln">   315</span>	func DeleteObject(obj objects.Hash) error {
<span id="L316" class="ln">   316</span>		var err error
<span id="L317" class="ln">   317</span>		switch Contains(obj) {
<span id="L318" class="ln">   318</span>		case INBOX:
<span id="L319" class="ln">   319</span>			fallthrough
<span id="L320" class="ln">   320</span>		case SENDBOX:
<span id="L321" class="ln">   321</span>			fallthrough
<span id="L322" class="ln">   322</span>		case OUTBOX:
<span id="L323" class="ln">   323</span>			err = LocalDB.Exec(&#34;DELETE FROM msg WHERE txid_hash=?&#34;, obj.GetBytes())
<span id="L324" class="ln">   324</span>		case ADDRESS:
<span id="L325" class="ln">   325</span>			err = LocalDB.Exec(&#34;DELETE FROM addressbook WHERE hash=?&#34;, obj.GetBytes())
<span id="L326" class="ln">   326</span>		default:
<span id="L327" class="ln">   327</span>			err = errors.New(&#34;Hash not found!&#34;)
<span id="L328" class="ln">   328</span>		}
<span id="L329" class="ln">   329</span>	
<span id="L330" class="ln">   330</span>		if err == nil {
<span id="L331" class="ln">   331</span>			Del(obj)
<span id="L332" class="ln">   332</span>		}
<span id="L333" class="ln">   333</span>	
<span id="L334" class="ln">   334</span>		return err
<span id="L335" class="ln">   335</span>	}
</pre><p><a href="/src/emp/local/localdb/dbAdmin.go?m=text">View as plain text</a></p>

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

