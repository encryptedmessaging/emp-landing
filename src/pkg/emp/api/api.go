<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/api/api.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/api/api.go</h1>




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
<span id="L12" class="ln">    12</span>	<span class="comment">// Package api provides a TCP server that fully implements the EMProtocol.</span>
<span id="L13" class="ln">    13</span>	package api
<span id="L14" class="ln">    14</span>	
<span id="L15" class="ln">    15</span>	import (
<span id="L16" class="ln">    16</span>		&#34;emp/db&#34;
<span id="L17" class="ln">    17</span>		&#34;emp/objects&#34;
<span id="L18" class="ln">    18</span>		&#34;fmt&#34;
<span id="L19" class="ln">    19</span>		&#34;quibit&#34;
<span id="L20" class="ln">    20</span>		&#34;runtime&#34;
<span id="L21" class="ln">    21</span>		&#34;time&#34;
<span id="L22" class="ln">    22</span>	)
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	<span class="comment">// Starts a new TCP Server wth configuration specified in ApiConfig. </span>
<span id="L25" class="ln">    25</span>	<span class="comment">// Server will terminate cleanly only when data is sent to the Quit channel.</span>
<span id="L26" class="ln">    26</span>	<span class="comment">// </span>
<span id="L27" class="ln">    27</span>	<span class="comment">// See (struct ApiConfig) for details.</span>
<span id="L28" class="ln">    28</span>	func Start(config *ApiConfig) {
<span id="L29" class="ln">    29</span>		var err error
<span id="L30" class="ln">    30</span>		var frame quibit.Frame
<span id="L31" class="ln">    31</span>	
<span id="L32" class="ln">    32</span>		defer quit(config)
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>		config.Log &lt;- &#34;Starting api...&#34;
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>		<span class="comment">// Start Database Services</span>
<span id="L37" class="ln">    37</span>		err = db.Initialize(config.Log, config.DbFile)
<span id="L38" class="ln">    38</span>		defer db.Cleanup()
<span id="L39" class="ln">    39</span>		if err != nil {
<span id="L40" class="ln">    40</span>			config.Log &lt;- fmt.Sprintf(&#34;Error initializing database: %s&#34;, err)
<span id="L41" class="ln">    41</span>			config.Log &lt;- &#34;Quit&#34;
<span id="L42" class="ln">    42</span>			return
<span id="L43" class="ln">    43</span>		}
<span id="L44" class="ln">    44</span>		config.LocalVersion.Timestamp = time.Now().Round(time.Second)
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>		locVersion := objects.MakeFrame(objects.VERSION, objects.REQUEST, &amp;config.LocalVersion)
<span id="L47" class="ln">    47</span>		for str, _ := range config.NodeList.Nodes {
<span id="L48" class="ln">    48</span>			locVersion.Peer = str
<span id="L49" class="ln">    49</span>			config.SendQueue &lt;- *locVersion
<span id="L50" class="ln">    50</span>		}
<span id="L51" class="ln">    51</span>	
<span id="L52" class="ln">    52</span>		<span class="comment">// Set Up Clocks</span>
<span id="L53" class="ln">    53</span>		second := time.Tick(2 * time.Second)
<span id="L54" class="ln">    54</span>		minute := time.Tick(time.Minute)
<span id="L55" class="ln">    55</span>	
<span id="L56" class="ln">    56</span>		for {
<span id="L57" class="ln">    57</span>			select {
<span id="L58" class="ln">    58</span>			case frame = &lt;-config.RecvQueue:
<span id="L59" class="ln">    59</span>				config.Log &lt;- fmt.Sprintf(&#34;Received %s frame...&#34;, CmdString(frame.Header.Command))
<span id="L60" class="ln">    60</span>				switch frame.Header.Command {
<span id="L61" class="ln">    61</span>				case objects.VERSION:
<span id="L62" class="ln">    62</span>					version := new(objects.Version)
<span id="L63" class="ln">    63</span>					err = version.FromBytes(frame.Payload)
<span id="L64" class="ln">    64</span>					if err != nil {
<span id="L65" class="ln">    65</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing version: %s&#34;, err)
<span id="L66" class="ln">    66</span>					} else {
<span id="L67" class="ln">    67</span>						fVERSION(config, frame, version)
<span id="L68" class="ln">    68</span>					}
<span id="L69" class="ln">    69</span>				case objects.PEER:
<span id="L70" class="ln">    70</span>					nodeList := new(objects.NodeList)
<span id="L71" class="ln">    71</span>					err = nodeList.FromBytes(frame.Payload)
<span id="L72" class="ln">    72</span>					if err != nil {
<span id="L73" class="ln">    73</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing peer list: %s&#34;, err)
<span id="L74" class="ln">    74</span>					} else {
<span id="L75" class="ln">    75</span>						fPEER(config, frame, nodeList)
<span id="L76" class="ln">    76</span>					}
<span id="L77" class="ln">    77</span>				case objects.OBJ:
<span id="L78" class="ln">    78</span>					obj := new(objects.Obj)
<span id="L79" class="ln">    79</span>					err = obj.FromBytes(frame.Payload)
<span id="L80" class="ln">    80</span>					if err != nil {
<span id="L81" class="ln">    81</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing obj list: %s&#34;, err)
<span id="L82" class="ln">    82</span>					} else {
<span id="L83" class="ln">    83</span>						fOBJ(config, frame, obj)
<span id="L84" class="ln">    84</span>					}
<span id="L85" class="ln">    85</span>				case objects.GETOBJ:
<span id="L86" class="ln">    86</span>					getObj := new(objects.Hash)
<span id="L87" class="ln">    87</span>					if len(frame.Payload) == 0 {
<span id="L88" class="ln">    88</span>						break
<span id="L89" class="ln">    89</span>					}
<span id="L90" class="ln">    90</span>					err = getObj.FromBytes(frame.Payload)
<span id="L91" class="ln">    91</span>					if err != nil {
<span id="L92" class="ln">    92</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing getobj hash: %s&#34;, err)
<span id="L93" class="ln">    93</span>					} else {
<span id="L94" class="ln">    94</span>						fGETOBJ(config, frame, getObj)
<span id="L95" class="ln">    95</span>					}
<span id="L96" class="ln">    96</span>				case objects.PUBKEY_REQUEST:
<span id="L97" class="ln">    97</span>					pubReq := new(objects.Hash)
<span id="L98" class="ln">    98</span>					err = pubReq.FromBytes(frame.Payload)
<span id="L99" class="ln">    99</span>					if err != nil {
<span id="L100" class="ln">   100</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing pubkey request hash: %s&#34;, err)
<span id="L101" class="ln">   101</span>					} else {
<span id="L102" class="ln">   102</span>						fPUBKEY_REQUEST(config, frame, pubReq)
<span id="L103" class="ln">   103</span>					}
<span id="L104" class="ln">   104</span>				case objects.PUBKEY:
<span id="L105" class="ln">   105</span>					pub := new(objects.EncryptedPubkey)
<span id="L106" class="ln">   106</span>					err = pub.FromBytes(frame.Payload)
<span id="L107" class="ln">   107</span>					if err != nil {
<span id="L108" class="ln">   108</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing pubkey: %s&#34;, err)
<span id="L109" class="ln">   109</span>					} else {
<span id="L110" class="ln">   110</span>						fPUBKEY(config, frame, pub)
<span id="L111" class="ln">   111</span>					}
<span id="L112" class="ln">   112</span>				case objects.MSG:
<span id="L113" class="ln">   113</span>					msg := new(objects.Message)
<span id="L114" class="ln">   114</span>					err = msg.FromBytes(frame.Payload)
<span id="L115" class="ln">   115</span>					if err != nil {
<span id="L116" class="ln">   116</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing message: %s&#34;, err)
<span id="L117" class="ln">   117</span>					} else {
<span id="L118" class="ln">   118</span>						fMSG(config, frame, msg)
<span id="L119" class="ln">   119</span>					}
<span id="L120" class="ln">   120</span>				case objects.PUB:
<span id="L121" class="ln">   121</span>					msg := new(objects.Message)
<span id="L122" class="ln">   122</span>					err = msg.FromBytes(frame.Payload)
<span id="L123" class="ln">   123</span>					if err != nil {
<span id="L124" class="ln">   124</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing publication: %s&#34;, err)
<span id="L125" class="ln">   125</span>					} else {
<span id="L126" class="ln">   126</span>						fPUB(config, frame, msg)
<span id="L127" class="ln">   127</span>					}
<span id="L128" class="ln">   128</span>				case objects.PURGE:
<span id="L129" class="ln">   129</span>					purge := new(objects.Purge)
<span id="L130" class="ln">   130</span>					err = purge.FromBytes(frame.Payload)
<span id="L131" class="ln">   131</span>					if err != nil {
<span id="L132" class="ln">   132</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing purge: %s&#34;, err)
<span id="L133" class="ln">   133</span>					} else {
<span id="L134" class="ln">   134</span>						fPURGE(config, frame, purge)
<span id="L135" class="ln">   135</span>					}
<span id="L136" class="ln">   136</span>				case objects.CHECKTXID:
<span id="L137" class="ln">   137</span>					chkTxid := new(objects.Hash)
<span id="L138" class="ln">   138</span>					if len(frame.Payload) == 0 {
<span id="L139" class="ln">   139</span>						break
<span id="L140" class="ln">   140</span>					}
<span id="L141" class="ln">   141</span>					err = chkTxid.FromBytes(frame.Payload)
<span id="L142" class="ln">   142</span>					if err != nil {
<span id="L143" class="ln">   143</span>						config.Log &lt;- fmt.Sprintf(&#34;Error parsing checktxid hash: %s&#34;, err)
<span id="L144" class="ln">   144</span>					} else {
<span id="L145" class="ln">   145</span>						fCHECKTXID(config, frame, chkTxid)
<span id="L146" class="ln">   146</span>					}
<span id="L147" class="ln">   147</span>				default:
<span id="L148" class="ln">   148</span>					config.Log &lt;- fmt.Sprintf(&#34;Received invalid frame for command: %d&#34;, frame.Header.Command)
<span id="L149" class="ln">   149</span>				}
<span id="L150" class="ln">   150</span>			case &lt;-config.Quit:
<span id="L151" class="ln">   151</span>				fmt.Println()
<span id="L152" class="ln">   152</span>				<span class="comment">// Dump Nodes to File</span>
<span id="L153" class="ln">   153</span>				DumpNodes(config)
<span id="L154" class="ln">   154</span>				return
<span id="L155" class="ln">   155</span>			case &lt;-second:
<span id="L156" class="ln">   156</span>				<span class="comment">// Reconnection Logic</span>
<span id="L157" class="ln">   157</span>				for key, node := range config.NodeList.Nodes {
<span id="L158" class="ln">   158</span>					peer := quibit.GetPeer(key)
<span id="L159" class="ln">   159</span>					if peer == nil || !peer.IsConnected() {
<span id="L160" class="ln">   160</span>						quibit.KillPeer(key)
<span id="L161" class="ln">   161</span>						if node.Attempts &gt;= 3 {
<span id="L162" class="ln">   162</span>							config.Log &lt;- fmt.Sprintf(&#34;Max connection attempts reached for %s, disconnecting...&#34;, key)
<span id="L163" class="ln">   163</span>							<span class="comment">// Max Attempts Reached, disconnect</span>
<span id="L164" class="ln">   164</span>							delete(config.NodeList.Nodes, key)
<span id="L165" class="ln">   165</span>						} else {
<span id="L166" class="ln">   166</span>							config.Log &lt;- fmt.Sprintf(&#34;Disconnected from peer %s, trying to reconnect...&#34;, key)
<span id="L167" class="ln">   167</span>							peer = new(quibit.Peer)
<span id="L168" class="ln">   168</span>							peer.IP = node.IP
<span id="L169" class="ln">   169</span>							peer.Port = node.Port
<span id="L170" class="ln">   170</span>							config.PeerQueue &lt;- *peer
<span id="L171" class="ln">   171</span>							runtime.Gosched()
<span id="L172" class="ln">   172</span>							peer = nil
<span id="L173" class="ln">   173</span>							node.Attempts++
<span id="L174" class="ln">   174</span>							config.NodeList.Nodes[key] = node
<span id="L175" class="ln">   175</span>							locVersion.Peer = key
<span id="L176" class="ln">   176</span>							config.SendQueue &lt;- *locVersion
<span id="L177" class="ln">   177</span>						}
<span id="L178" class="ln">   178</span>					}
<span id="L179" class="ln">   179</span>				}
<span id="L180" class="ln">   180</span>	
<span id="L181" class="ln">   181</span>				if len(config.NodeList.Nodes) &lt; 1 {
<span id="L182" class="ln">   182</span>					config.Log &lt;- &#34;All connections lost, re-bootstrapping...&#34;
<span id="L183" class="ln">   183</span>	
<span id="L184" class="ln">   184</span>					for i, str := range config.Bootstrap {
<span id="L185" class="ln">   185</span>						if i &gt;= bufLen {
<span id="L186" class="ln">   186</span>							break
<span id="L187" class="ln">   187</span>						}
<span id="L188" class="ln">   188</span>	
<span id="L189" class="ln">   189</span>						p := new(quibit.Peer)
<span id="L190" class="ln">   190</span>						n := new(objects.Node)
<span id="L191" class="ln">   191</span>						err := n.FromString(str)
<span id="L192" class="ln">   192</span>						if err != nil {
<span id="L193" class="ln">   193</span>							fmt.Println(&#34;Error Decoding Peer &#34;, str, &#34;: &#34;, err)
<span id="L194" class="ln">   194</span>							continue
<span id="L195" class="ln">   195</span>						}
<span id="L196" class="ln">   196</span>	
<span id="L197" class="ln">   197</span>						p.IP = n.IP
<span id="L198" class="ln">   198</span>						p.Port = n.Port
<span id="L199" class="ln">   199</span>						config.PeerQueue &lt;- *p
<span id="L200" class="ln">   200</span>						runtime.Gosched()
<span id="L201" class="ln">   201</span>						config.NodeList.Nodes[n.String()] = *n
<span id="L202" class="ln">   202</span>					}
<span id="L203" class="ln">   203</span>	
<span id="L204" class="ln">   204</span>					for str, _ := range config.NodeList.Nodes {
<span id="L205" class="ln">   205</span>						locVersion.Peer = str
<span id="L206" class="ln">   206</span>						config.SendQueue &lt;- *locVersion
<span id="L207" class="ln">   207</span>					}
<span id="L208" class="ln">   208</span>				}
<span id="L209" class="ln">   209</span>			case &lt;-minute:
<span id="L210" class="ln">   210</span>				<span class="comment">// Dump old messages</span>
<span id="L211" class="ln">   211</span>				err = db.SweepMessages(30 * 24 * time.Hour)
<span id="L212" class="ln">   212</span>				if err != nil {
<span id="L213" class="ln">   213</span>					config.Log &lt;- fmt.Sprintf(&#34;Error Sweeping Messages: %s&#34;, err)
<span id="L214" class="ln">   214</span>				}
<span id="L215" class="ln">   215</span>			}
<span id="L216" class="ln">   216</span>		}
<span id="L217" class="ln">   217</span>	
<span id="L218" class="ln">   218</span>		<span class="comment">// Should NEVER get here!</span>
<span id="L219" class="ln">   219</span>		panic(&#34;Must&#39;ve been a cosmic ray!&#34;)
<span id="L220" class="ln">   220</span>	}
<span id="L221" class="ln">   221</span>	
<span id="L222" class="ln">   222</span>	func quit(config *ApiConfig) {
<span id="L223" class="ln">   223</span>		config.Log &lt;- &#34;Quit&#34;
<span id="L224" class="ln">   224</span>	}
</pre><p><a href="/src/pkg/emp/api/api.go?m=text">View as plain text</a></p>

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

