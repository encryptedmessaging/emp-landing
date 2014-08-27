<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/local/localapi/rpcmsg.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/local/localapi/rpcmsg.go</h1>




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
<span id="L12" class="ln">    12</span>	package localapi
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;crypto/ecdsa&#34;
<span id="L16" class="ln">    16</span>		&#34;crypto/rand&#34;
<span id="L17" class="ln">    17</span>		&#34;emp/encryption&#34;
<span id="L18" class="ln">    18</span>		&#34;emp/local/localdb&#34;
<span id="L19" class="ln">    19</span>		&#34;emp/objects&#34;
<span id="L20" class="ln">    20</span>		&#34;errors&#34;
<span id="L21" class="ln">    21</span>		&#34;fmt&#34;
<span id="L22" class="ln">    22</span>		&#34;math/big&#34;
<span id="L23" class="ln">    23</span>		&#34;net/http&#34;
<span id="L24" class="ln">    24</span>		&#34;time&#34;
<span id="L25" class="ln">    25</span>	)
<span id="L26" class="ln">    26</span>	
<span id="L27" class="ln">    27</span>	type SendMsg struct {
<span id="L28" class="ln">    28</span>		Sender    string `json:&#34;sender&#34;`
<span id="L29" class="ln">    29</span>		Recipient string `json:&#34;recipient&#34;`
<span id="L30" class="ln">    30</span>		Subject   string `json:&#34;subject&#34;`
<span id="L31" class="ln">    31</span>		Plaintext string `json:&#34;content&#34;`
<span id="L32" class="ln">    32</span>	}
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>	type SendResponse struct {
<span id="L35" class="ln">    35</span>		TxidHash []byte `json:&#34;txid_hash&#34;`
<span id="L36" class="ln">    36</span>		IsSent   bool   `json:&#34;sent&#34;`
<span id="L37" class="ln">    37</span>	}
<span id="L38" class="ln">    38</span>	
<span id="L39" class="ln">    39</span>	type PubMsg struct {
<span id="L40" class="ln">    40</span>		Sender    string `json:&#34;sender&#34;`
<span id="L41" class="ln">    41</span>		Subject   string `json:&#34;subject&#34;`
<span id="L42" class="ln">    42</span>		Plaintext string `json:&#34;content&#34;`
<span id="L43" class="ln">    43</span>	}
<span id="L44" class="ln">    44</span>	
<span id="L45" class="ln">    45</span>	type RawMsg struct {
<span id="L46" class="ln">    46</span>		Message objects.Message `json:&#34;message&#34;`
<span id="L47" class="ln">    47</span>		SendAddress string `json:&#34;sender&#34;`
<span id="L48" class="ln">    48</span>		Subscription bool `json:&#34;is_subscription&#34;`
<span id="L49" class="ln">    49</span>	}
<span id="L50" class="ln">    50</span>	
<span id="L51" class="ln">    51</span>	func (service *EMPService) PublishMessage(r *http.Request, args *SendMsg, reply *SendResponse) error {
<span id="L52" class="ln">    52</span>		if !basicAuth(service.Config, r) {
<span id="L53" class="ln">    53</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L54" class="ln">    54</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L55" class="ln">    55</span>		}
<span id="L56" class="ln">    56</span>	
<span id="L57" class="ln">    57</span>		<span class="comment">// Nil Check</span>
<span id="L58" class="ln">    58</span>		if len(args.Sender) == 0 || len(args.Plaintext) == 0 {
<span id="L59" class="ln">    59</span>			return errors.New(&#34;All fields required.&#34;)
<span id="L60" class="ln">    60</span>		}
<span id="L61" class="ln">    61</span>	
<span id="L62" class="ln">    62</span>		var err error
<span id="L63" class="ln">    63</span>	
<span id="L64" class="ln">    64</span>		<span class="comment">// Get Addresses</span>
<span id="L65" class="ln">    65</span>		sendAddr := encryption.StringToAddress(args.Sender)
<span id="L66" class="ln">    66</span>		if len(sendAddr) == 0 {
<span id="L67" class="ln">    67</span>			return errors.New(&#34;Invalid sender address!&#34;)
<span id="L68" class="ln">    68</span>		}
<span id="L69" class="ln">    69</span>	
<span id="L70" class="ln">    70</span>		sender, err := localdb.GetAddressDetail(objects.MakeHash(sendAddr))
<span id="L71" class="ln">    71</span>		if err != nil {
<span id="L72" class="ln">    72</span>			return errors.New(fmt.Sprintf(&#34;Error pulling send address from Database: %s&#34;, err))
<span id="L73" class="ln">    73</span>		}
<span id="L74" class="ln">    74</span>		if sender.Privkey == nil {
<span id="L75" class="ln">    75</span>			return errors.New(&#34;Private Key Required to Publish Message&#34;)
<span id="L76" class="ln">    76</span>		}
<span id="L77" class="ln">    77</span>	
<span id="L78" class="ln">    78</span>		<span class="comment">// Create New Message</span>
<span id="L79" class="ln">    79</span>		msg := new(objects.FullMessage)
<span id="L80" class="ln">    80</span>		msg.Decrypted = new(objects.DecryptedMessage)
<span id="L81" class="ln">    81</span>		msg.Encrypted = nil
<span id="L82" class="ln">    82</span>	
<span id="L83" class="ln">    83</span>		txid := make([]byte, len(msg.Decrypted.Txid), cap(msg.Decrypted.Txid))
<span id="L84" class="ln">    84</span>	
<span id="L85" class="ln">    85</span>		<span class="comment">// Fill out decrypted message</span>
<span id="L86" class="ln">    86</span>		n, err := rand.Read(txid)
<span id="L87" class="ln">    87</span>		if n &lt; len(msg.Decrypted.Txid[:]) || err != nil {
<span id="L88" class="ln">    88</span>			return errors.New(fmt.Sprintf(&#34;Problem with random reader: %s&#34;, err))
<span id="L89" class="ln">    89</span>		}
<span id="L90" class="ln">    90</span>		copy(msg.Decrypted.Pubkey[:], sender.Pubkey)
<span id="L91" class="ln">    91</span>		msg.Decrypted.Subject = args.Subject
<span id="L92" class="ln">    92</span>		msg.Decrypted.MimeType = &#34;text/plain&#34;
<span id="L93" class="ln">    93</span>		msg.Decrypted.Content = args.Plaintext
<span id="L94" class="ln">    94</span>		msg.Decrypted.Length = uint32(len(msg.Decrypted.Content))
<span id="L95" class="ln">    95</span>	
<span id="L96" class="ln">    96</span>		<span class="comment">// Fill Out Meta Message (save timestamp)</span>
<span id="L97" class="ln">    97</span>		msg.MetaMessage.Purged = false
<span id="L98" class="ln">    98</span>		msg.MetaMessage.TxidHash = objects.MakeHash(txid)
<span id="L99" class="ln">    99</span>		msg.MetaMessage.Sender = sender.String
<span id="L100" class="ln">   100</span>		msg.MetaMessage.Recipient = &#34;&lt;Subscription Message&gt;&#34;
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		<span class="comment">// Get Signature</span>
<span id="L103" class="ln">   103</span>		priv := new(ecdsa.PrivateKey)
<span id="L104" class="ln">   104</span>		priv.PublicKey.Curve = encryption.GetCurve()
<span id="L105" class="ln">   105</span>		priv.D = new(big.Int)
<span id="L106" class="ln">   106</span>		priv.D.SetBytes(sender.Privkey)
<span id="L107" class="ln">   107</span>	
<span id="L108" class="ln">   108</span>		sign := msg.Decrypted.GetBytes()
<span id="L109" class="ln">   109</span>		sign = sign[:len(sign)-65]
<span id="L110" class="ln">   110</span>		signHash := objects.MakeHash(sign)
<span id="L111" class="ln">   111</span>	
<span id="L112" class="ln">   112</span>		x, y, err := ecdsa.Sign(rand.Reader, priv, signHash.GetBytes())
<span id="L113" class="ln">   113</span>		if err != nil {
<span id="L114" class="ln">   114</span>			return err
<span id="L115" class="ln">   115</span>		}
<span id="L116" class="ln">   116</span>	
<span id="L117" class="ln">   117</span>		copy(msg.Decrypted.Signature[:], encryption.MarshalPubkey(x, y))
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>		<span class="comment">// Send message and add to sendbox...</span>
<span id="L120" class="ln">   120</span>		msg.Encrypted = encryption.EncryptPub(service.Config.Log, sender.Privkey, string(msg.Decrypted.GetBytes()))
<span id="L121" class="ln">   121</span>		msg.MetaMessage.Timestamp = time.Now().Round(time.Second)
<span id="L122" class="ln">   122</span>	
<span id="L123" class="ln">   123</span>		<span class="comment">// Now Add Txid</span>
<span id="L124" class="ln">   124</span>		copy(msg.Decrypted.Txid[:], txid)
<span id="L125" class="ln">   125</span>	
<span id="L126" class="ln">   126</span>		err = localdb.AddUpdateMessage(msg, localdb.SENDBOX)
<span id="L127" class="ln">   127</span>		if err != nil {
<span id="L128" class="ln">   128</span>			return err
<span id="L129" class="ln">   129</span>		}
<span id="L130" class="ln">   130</span>	
<span id="L131" class="ln">   131</span>		sendMsg := new(objects.Message)
<span id="L132" class="ln">   132</span>		sendMsg.TxidHash = msg.MetaMessage.TxidHash
<span id="L133" class="ln">   133</span>		sendMsg.AddrHash = objects.MakeHash(sender.Address)
<span id="L134" class="ln">   134</span>		sendMsg.Timestamp = msg.MetaMessage.Timestamp
<span id="L135" class="ln">   135</span>		sendMsg.Content = *msg.Encrypted
<span id="L136" class="ln">   136</span>	
<span id="L137" class="ln">   137</span>		service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.PUB, objects.BROADCAST, sendMsg)
<span id="L138" class="ln">   138</span>	
<span id="L139" class="ln">   139</span>		reply.IsSent = true
<span id="L140" class="ln">   140</span>	
<span id="L141" class="ln">   141</span>		<span class="comment">// Finish by setting msg&#39;s txid</span>
<span id="L142" class="ln">   142</span>		reply.TxidHash = msg.MetaMessage.TxidHash.GetBytes()
<span id="L143" class="ln">   143</span>		return nil
<span id="L144" class="ln">   144</span>	
<span id="L145" class="ln">   145</span>	}
<span id="L146" class="ln">   146</span>	
<span id="L147" class="ln">   147</span>	func (service *EMPService) PurgeMessage(r *http.Request, args *[]byte, reply *NilParam) error {
<span id="L148" class="ln">   148</span>		if len(*args) != 16 {
<span id="L149" class="ln">   149</span>			return errors.New(&#34;Invalid Txid: Bad Length&#34;)
<span id="L150" class="ln">   150</span>		}
<span id="L151" class="ln">   151</span>		
<span id="L152" class="ln">   152</span>		txidHash := objects.MakeHash(*args)
<span id="L153" class="ln">   153</span>	
<span id="L154" class="ln">   154</span>		if (localdb.Contains(txidHash) &lt;= localdb.SENDBOX) {
<span id="L155" class="ln">   155</span>			msg, err := localdb.GetMessageDetail(txidHash)
<span id="L156" class="ln">   156</span>			if (err != nil) {
<span id="L157" class="ln">   157</span>				return errors.New(fmt.Sprintf(&#34;Problem Retrieving Message: %s&#34;, err))
<span id="L158" class="ln">   158</span>			}
<span id="L159" class="ln">   159</span>			msg.MetaMessage.Purged = true
<span id="L160" class="ln">   160</span>			localdb.AddUpdateMessage(msg, -1)
<span id="L161" class="ln">   161</span>	
<span id="L162" class="ln">   162</span>			<span class="comment">// Send Purge Request</span>
<span id="L163" class="ln">   163</span>			purge := new(objects.Purge)
<span id="L164" class="ln">   164</span>			purge.Txid = msg.Decrypted.Txid
<span id="L165" class="ln">   165</span>	
<span id="L166" class="ln">   166</span>			service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.PURGE, objects.BROADCAST, purge)
<span id="L167" class="ln">   167</span>	
<span id="L168" class="ln">   168</span>			return nil
<span id="L169" class="ln">   169</span>		}
<span id="L170" class="ln">   170</span>	
<span id="L171" class="ln">   171</span>		return errors.New(&#34;Txid Not Found&#34;)
<span id="L172" class="ln">   172</span>	}
<span id="L173" class="ln">   173</span>	
<span id="L174" class="ln">   174</span>	func (service *EMPService) DeleteMessage(r *http.Request, args *[]byte, reply *NilParam) error {
<span id="L175" class="ln">   175</span>		txidHash := new(objects.Hash)
<span id="L176" class="ln">   176</span>		txidHash.FromBytes(*args)
<span id="L177" class="ln">   177</span>	
<span id="L178" class="ln">   178</span>		return localdb.DeleteMessage(txidHash)
<span id="L179" class="ln">   179</span>	}
<span id="L180" class="ln">   180</span>	
<span id="L181" class="ln">   181</span>	func (service *EMPService) SendRawMsg(r *http.Request, args *RawMsg, reply *NilParam) error {
<span id="L182" class="ln">   182</span>		if (args == nil) {
<span id="L183" class="ln">   183</span>			return errors.New(&#34;Cannot work with nil message object!&#34;)
<span id="L184" class="ln">   184</span>		}
<span id="L185" class="ln">   185</span>	
<span id="L186" class="ln">   186</span>		detail, err := localdb.GetAddressDetail(args.Message.AddrHash)
<span id="L187" class="ln">   187</span>		if (err != nil) {
<span id="L188" class="ln">   188</span>			return err
<span id="L189" class="ln">   189</span>		}
<span id="L190" class="ln">   190</span>	
<span id="L191" class="ln">   191</span>		<span class="comment">// Create New Message</span>
<span id="L192" class="ln">   192</span>		msg := new(objects.FullMessage)
<span id="L193" class="ln">   193</span>		msg.Decrypted = nil
<span id="L194" class="ln">   194</span>		msg.Encrypted = new(encryption.EncryptedMessage)
<span id="L195" class="ln">   195</span>		msg.Encrypted.FromBytes(args.Message.Content.GetBytes())
<span id="L196" class="ln">   196</span>	
<span id="L197" class="ln">   197</span>		<span class="comment">// Fill Out Meta Message (save timestamp)</span>
<span id="L198" class="ln">   198</span>		msg.MetaMessage.Purged = false
<span id="L199" class="ln">   199</span>	
<span id="L200" class="ln">   200</span>		msg.MetaMessage.TxidHash.FromBytes(args.Message.TxidHash.GetBytes())
<span id="L201" class="ln">   201</span>	
<span id="L202" class="ln">   202</span>		if (args.Subscription) {
<span id="L203" class="ln">   203</span>			msg.MetaMessage.Sender = detail.String
<span id="L204" class="ln">   204</span>			msg.MetaMessage.Recipient = &#34;&#34;
<span id="L205" class="ln">   205</span>		} else {
<span id="L206" class="ln">   206</span>			msg.MetaMessage.Sender = args.SendAddress
<span id="L207" class="ln">   207</span>			msg.MetaMessage.Recipient = detail.String
<span id="L208" class="ln">   208</span>		} 
<span id="L209" class="ln">   209</span>		
<span id="L210" class="ln">   210</span>		err = localdb.AddUpdateMessage(msg, localdb.SENDBOX)
<span id="L211" class="ln">   211</span>		if (err != nil) {
<span id="L212" class="ln">   212</span>			return err
<span id="L213" class="ln">   213</span>		}
<span id="L214" class="ln">   214</span>	
<span id="L215" class="ln">   215</span>	
<span id="L216" class="ln">   216</span>		if (args.Subscription) {
<span id="L217" class="ln">   217</span>			service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.PUB, objects.BROADCAST, &amp;(args.Message))
<span id="L218" class="ln">   218</span>		} else {
<span id="L219" class="ln">   219</span>			service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.MSG, objects.BROADCAST, &amp;(args.Message))
<span id="L220" class="ln">   220</span>		}
<span id="L221" class="ln">   221</span>	
<span id="L222" class="ln">   222</span>		return nil
<span id="L223" class="ln">   223</span>	}
<span id="L224" class="ln">   224</span>	
<span id="L225" class="ln">   225</span>	func (service *EMPService) SendMessage(r *http.Request, args *SendMsg, reply *SendResponse) error {
<span id="L226" class="ln">   226</span>		if !basicAuth(service.Config, r) {
<span id="L227" class="ln">   227</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L228" class="ln">   228</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L229" class="ln">   229</span>		}
<span id="L230" class="ln">   230</span>	
<span id="L231" class="ln">   231</span>		<span class="comment">// Nil Check</span>
<span id="L232" class="ln">   232</span>		if len(args.Sender) == 0 || len(args.Recipient) == 0 || len(args.Plaintext) == 0 {
<span id="L233" class="ln">   233</span>			return errors.New(&#34;All fields required except signature.&#34;)
<span id="L234" class="ln">   234</span>		}
<span id="L235" class="ln">   235</span>	
<span id="L236" class="ln">   236</span>		var err error
<span id="L237" class="ln">   237</span>	
<span id="L238" class="ln">   238</span>		<span class="comment">// Get Addresses</span>
<span id="L239" class="ln">   239</span>		sendAddr := encryption.StringToAddress(args.Sender)
<span id="L240" class="ln">   240</span>		if len(sendAddr) == 0 {
<span id="L241" class="ln">   241</span>			return errors.New(&#34;Invalid sender address!&#34;)
<span id="L242" class="ln">   242</span>		}
<span id="L243" class="ln">   243</span>	
<span id="L244" class="ln">   244</span>		recvAddr := encryption.StringToAddress(args.Recipient)
<span id="L245" class="ln">   245</span>		if len(recvAddr) == 0 {
<span id="L246" class="ln">   246</span>			return errors.New(&#34;Invalid recipient address!&#34;)
<span id="L247" class="ln">   247</span>		}
<span id="L248" class="ln">   248</span>	
<span id="L249" class="ln">   249</span>		sender, err := localdb.GetAddressDetail(objects.MakeHash(sendAddr))
<span id="L250" class="ln">   250</span>		if err != nil {
<span id="L251" class="ln">   251</span>			return errors.New(fmt.Sprintf(&#34;Error pulling send address from Database: %s&#34;, err))
<span id="L252" class="ln">   252</span>		}
<span id="L253" class="ln">   253</span>		if sender.Pubkey == nil {
<span id="L254" class="ln">   254</span>			sender.Pubkey = checkPubkey(service.Config, objects.MakeHash(sendAddr))
<span id="L255" class="ln">   255</span>			if sender.Pubkey == nil {
<span id="L256" class="ln">   256</span>				return errors.New(&#34;Sender&#39;s Public Key is required to send message!&#34;)
<span id="L257" class="ln">   257</span>			}
<span id="L258" class="ln">   258</span>		}
<span id="L259" class="ln">   259</span>		if sender.Privkey == nil {
<span id="L260" class="ln">   260</span>			return errors.New(&#34;SendMsg() requires a stored private key. Use SendRawMsg() instead.&#34;)
<span id="L261" class="ln">   261</span>		}
<span id="L262" class="ln">   262</span>	
<span id="L263" class="ln">   263</span>		recipient, err := localdb.GetAddressDetail(objects.MakeHash(recvAddr))
<span id="L264" class="ln">   264</span>		if err != nil {
<span id="L265" class="ln">   265</span>			return errors.New(fmt.Sprintf(&#34;Error pulling recipient address from Database: %s&#34;, err))
<span id="L266" class="ln">   266</span>		}
<span id="L267" class="ln">   267</span>	
<span id="L268" class="ln">   268</span>		<span class="comment">// Create New Message</span>
<span id="L269" class="ln">   269</span>		msg := new(objects.FullMessage)
<span id="L270" class="ln">   270</span>		msg.Decrypted = new(objects.DecryptedMessage)
<span id="L271" class="ln">   271</span>		msg.Encrypted = nil
<span id="L272" class="ln">   272</span>	
<span id="L273" class="ln">   273</span>		<span class="comment">// Fill out decrypted message</span>
<span id="L274" class="ln">   274</span>		n, err := rand.Read(msg.Decrypted.Txid[:])
<span id="L275" class="ln">   275</span>		if n &lt; len(msg.Decrypted.Txid[:]) || err != nil {
<span id="L276" class="ln">   276</span>			return errors.New(fmt.Sprintf(&#34;Problem with random reader: %s&#34;, err))
<span id="L277" class="ln">   277</span>		}
<span id="L278" class="ln">   278</span>		copy(msg.Decrypted.Pubkey[:], sender.Pubkey)
<span id="L279" class="ln">   279</span>		msg.Decrypted.Subject = args.Subject
<span id="L280" class="ln">   280</span>		msg.Decrypted.MimeType = &#34;text/plain&#34;
<span id="L281" class="ln">   281</span>		msg.Decrypted.Content = args.Plaintext
<span id="L282" class="ln">   282</span>		msg.Decrypted.Length = uint32(len(msg.Decrypted.Content))
<span id="L283" class="ln">   283</span>	
<span id="L284" class="ln">   284</span>		<span class="comment">// Fill Out Meta Message (save timestamp)</span>
<span id="L285" class="ln">   285</span>		msg.MetaMessage.Purged = false
<span id="L286" class="ln">   286</span>		msg.MetaMessage.TxidHash = objects.MakeHash(msg.Decrypted.Txid[:])
<span id="L287" class="ln">   287</span>		msg.MetaMessage.Sender = sender.String
<span id="L288" class="ln">   288</span>		msg.MetaMessage.Recipient = recipient.String
<span id="L289" class="ln">   289</span>	
<span id="L290" class="ln">   290</span>		<span class="comment">// Get Signature</span>
<span id="L291" class="ln">   291</span>		priv := new(ecdsa.PrivateKey)
<span id="L292" class="ln">   292</span>		priv.PublicKey.Curve = encryption.GetCurve()
<span id="L293" class="ln">   293</span>		priv.D = new(big.Int)
<span id="L294" class="ln">   294</span>		priv.D.SetBytes(sender.Privkey)
<span id="L295" class="ln">   295</span>	
<span id="L296" class="ln">   296</span>		sign := msg.Decrypted.GetBytes()
<span id="L297" class="ln">   297</span>		sign = sign[:len(sign)-65]
<span id="L298" class="ln">   298</span>		signHash := objects.MakeHash(sign)
<span id="L299" class="ln">   299</span>	
<span id="L300" class="ln">   300</span>		x, y, err := ecdsa.Sign(rand.Reader, priv, signHash.GetBytes())
<span id="L301" class="ln">   301</span>		if err != nil {
<span id="L302" class="ln">   302</span>			return err
<span id="L303" class="ln">   303</span>		}
<span id="L304" class="ln">   304</span>	
<span id="L305" class="ln">   305</span>		copy(msg.Decrypted.Signature[:], encryption.MarshalPubkey(x, y))
<span id="L306" class="ln">   306</span>	
<span id="L307" class="ln">   307</span>		<span class="comment">// Check for pubkey</span>
<span id="L308" class="ln">   308</span>		if recipient.Pubkey == nil {
<span id="L309" class="ln">   309</span>			recipient.Pubkey = checkPubkey(service.Config, objects.MakeHash(recipient.Address))
<span id="L310" class="ln">   310</span>		}
<span id="L311" class="ln">   311</span>	
<span id="L312" class="ln">   312</span>		if recipient.Pubkey == nil {
<span id="L313" class="ln">   313</span>			reply.IsSent = false
<span id="L314" class="ln">   314</span>			<span class="comment">// Add message to outbox...</span>
<span id="L315" class="ln">   315</span>			err = localdb.AddUpdateMessage(msg, localdb.OUTBOX)
<span id="L316" class="ln">   316</span>			if err != nil {
<span id="L317" class="ln">   317</span>				return err
<span id="L318" class="ln">   318</span>			}
<span id="L319" class="ln">   319</span>	
<span id="L320" class="ln">   320</span>		} else {
<span id="L321" class="ln">   321</span>			<span class="comment">// Send message and add to sendbox...</span>
<span id="L322" class="ln">   322</span>			msg.Encrypted = encryption.Encrypt(service.Config.Log, recipient.Pubkey, string(msg.Decrypted.GetBytes()))
<span id="L323" class="ln">   323</span>			msg.MetaMessage.Timestamp = time.Now().Round(time.Second)
<span id="L324" class="ln">   324</span>	
<span id="L325" class="ln">   325</span>			err = localdb.AddUpdateMessage(msg, localdb.SENDBOX)
<span id="L326" class="ln">   326</span>			if err != nil {
<span id="L327" class="ln">   327</span>				return err
<span id="L328" class="ln">   328</span>			}
<span id="L329" class="ln">   329</span>	
<span id="L330" class="ln">   330</span>			sendMsg := new(objects.Message)
<span id="L331" class="ln">   331</span>			sendMsg.TxidHash = msg.MetaMessage.TxidHash
<span id="L332" class="ln">   332</span>			sendMsg.AddrHash = objects.MakeHash(recipient.Address)
<span id="L333" class="ln">   333</span>			sendMsg.Timestamp = msg.MetaMessage.Timestamp
<span id="L334" class="ln">   334</span>			sendMsg.Content = *msg.Encrypted
<span id="L335" class="ln">   335</span>	
<span id="L336" class="ln">   336</span>			service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.MSG, objects.BROADCAST, sendMsg)
<span id="L337" class="ln">   337</span>	
<span id="L338" class="ln">   338</span>			reply.IsSent = true
<span id="L339" class="ln">   339</span>		}
<span id="L340" class="ln">   340</span>	
<span id="L341" class="ln">   341</span>		<span class="comment">// Finish by setting msg&#39;s txid</span>
<span id="L342" class="ln">   342</span>		reply.TxidHash = msg.MetaMessage.TxidHash.GetBytes()
<span id="L343" class="ln">   343</span>		return nil
<span id="L344" class="ln">   344</span>	}
<span id="L345" class="ln">   345</span>	
<span id="L346" class="ln">   346</span>	func (service *EMPService) ListMessagesBySender(r *http.Request, args *string, reply *[]objects.MetaMessage) error {
<span id="L347" class="ln">   347</span>		if !basicAuth(service.Config, r) {
<span id="L348" class="ln">   348</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L349" class="ln">   349</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L350" class="ln">   350</span>		}
<span id="L351" class="ln">   351</span>	
<span id="L352" class="ln">   352</span>		*reply = localdb.GetBySender(*args)
<span id="L353" class="ln">   353</span>		return nil
<span id="L354" class="ln">   354</span>	}
<span id="L355" class="ln">   355</span>	
<span id="L356" class="ln">   356</span>	func (service *EMPService) ListMessagesByRecpient(r *http.Request, args *string, reply *[]objects.MetaMessage) error {
<span id="L357" class="ln">   357</span>		if !basicAuth(service.Config, r) {
<span id="L358" class="ln">   358</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L359" class="ln">   359</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L360" class="ln">   360</span>		}
<span id="L361" class="ln">   361</span>	
<span id="L362" class="ln">   362</span>		*reply = localdb.GetByRecipient(*args)
<span id="L363" class="ln">   363</span>		return nil
<span id="L364" class="ln">   364</span>	}
<span id="L365" class="ln">   365</span>	
<span id="L366" class="ln">   366</span>	func (service *EMPService) Inbox(r *http.Request, args *NilParam, reply *[]objects.MetaMessage) error {
<span id="L367" class="ln">   367</span>		if !basicAuth(service.Config, r) {
<span id="L368" class="ln">   368</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L369" class="ln">   369</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L370" class="ln">   370</span>		}
<span id="L371" class="ln">   371</span>	
<span id="L372" class="ln">   372</span>		*reply = localdb.GetBox(localdb.INBOX)
<span id="L373" class="ln">   373</span>		return nil
<span id="L374" class="ln">   374</span>	}
<span id="L375" class="ln">   375</span>	
<span id="L376" class="ln">   376</span>	func (service *EMPService) Outbox(r *http.Request, args *NilParam, reply *[]objects.MetaMessage) error {
<span id="L377" class="ln">   377</span>		if !basicAuth(service.Config, r) {
<span id="L378" class="ln">   378</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L379" class="ln">   379</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L380" class="ln">   380</span>		}
<span id="L381" class="ln">   381</span>	
<span id="L382" class="ln">   382</span>		*reply = localdb.GetBox(localdb.OUTBOX)
<span id="L383" class="ln">   383</span>		return nil
<span id="L384" class="ln">   384</span>	}
<span id="L385" class="ln">   385</span>	
<span id="L386" class="ln">   386</span>	func (service *EMPService) Sendbox(r *http.Request, args *NilParam, reply *[]objects.MetaMessage) error {
<span id="L387" class="ln">   387</span>		if !basicAuth(service.Config, r) {
<span id="L388" class="ln">   388</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L389" class="ln">   389</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L390" class="ln">   390</span>		}
<span id="L391" class="ln">   391</span>	
<span id="L392" class="ln">   392</span>		*reply = localdb.GetBox(localdb.SENDBOX)
<span id="L393" class="ln">   393</span>		return nil
<span id="L394" class="ln">   394</span>	}
<span id="L395" class="ln">   395</span>	
<span id="L396" class="ln">   396</span>	func (service *EMPService) GetEncrypted(r *http.Request, args *[]byte, reply *encryption.EncryptedMessage) error {
<span id="L397" class="ln">   397</span>		if !basicAuth(service.Config, r) {
<span id="L398" class="ln">   398</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L399" class="ln">   399</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L400" class="ln">   400</span>		}
<span id="L401" class="ln">   401</span>	
<span id="L402" class="ln">   402</span>		var txidHash objects.Hash
<span id="L403" class="ln">   403</span>		txidHash.FromBytes(*args)
<span id="L404" class="ln">   404</span>	
<span id="L405" class="ln">   405</span>		<span class="comment">// Get Message from Database</span>
<span id="L406" class="ln">   406</span>		msg, err := localdb.GetMessageDetail(txidHash)
<span id="L407" class="ln">   407</span>		if err != nil {
<span id="L408" class="ln">   408</span>			return err
<span id="L409" class="ln">   409</span>		}
<span id="L410" class="ln">   410</span>	
<span id="L411" class="ln">   411</span>		*reply = *msg.Encrypted
<span id="L412" class="ln">   412</span>		return nil
<span id="L413" class="ln">   413</span>	}
<span id="L414" class="ln">   414</span>	
<span id="L415" class="ln">   415</span>	func (service *EMPService) OpenMessage(r *http.Request, args *[]byte, reply *objects.FullMessage) error {
<span id="L416" class="ln">   416</span>		if !basicAuth(service.Config, r) {
<span id="L417" class="ln">   417</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L418" class="ln">   418</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L419" class="ln">   419</span>		}
<span id="L420" class="ln">   420</span>	
<span id="L421" class="ln">   421</span>		var txidHash objects.Hash
<span id="L422" class="ln">   422</span>		txidHash.FromBytes(*args)
<span id="L423" class="ln">   423</span>	
<span id="L424" class="ln">   424</span>		<span class="comment">// Get Message from Database</span>
<span id="L425" class="ln">   425</span>		msg, err := localdb.GetMessageDetail(txidHash)
<span id="L426" class="ln">   426</span>		if err != nil {
<span id="L427" class="ln">   427</span>			return err
<span id="L428" class="ln">   428</span>		}
<span id="L429" class="ln">   429</span>	
<span id="L430" class="ln">   430</span>		if msg.Encrypted == nil {
<span id="L431" class="ln">   431</span>			*reply = *msg
<span id="L432" class="ln">   432</span>			return nil
<span id="L433" class="ln">   433</span>		}
<span id="L434" class="ln">   434</span>	
<span id="L435" class="ln">   435</span>		<span class="comment">// If not decrypted, decrypt message and purge</span>
<span id="L436" class="ln">   436</span>		if msg.Decrypted == nil {
<span id="L437" class="ln">   437</span>			recipient, err := localdb.GetAddressDetail(objects.MakeHash(encryption.StringToAddress(msg.MetaMessage.Recipient)))
<span id="L438" class="ln">   438</span>			if err != nil {
<span id="L439" class="ln">   439</span>				return err
<span id="L440" class="ln">   440</span>			}
<span id="L441" class="ln">   441</span>	
<span id="L442" class="ln">   442</span>			if recipient.Privkey == nil {
<span id="L443" class="ln">   443</span>				*reply = *msg
<span id="L444" class="ln">   444</span>				return nil
<span id="L445" class="ln">   445</span>			}
<span id="L446" class="ln">   446</span>	
<span id="L447" class="ln">   447</span>			<span class="comment">// Decrypt Message</span>
<span id="L448" class="ln">   448</span>			decrypted := encryption.Decrypt(service.Config.Log, recipient.Privkey, msg.Encrypted)
<span id="L449" class="ln">   449</span>			if len(decrypted) == 0 {
<span id="L450" class="ln">   450</span>				*reply = *msg
<span id="L451" class="ln">   451</span>				return nil
<span id="L452" class="ln">   452</span>			}
<span id="L453" class="ln">   453</span>			msg.Decrypted = new(objects.DecryptedMessage)
<span id="L454" class="ln">   454</span>			msg.Decrypted.FromBytes(decrypted)
<span id="L455" class="ln">   455</span>	
<span id="L456" class="ln">   456</span>			<span class="comment">// Update Sender</span>
<span id="L457" class="ln">   457</span>	
<span id="L458" class="ln">   458</span>			
<span id="L459" class="ln">   459</span>			x, y := encryption.UnmarshalPubkey(msg.Decrypted.Pubkey[:])
<span id="L460" class="ln">   460</span>			address := encryption.GetAddress(service.Config.Log, x, y)
<span id="L461" class="ln">   461</span>			addrStr := encryption.AddressToString(address)
<span id="L462" class="ln">   462</span>			addrHash := objects.MakeHash(address)
<span id="L463" class="ln">   463</span>	
<span id="L464" class="ln">   464</span>			detail, _ := localdb.GetAddressDetail(addrHash)
<span id="L465" class="ln">   465</span>			if detail == nil {
<span id="L466" class="ln">   466</span>				detail = new(objects.AddressDetail)
<span id="L467" class="ln">   467</span>			}
<span id="L468" class="ln">   468</span>			detail.Address = address
<span id="L469" class="ln">   469</span>			detail.String = addrStr
<span id="L470" class="ln">   470</span>			detail.Pubkey = msg.Decrypted.Pubkey[:]
<span id="L471" class="ln">   471</span>	
<span id="L472" class="ln">   472</span>			localdb.AddUpdateAddress(detail)
<span id="L473" class="ln">   473</span>			msg.MetaMessage.Sender = detail.String
<span id="L474" class="ln">   474</span>	
<span id="L475" class="ln">   475</span>			<span class="comment">// Send Purge Request</span>
<span id="L476" class="ln">   476</span>			purge := new(objects.Purge)
<span id="L477" class="ln">   477</span>			purge.Txid = msg.Decrypted.Txid
<span id="L478" class="ln">   478</span>	
<span id="L479" class="ln">   479</span>			service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.PURGE, objects.BROADCAST, purge)
<span id="L480" class="ln">   480</span>			msg.MetaMessage.Purged = true
<span id="L481" class="ln">   481</span>	
<span id="L482" class="ln">   482</span>			localdb.AddUpdateMessage(msg, localdb.Contains(msg.MetaMessage.TxidHash))
<span id="L483" class="ln">   483</span>		} else {
<span id="L484" class="ln">   484</span>			if msg.MetaMessage.Purged == false &amp;&amp; localdb.Contains(txidHash) == localdb.INBOX {
<span id="L485" class="ln">   485</span>				msg.MetaMessage.Purged = true
<span id="L486" class="ln">   486</span>				localdb.AddUpdateMessage(msg, localdb.Contains(msg.MetaMessage.TxidHash))
<span id="L487" class="ln">   487</span>			}
<span id="L488" class="ln">   488</span>		}
<span id="L489" class="ln">   489</span>	
<span id="L490" class="ln">   490</span>		*reply = *msg
<span id="L491" class="ln">   491</span>		return nil
<span id="L492" class="ln">   492</span>	}
</pre><p><a href="/src/pkg/emp/local/localapi/rpcmsg.go?m=text">View as plain text</a></p>

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

