<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/quibit/quibit.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/quibit/quibit.go</h1>




<div id="nav"></div>


<pre><span id="L1" class="ln">     1</span>	<span class="comment">/**
</span><span id="L2" class="ln">     2</span>	<span class="comment">    Copyright 2014 JARST, LLC
</span><span id="L3" class="ln">     3</span>	<span class="comment">    
</span><span id="L4" class="ln">     4</span>	<span class="comment">    This file is part of Quibit.
</span><span id="L5" class="ln">     5</span>	<span class="comment">
</span><span id="L6" class="ln">     6</span>	<span class="comment">    Quibit is distributed in the hope that it will be useful,
</span><span id="L7" class="ln">     7</span>	<span class="comment">    but WITHOUT ANY WARRANTY; without even the implied warranty of
</span><span id="L8" class="ln">     8</span>	<span class="comment">    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
</span><span id="L9" class="ln">     9</span>	<span class="comment">    LICENSE file for details.
</span><span id="L10" class="ln">    10</span>	<span class="comment">**/</span>
<span id="L11" class="ln">    11</span>	
<span id="L12" class="ln">    12</span>	<span class="comment">// Package quibit provides basic Peer-To-Peer asynchronous network</span>
<span id="L13" class="ln">    13</span>	<span class="comment">// functionality and peer management.</span>
<span id="L14" class="ln">    14</span>	package quibit
<span id="L15" class="ln">    15</span>	
<span id="L16" class="ln">    16</span>	import (
<span id="L17" class="ln">    17</span>		&#34;fmt&#34;
<span id="L18" class="ln">    18</span>	)
<span id="L19" class="ln">    19</span>	
<span id="L20" class="ln">    20</span>	var peerList map[string]Peer
<span id="L21" class="ln">    21</span>	var quit chan bool
<span id="L22" class="ln">    22</span>	var incomingConnections int
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	
<span id="L25" class="ln">    25</span>	<span class="comment">//Message Types:</span>
<span id="L26" class="ln">    26</span>	<span class="comment">//</span>
<span id="L27" class="ln">    27</span>	<span class="comment">//Broadcast goes to all connected peers except for the Peer specified in the Frame.</span>
<span id="L28" class="ln">    28</span>	<span class="comment">//</span>
<span id="L29" class="ln">    29</span>	<span class="comment">//Request and Reply go to the Peer specified in the Frame.</span>
<span id="L30" class="ln">    30</span>	const (
<span id="L31" class="ln">    31</span>		BROADCAST = iota
<span id="L32" class="ln">    32</span>		REQUEST   = iota
<span id="L33" class="ln">    33</span>		REPLY     = iota
<span id="L34" class="ln">    34</span>	)
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>	
<span id="L37" class="ln">    37</span>	<span class="comment">//Initialize the Quibit Service</span>
<span id="L38" class="ln">    38</span>	<span class="comment">//</span>
<span id="L39" class="ln">    39</span>	<span class="comment">//Frames from the network will be sent to recvChan, and includes the sending peer</span>
<span id="L40" class="ln">    40</span>	<span class="comment">//</span>
<span id="L41" class="ln">    41</span>	<span class="comment">//Frames for the network should be sent to sendChan, and include the receiving peer</span>
<span id="L42" class="ln">    42</span>	<span class="comment">//</span>
<span id="L43" class="ln">    43</span>	<span class="comment">//New Peers for connecting should be sent to peerChan.</span>
<span id="L44" class="ln">    44</span>	<span class="comment">//</span>
<span id="L45" class="ln">    45</span>	<span class="comment">//A local server will be started on the port specified by &#34;port&#34;</span>
<span id="L46" class="ln">    46</span>	<span class="comment">//</span>
<span id="L47" class="ln">    47</span>	<span class="comment">//If an error is returned, than neither the server or mux has been started.</span>
<span id="L48" class="ln">    48</span>	func Initialize(log chan string, recvChan, sendChan chan Frame, peerChan chan Peer, port uint16) error {
<span id="L49" class="ln">    49</span>		var err error
<span id="L50" class="ln">    50</span>	
<span id="L51" class="ln">    51</span>		incomingConnections = 0
<span id="L52" class="ln">    52</span>	
<span id="L53" class="ln">    53</span>		err = initServer(recvChan, peerChan, fmt.Sprintf(&#34;:%d&#34;, port))
<span id="L54" class="ln">    54</span>		if err != nil {
<span id="L55" class="ln">    55</span>			return err
<span id="L56" class="ln">    56</span>		}
<span id="L57" class="ln">    57</span>	
<span id="L58" class="ln">    58</span>		peerList = make(map[string]Peer)
<span id="L59" class="ln">    59</span>		quit = make(chan bool)
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>		go mux(recvChan, sendChan, peerChan, quit, log)
<span id="L62" class="ln">    62</span>		return nil
<span id="L63" class="ln">    63</span>	}
<span id="L64" class="ln">    64</span>	
<span id="L65" class="ln">    65</span>	
<span id="L66" class="ln">    66</span>	<span class="comment">//Cleanup the Quibit Service</span>
<span id="L67" class="ln">    67</span>	<span class="comment">//</span>
<span id="L68" class="ln">    68</span>	<span class="comment">//End the mux and server routines, and Disconnect from all peers.</span>
<span id="L69" class="ln">    69</span>	func Cleanup() {
<span id="L70" class="ln">    70</span>		quit &lt;- true
<span id="L71" class="ln">    71</span>	}
<span id="L72" class="ln">    72</span>	
<span id="L73" class="ln">    73</span>	
<span id="L74" class="ln">    74</span>	<span class="comment">//KillPeer Force Disconnects from a Peer</span>
<span id="L75" class="ln">    75</span>	<span class="comment">//</span>
<span id="L76" class="ln">    76</span>	<span class="comment">//All incoming data is dropped and the peer is removed from the Peer List</span>
<span id="L77" class="ln">    77</span>	func KillPeer(p string) {
<span id="L78" class="ln">    78</span>		peer, ok := peerList[p]
<span id="L79" class="ln">    79</span>		if ok {
<span id="L80" class="ln">    80</span>			peer.Disconnect()
<span id="L81" class="ln">    81</span>			delete(peerList, p)
<span id="L82" class="ln">    82</span>		}
<span id="L83" class="ln">    83</span>	}
<span id="L84" class="ln">    84</span>	
<span id="L85" class="ln">    85</span>	
<span id="L86" class="ln">    86</span>	<span class="comment">//Get Peer associated with the given &lt;IP&gt;:&lt;Host&gt; string</span>
<span id="L87" class="ln">    87</span>	<span class="comment">//</span>
<span id="L88" class="ln">    88</span>	<span class="comment">//&lt;nil&gt; Signifies a disconnected or unknown peer.</span>
<span id="L89" class="ln">    89</span>	func GetPeer(p string) *Peer {
<span id="L90" class="ln">    90</span>		peer, ok := peerList[p]
<span id="L91" class="ln">    91</span>		if ok {
<span id="L92" class="ln">    92</span>			return &amp;peer
<span id="L93" class="ln">    93</span>		} else {
<span id="L94" class="ln">    94</span>			return nil
<span id="L95" class="ln">    95</span>		}
<span id="L96" class="ln">    96</span>	}
<span id="L97" class="ln">    97</span>	
<span id="L98" class="ln">    98</span>	
<span id="L99" class="ln">    99</span>	<span class="comment">//Status returns the Current Connection Status</span>
<span id="L100" class="ln">   100</span>	<span class="comment">//</span>
<span id="L101" class="ln">   101</span>	<span class="comment">// Returns 0 on disconnected.  </span>
<span id="L102" class="ln">   102</span>	<span class="comment">// Returns 1 on Client Connection (Outgoing Connections Only).  </span>
<span id="L103" class="ln">   103</span>	<span class="comment">// Returns 2 On Full Connection (Incoming and Outgoing Connections).</span>
<span id="L104" class="ln">   104</span>	func Status() int {
<span id="L105" class="ln">   105</span>		if len(peerList) &lt; 1 {
<span id="L106" class="ln">   106</span>			return 0
<span id="L107" class="ln">   107</span>		}
<span id="L108" class="ln">   108</span>		if incomingConnections &lt; 1 {
<span id="L109" class="ln">   109</span>			return 1
<span id="L110" class="ln">   110</span>		}
<span id="L111" class="ln">   111</span>		return 2
<span id="L112" class="ln">   112</span>	}
<span id="L113" class="ln">   113</span>	
<span id="L114" class="ln">   114</span>	func mux(recvChan, sendChan chan Frame, peerChan chan Peer, quit chan bool, log chan string) {
<span id="L115" class="ln">   115</span>		var frame Frame
<span id="L116" class="ln">   116</span>		var peer Peer
<span id="L117" class="ln">   117</span>		var err error
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>		for {
<span id="L120" class="ln">   120</span>			select {
<span id="L121" class="ln">   121</span>			case frame = &lt;-sendChan:
<span id="L122" class="ln">   122</span>				<span class="comment">// Received frame to send to peer(s)</span>
<span id="L123" class="ln">   123</span>				if frame.Header.Type == BROADCAST {
<span id="L124" class="ln">   124</span>					<span class="comment">// Send to all peers</span>
<span id="L125" class="ln">   125</span>					for key, p := range peerList {
<span id="L126" class="ln">   126</span>						if key == frame.Peer {
<span id="L127" class="ln">   127</span>							<span class="comment">// Exclude peer in message</span>
<span id="L128" class="ln">   128</span>							continue
<span id="L129" class="ln">   129</span>						}
<span id="L130" class="ln">   130</span>	
<span id="L131" class="ln">   131</span>						err = p.sendFrame(frame)
<span id="L132" class="ln">   132</span>						if err != nil {
<span id="L133" class="ln">   133</span>							if err.Error() != QuibitError(eHEADER).Error() {
<span id="L134" class="ln">   134</span>								<span class="comment">// Disconnect from Peer</span>
<span id="L135" class="ln">   135</span>								p.Disconnect()
<span id="L136" class="ln">   136</span>								delete(peerList, key)
<span id="L137" class="ln">   137</span>							}
<span id="L138" class="ln">   138</span>							<span class="comment">// Malformed header, break out of for loop</span>
<span id="L139" class="ln">   139</span>							log &lt;- fmt.Sprintln(&#34;Error sending frame: &#34;, err)
<span id="L140" class="ln">   140</span>						}
<span id="L141" class="ln">   141</span>					}
<span id="L142" class="ln">   142</span>				} else {
<span id="L143" class="ln">   143</span>					<span class="comment">// Send to one peer</span>
<span id="L144" class="ln">   144</span>					if frame.Peer == &#34;&#34; {
<span id="L145" class="ln">   145</span>						<span class="comment">// Error, can&#39;t broadcast a non-broadcast message</span>
<span id="L146" class="ln">   146</span>						break
<span id="L147" class="ln">   147</span>					}
<span id="L148" class="ln">   148</span>					p, ok := peerList[frame.Peer]
<span id="L149" class="ln">   149</span>					if ok {
<span id="L150" class="ln">   150</span>						err = p.sendFrame(frame)
<span id="L151" class="ln">   151</span>						if err != nil {
<span id="L152" class="ln">   152</span>							if err.Error() != QuibitError(eHEADER).Error() {
<span id="L153" class="ln">   153</span>								<span class="comment">// Disconnect from Peer</span>
<span id="L154" class="ln">   154</span>								p.Disconnect()
<span id="L155" class="ln">   155</span>								delete(peerList, frame.Peer)
<span id="L156" class="ln">   156</span>							}
<span id="L157" class="ln">   157</span>							<span class="comment">// Malformed header</span>
<span id="L158" class="ln">   158</span>							log &lt;- fmt.Sprintln(&#34;Malformed header in frame!&#34;)
<span id="L159" class="ln">   159</span>						}
<span id="L160" class="ln">   160</span>					} else {
<span id="L161" class="ln">   161</span>						log &lt;- fmt.Sprintln(&#34;Peer not found: &#34;, frame.Peer)
<span id="L162" class="ln">   162</span>					}
<span id="L163" class="ln">   163</span>				}
<span id="L164" class="ln">   164</span>	
<span id="L165" class="ln">   165</span>			case peer = &lt;-peerChan:
<span id="L166" class="ln">   166</span>				<span class="comment">// Received a new peer to connect to...</span>
<span id="L167" class="ln">   167</span>				err = peer.connect()
<span id="L168" class="ln">   168</span>				if err == nil {
<span id="L169" class="ln">   169</span>					peerList[peer.String()] = peer
<span id="L170" class="ln">   170</span>	
<span id="L171" class="ln">   171</span>					<span class="comment">// Prevent overwriting...</span>
<span id="L172" class="ln">   172</span>					rawPeer := new(Peer)
<span id="L173" class="ln">   173</span>					*rawPeer = peer
<span id="L174" class="ln">   174</span>					go rawPeer.receive(recvChan, log)
<span id="L175" class="ln">   175</span>	
<span id="L176" class="ln">   176</span>				} else {
<span id="L177" class="ln">   177</span>					log &lt;- fmt.Sprintln(&#34;Error adding peer: &#34;, err)
<span id="L178" class="ln">   178</span>				}
<span id="L179" class="ln">   179</span>			case &lt;-quit:
<span id="L180" class="ln">   180</span>				for _, p := range peerList {
<span id="L181" class="ln">   181</span>					p.Disconnect()
<span id="L182" class="ln">   182</span>				}
<span id="L183" class="ln">   183</span>				return
<span id="L184" class="ln">   184</span>			} <span class="comment">// End select</span>
<span id="L185" class="ln">   185</span>		} <span class="comment">// End for</span>
<span id="L186" class="ln">   186</span>	} <span class="comment">// End mux()</span>
</pre><p><a href="/src/pkg/quibit/quibit.go?m=text">View as plain text</a></p>

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

