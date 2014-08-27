<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/quibit/peer.go - The Go Programming Language</title>

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
  <h1>Source file src/quibit/peer.go</h1>




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
<span id="L12" class="ln">    12</span>	
<span id="L13" class="ln">    13</span>	package quibit
<span id="L14" class="ln">    14</span>	
<span id="L15" class="ln">    15</span>	import (
<span id="L16" class="ln">    16</span>		&#34;fmt&#34;
<span id="L17" class="ln">    17</span>		&#34;net&#34;
<span id="L18" class="ln">    18</span>		&#34;strconv&#34;
<span id="L19" class="ln">    19</span>		&#34;time&#34;
<span id="L20" class="ln">    20</span>	)
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	type Peer struct {
<span id="L23" class="ln">    23</span>		IP   net.IP <span class="comment">// Standard 16-byte IP Address</span>
<span id="L24" class="ln">    24</span>		Port uint16 <span class="comment">// Standard 2-byte TCP Port</span>
<span id="L25" class="ln">    25</span>		conn *net.Conn
<span id="L26" class="ln">    26</span>		external bool
<span id="L27" class="ln">    27</span>	}
<span id="L28" class="ln">    28</span>	
<span id="L29" class="ln">    29</span>	func peerFromConn(conn *net.Conn, external bool) Peer {
<span id="L30" class="ln">    30</span>		var p Peer
<span id="L31" class="ln">    31</span>		if conn == nil {
<span id="L32" class="ln">    32</span>			return p
<span id="L33" class="ln">    33</span>		}
<span id="L34" class="ln">    34</span>		
<span id="L35" class="ln">    35</span>		addr := (*conn).RemoteAddr()
<span id="L36" class="ln">    36</span>		if addr.Network() != &#34;tcp&#34; {
<span id="L37" class="ln">    37</span>			return p
<span id="L38" class="ln">    38</span>		}
<span id="L39" class="ln">    39</span>		ip, portStr, err := net.SplitHostPort(addr.String())
<span id="L40" class="ln">    40</span>		port, _ := strconv.Atoi(portStr)
<span id="L41" class="ln">    41</span>		if err != nil {
<span id="L42" class="ln">    42</span>			return p
<span id="L43" class="ln">    43</span>		}
<span id="L44" class="ln">    44</span>	
<span id="L45" class="ln">    45</span>		<span class="comment">// Create New Peer</span>
<span id="L46" class="ln">    46</span>		p.IP = net.ParseIP(ip)
<span id="L47" class="ln">    47</span>		p.Port = uint16(port)
<span id="L48" class="ln">    48</span>	
<span id="L49" class="ln">    49</span>		p.conn = conn
<span id="L50" class="ln">    50</span>		p.external = external
<span id="L51" class="ln">    51</span>		return p
<span id="L52" class="ln">    52</span>	}
<span id="L53" class="ln">    53</span>	
<span id="L54" class="ln">    54</span>	<span class="comment">// String returns the &lt;IP&gt;:&lt;Host&gt; Identifier of Peer</span>
<span id="L55" class="ln">    55</span>	func (p *Peer) String() string {
<span id="L56" class="ln">    56</span>		if p == nil {
<span id="L57" class="ln">    57</span>			return &#34;&#34;
<span id="L58" class="ln">    58</span>		}
<span id="L59" class="ln">    59</span>		return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
<span id="L60" class="ln">    60</span>	}
<span id="L61" class="ln">    61</span>	
<span id="L62" class="ln">    62</span>	<span class="comment">// Returns the connection status of Peer</span>
<span id="L63" class="ln">    63</span>	func (p *Peer) IsConnected() bool {
<span id="L64" class="ln">    64</span>		if p == nil {
<span id="L65" class="ln">    65</span>			return false
<span id="L66" class="ln">    66</span>		}
<span id="L67" class="ln">    67</span>		return (p.conn != nil)
<span id="L68" class="ln">    68</span>	}
<span id="L69" class="ln">    69</span>	
<span id="L70" class="ln">    70</span>	func (p *Peer) connect() error {
<span id="L71" class="ln">    71</span>		<span class="comment">// Check for sane peer object</span>
<span id="L72" class="ln">    72</span>		if p == nil {
<span id="L73" class="ln">    73</span>			return QuibitError(eNILOBJ)
<span id="L74" class="ln">    74</span>		}
<span id="L75" class="ln">    75</span>		if p.conn != nil {
<span id="L76" class="ln">    76</span>			return nil
<span id="L77" class="ln">    77</span>		}
<span id="L78" class="ln">    78</span>	
<span id="L79" class="ln">    79</span>		var err error
<span id="L80" class="ln">    80</span>		lConn, err := net.DialTimeout(&#34;tcp&#34;, p.String(), 10 * time.Second)
<span id="L81" class="ln">    81</span>		p.conn = &amp;lConn
<span id="L82" class="ln">    82</span>		if err != nil {
<span id="L83" class="ln">    83</span>			p.conn = nil
<span id="L84" class="ln">    84</span>			return err
<span id="L85" class="ln">    85</span>		}
<span id="L86" class="ln">    86</span>	
<span id="L87" class="ln">    87</span>		<span class="comment">// Set Keep-Alives</span>
<span id="L88" class="ln">    88</span>		(*p.conn).(*net.TCPConn).SetKeepAlive(true)
<span id="L89" class="ln">    89</span>		(*p.conn).(*net.TCPConn).SetKeepAlivePeriod(time.Second)
<span id="L90" class="ln">    90</span>	
<span id="L91" class="ln">    91</span>		if p.external {
<span id="L92" class="ln">    92</span>			incomingConnections++
<span id="L93" class="ln">    93</span>		}
<span id="L94" class="ln">    94</span>	
<span id="L95" class="ln">    95</span>	
<span id="L96" class="ln">    96</span>		return nil
<span id="L97" class="ln">    97</span>	}
<span id="L98" class="ln">    98</span>	
<span id="L99" class="ln">    99</span>	<span class="comment">// Forces Disconnect from Peer and closes all incoming connections.</span>
<span id="L100" class="ln">   100</span>	func (p *Peer) Disconnect() {
<span id="L101" class="ln">   101</span>		if p.conn == nil {
<span id="L102" class="ln">   102</span>			return
<span id="L103" class="ln">   103</span>		}
<span id="L104" class="ln">   104</span>		(*p.conn).Close()
<span id="L105" class="ln">   105</span>		p.conn = nil
<span id="L106" class="ln">   106</span>		if p.external {
<span id="L107" class="ln">   107</span>			incomingConnections--
<span id="L108" class="ln">   108</span>		}
<span id="L109" class="ln">   109</span>	}
<span id="L110" class="ln">   110</span>	
<span id="L111" class="ln">   111</span>	func (p *Peer) sendFrame(frame Frame) error {
<span id="L112" class="ln">   112</span>		if p == nil {
<span id="L113" class="ln">   113</span>			return QuibitError(eNILOBJ)
<span id="L114" class="ln">   114</span>		}
<span id="L115" class="ln">   115</span>		if p.conn == nil {
<span id="L116" class="ln">   116</span>			return QuibitError(eNILOBJ)
<span id="L117" class="ln">   117</span>		}
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>		var err error
<span id="L120" class="ln">   120</span>	
<span id="L121" class="ln">   121</span>		headerBytes, err := frame.Header.ToBytes()
<span id="L122" class="ln">   122</span>		if err != nil {
<span id="L123" class="ln">   123</span>			return QuibitError(eHEADER)
<span id="L124" class="ln">   124</span>		}
<span id="L125" class="ln">   125</span>	
<span id="L126" class="ln">   126</span>		_, err = (*p.conn).Write(append(headerBytes, frame.Payload...))
<span id="L127" class="ln">   127</span>		if err != nil {
<span id="L128" class="ln">   128</span>			return err
<span id="L129" class="ln">   129</span>		}
<span id="L130" class="ln">   130</span>	
<span id="L131" class="ln">   131</span>		return nil
<span id="L132" class="ln">   132</span>	}
<span id="L133" class="ln">   133</span>	
<span id="L134" class="ln">   134</span>	func (p *Peer) receive(recvChan chan Frame, log chan string) {
<span id="L135" class="ln">   135</span>		if p.conn == nil {
<span id="L136" class="ln">   136</span>			return
<span id="L137" class="ln">   137</span>		}
<span id="L138" class="ln">   138</span>		for {
<span id="L139" class="ln">   139</span>			<span class="comment">// So now we have a connection.  Let&#39;s start receiving</span>
<span id="L140" class="ln">   140</span>			frame, err := recvAll(*p.conn, log)
<span id="L141" class="ln">   141</span>	
<span id="L142" class="ln">   142</span>			if err != nil {
<span id="L143" class="ln">   143</span>				log &lt;- fmt.Sprintln(&#34;Error receiving from Peer: &#34;, err)
<span id="L144" class="ln">   144</span>				KillPeer(p.String())
<span id="L145" class="ln">   145</span>				break
<span id="L146" class="ln">   146</span>			} else {
<span id="L147" class="ln">   147</span>				frame.Peer = p.String()
<span id="L148" class="ln">   148</span>				recvChan &lt;- frame
<span id="L149" class="ln">   149</span>			}
<span id="L150" class="ln">   150</span>		} <span class="comment">// End for</span>
<span id="L151" class="ln">   151</span>	} <span class="comment">// End receive()</span>
</pre><p><a href="/src/quibit/peer.go?m=text">View as plain text</a></p>

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

