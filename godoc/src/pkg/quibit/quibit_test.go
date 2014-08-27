<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/quibit/quibit_test.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/quibit/quibit_test.go</h1>




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
<span id="L18" class="ln">    18</span>		&#34;testing&#34;
<span id="L19" class="ln">    19</span>		&#34;time&#34;
<span id="L20" class="ln">    20</span>	)
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	<span class="comment">// TestAcceptence is an all-encompassing end-to-end acceptance test. </span>
<span id="L23" class="ln">    23</span>	<span class="comment">// Includes the testing of the Initialize Function, Incoming Connections,</span>
<span id="L24" class="ln">    24</span>	<span class="comment">// and Peer Management functions.</span>
<span id="L25" class="ln">    25</span>	func TestAcceptance(t *testing.T) {
<span id="L26" class="ln">    26</span>		log := make(chan string, 100)
<span id="L27" class="ln">    27</span>		recvChan := make(chan Frame, 10)
<span id="L28" class="ln">    28</span>		sendChan := make(chan Frame, 10)
<span id="L29" class="ln">    29</span>		peerChan := make(chan Peer)
<span id="L30" class="ln">    30</span>		port := uint16(4444)
<span id="L31" class="ln">    31</span>	
<span id="L32" class="ln">    32</span>		<span class="comment">// Initialize Quibit</span>
<span id="L33" class="ln">    33</span>		err := Initialize(log, recvChan, sendChan, peerChan, port)
<span id="L34" class="ln">    34</span>		if err != nil {
<span id="L35" class="ln">    35</span>			fmt.Println(&#34;ERROR INITIALIZING! &#34;, err)
<span id="L36" class="ln">    36</span>			t.FailNow()
<span id="L37" class="ln">    37</span>		}
<span id="L38" class="ln">    38</span>	
<span id="L39" class="ln">    39</span>		<span class="comment">// Test 1: Manual Connection, look for receive</span>
<span id="L40" class="ln">    40</span>		conn, err := net.Dial(&#34;tcp&#34;, &#34;127.0.0.1:4444&#34;)
<span id="L41" class="ln">    41</span>		if err != nil {
<span id="L42" class="ln">    42</span>			fmt.Println(&#34;Error connecting: &#34;, err)
<span id="L43" class="ln">    43</span>			t.FailNow()
<span id="L44" class="ln">    44</span>		}
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>		time.Sleep(time.Millisecond)
<span id="L47" class="ln">    47</span>		if len(peerList) == 0 {
<span id="L48" class="ln">    48</span>			fmt.Println(&#34;Not in peer list!&#34;)
<span id="L49" class="ln">    49</span>			t.FailNow()
<span id="L50" class="ln">    50</span>		}
<span id="L51" class="ln">    51</span>	
<span id="L52" class="ln">    52</span>		data := []byte{&#39;a&#39;, &#39;b&#39;, &#39;c&#39;, &#39;d&#39;}
<span id="L53" class="ln">    53</span>		frame := new(Frame)
<span id="L54" class="ln">    54</span>		frame.Configure(data, 1, 1)
<span id="L55" class="ln">    55</span>	
<span id="L56" class="ln">    56</span>		buf, _ := frame.Header.ToBytes()
<span id="L57" class="ln">    57</span>		_, err = conn.Write(buf)
<span id="L58" class="ln">    58</span>		if err != nil {
<span id="L59" class="ln">    59</span>			fmt.Println(&#34;Error writing header: &#34;, err)
<span id="L60" class="ln">    60</span>			t.FailNow()
<span id="L61" class="ln">    61</span>		}
<span id="L62" class="ln">    62</span>	
<span id="L63" class="ln">    63</span>		_, err = conn.Write(frame.Payload)
<span id="L64" class="ln">    64</span>		if err != nil {
<span id="L65" class="ln">    65</span>			fmt.Println(&#34;Error writing payload: &#34;, err)
<span id="L66" class="ln">    66</span>			t.FailNow()
<span id="L67" class="ln">    67</span>		}
<span id="L68" class="ln">    68</span>	
<span id="L69" class="ln">    69</span>		frame2 := &lt;-recvChan
<span id="L70" class="ln">    70</span>		if string(frame2.Payload) != string(data) {
<span id="L71" class="ln">    71</span>			fmt.Println(&#34;Bad frame! &#34;, frame2)
<span id="L72" class="ln">    72</span>			t.FailNow()
<span id="L73" class="ln">    73</span>		}
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		if frame2.Peer != conn.LocalAddr().String() {
<span id="L76" class="ln">    76</span>			fmt.Println(&#34;Peer doesn&#39;t match! &#34;, frame2.Peer, conn.LocalAddr().String())
<span id="L77" class="ln">    77</span>			t.FailNow()
<span id="L78" class="ln">    78</span>		}
<span id="L79" class="ln">    79</span>	
<span id="L80" class="ln">    80</span>		<span class="comment">// Test 2: Send, look for manual receive</span>
<span id="L81" class="ln">    81</span>		sendChan &lt;- frame2
<span id="L82" class="ln">    82</span>		time.Sleep(time.Millisecond)
<span id="L83" class="ln">    83</span>	
<span id="L84" class="ln">    84</span>		<span class="comment">// So now we have a connection.  Let&#39;s shake hands.</span>
<span id="L85" class="ln">    85</span>		frame3, err := recvAll(conn, log)
<span id="L86" class="ln">    86</span>	
<span id="L87" class="ln">    87</span>		if err != nil {
<span id="L88" class="ln">    88</span>			fmt.Println(&#34;Error Receiving Frame 3... &#34;, err)
<span id="L89" class="ln">    89</span>			t.FailNow()
<span id="L90" class="ln">    90</span>		}
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>		if string(frame3.Payload) != string(data) {
<span id="L93" class="ln">    93</span>			fmt.Println(&#34;Bad frame! &#34;, frame3)
<span id="L94" class="ln">    94</span>			t.FailNow()
<span id="L95" class="ln">    95</span>		}
<span id="L96" class="ln">    96</span>	
<span id="L97" class="ln">    97</span>		conn.Close()
<span id="L98" class="ln">    98</span>		Cleanup()
<span id="L99" class="ln">    99</span>	}
</pre><p><a href="/src/pkg/quibit/quibit_test.go?m=text">View as plain text</a></p>

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

