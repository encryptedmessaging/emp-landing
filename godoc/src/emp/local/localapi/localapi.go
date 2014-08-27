<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/local/localapi/localapi.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/local/localapi/localapi.go</h1>




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
<span id="L12" class="ln">    12</span>	<span class="comment">// Package LocalAPI provides the RPC-API that allows management of Addresses and Messages</span>
<span id="L13" class="ln">    13</span>	<span class="comment">// in EMPLocal.</span>
<span id="L14" class="ln">    14</span>	package localapi
<span id="L15" class="ln">    15</span>	
<span id="L16" class="ln">    16</span>	import (
<span id="L17" class="ln">    17</span>		&#34;emp/api&#34;
<span id="L18" class="ln">    18</span>		&#34;emp/db&#34;
<span id="L19" class="ln">    19</span>		&#34;emp/encryption&#34;
<span id="L20" class="ln">    20</span>		&#34;emp/local/localdb&#34;
<span id="L21" class="ln">    21</span>		&#34;emp/objects&#34;
<span id="L22" class="ln">    22</span>		&#34;encoding/base64&#34;
<span id="L23" class="ln">    23</span>		&#34;errors&#34;
<span id="L24" class="ln">    24</span>		&#34;fmt&#34;
<span id="L25" class="ln">    25</span>		&#34;github.com/gorilla/rpc&#34;
<span id="L26" class="ln">    26</span>		&#34;github.com/gorilla/rpc/json&#34;
<span id="L27" class="ln">    27</span>		&#34;net&#34;
<span id="L28" class="ln">    28</span>		&#34;net/http&#34;
<span id="L29" class="ln">    29</span>		&#34;time&#34;
<span id="L30" class="ln">    30</span>	)
<span id="L31" class="ln">    31</span>	
<span id="L32" class="ln">    32</span>	type EMPService struct {
<span id="L33" class="ln">    33</span>		Config *api.ApiConfig
<span id="L34" class="ln">    34</span>	}
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>	type NilParam struct{}
<span id="L37" class="ln">    37</span>	
<span id="L38" class="ln">    38</span>	func (s *EMPService) Version(r *http.Request, args *NilParam, reply *objects.Version) error {
<span id="L39" class="ln">    39</span>		if !basicAuth(s.Config, r) {
<span id="L40" class="ln">    40</span>			s.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L41" class="ln">    41</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L42" class="ln">    42</span>		}
<span id="L43" class="ln">    43</span>	
<span id="L44" class="ln">    44</span>		*reply = s.Config.LocalVersion
<span id="L45" class="ln">    45</span>		return nil
<span id="L46" class="ln">    46</span>	}
<span id="L47" class="ln">    47</span>	
<span id="L48" class="ln">    48</span>	func basicAuth(config *api.ApiConfig, r *http.Request) bool {
<span id="L49" class="ln">    49</span>		if config == nil || r == nil {
<span id="L50" class="ln">    50</span>			return false
<span id="L51" class="ln">    51</span>		}
<span id="L52" class="ln">    52</span>	
<span id="L53" class="ln">    53</span>		auth := r.Header.Get(&#34;Authorization&#34;)
<span id="L54" class="ln">    54</span>	
<span id="L55" class="ln">    55</span>		auth2 := &#34;Basic &#34; + base64.StdEncoding.EncodeToString([]byte(config.RPCUser+&#34;:&#34;+config.RPCPass))
<span id="L56" class="ln">    56</span>	
<span id="L57" class="ln">    57</span>		return (auth == auth2)
<span id="L58" class="ln">    58</span>	}
<span id="L59" class="ln">    59</span>	
<span id="L60" class="ln">    60</span>	func Initialize(config *api.ApiConfig) error {
<span id="L61" class="ln">    61</span>	
<span id="L62" class="ln">    62</span>		e := localdb.Initialize(config.Log, config.LocalDB)
<span id="L63" class="ln">    63</span>	
<span id="L64" class="ln">    64</span>		if e != nil {
<span id="L65" class="ln">    65</span>			return e
<span id="L66" class="ln">    66</span>		}
<span id="L67" class="ln">    67</span>	
<span id="L68" class="ln">    68</span>		s := rpc.NewServer()
<span id="L69" class="ln">    69</span>		s.RegisterCodec(json.NewCodec(), &#34;application/json&#34;)
<span id="L70" class="ln">    70</span>		service := new(EMPService)
<span id="L71" class="ln">    71</span>		service.Config = config
<span id="L72" class="ln">    72</span>		s.RegisterService(service, &#34;EMPService&#34;)
<span id="L73" class="ln">    73</span>	
<span id="L74" class="ln">    74</span>		<span class="comment">// Register RPC Services</span>
<span id="L75" class="ln">    75</span>		http.Handle(&#34;/rpc&#34;, s)
<span id="L76" class="ln">    76</span>	
<span id="L77" class="ln">    77</span>		<span class="comment">// Register JS Client</span>
<span id="L78" class="ln">    78</span>		http.Handle(&#34;/&#34;, http.FileServer(http.Dir(config.HttpRoot)))
<span id="L79" class="ln">    79</span>	
<span id="L80" class="ln">    80</span>		l, e := net.Listen(&#34;tcp&#34;, fmt.Sprintf(&#34;:%d&#34;, config.RPCPort))
<span id="L81" class="ln">    81</span>		if e != nil {
<span id="L82" class="ln">    82</span>			config.Log &lt;- fmt.Sprintf(&#34;RPC Listen Error: %s&#34;, e)
<span id="L83" class="ln">    83</span>			return e
<span id="L84" class="ln">    84</span>		}
<span id="L85" class="ln">    85</span>	
<span id="L86" class="ln">    86</span>		go http.Serve(l, nil)
<span id="L87" class="ln">    87</span>	
<span id="L88" class="ln">    88</span>		go register(config)
<span id="L89" class="ln">    89</span>	
<span id="L90" class="ln">    90</span>		portStr := fmt.Sprintf(&#34;:%d&#34;, config.RPCPort)
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>		config.Log &lt;- fmt.Sprintf(&#34;Started RPC Server on: %s&#34;, portStr)
<span id="L93" class="ln">    93</span>		return nil
<span id="L94" class="ln">    94</span>	}
<span id="L95" class="ln">    95</span>	
<span id="L96" class="ln">    96</span>	func Cleanup() {
<span id="L97" class="ln">    97</span>		localdb.Cleanup()
<span id="L98" class="ln">    98</span>	}
<span id="L99" class="ln">    99</span>	
<span id="L100" class="ln">   100</span>	<span class="comment">// Handle Pubkey, Message, and Purge Registration</span>
<span id="L101" class="ln">   101</span>	func register(config *api.ApiConfig) {
<span id="L102" class="ln">   102</span>		var message objects.Message
<span id="L103" class="ln">   103</span>		var txid [16]byte
<span id="L104" class="ln">   104</span>	
<span id="L105" class="ln">   105</span>		for {
<span id="L106" class="ln">   106</span>			select {
<span id="L107" class="ln">   107</span>			case pubHash := &lt;-config.PubkeyRegister:
<span id="L108" class="ln">   108</span>	
<span id="L109" class="ln">   109</span>				<span class="comment">// Check if pubkey is in database...</span>
<span id="L110" class="ln">   110</span>				pubkey := checkPubkey(config, pubHash)
<span id="L111" class="ln">   111</span>	
<span id="L112" class="ln">   112</span>				if pubkey == nil {
<span id="L113" class="ln">   113</span>					break
<span id="L114" class="ln">   114</span>				}
<span id="L115" class="ln">   115</span>	
<span id="L116" class="ln">   116</span>				outbox := localdb.GetBox(localdb.OUTBOX)
<span id="L117" class="ln">   117</span>				for _, metamsg := range outbox {
<span id="L118" class="ln">   118</span>					recvHash := objects.MakeHash([]byte(metamsg.Recipient))
<span id="L119" class="ln">   119</span>					if string(pubHash.GetBytes()) == string(recvHash.GetBytes()) {
<span id="L120" class="ln">   120</span>						<span class="comment">// Send message and move to sendbox</span>
<span id="L121" class="ln">   121</span>						msg, err := localdb.GetMessageDetail(metamsg.TxidHash)
<span id="L122" class="ln">   122</span>						if err != nil {
<span id="L123" class="ln">   123</span>							config.Log &lt;- err.Error()
<span id="L124" class="ln">   124</span>							break
<span id="L125" class="ln">   125</span>						}
<span id="L126" class="ln">   126</span>						msg.Encrypted = encryption.Encrypt(config.Log, pubkey, string(msg.Decrypted.GetBytes()))
<span id="L127" class="ln">   127</span>						msg.MetaMessage.Timestamp = time.Now().Round(time.Second)
<span id="L128" class="ln">   128</span>						err = localdb.AddUpdateMessage(msg, localdb.SENDBOX)
<span id="L129" class="ln">   129</span>						if err != nil {
<span id="L130" class="ln">   130</span>							config.Log &lt;- err.Error()
<span id="L131" class="ln">   131</span>							break
<span id="L132" class="ln">   132</span>						}
<span id="L133" class="ln">   133</span>	
<span id="L134" class="ln">   134</span>						sendMsg := new(objects.Message)
<span id="L135" class="ln">   135</span>						sendMsg.Timestamp = msg.MetaMessage.Timestamp
<span id="L136" class="ln">   136</span>						sendMsg.TxidHash = msg.MetaMessage.TxidHash
<span id="L137" class="ln">   137</span>						sendMsg.AddrHash = recvHash
<span id="L138" class="ln">   138</span>						sendMsg.Content = *msg.Encrypted
<span id="L139" class="ln">   139</span>	
<span id="L140" class="ln">   140</span>						config.RecvQueue &lt;- *objects.MakeFrame(objects.MSG, objects.BROADCAST, sendMsg)
<span id="L141" class="ln">   141</span>					}
<span id="L142" class="ln">   142</span>				}
<span id="L143" class="ln">   143</span>	
<span id="L144" class="ln">   144</span>			case message = &lt;-config.MessageRegister:
<span id="L145" class="ln">   145</span>				<span class="comment">// If address is registered, store message in inbox</span>
<span id="L146" class="ln">   146</span>				detail, err := localdb.GetAddressDetail(message.AddrHash)
<span id="L147" class="ln">   147</span>				if err != nil {
<span id="L148" class="ln">   148</span>					config.Log &lt;- &#34;Message address not in database...&#34;
<span id="L149" class="ln">   149</span>					break
<span id="L150" class="ln">   150</span>				}
<span id="L151" class="ln">   151</span>				if !detail.IsRegistered {
<span id="L152" class="ln">   152</span>					config.Log &lt;- &#34;Message not for registered address...&#34;
<span id="L153" class="ln">   153</span>					break
<span id="L154" class="ln">   154</span>				}
<span id="L155" class="ln">   155</span>	
<span id="L156" class="ln">   156</span>				config.Log &lt;- &#34;Registering new encrypted message...&#34;
<span id="L157" class="ln">   157</span>	
<span id="L158" class="ln">   158</span>				msg := new(objects.FullMessage)
<span id="L159" class="ln">   159</span>				msg.MetaMessage.TxidHash = message.TxidHash
<span id="L160" class="ln">   160</span>				msg.MetaMessage.Timestamp = message.Timestamp
<span id="L161" class="ln">   161</span>				msg.MetaMessage.Recipient = detail.String
<span id="L162" class="ln">   162</span>				msg.Encrypted = &amp;message.Content
<span id="L163" class="ln">   163</span>	
<span id="L164" class="ln">   164</span>				err = localdb.AddUpdateMessage(msg, localdb.INBOX)
<span id="L165" class="ln">   165</span>				if err != nil {
<span id="L166" class="ln">   166</span>					config.Log &lt;- err.Error()
<span id="L167" class="ln">   167</span>				}
<span id="L168" class="ln">   168</span>			case message = &lt;-config.PubRegister:
<span id="L169" class="ln">   169</span>				<span class="comment">// If address is registered, store message in inbox</span>
<span id="L170" class="ln">   170</span>				detail, err := localdb.GetAddressDetail(message.AddrHash)
<span id="L171" class="ln">   171</span>				if err != nil {
<span id="L172" class="ln">   172</span>					config.Log &lt;- &#34;Message address not in database...&#34;
<span id="L173" class="ln">   173</span>					break
<span id="L174" class="ln">   174</span>				}
<span id="L175" class="ln">   175</span>				if !detail.IsSubscribed {
<span id="L176" class="ln">   176</span>					config.Log &lt;- &#34;Not Subscribed to Address...&#34;
<span id="L177" class="ln">   177</span>					break
<span id="L178" class="ln">   178</span>				}
<span id="L179" class="ln">   179</span>	
<span id="L180" class="ln">   180</span>				config.Log &lt;- &#34;Registering new publication...&#34;
<span id="L181" class="ln">   181</span>	
<span id="L182" class="ln">   182</span>				msg := new(objects.FullMessage)
<span id="L183" class="ln">   183</span>				msg.MetaMessage.TxidHash = message.TxidHash
<span id="L184" class="ln">   184</span>				msg.MetaMessage.Timestamp = message.Timestamp
<span id="L185" class="ln">   185</span>				msg.MetaMessage.Sender = detail.String
<span id="L186" class="ln">   186</span>				msg.MetaMessage.Recipient = &#34;&lt;Subscription Message&gt;&#34;
<span id="L187" class="ln">   187</span>				msg.Encrypted = &amp;message.Content
<span id="L188" class="ln">   188</span>	
<span id="L189" class="ln">   189</span>				msg.Decrypted = new(objects.DecryptedMessage)
<span id="L190" class="ln">   190</span>				msg.Decrypted.FromBytes(encryption.DecryptPub(config.Log, detail.Pubkey, msg.Encrypted))
<span id="L191" class="ln">   191</span>	
<span id="L192" class="ln">   192</span>				err = localdb.AddUpdateMessage(msg, localdb.INBOX)
<span id="L193" class="ln">   193</span>				if err != nil {
<span id="L194" class="ln">   194</span>					config.Log &lt;- err.Error()
<span id="L195" class="ln">   195</span>				}
<span id="L196" class="ln">   196</span>			case txid = &lt;-config.PurgeRegister:
<span id="L197" class="ln">   197</span>				<span class="comment">// If Message in database, mark as purged</span>
<span id="L198" class="ln">   198</span>				detail, err := localdb.GetMessageDetail(objects.MakeHash(txid[:]))
<span id="L199" class="ln">   199</span>				if err != nil {
<span id="L200" class="ln">   200</span>					break
<span id="L201" class="ln">   201</span>				}
<span id="L202" class="ln">   202</span>				detail.MetaMessage.Purged = true
<span id="L203" class="ln">   203</span>				err = localdb.AddUpdateMessage(detail, -1)
<span id="L204" class="ln">   204</span>				if err != nil {
<span id="L205" class="ln">   205</span>					config.Log &lt;- fmt.Sprintf(&#34;Error registering purge: %s&#34;, err)
<span id="L206" class="ln">   206</span>				}
<span id="L207" class="ln">   207</span>			} <span class="comment">// End select</span>
<span id="L208" class="ln">   208</span>		} <span class="comment">// End for</span>
<span id="L209" class="ln">   209</span>	} <span class="comment">// End register</span>
<span id="L210" class="ln">   210</span>	
<span id="L211" class="ln">   211</span>	func checkPubkey(config *api.ApiConfig, addrHash objects.Hash) []byte {
<span id="L212" class="ln">   212</span>	
<span id="L213" class="ln">   213</span>		<span class="comment">// First check local DB</span>
<span id="L214" class="ln">   214</span>		detail, err := localdb.GetAddressDetail(addrHash)
<span id="L215" class="ln">   215</span>		if err != nil {
<span id="L216" class="ln">   216</span>			<span class="comment">// If not in database, won&#39;t be able to decrypt anyway!</span>
<span id="L217" class="ln">   217</span>			return nil
<span id="L218" class="ln">   218</span>		}
<span id="L219" class="ln">   219</span>		if len(detail.Pubkey) &gt; 0 {
<span id="L220" class="ln">   220</span>			if db.Contains(addrHash) != db.PUBKEY {
<span id="L221" class="ln">   221</span>				enc := new(objects.EncryptedPubkey)
<span id="L222" class="ln">   222</span>	
<span id="L223" class="ln">   223</span>				enc.IV, enc.Payload, _ = encryption.SymmetricEncrypt(detail.Address, string(detail.Pubkey))
<span id="L224" class="ln">   224</span>				enc.AddrHash = objects.MakeHash(detail.Address)
<span id="L225" class="ln">   225</span>	
<span id="L226" class="ln">   226</span>				config.RecvQueue &lt;- *objects.MakeFrame(objects.PUBKEY, objects.BROADCAST, enc)
<span id="L227" class="ln">   227</span>			}
<span id="L228" class="ln">   228</span>			return detail.Pubkey
<span id="L229" class="ln">   229</span>		}
<span id="L230" class="ln">   230</span>	
<span id="L231" class="ln">   231</span>		<span class="comment">// If not there, check local database</span>
<span id="L232" class="ln">   232</span>		if db.Contains(addrHash) == db.PUBKEY {
<span id="L233" class="ln">   233</span>			enc := db.GetPubkey(config.Log, addrHash)
<span id="L234" class="ln">   234</span>	
<span id="L235" class="ln">   235</span>			pubkey := encryption.SymmetricDecrypt(enc.IV, detail.Address, enc.Payload)
<span id="L236" class="ln">   236</span>			pubkey = pubkey[:65]
<span id="L237" class="ln">   237</span>	
<span id="L238" class="ln">   238</span>			<span class="comment">// Check public Key</span>
<span id="L239" class="ln">   239</span>			x, y := encryption.UnmarshalPubkey(pubkey)
<span id="L240" class="ln">   240</span>			if x == nil {
<span id="L241" class="ln">   241</span>				config.Log &lt;- &#34;Decrypted Public Key Invalid&#34;
<span id="L242" class="ln">   242</span>				return nil
<span id="L243" class="ln">   243</span>			}
<span id="L244" class="ln">   244</span>	
<span id="L245" class="ln">   245</span>			address2 := encryption.GetAddress(config.Log, x, y)
<span id="L246" class="ln">   246</span>			if string(detail.Address) != string(address2) {
<span id="L247" class="ln">   247</span>				config.Log &lt;- &#34;Decrypted Public Key doesn&#39;t match provided address!&#34;
<span id="L248" class="ln">   248</span>				return nil
<span id="L249" class="ln">   249</span>			}
<span id="L250" class="ln">   250</span>	
<span id="L251" class="ln">   251</span>			detail.Pubkey = pubkey
<span id="L252" class="ln">   252</span>			err := localdb.AddUpdateAddress(detail)
<span id="L253" class="ln">   253</span>			if err != nil {
<span id="L254" class="ln">   254</span>				config.Log &lt;- &#34;Error adding pubkey to local database!&#34;
<span id="L255" class="ln">   255</span>				return nil
<span id="L256" class="ln">   256</span>			}
<span id="L257" class="ln">   257</span>	
<span id="L258" class="ln">   258</span>			return pubkey
<span id="L259" class="ln">   259</span>		}
<span id="L260" class="ln">   260</span>	
<span id="L261" class="ln">   261</span>		<span class="comment">// If not there, send a pubkey request</span>
<span id="L262" class="ln">   262</span>		config.RecvQueue &lt;- *objects.MakeFrame(objects.PUBKEY_REQUEST, objects.BROADCAST, &amp;addrHash)
<span id="L263" class="ln">   263</span>		return nil
<span id="L264" class="ln">   264</span>	}
</pre><p><a href="/src/emp/local/localapi/localapi.go?m=text">View as plain text</a></p>

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

