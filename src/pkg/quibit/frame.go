<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/quibit/frame.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/quibit/frame.go</h1>




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
<span id="L15" class="ln">    15</span>	type Frame struct {
<span id="L16" class="ln">    16</span>		Peer    string <span class="comment">// Peer who sent or will receive this frame</span>
<span id="L17" class="ln">    17</span>		Header  Header <span class="comment">// Header associated with frame</span>
<span id="L18" class="ln">    18</span>		Payload []byte <span class="comment">// Serialize frame payload</span>
<span id="L19" class="ln">    19</span>	}
<span id="L20" class="ln">    20</span>	
<span id="L21" class="ln">    21</span>	<span class="comment">// Configure Frame f with a proper header for</span>
<span id="L22" class="ln">    22</span>	<span class="comment">// payload &#34;data&#34; interpreted as &#34;command&#34;</span>
<span id="L23" class="ln">    23</span>	func (f *Frame) Configure(data []byte, command, t uint8) {
<span id="L24" class="ln">    24</span>		var h Header
<span id="L25" class="ln">    25</span>		h.Configure(data)
<span id="L26" class="ln">    26</span>		h.Command = command
<span id="L27" class="ln">    27</span>		h.Type = t
<span id="L28" class="ln">    28</span>		f.Header = h
<span id="L29" class="ln">    29</span>		f.Payload = data
<span id="L30" class="ln">    30</span>	}
</pre><p><a href="/src/pkg/quibit/frame.go?m=text">View as plain text</a></p>

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

