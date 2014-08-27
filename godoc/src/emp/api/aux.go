<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/api/aux.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/api/aux.go</h1>




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
<span id="L12" class="ln">    12</span>	package api
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;emp/db&#34;
<span id="L16" class="ln">    16</span>		&#34;emp/objects&#34;
<span id="L17" class="ln">    17</span>		&#34;fmt&#34;
<span id="L18" class="ln">    18</span>		&#34;quibit&#34;
<span id="L19" class="ln">    19</span>		&#34;runtime&#34;
<span id="L20" class="ln">    20</span>		&#34;time&#34;
<span id="L21" class="ln">    21</span>	)
<span id="L22" class="ln">    22</span>	
<span id="L23" class="ln">    23</span>	<span class="comment">// Handle a Version Request or Reply</span>
<span id="L24" class="ln">    24</span>	func fVERSION(config *ApiConfig, frame quibit.Frame, version *objects.Version) {
<span id="L25" class="ln">    25</span>	
<span id="L26" class="ln">    26</span>		<span class="comment">// Verify not objects.BROADCAST</span>
<span id="L27" class="ln">    27</span>		if frame.Header.Type == objects.BROADCAST {
<span id="L28" class="ln">    28</span>			<span class="comment">// SHUN THE NODE! SHUN IT WITH FIRE!</span>
<span id="L29" class="ln">    29</span>			config.Log &lt;- &#34;Node sent a version message as a broadcast. Disconnecting...&#34;
<span id="L30" class="ln">    30</span>			quibit.KillPeer(frame.Peer)
<span id="L31" class="ln">    31</span>			return
<span id="L32" class="ln">    32</span>		}
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>		<span class="comment">// Verify Protcol Version, else Disconnect</span>
<span id="L35" class="ln">    35</span>		if version.Version != objects.LOCAL_VERSION {
<span id="L36" class="ln">    36</span>			config.Log &lt;- fmt.Sprintf(&#34;Peer protocol version does not match local version: %d&#34;, version.Version)
<span id="L37" class="ln">    37</span>			quibit.KillPeer(frame.Peer)
<span id="L38" class="ln">    38</span>			return
<span id="L39" class="ln">    39</span>		}
<span id="L40" class="ln">    40</span>	
<span id="L41" class="ln">    41</span>		<span class="comment">// Verify Timestamp (5 minute window), else Disconnect</span>
<span id="L42" class="ln">    42</span>		dur := time.Since(version.Timestamp)
<span id="L43" class="ln">    43</span>		if dur.Minutes()+5 &gt; 10 {
<span id="L44" class="ln">    44</span>			config.Log &lt;- fmt.Sprintf(&#34;Peer timestamp too far off local time: %s&#34;, dur.String())
<span id="L45" class="ln">    45</span>			quibit.KillPeer(frame.Peer)
<span id="L46" class="ln">    46</span>			return
<span id="L47" class="ln">    47</span>		}
<span id="L48" class="ln">    48</span>	
<span id="L49" class="ln">    49</span>		<span class="comment">// If backbone node, verify IP</span>
<span id="L50" class="ln">    50</span>		backbone := false
<span id="L51" class="ln">    51</span>		for _, b := range []byte(version.IpAddress) {
<span id="L52" class="ln">    52</span>			if b != 0 {
<span id="L53" class="ln">    53</span>				backbone = true
<span id="L54" class="ln">    54</span>			}
<span id="L55" class="ln">    55</span>		}
<span id="L56" class="ln">    56</span>	
<span id="L57" class="ln">    57</span>		if backbone {
<span id="L58" class="ln">    58</span>			testIP := quibit.GetPeer(frame.Peer).IP
<span id="L59" class="ln">    59</span>			if version.IpAddress.String() != testIP.String() {
<span id="L60" class="ln">    60</span>				config.Log &lt;- fmt.Sprintf(&#34;Backbone node broadcast incorrect IP: %s&#34;, version.IpAddress.String())
<span id="L61" class="ln">    61</span>				quibit.KillPeer(frame.Peer)
<span id="L62" class="ln">    62</span>				return
<span id="L63" class="ln">    63</span>			}
<span id="L64" class="ln">    64</span>	
<span id="L65" class="ln">    65</span>			<span class="comment">// Add to Master Node List</span>
<span id="L66" class="ln">    66</span>			var node objects.Node
<span id="L67" class="ln">    67</span>	
<span id="L68" class="ln">    68</span>			node.IP = version.IpAddress
<span id="L69" class="ln">    69</span>			node.Port = version.Port
<span id="L70" class="ln">    70</span>			node.LastSeen = time.Now().Round(time.Second)
<span id="L71" class="ln">    71</span>			config.NodeList.Nodes[node.String()] = node
<span id="L72" class="ln">    72</span>		}
<span id="L73" class="ln">    73</span>	
<span id="L74" class="ln">    74</span>		var sending *quibit.Frame
<span id="L75" class="ln">    75</span>		if frame.Header.Type == objects.REQUEST {
<span id="L76" class="ln">    76</span>			<span class="comment">// If a objects.REQUEST, send local version as a objects.REPLY</span>
<span id="L77" class="ln">    77</span>			config.LocalVersion.Timestamp = time.Now().Round(time.Second)
<span id="L78" class="ln">    78</span>			sending = objects.MakeFrame(objects.VERSION, objects.REPLY, &amp;config.LocalVersion)
<span id="L79" class="ln">    79</span>		} else {
<span id="L80" class="ln">    80</span>			<span class="comment">// If a objects.REPLY, send a peer list as a objects.REQUEST</span>
<span id="L81" class="ln">    81</span>			sending = objects.MakeFrame(objects.PEER, objects.REQUEST, &amp;config.NodeList)
<span id="L82" class="ln">    82</span>		}
<span id="L83" class="ln">    83</span>		sending.Peer = frame.Peer
<span id="L84" class="ln">    84</span>		config.SendQueue &lt;- *sending
<span id="L85" class="ln">    85</span>	} <span class="comment">// End fVERSION</span>
<span id="L86" class="ln">    86</span>	
<span id="L87" class="ln">    87</span>	<span class="comment">// Handle Peer List Requests or Replies</span>
<span id="L88" class="ln">    88</span>	func fPEER(config *ApiConfig, frame quibit.Frame, nodeList *objects.NodeList) {
<span id="L89" class="ln">    89</span>	
<span id="L90" class="ln">    90</span>		<span class="comment">// Verify not objects.BROADCAST</span>
<span id="L91" class="ln">    91</span>		if frame.Header.Type == objects.BROADCAST {
<span id="L92" class="ln">    92</span>			<span class="comment">// SHUN THE NODE! SHUN IT WITH FIRE!</span>
<span id="L93" class="ln">    93</span>			config.Log &lt;- &#34;Node sent a peer frame as a broadcast. Disconnecting...&#34;
<span id="L94" class="ln">    94</span>			quibit.KillPeer(frame.Peer)
<span id="L95" class="ln">    95</span>			return
<span id="L96" class="ln">    96</span>		}
<span id="L97" class="ln">    97</span>	
<span id="L98" class="ln">    98</span>		var sending *quibit.Frame
<span id="L99" class="ln">    99</span>		if frame.Header.Type == objects.REQUEST {
<span id="L100" class="ln">   100</span>			<span class="comment">// If a objects.REQUEST, send back peer objects.REPLY</span>
<span id="L101" class="ln">   101</span>			sending = objects.MakeFrame(objects.PEER, objects.REPLY, &amp;config.NodeList)
<span id="L102" class="ln">   102</span>		} else {
<span id="L103" class="ln">   103</span>			<span class="comment">// If a objects.REPLY, send an object list as a objects.REQUEST</span>
<span id="L104" class="ln">   104</span>			sending = objects.MakeFrame(objects.OBJ, objects.REQUEST, db.ObjList())
<span id="L105" class="ln">   105</span>		}
<span id="L106" class="ln">   106</span>	
<span id="L107" class="ln">   107</span>		sending.Peer = frame.Peer
<span id="L108" class="ln">   108</span>	
<span id="L109" class="ln">   109</span>		config.SendQueue &lt;- *sending
<span id="L110" class="ln">   110</span>	
<span id="L111" class="ln">   111</span>		if nodeList != nil {
<span id="L112" class="ln">   112</span>	
<span id="L113" class="ln">   113</span>			<span class="comment">// Merge incoming list with current list</span>
<span id="L114" class="ln">   114</span>			for key, node := range nodeList.Nodes {
<span id="L115" class="ln">   115</span>				if node.IP.String() == config.LocalVersion.IpAddress.String() {
<span id="L116" class="ln">   116</span>					continue
<span id="L117" class="ln">   117</span>				}
<span id="L118" class="ln">   118</span>				_, ok := config.NodeList.Nodes[key]
<span id="L119" class="ln">   119</span>				if !ok {
<span id="L120" class="ln">   120</span>					config.NodeList.Nodes[key] = node
<span id="L121" class="ln">   121</span>					p := new(quibit.Peer)
<span id="L122" class="ln">   122</span>					p.IP = node.IP
<span id="L123" class="ln">   123</span>					p.Port = node.Port
<span id="L124" class="ln">   124</span>					config.PeerQueue &lt;- *p
<span id="L125" class="ln">   125</span>					runtime.Gosched()
<span id="L126" class="ln">   126</span>	
<span id="L127" class="ln">   127</span>					newVer := objects.MakeFrame(objects.VERSION, objects.REQUEST, &amp;config.LocalVersion)
<span id="L128" class="ln">   128</span>					newVer.Peer = p.String()
<span id="L129" class="ln">   129</span>	
<span id="L130" class="ln">   130</span>					config.SendQueue &lt;- *newVer
<span id="L131" class="ln">   131</span>				} <span class="comment">// End if</span>
<span id="L132" class="ln">   132</span>			} <span class="comment">// End for</span>
<span id="L133" class="ln">   133</span>	
<span id="L134" class="ln">   134</span>		}
<span id="L135" class="ln">   135</span>	} <span class="comment">// End fPEER</span>
<span id="L136" class="ln">   136</span>	
<span id="L137" class="ln">   137</span>	<span class="comment">// Handle Object Vector Requests or Replies</span>
<span id="L138" class="ln">   138</span>	func fOBJ(config *ApiConfig, frame quibit.Frame, obj *objects.Obj) {
<span id="L139" class="ln">   139</span>		var sending *quibit.Frame
<span id="L140" class="ln">   140</span>	
<span id="L141" class="ln">   141</span>		<span class="comment">// Verify not objects.BROADCAST</span>
<span id="L142" class="ln">   142</span>		if frame.Header.Type == objects.BROADCAST {
<span id="L143" class="ln">   143</span>			<span class="comment">// SHUN THE NODE! SHUN IT WITH FIRE!</span>
<span id="L144" class="ln">   144</span>			config.Log &lt;- &#34;Node sent an obj frame as a broadcast. Disconnecting...&#34;
<span id="L145" class="ln">   145</span>			quibit.KillPeer(frame.Peer)
<span id="L146" class="ln">   146</span>			return
<span id="L147" class="ln">   147</span>		}
<span id="L148" class="ln">   148</span>	
<span id="L149" class="ln">   149</span>		if frame.Header.Type == objects.REQUEST {
<span id="L150" class="ln">   150</span>			<span class="comment">// If a objects.REQUEST, send local object list as objects.REPLY</span>
<span id="L151" class="ln">   151</span>			sending = objects.MakeFrame(objects.OBJ, objects.REPLY, db.ObjList())
<span id="L152" class="ln">   152</span>			sending.Peer = frame.Peer
<span id="L153" class="ln">   153</span>			config.SendQueue &lt;- *sending
<span id="L154" class="ln">   154</span>		}
<span id="L155" class="ln">   155</span>	
<span id="L156" class="ln">   156</span>		<span class="comment">// For each object in object list:</span>
<span id="L157" class="ln">   157</span>		<span class="comment">// If object not stored locally, send GETOBJ objects.REQUEST</span>
<span id="L158" class="ln">   158</span>		for _, hash := range obj.HashList {
<span id="L159" class="ln">   159</span>			if db.Contains(hash) == db.NOTFOUND {
<span id="L160" class="ln">   160</span>				sending = objects.MakeFrame(objects.GETOBJ, objects.REQUEST, &amp;hash)
<span id="L161" class="ln">   161</span>				sending.Peer = frame.Peer
<span id="L162" class="ln">   162</span>				config.SendQueue &lt;- *sending
<span id="L163" class="ln">   163</span>			} else if db.Contains(hash) == db.MSG {
<span id="L164" class="ln">   164</span>				<span class="comment">// Check for purge</span>
<span id="L165" class="ln">   165</span>				sending = objects.MakeFrame(objects.CHECKTXID, objects.REQUEST, &amp;hash)
<span id="L166" class="ln">   166</span>				sending.Peer = frame.Peer
<span id="L167" class="ln">   167</span>				config.SendQueue &lt;- *sending
<span id="L168" class="ln">   168</span>			}
<span id="L169" class="ln">   169</span>		}
<span id="L170" class="ln">   170</span>	}
<span id="L171" class="ln">   171</span>	
<span id="L172" class="ln">   172</span>	<span class="comment">// Handle Object Detail Requests</span>
<span id="L173" class="ln">   173</span>	func fGETOBJ(config *ApiConfig, frame quibit.Frame, hash *objects.Hash) {
<span id="L174" class="ln">   174</span>		<span class="comment">// Verify not objects.BROADCAST</span>
<span id="L175" class="ln">   175</span>		if frame.Header.Type == objects.BROADCAST {
<span id="L176" class="ln">   176</span>			<span class="comment">// SHUN THE NODE! SHUN IT WITH FIRE!</span>
<span id="L177" class="ln">   177</span>			config.Log &lt;- &#34;Node sent a getobj message as a broadcast. Disconnecting...&#34;
<span id="L178" class="ln">   178</span>			quibit.KillPeer(frame.Peer)
<span id="L179" class="ln">   179</span>			return
<span id="L180" class="ln">   180</span>		}
<span id="L181" class="ln">   181</span>	
<span id="L182" class="ln">   182</span>		<span class="comment">// If object stored locally, send object as a objects.REPLY</span>
<span id="L183" class="ln">   183</span>		var sending *quibit.Frame
<span id="L184" class="ln">   184</span>		if frame.Header.Type == objects.REQUEST {
<span id="L185" class="ln">   185</span>			switch db.Contains(*hash) {
<span id="L186" class="ln">   186</span>			case db.PUBKEY:
<span id="L187" class="ln">   187</span>				sending = objects.MakeFrame(objects.PUBKEY, objects.REPLY, db.GetPubkey(config.Log, *hash))
<span id="L188" class="ln">   188</span>			case db.PURGE:
<span id="L189" class="ln">   189</span>				sending = objects.MakeFrame(objects.PURGE, objects.REPLY, db.GetPurge(config.Log, *hash))
<span id="L190" class="ln">   190</span>			case db.MSG:
<span id="L191" class="ln">   191</span>				message := db.GetMessage(config.Log, *hash)
<span id="L192" class="ln">   192</span>				if message != nil {
<span id="L193" class="ln">   193</span>					sending = objects.MakeFrame(objects.MSG, objects.REPLY, message)
<span id="L194" class="ln">   194</span>				} else {
<span id="L195" class="ln">   195</span>					config.Log &lt;- &#34;Error pulling message from database!&#34;
<span id="L196" class="ln">   196</span>				}
<span id="L197" class="ln">   197</span>			case db.PUB:
<span id="L198" class="ln">   198</span>				message := db.GetMessage(config.Log, *hash)
<span id="L199" class="ln">   199</span>				if message != nil {
<span id="L200" class="ln">   200</span>					sending = objects.MakeFrame(objects.PUB, objects.REPLY, message)
<span id="L201" class="ln">   201</span>				} else {
<span id="L202" class="ln">   202</span>					config.Log &lt;- &#34;Error pulling publication from database!&#34;
<span id="L203" class="ln">   203</span>				}
<span id="L204" class="ln">   204</span>			case db.PUBKEYRQ:
<span id="L205" class="ln">   205</span>				sending = objects.MakeFrame(objects.PUBKEY_REQUEST, objects.REPLY, hash)
<span id="L206" class="ln">   206</span>			default:
<span id="L207" class="ln">   207</span>				sending = objects.MakeFrame(objects.GETOBJ, objects.REPLY, new(objects.NilPayload))
<span id="L208" class="ln">   208</span>			} <span class="comment">// End switch</span>
<span id="L209" class="ln">   209</span>			if (sending == nil) {
<span id="L210" class="ln">   210</span>				return
<span id="L211" class="ln">   211</span>			}
<span id="L212" class="ln">   212</span>			sending.Peer = frame.Peer
<span id="L213" class="ln">   213</span>			config.SendQueue &lt;- *sending
<span id="L214" class="ln">   214</span>		} <span class="comment">// End if</span>
<span id="L215" class="ln">   215</span>	} <span class="comment">// End fGETOBJ</span>
<span id="L216" class="ln">   216</span>	
<span id="L217" class="ln">   217</span>	<span class="comment">// Handle Public Key Request Broadcasts</span>
<span id="L218" class="ln">   218</span>	func fPUBKEY_REQUEST(config *ApiConfig, frame quibit.Frame, pubHash *objects.Hash) {
<span id="L219" class="ln">   219</span>		<span class="comment">// Check Hash in Object List</span>
<span id="L220" class="ln">   220</span>		var sending quibit.Frame
<span id="L221" class="ln">   221</span>	
<span id="L222" class="ln">   222</span>		switch db.Contains(*pubHash) {
<span id="L223" class="ln">   223</span>		<span class="comment">// If request is Not in List, store the request</span>
<span id="L224" class="ln">   224</span>		case db.NOTFOUND:
<span id="L225" class="ln">   225</span>			<span class="comment">// If a objects.BROADCAST, send out another objects.BROADCAST</span>
<span id="L226" class="ln">   226</span>			db.Add(*pubHash, db.PUBKEYRQ)
<span id="L227" class="ln">   227</span>			if frame.Header.Type == objects.BROADCAST {
<span id="L228" class="ln">   228</span>				sending = *objects.MakeFrame(objects.PUBKEY_REQUEST, objects.BROADCAST, pubHash)
<span id="L229" class="ln">   229</span>				sending.Peer = frame.Peer
<span id="L230" class="ln">   230</span>				config.SendQueue &lt;- sending
<span id="L231" class="ln">   231</span>			}
<span id="L232" class="ln">   232</span>	
<span id="L233" class="ln">   233</span>		<span class="comment">// If request is a Public Key in List:</span>
<span id="L234" class="ln">   234</span>		case db.PUBKEY:
<span id="L235" class="ln">   235</span>			<span class="comment">// Send out the PUBKEY as a objects.BROADCAST</span>
<span id="L236" class="ln">   236</span>			sending = *objects.MakeFrame(objects.PUBKEY, objects.BROADCAST, db.GetPubkey(config.Log, *pubHash))
<span id="L237" class="ln">   237</span>			sending.Peer = frame.Peer
<span id="L238" class="ln">   238</span>			config.SendQueue &lt;- sending
<span id="L239" class="ln">   239</span>		}
<span id="L240" class="ln">   240</span>	}
<span id="L241" class="ln">   241</span>	
<span id="L242" class="ln">   242</span>	<span class="comment">// Handle Public Key Broadcasts</span>
<span id="L243" class="ln">   243</span>	func fPUBKEY(config *ApiConfig, frame quibit.Frame, pubkey *objects.EncryptedPubkey) {
<span id="L244" class="ln">   244</span>		<span class="comment">// Check Hash in Object List</span>
<span id="L245" class="ln">   245</span>		switch db.Contains(pubkey.AddrHash) {
<span id="L246" class="ln">   246</span>		<span class="comment">// If request is a Pubkey Request, remove the pubkey request</span>
<span id="L247" class="ln">   247</span>		case db.PUBKEYRQ:
<span id="L248" class="ln">   248</span>			db.Delete(pubkey.AddrHash)
<span id="L249" class="ln">   249</span>			fallthrough
<span id="L250" class="ln">   250</span>		case db.NOTFOUND:
<span id="L251" class="ln">   251</span>			<span class="comment">// Add Pubkey to database</span>
<span id="L252" class="ln">   252</span>			err := db.AddPubkey(config.Log, *pubkey)
<span id="L253" class="ln">   253</span>			if err != nil {
<span id="L254" class="ln">   254</span>				config.Log &lt;- fmt.Sprintf(&#34;Error adding pubkey to database: %s&#34;, err)
<span id="L255" class="ln">   255</span>				break
<span id="L256" class="ln">   256</span>			}
<span id="L257" class="ln">   257</span>			<span class="comment">// If a objects.BROADCAST, send a objects.BROADCAST</span>
<span id="L258" class="ln">   258</span>			if frame.Header.Type == objects.BROADCAST {
<span id="L259" class="ln">   259</span>				sending := *objects.MakeFrame(objects.PUBKEY, objects.BROADCAST, pubkey)
<span id="L260" class="ln">   260</span>				sending.Peer = frame.Peer
<span id="L261" class="ln">   261</span>				config.SendQueue &lt;- sending
<span id="L262" class="ln">   262</span>			}
<span id="L263" class="ln">   263</span>	
<span id="L264" class="ln">   264</span>			config.PubkeyRegister &lt;- pubkey.AddrHash
<span id="L265" class="ln">   265</span>		}
<span id="L266" class="ln">   266</span>	} <span class="comment">// End fPUBKEY</span>
<span id="L267" class="ln">   267</span>	
<span id="L268" class="ln">   268</span>	<span class="comment">// Handle Encrypted Message Broadcasts</span>
<span id="L269" class="ln">   269</span>	func fMSG(config *ApiConfig, frame quibit.Frame, msg *objects.Message) {
<span id="L270" class="ln">   270</span>		var sending quibit.Frame
<span id="L271" class="ln">   271</span>		<span class="comment">// Check Hash in Object List</span>
<span id="L272" class="ln">   272</span>		switch db.Contains(msg.TxidHash) {
<span id="L273" class="ln">   273</span>		<span class="comment">// If Not in List, Store and objects.BROADCAST</span>
<span id="L274" class="ln">   274</span>		case db.NOTFOUND:
<span id="L275" class="ln">   275</span>			err := db.AddMessage(config.Log, msg)
<span id="L276" class="ln">   276</span>			if err != nil {
<span id="L277" class="ln">   277</span>				config.Log &lt;- fmt.Sprintf(&#34;Error adding message to database: %s&#34;, err)
<span id="L278" class="ln">   278</span>				break
<span id="L279" class="ln">   279</span>			}
<span id="L280" class="ln">   280</span>			<span class="comment">// Re-broadcast unpurged message</span>
<span id="L281" class="ln">   281</span>			sending = *objects.MakeFrame(objects.MSG, objects.BROADCAST, msg)
<span id="L282" class="ln">   282</span>			sending.Peer = frame.Peer
<span id="L283" class="ln">   283</span>			config.SendQueue &lt;- sending
<span id="L284" class="ln">   284</span>	
<span id="L285" class="ln">   285</span>			config.Log &lt;- &#34;Registering message...&#34;
<span id="L286" class="ln">   286</span>			config.MessageRegister &lt;- *msg
<span id="L287" class="ln">   287</span>	
<span id="L288" class="ln">   288</span>		<span class="comment">// If found as PURGE, reply with PURGE</span>
<span id="L289" class="ln">   289</span>		case db.PURGE:
<span id="L290" class="ln">   290</span>			config.Log &lt;- &#34;Received already-purged message!&#34;
<span id="L291" class="ln">   291</span>			sending = *objects.MakeFrame(objects.PURGE, objects.REPLY, db.GetPurge(config.Log, msg.TxidHash))
<span id="L292" class="ln">   292</span>			sending.Peer = frame.Peer
<span id="L293" class="ln">   293</span>			config.SendQueue &lt;- sending
<span id="L294" class="ln">   294</span>		}
<span id="L295" class="ln">   295</span>	} <span class="comment">// End fMSG</span>
<span id="L296" class="ln">   296</span>	
<span id="L297" class="ln">   297</span>	<span class="comment">// Handle Encrypted Publication Broadcasts</span>
<span id="L298" class="ln">   298</span>	func fPUB(config *ApiConfig, frame quibit.Frame, msg *objects.Message) {
<span id="L299" class="ln">   299</span>		var sending quibit.Frame
<span id="L300" class="ln">   300</span>		<span class="comment">// Check Hash in Object List</span>
<span id="L301" class="ln">   301</span>		switch db.Contains(msg.TxidHash) {
<span id="L302" class="ln">   302</span>		<span class="comment">// If Not in List, Store and objects.BROADCAST</span>
<span id="L303" class="ln">   303</span>		case db.NOTFOUND:
<span id="L304" class="ln">   304</span>			err := db.AddPub(config.Log, msg)
<span id="L305" class="ln">   305</span>			if err != nil {
<span id="L306" class="ln">   306</span>				config.Log &lt;- fmt.Sprintf(&#34;Error adding publication to database: %s&#34;, err)
<span id="L307" class="ln">   307</span>				break
<span id="L308" class="ln">   308</span>			}
<span id="L309" class="ln">   309</span>			<span class="comment">// Re-broadcast</span>
<span id="L310" class="ln">   310</span>				sending = *objects.MakeFrame(objects.PUB, objects.BROADCAST, msg)
<span id="L311" class="ln">   311</span>				sending.Peer = frame.Peer
<span id="L312" class="ln">   312</span>				config.SendQueue &lt;- sending
<span id="L313" class="ln">   313</span>			config.Log &lt;- &#34;Registering publication...&#34;
<span id="L314" class="ln">   314</span>			config.PubRegister &lt;- *msg
<span id="L315" class="ln">   315</span>	
<span id="L316" class="ln">   316</span>		<span class="comment">// If found as PURGE, reply with PURGE</span>
<span id="L317" class="ln">   317</span>		case db.PURGE:
<span id="L318" class="ln">   318</span>			config.Log &lt;- &#34;Received already-purged publication!&#34;
<span id="L319" class="ln">   319</span>			sending = *objects.MakeFrame(objects.PURGE, objects.REPLY, db.GetPurge(config.Log, msg.TxidHash))
<span id="L320" class="ln">   320</span>			sending.Peer = frame.Peer
<span id="L321" class="ln">   321</span>			config.SendQueue &lt;- sending
<span id="L322" class="ln">   322</span>		}
<span id="L323" class="ln">   323</span>	} <span class="comment">// End fMSG</span>
<span id="L324" class="ln">   324</span>	
<span id="L325" class="ln">   325</span>	<span class="comment">// Handle Purge Broadcasts</span>
<span id="L326" class="ln">   326</span>	func fPURGE(config *ApiConfig, frame quibit.Frame, purge *objects.Purge) {
<span id="L327" class="ln">   327</span>		var err error
<span id="L328" class="ln">   328</span>		txidHash := objects.MakeHash(purge.Txid[:])
<span id="L329" class="ln">   329</span>	
<span id="L330" class="ln">   330</span>		<span class="comment">// Check Hash in Object List</span>
<span id="L331" class="ln">   331</span>		switch db.Contains(txidHash) {
<span id="L332" class="ln">   332</span>		<span class="comment">// Delete Stored Messages</span>
<span id="L333" class="ln">   333</span>		case db.PUB:
<span id="L334" class="ln">   334</span>			fallthrough
<span id="L335" class="ln">   335</span>		case db.MSG:
<span id="L336" class="ln">   336</span>			err = db.RemoveHash(config.Log, txidHash)
<span id="L337" class="ln">   337</span>			if err != nil {
<span id="L338" class="ln">   338</span>				config.Log &lt;- fmt.Sprintf(&#34;Error removing message/publication from database: %s&#34;, err)
<span id="L339" class="ln">   339</span>				break
<span id="L340" class="ln">   340</span>			}
<span id="L341" class="ln">   341</span>			fallthrough
<span id="L342" class="ln">   342</span>		<span class="comment">// Add to database</span>
<span id="L343" class="ln">   343</span>		case db.NOTFOUND:
<span id="L344" class="ln">   344</span>			err = db.AddPurge(config.Log, *purge)
<span id="L345" class="ln">   345</span>			if err != nil {
<span id="L346" class="ln">   346</span>				config.Log &lt;- fmt.Sprintf(&#34;Error adding purge to database: &#34;, err)
<span id="L347" class="ln">   347</span>				break
<span id="L348" class="ln">   348</span>			}
<span id="L349" class="ln">   349</span>	
<span id="L350" class="ln">   350</span>			<span class="comment">// Re-objects.BROADCAST if necessary</span>
<span id="L351" class="ln">   351</span>				sending := *objects.MakeFrame(objects.PURGE, objects.BROADCAST, purge)
<span id="L352" class="ln">   352</span>				sending.Peer = frame.Peer
<span id="L353" class="ln">   353</span>				config.SendQueue &lt;- sending
<span id="L354" class="ln">   354</span>			config.PurgeRegister &lt;- purge.Txid
<span id="L355" class="ln">   355</span>		} <span class="comment">// End Switch</span>
<span id="L356" class="ln">   356</span>	} <span class="comment">// End fPURGE</span>
<span id="L357" class="ln">   357</span>	
<span id="L358" class="ln">   358</span>	func fCHECKTXID(config *ApiConfig, frame quibit.Frame, hash *objects.Hash) {
<span id="L359" class="ln">   359</span>		<span class="comment">// Verify not objects.BROADCAST</span>
<span id="L360" class="ln">   360</span>		if frame.Header.Type == objects.BROADCAST {
<span id="L361" class="ln">   361</span>			<span class="comment">// SHUN THE NODE! SHUN IT WITH FIRE!</span>
<span id="L362" class="ln">   362</span>			config.Log &lt;- &#34;Node sent a checktxid frame as a broadcast. Disconnecting...&#34;
<span id="L363" class="ln">   363</span>			quibit.KillPeer(frame.Peer)
<span id="L364" class="ln">   364</span>			return
<span id="L365" class="ln">   365</span>		}
<span id="L366" class="ln">   366</span>	
<span id="L367" class="ln">   367</span>		<span class="comment">// If object stored locally, send object as a objects.REPLY</span>
<span id="L368" class="ln">   368</span>		var sending *quibit.Frame
<span id="L369" class="ln">   369</span>		if frame.Header.Type == objects.REQUEST {
<span id="L370" class="ln">   370</span>			if db.Contains(*hash) == db.PURGE {
<span id="L371" class="ln">   371</span>				sending = objects.MakeFrame(objects.PURGE, objects.REPLY, db.GetPurge(config.Log, *hash))
<span id="L372" class="ln">   372</span>				sending.Peer = frame.Peer
<span id="L373" class="ln">   373</span>				config.SendQueue &lt;- *sending
<span id="L374" class="ln">   374</span>			} else {
<span id="L375" class="ln">   375</span>				sending = objects.MakeFrame(objects.CHECKTXID, objects.REPLY, new(objects.NilPayload))
<span id="L376" class="ln">   376</span>				sending.Peer = frame.Peer
<span id="L377" class="ln">   377</span>				config.SendQueue &lt;- *sending
<span id="L378" class="ln">   378</span>			}
<span id="L379" class="ln">   379</span>		}
<span id="L380" class="ln">   380</span>	} <span class="comment">// End fCHECKTXID</span>
</pre><p><a href="/src/emp/api/aux.go?m=text">View as plain text</a></p>

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
