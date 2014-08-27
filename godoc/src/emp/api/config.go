<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/api/config.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/api/config.go</h1>




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
<span id="L15" class="ln">    15</span>		&#34;emp/objects&#34;
<span id="L16" class="ln">    16</span>		&#34;fmt&#34;
<span id="L17" class="ln">    17</span>		&#34;github.com/BurntSushi/toml&#34;
<span id="L18" class="ln">    18</span>		&#34;net&#34;
<span id="L19" class="ln">    19</span>		&#34;os&#34;
<span id="L20" class="ln">    20</span>		&#34;bufio&#34;
<span id="L21" class="ln">    21</span>		&#34;os/user&#34;
<span id="L22" class="ln">    22</span>		&#34;quibit&#34;
<span id="L23" class="ln">    23</span>		&#34;time&#34;
<span id="L24" class="ln">    24</span>		&#34;io/ioutil&#34;
<span id="L25" class="ln">    25</span>	)
<span id="L26" class="ln">    26</span>	
<span id="L27" class="ln">    27</span>	var confDir string;
<span id="L28" class="ln">    28</span>	
<span id="L29" class="ln">    29</span>	type ApiConfig struct {
<span id="L30" class="ln">    30</span>		<span class="comment">// Network Channels</span>
<span id="L31" class="ln">    31</span>		RecvQueue chan quibit.Frame <span class="comment">// Send frames here to be handled by the Running API</span>
<span id="L32" class="ln">    32</span>		SendQueue chan quibit.Frame <span class="comment">// Frames to be broadcast to the network are sent here</span>
<span id="L33" class="ln">    33</span>		PeerQueue chan quibit.Peer <span class="comment">// New peers to connect to are sent here</span>
<span id="L34" class="ln">    34</span>	
<span id="L35" class="ln">    35</span>		<span class="comment">// Local Logic</span>
<span id="L36" class="ln">    36</span>		DbFile       string <span class="comment">// Inventory File relative to Config Directory</span>
<span id="L37" class="ln">    37</span>		LocalDB      string <span class="comment">// EMPLocal Database relative to Config Directory</span>
<span id="L38" class="ln">    38</span>		NodeFile     string <span class="comment">// File to store list of &lt;IP&gt;:&lt;Host&gt; Strings</span>
<span id="L39" class="ln">    39</span>		NodeList     objects.NodeList <span class="comment">// Active list of connected backbone nodes.</span>
<span id="L40" class="ln">    40</span>		LocalVersion objects.Version <span class="comment">// Local version broadcast to nodes upon connection</span>
<span id="L41" class="ln">    41</span>		Bootstrap    []string <span class="comment">// List of bootstrap nodes to use when all other nodes are disconnected.</span>
<span id="L42" class="ln">    42</span>	
<span id="L43" class="ln">    43</span>		<span class="comment">// Local Register</span>
<span id="L44" class="ln">    44</span>		PubkeyRegister  chan objects.Hash <span class="comment">// Identifiers for incoming encrypted public keys are sent here.</span>
<span id="L45" class="ln">    45</span>		MessageRegister chan objects.Message <span class="comment">// Incoming basic messages are copied here.</span>
<span id="L46" class="ln">    46</span>		PubRegister     chan objects.Message <span class="comment">// Incoming published messages are copied here.</span>
<span id="L47" class="ln">    47</span>		PurgeRegister   chan [16]byte <span class="comment">// Incomping purge tokens are copied here.</span>
<span id="L48" class="ln">    48</span>	
<span id="L49" class="ln">    49</span>		<span class="comment">// Administration</span>
<span id="L50" class="ln">    50</span>		Log  chan string <span class="comment">// Error messages are sent here. WILL BLOCK if not </span>
<span id="L51" class="ln">    51</span>		Quit chan os.Signal <span class="comment">// Send data here to cleanly quit the API Server</span>
<span id="L52" class="ln">    52</span>	
<span id="L53" class="ln">    53</span>		<span class="comment">// Network</span>
<span id="L54" class="ln">    54</span>		RPCPort uint16 <span class="comment">// Port to run RPC API and EMPLocal Client</span>
<span id="L55" class="ln">    55</span>		RPCUser string <span class="comment">// Username for RPC server</span>
<span id="L56" class="ln">    56</span>		RPCPass string <span class="comment">// Password for RPC Server</span>
<span id="L57" class="ln">    57</span>	
<span id="L58" class="ln">    58</span>		HttpRoot string <span class="comment">// HTML Root of EMPLocal Client</span>
<span id="L59" class="ln">    59</span>	}
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>	<span class="comment">// Returns Human-Readable string for a specific EMP command.</span>
<span id="L62" class="ln">    62</span>	func CmdString(cmd uint8) string {
<span id="L63" class="ln">    63</span>		var ret string
<span id="L64" class="ln">    64</span>	
<span id="L65" class="ln">    65</span>		switch cmd {
<span id="L66" class="ln">    66</span>		case objects.VERSION:
<span id="L67" class="ln">    67</span>			ret = &#34;version&#34;
<span id="L68" class="ln">    68</span>		case objects.PEER:
<span id="L69" class="ln">    69</span>			ret = &#34;peer list&#34;
<span id="L70" class="ln">    70</span>		case objects.OBJ:
<span id="L71" class="ln">    71</span>			ret = &#34;object vector&#34;
<span id="L72" class="ln">    72</span>		case objects.GETOBJ:
<span id="L73" class="ln">    73</span>			ret = &#34;object request&#34;
<span id="L74" class="ln">    74</span>		case objects.PUBKEY_REQUEST:
<span id="L75" class="ln">    75</span>			ret = &#34;public key request&#34;
<span id="L76" class="ln">    76</span>		case objects.PUBKEY:
<span id="L77" class="ln">    77</span>			ret = &#34;public key&#34;
<span id="L78" class="ln">    78</span>		case objects.MSG:
<span id="L79" class="ln">    79</span>			ret = &#34;encrypted message&#34;
<span id="L80" class="ln">    80</span>		case objects.PUB:
<span id="L81" class="ln">    81</span>			ret = &#34;encrypted publication&#34;
<span id="L82" class="ln">    82</span>		case objects.PURGE:
<span id="L83" class="ln">    83</span>			ret = &#34;purge notification&#34;
<span id="L84" class="ln">    84</span>		case objects.CHECKTXID:
<span id="L85" class="ln">    85</span>			ret = &#34;purge check&#34;
<span id="L86" class="ln">    86</span>		default:
<span id="L87" class="ln">    87</span>			ret = &#34;unknown&#34;
<span id="L88" class="ln">    88</span>		}
<span id="L89" class="ln">    89</span>	
<span id="L90" class="ln">    90</span>		return ret
<span id="L91" class="ln">    91</span>	}
<span id="L92" class="ln">    92</span>	
<span id="L93" class="ln">    93</span>	const (
<span id="L94" class="ln">    94</span>		bufLen = 10
<span id="L95" class="ln">    95</span>	)
<span id="L96" class="ln">    96</span>	
<span id="L97" class="ln">    97</span>	type tomlConfig struct {
<span id="L98" class="ln">    98</span>		Inventory string `toml:&#34;inventory&#34;`
<span id="L99" class="ln">    99</span>		Local     string `toml:&#34;local&#34;`
<span id="L100" class="ln">   100</span>		Nodes     string `toml:&#34;nodes&#34;`
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		IP   string
<span id="L103" class="ln">   103</span>		Port uint16
<span id="L104" class="ln">   104</span>	
<span id="L105" class="ln">   105</span>		Peers []string `toml:&#34;bootstrap&#34;`
<span id="L106" class="ln">   106</span>	
<span id="L107" class="ln">   107</span>		RPCConf rpcConf `toml:&#34;rpc&#34;`
<span id="L108" class="ln">   108</span>	}
<span id="L109" class="ln">   109</span>	
<span id="L110" class="ln">   110</span>	type rpcConf struct {
<span id="L111" class="ln">   111</span>		User  string
<span id="L112" class="ln">   112</span>		Pass  string
<span id="L113" class="ln">   113</span>		Port  uint16
<span id="L114" class="ln">   114</span>		Local string `toml:&#34;local_client&#34;`
<span id="L115" class="ln">   115</span>	}
<span id="L116" class="ln">   116</span>	
<span id="L117" class="ln">   117</span>	<span class="comment">// Set Config Directory where databases and configuration are stored.</span>
<span id="L118" class="ln">   118</span>	func SetConfDir(conf string) {
<span id="L119" class="ln">   119</span>		confDir = conf
<span id="L120" class="ln">   120</span>	}
<span id="L121" class="ln">   121</span>	
<span id="L122" class="ln">   122</span>	<span class="comment">// Get Config Directory: Defaults to $(HOME)/.config/emp/</span>
<span id="L123" class="ln">   123</span>	func GetConfDir() string {
<span id="L124" class="ln">   124</span>		if len(confDir) != 0 {
<span id="L125" class="ln">   125</span>			return confDir
<span id="L126" class="ln">   126</span>		}
<span id="L127" class="ln">   127</span>	
<span id="L128" class="ln">   128</span>		usr, err := user.Current()
<span id="L129" class="ln">   129</span>		if err != nil {
<span id="L130" class="ln">   130</span>			return &#34;./&#34;
<span id="L131" class="ln">   131</span>		}
<span id="L132" class="ln">   132</span>	
<span id="L133" class="ln">   133</span>		return usr.HomeDir + &#34;/.config/emp/&#34;
<span id="L134" class="ln">   134</span>	}
<span id="L135" class="ln">   135</span>	
<span id="L136" class="ln">   136</span>	<span class="comment">// Generate new config from configuration file. File provided as an Absolute Path.</span>
<span id="L137" class="ln">   137</span>	func GetConfig(confFile string) *ApiConfig {
<span id="L138" class="ln">   138</span>	
<span id="L139" class="ln">   139</span>		var tomlConf tomlConfig
<span id="L140" class="ln">   140</span>	
<span id="L141" class="ln">   141</span>		if _, err := toml.DecodeFile(confFile, &amp;tomlConf); err != nil {
<span id="L142" class="ln">   142</span>			fmt.Println(&#34;Config Error: &#34;, err)
<span id="L143" class="ln">   143</span>			return nil
<span id="L144" class="ln">   144</span>		}
<span id="L145" class="ln">   145</span>	
<span id="L146" class="ln">   146</span>		config := new(ApiConfig)
<span id="L147" class="ln">   147</span>	
<span id="L148" class="ln">   148</span>		<span class="comment">// Network Channels</span>
<span id="L149" class="ln">   149</span>		config.RecvQueue = make(chan quibit.Frame, bufLen)
<span id="L150" class="ln">   150</span>		config.SendQueue = make(chan quibit.Frame, bufLen)
<span id="L151" class="ln">   151</span>		config.PeerQueue = make(chan quibit.Peer, bufLen)
<span id="L152" class="ln">   152</span>	
<span id="L153" class="ln">   153</span>		<span class="comment">// Local Logic</span>
<span id="L154" class="ln">   154</span>		config.DbFile = GetConfDir() + tomlConf.Inventory
<span id="L155" class="ln">   155</span>		config.LocalDB = GetConfDir() + tomlConf.Local
<span id="L156" class="ln">   156</span>		if len(config.DbFile) == 0 || len(config.LocalDB) == 0 {
<span id="L157" class="ln">   157</span>			fmt.Println(&#34;Database file not found in config!&#34;)
<span id="L158" class="ln">   158</span>			return nil
<span id="L159" class="ln">   159</span>		}
<span id="L160" class="ln">   160</span>		config.NodeFile = GetConfDir() + tomlConf.Nodes
<span id="L161" class="ln">   161</span>	
<span id="L162" class="ln">   162</span>		config.LocalVersion.Port = tomlConf.Port
<span id="L163" class="ln">   163</span>		if tomlConf.IP != &#34;0.0.0.0&#34; {
<span id="L164" class="ln">   164</span>			config.LocalVersion.IpAddress = net.ParseIP(tomlConf.IP)
<span id="L165" class="ln">   165</span>		}
<span id="L166" class="ln">   166</span>		config.LocalVersion.Timestamp = time.Now().Round(time.Second)
<span id="L167" class="ln">   167</span>		config.LocalVersion.Version = objects.LOCAL_VERSION
<span id="L168" class="ln">   168</span>		config.LocalVersion.UserAgent = objects.LOCAL_USER
<span id="L169" class="ln">   169</span>	
<span id="L170" class="ln">   170</span>		<span class="comment">// RPC</span>
<span id="L171" class="ln">   171</span>		config.RPCPort = tomlConf.RPCConf.Port
<span id="L172" class="ln">   172</span>		config.RPCUser = tomlConf.RPCConf.User
<span id="L173" class="ln">   173</span>		config.RPCPass = tomlConf.RPCConf.Pass
<span id="L174" class="ln">   174</span>		config.HttpRoot = GetConfDir() + tomlConf.RPCConf.Local
<span id="L175" class="ln">   175</span>	
<span id="L176" class="ln">   176</span>		<span class="comment">// Local Registers</span>
<span id="L177" class="ln">   177</span>		config.PubkeyRegister  = make(chan objects.Hash, bufLen)
<span id="L178" class="ln">   178</span>		config.MessageRegister = make(chan objects.Message, bufLen)
<span id="L179" class="ln">   179</span>		config.PubRegister     = make(chan objects.Message, bufLen)
<span id="L180" class="ln">   180</span>		config.PurgeRegister   = make(chan [16]byte, bufLen)
<span id="L181" class="ln">   181</span>	
<span id="L182" class="ln">   182</span>		<span class="comment">// Administration</span>
<span id="L183" class="ln">   183</span>		config.Log = make(chan string, bufLen)
<span id="L184" class="ln">   184</span>		config.Quit = make(chan os.Signal, 1)
<span id="L185" class="ln">   185</span>	
<span id="L186" class="ln">   186</span>		<span class="comment">// Initialize Map</span>
<span id="L187" class="ln">   187</span>		config.NodeList.Nodes = make(map[string]objects.Node)
<span id="L188" class="ln">   188</span>	
<span id="L189" class="ln">   189</span>		<span class="comment">// Bootstrap Nodes</span>
<span id="L190" class="ln">   190</span>		config.Bootstrap = make([]string, len(tomlConf.Peers), cap(tomlConf.Peers))
<span id="L191" class="ln">   191</span>		copy(config.Bootstrap, tomlConf.Peers)
<span id="L192" class="ln">   192</span>		for i, str := range tomlConf.Peers {
<span id="L193" class="ln">   193</span>			if i &gt;= bufLen {
<span id="L194" class="ln">   194</span>				break
<span id="L195" class="ln">   195</span>			}
<span id="L196" class="ln">   196</span>	
<span id="L197" class="ln">   197</span>			p := new(quibit.Peer)
<span id="L198" class="ln">   198</span>			n := new(objects.Node)
<span id="L199" class="ln">   199</span>			err := n.FromString(str)
<span id="L200" class="ln">   200</span>			if err != nil {
<span id="L201" class="ln">   201</span>				fmt.Println(&#34;Error Decoding Peer &#34;, str, &#34;: &#34;, err)
<span id="L202" class="ln">   202</span>				continue
<span id="L203" class="ln">   203</span>			}
<span id="L204" class="ln">   204</span>	
<span id="L205" class="ln">   205</span>			p.IP = n.IP
<span id="L206" class="ln">   206</span>			p.Port = n.Port
<span id="L207" class="ln">   207</span>			config.PeerQueue &lt;- *p
<span id="L208" class="ln">   208</span>			config.NodeList.Nodes[n.String()] = *n
<span id="L209" class="ln">   209</span>		}
<span id="L210" class="ln">   210</span>	
<span id="L211" class="ln">   211</span>		<span class="comment">// Pull Nodes from node file</span>
<span id="L212" class="ln">   212</span>		if len(config.NodeFile) &gt; 0 {
<span id="L213" class="ln">   213</span>			ReadNodes(config)
<span id="L214" class="ln">   214</span>		}
<span id="L215" class="ln">   215</span>	
<span id="L216" class="ln">   216</span>		return config
<span id="L217" class="ln">   217</span>	}
<span id="L218" class="ln">   218</span>	
<span id="L219" class="ln">   219</span>	<span class="comment">// Load and connect to all nodes from the NodeFile found in the ApiConfig.</span>
<span id="L220" class="ln">   220</span>	func ReadNodes(config *ApiConfig) {
<span id="L221" class="ln">   221</span>		file, err := os.Open(config.NodeFile);
<span id="L222" class="ln">   222</span>		defer file.Close()
<span id="L223" class="ln">   223</span>		if err != nil {
<span id="L224" class="ln">   224</span>			fmt.Println(&#34;Could not open node file: &#34;, err)
<span id="L225" class="ln">   225</span>		}
<span id="L226" class="ln">   226</span>	
<span id="L227" class="ln">   227</span>		var count int
<span id="L228" class="ln">   228</span>	
<span id="L229" class="ln">   229</span>		scanner := bufio.NewScanner(file)
<span id="L230" class="ln">   230</span>	
<span id="L231" class="ln">   231</span>		for scanner.Scan() {
<span id="L232" class="ln">   232</span>	    	str := scanner.Text()
<span id="L233" class="ln">   233</span>	    	if len(str) &lt; 0 || str == &#34;&lt;nil&gt;&#34; {
<span id="L234" class="ln">   234</span>	    		continue
<span id="L235" class="ln">   235</span>	    	}
<span id="L236" class="ln">   236</span>	
<span id="L237" class="ln">   237</span>			p := new(quibit.Peer)
<span id="L238" class="ln">   238</span>			n := new(objects.Node)
<span id="L239" class="ln">   239</span>			err = n.FromString(str)
<span id="L240" class="ln">   240</span>			if err != nil {
<span id="L241" class="ln">   241</span>				fmt.Println(&#34;Error Decoding Peer &#34;, str, &#34;: &#34;, err)
<span id="L242" class="ln">   242</span>				continue
<span id="L243" class="ln">   243</span>			}
<span id="L244" class="ln">   244</span>	
<span id="L245" class="ln">   245</span>			p.IP = n.IP
<span id="L246" class="ln">   246</span>			p.Port = n.Port
<span id="L247" class="ln">   247</span>			config.PeerQueue &lt;- *p
<span id="L248" class="ln">   248</span>			config.NodeList.Nodes[n.String()] = *n
<span id="L249" class="ln">   249</span>			count++
<span id="L250" class="ln">   250</span>		}
<span id="L251" class="ln">   251</span>		fmt.Println(count, &#34;nodes pulled from node file.&#34;)
<span id="L252" class="ln">   252</span>	}
<span id="L253" class="ln">   253</span>	
<span id="L254" class="ln">   254</span>	<span class="comment">// Dump all nodes in config.NodeList to config.NodeFile.</span>
<span id="L255" class="ln">   255</span>	func DumpNodes(config *ApiConfig) {
<span id="L256" class="ln">   256</span>		if config == nil {
<span id="L257" class="ln">   257</span>			return
<span id="L258" class="ln">   258</span>		}
<span id="L259" class="ln">   259</span>		if len(config.NodeFile) &lt; 1 {
<span id="L260" class="ln">   260</span>			return
<span id="L261" class="ln">   261</span>		}
<span id="L262" class="ln">   262</span>		writeBytes := make([]byte, 0, 0)
<span id="L263" class="ln">   263</span>	
<span id="L264" class="ln">   264</span>		for key, _ := range config.NodeList.Nodes {
<span id="L265" class="ln">   265</span>			if quibit.GetPeer(key).IsConnected() {
<span id="L266" class="ln">   266</span>				writeBytes = append(writeBytes, key...)
<span id="L267" class="ln">   267</span>				writeBytes = append(writeBytes, byte(&#39;\n&#39;))
<span id="L268" class="ln">   268</span>			}
<span id="L269" class="ln">   269</span>		}
<span id="L270" class="ln">   270</span>	
<span id="L271" class="ln">   271</span>		err := ioutil.WriteFile(config.NodeFile, writeBytes, 0644)
<span id="L272" class="ln">   272</span>		if err != nil {
<span id="L273" class="ln">   273</span>			fmt.Println(&#34;Error writing peers to file: &#34;, err)
<span id="L274" class="ln">   274</span>		}
<span id="L275" class="ln">   275</span>	}
</pre><p><a href="/src/emp/api/config.go?m=text">View as plain text</a></p>

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

