<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/api/api_test.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/api/api_test.go</h1>




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
<span id="L17" class="ln">    17</span>		&#34;os&#34;
<span id="L18" class="ln">    18</span>		&#34;os/exec&#34;
<span id="L19" class="ln">    19</span>		&#34;quibit&#34;
<span id="L20" class="ln">    20</span>		&#34;testing&#34;
<span id="L21" class="ln">    21</span>		&#34;time&#34;
<span id="L22" class="ln">    22</span>	)
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	func initialize() *ApiConfig {
<span id="L25" class="ln">    25</span>		config := new(ApiConfig)
<span id="L26" class="ln">    26</span>	
<span id="L27" class="ln">    27</span>		<span class="comment">// Network Channels</span>
<span id="L28" class="ln">    28</span>		config.RecvQueue = make(chan quibit.Frame)
<span id="L29" class="ln">    29</span>		config.SendQueue = make(chan quibit.Frame)
<span id="L30" class="ln">    30</span>		config.PeerQueue = make(chan quibit.Peer)
<span id="L31" class="ln">    31</span>	
<span id="L32" class="ln">    32</span>		<span class="comment">// Local Logic</span>
<span id="L33" class="ln">    33</span>		config.DbFile = &#34;testdb.db&#34;
<span id="L34" class="ln">    34</span>	
<span id="L35" class="ln">    35</span>		config.LocalVersion.Version = 1
<span id="L36" class="ln">    36</span>		config.LocalVersion.Timestamp = time.Now().Round(time.Second)
<span id="L37" class="ln">    37</span>		config.LocalVersion.Port = 4444
<span id="L38" class="ln">    38</span>		config.LocalVersion.UserAgent = &#34;strongmsg v0.1&#34;
<span id="L39" class="ln">    39</span>	
<span id="L40" class="ln">    40</span>		<span class="comment">// Administration</span>
<span id="L41" class="ln">    41</span>		config.Log = make(chan string, 100)
<span id="L42" class="ln">    42</span>		config.Quit = make(chan os.Signal, 1)
<span id="L43" class="ln">    43</span>	
<span id="L44" class="ln">    44</span>		go Start(config)
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>		return config
<span id="L47" class="ln">    47</span>	}
<span id="L48" class="ln">    48</span>	
<span id="L49" class="ln">    49</span>	func cleanup(config *ApiConfig) {
<span id="L50" class="ln">    50</span>		var s os.Signal
<span id="L51" class="ln">    51</span>		config.Quit &lt;- s
<span id="L52" class="ln">    52</span>	
<span id="L53" class="ln">    53</span>		str := &lt;-config.Log
<span id="L54" class="ln">    54</span>		for str != &#34;Quit&#34; {
<span id="L55" class="ln">    55</span>			fmt.Println(str)
<span id="L56" class="ln">    56</span>			str = &lt;-config.Log
<span id="L57" class="ln">    57</span>		}
<span id="L58" class="ln">    58</span>	
<span id="L59" class="ln">    59</span>		exec.Command(&#34;rm&#34;, &#34;testdb.db&#34;).Run()
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>	}
<span id="L62" class="ln">    62</span>	
<span id="L63" class="ln">    63</span>	func TestHandshake(t *testing.T) {
<span id="L64" class="ln">    64</span>		config := initialize()
<span id="L65" class="ln">    65</span>	
<span id="L66" class="ln">    66</span>		var frame quibit.Frame
<span id="L67" class="ln">    67</span>		var err error
<span id="L68" class="ln">    68</span>	
<span id="L69" class="ln">    69</span>		<span class="comment">// Test Version</span>
<span id="L70" class="ln">    70</span>		frame = *objects.MakeFrame(objects.VERSION, objects.REQUEST, &amp;config.LocalVersion)
<span id="L71" class="ln">    71</span>		frame.Peer = &#34;127.0.0.1:4444&#34;
<span id="L72" class="ln">    72</span>	
<span id="L73" class="ln">    73</span>		config.RecvQueue &lt;- frame
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		frame = &lt;-config.SendQueue
<span id="L76" class="ln">    76</span>	
<span id="L77" class="ln">    77</span>		if frame.Header.Command != objects.VERSION || frame.Header.Type != objects.REPLY {
<span id="L78" class="ln">    78</span>			fmt.Println(&#34;Frame is not a proper reply to a version request: &#34;, frame.Header)
<span id="L79" class="ln">    79</span>			t.FailNow()
<span id="L80" class="ln">    80</span>		}
<span id="L81" class="ln">    81</span>	
<span id="L82" class="ln">    82</span>		version := new(objects.Version)
<span id="L83" class="ln">    83</span>		err = version.FromBytes(frame.Payload)
<span id="L84" class="ln">    84</span>		if err != nil {
<span id="L85" class="ln">    85</span>			fmt.Println(&#34;Error parsing version reply: &#34;, err)
<span id="L86" class="ln">    86</span>			t.FailNow()
<span id="L87" class="ln">    87</span>		}
<span id="L88" class="ln">    88</span>	
<span id="L89" class="ln">    89</span>		<span class="comment">// Test Peer</span>
<span id="L90" class="ln">    90</span>		frame = *objects.MakeFrame(objects.PEER, objects.REQUEST, &amp;config.NodeList)
<span id="L91" class="ln">    91</span>		frame.Peer = &#34;127.0.0.1:4444&#34;
<span id="L92" class="ln">    92</span>	
<span id="L93" class="ln">    93</span>		config.RecvQueue &lt;- frame
<span id="L94" class="ln">    94</span>	
<span id="L95" class="ln">    95</span>		frame = &lt;-config.SendQueue
<span id="L96" class="ln">    96</span>	
<span id="L97" class="ln">    97</span>		if frame.Header.Command != objects.PEER || frame.Header.Type != objects.REPLY || frame.Header.Length != 0 {
<span id="L98" class="ln">    98</span>			fmt.Println(&#34;Frame is not a proper reply to a peer request: &#34;, frame.Header)
<span id="L99" class="ln">    99</span>			t.FailNow()
<span id="L100" class="ln">   100</span>		}
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		<span class="comment">// Test Obj</span>
<span id="L103" class="ln">   103</span>		frame = *objects.MakeFrame(objects.OBJ, objects.REQUEST, &amp;config.NodeList)
<span id="L104" class="ln">   104</span>		frame.Peer = &#34;127.0.0.1:4444&#34;
<span id="L105" class="ln">   105</span>	
<span id="L106" class="ln">   106</span>		config.RecvQueue &lt;- frame
<span id="L107" class="ln">   107</span>	
<span id="L108" class="ln">   108</span>		frame = &lt;-config.SendQueue
<span id="L109" class="ln">   109</span>	
<span id="L110" class="ln">   110</span>		if frame.Header.Command != objects.OBJ || frame.Header.Type != objects.REPLY || frame.Header.Length != 0 {
<span id="L111" class="ln">   111</span>			fmt.Println(&#34;Frame is not a proper reply to a peer request: &#34;, frame.Header)
<span id="L112" class="ln">   112</span>			t.FailNow()
<span id="L113" class="ln">   113</span>		}
<span id="L114" class="ln">   114</span>	
<span id="L115" class="ln">   115</span>		cleanup(config)
<span id="L116" class="ln">   116</span>	}
</pre><p><a href="/src/emp/api/api_test.go?m=text">View as plain text</a></p>

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

