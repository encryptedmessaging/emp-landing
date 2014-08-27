<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/quibit/recvUtil.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/quibit/recvUtil.go</h1>




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
<span id="L12" class="ln">    12</span>	package quibit
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;bytes&#34;
<span id="L16" class="ln">    16</span>		&#34;crypto/sha512&#34;
<span id="L17" class="ln">    17</span>		&#34;encoding/binary&#34;
<span id="L18" class="ln">    18</span>		&#34;errors&#34;
<span id="L19" class="ln">    19</span>		&#34;net&#34;
<span id="L20" class="ln">    20</span>		&#34;time&#34;
<span id="L21" class="ln">    21</span>		&#34;io&#34;
<span id="L22" class="ln">    22</span>	)
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	const (
<span id="L25" class="ln">    25</span>		MAGIC = 6667787
<span id="L26" class="ln">    26</span>	)
<span id="L27" class="ln">    27</span>	
<span id="L28" class="ln">    28</span>	func recvAll(conn net.Conn, log chan string) (Frame, error) {
<span id="L29" class="ln">    29</span>		<span class="comment">// ret val</span>
<span id="L30" class="ln">    30</span>		var h Header
<span id="L31" class="ln">    31</span>		var t time.Time
<span id="L32" class="ln">    32</span>		var frame Frame
<span id="L33" class="ln">    33</span>		<span class="comment">// a buffer for decoing</span>
<span id="L34" class="ln">    34</span>		var headerBuffer bytes.Buffer
<span id="L35" class="ln">    35</span>		for {
<span id="L36" class="ln">    36</span>			headerSize := int(binary.Size(h))
<span id="L37" class="ln">    37</span>			<span class="comment">// Byte slice for moving to buffer</span>
<span id="L38" class="ln">    38</span>			buffer := make([]byte, headerSize)
<span id="L39" class="ln">    39</span>			if conn == nil {
<span id="L40" class="ln">    40</span>				return frame, errors.New(&#34;Nil connection&#34;)
<span id="L41" class="ln">    41</span>			}
<span id="L42" class="ln">    42</span>			conn.SetReadDeadline(t)
<span id="L43" class="ln">    43</span>			n, err := io.ReadFull(conn, buffer)
<span id="L44" class="ln">    44</span>			if err != nil {
<span id="L45" class="ln">    45</span>				log &lt;- err.Error()
<span id="L46" class="ln">    46</span>				return frame, err
<span id="L47" class="ln">    47</span>			}
<span id="L48" class="ln">    48</span>			if n &gt; 0 {
<span id="L49" class="ln">    49</span>				<span class="comment">// Add to header buffer</span>
<span id="L50" class="ln">    50</span>				headerBuffer.Write(buffer)
<span id="L51" class="ln">    51</span>				<span class="comment">// Check to see if we have the whole header</span>
<span id="L52" class="ln">    52</span>				if len(headerBuffer.Bytes()) != headerSize {
<span id="L53" class="ln">    53</span>					return frame, errors.New(&#34;Incorrect header size...&#34;)
<span id="L54" class="ln">    54</span>				}
<span id="L55" class="ln">    55</span>				h.FromBytes(headerBuffer.Bytes())
<span id="L56" class="ln">    56</span>				if h.Magic != MAGIC {
<span id="L57" class="ln">    57</span>					return frame, errors.New(&#34;Incorrect Magic Number!&#34;)
<span id="L58" class="ln">    58</span>				}
<span id="L59" class="ln">    59</span>				frame.Header = h
<span id="L60" class="ln">    60</span>				break
<span id="L61" class="ln">    61</span>			}
<span id="L62" class="ln">    62</span>		}
<span id="L63" class="ln">    63</span>	
<span id="L64" class="ln">    64</span>		payload := make([]byte, h.Length)
<span id="L65" class="ln">    65</span>		var payloadBuffer bytes.Buffer
<span id="L66" class="ln">    66</span>		if h.Length &lt; 1 {
<span id="L67" class="ln">    67</span>			frame.Payload = nil
<span id="L68" class="ln">    68</span>			frame.Header = h
<span id="L69" class="ln">    69</span>			return frame, nil
<span id="L70" class="ln">    70</span>		}
<span id="L71" class="ln">    71</span>		for {
<span id="L72" class="ln">    72</span>			<span class="comment">// store in byte array</span>
<span id="L73" class="ln">    73</span>			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
<span id="L74" class="ln">    74</span>			n, err := io.ReadFull(conn, payload)
<span id="L75" class="ln">    75</span>			if err != nil {
<span id="L76" class="ln">    76</span>				return frame, err
<span id="L77" class="ln">    77</span>			}
<span id="L78" class="ln">    78</span>			if n &gt; 1 {
<span id="L79" class="ln">    79</span>				<span class="comment">// write to buffer</span>
<span id="L80" class="ln">    80</span>				payloadBuffer.Write(payload)
<span id="L81" class="ln">    81</span>				<span class="comment">// Check to see if we have whole payload</span>
<span id="L82" class="ln">    82</span>				if len(payloadBuffer.Bytes()) == int(h.Length) {
<span id="L83" class="ln">    83</span>					<span class="comment">// Verify checksum</span>
<span id="L84" class="ln">    84</span>					frame.Payload = payloadBuffer.Bytes()
<span id="L85" class="ln">    85</span>					frame.Header = h
<span id="L86" class="ln">    86</span>					if h.Checksum != sha512.Sum384(payloadBuffer.Bytes()) {
<span id="L87" class="ln">    87</span>						return frame, errors.New(&#34;Incorrect Checksum&#34;)
<span id="L88" class="ln">    88</span>					}
<span id="L89" class="ln">    89</span>					return frame, nil
<span id="L90" class="ln">    90</span>				}
<span id="L91" class="ln">    91</span>			}
<span id="L92" class="ln">    92</span>		}
<span id="L93" class="ln">    93</span>		<span class="comment">//Should never end here</span>
<span id="L94" class="ln">    94</span>		panic(&#34;RECV PAYLOAD&#34;)
<span id="L95" class="ln">    95</span>		return frame, nil
<span id="L96" class="ln">    96</span>	}
</pre><p><a href="/src/pkg/quibit/recvUtil.go?m=text">View as plain text</a></p>

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

