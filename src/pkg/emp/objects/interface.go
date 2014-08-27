<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/objects/interface.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/objects/interface.go</h1>




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
<span id="L12" class="ln">    12</span>	<span class="comment">// Package objects provides all serializeable data structures within the EMProtocol.</span>
<span id="L13" class="ln">    13</span>	package objects
<span id="L14" class="ln">    14</span>	
<span id="L15" class="ln">    15</span>	import (
<span id="L16" class="ln">    16</span>		&#34;quibit&#34;
<span id="L17" class="ln">    17</span>	)
<span id="L18" class="ln">    18</span>	
<span id="L19" class="ln">    19</span>	type Serializer interface {
<span id="L20" class="ln">    20</span>		GetBytes() []byte
<span id="L21" class="ln">    21</span>		FromBytes([]byte) error
<span id="L22" class="ln">    22</span>	}
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	type NilPayload bool
<span id="L25" class="ln">    25</span>	
<span id="L26" class="ln">    26</span>	func (n *NilPayload) GetBytes() []byte {
<span id="L27" class="ln">    27</span>		return nil
<span id="L28" class="ln">    28</span>	}
<span id="L29" class="ln">    29</span>	
<span id="L30" class="ln">    30</span>	func (n *NilPayload) FromBytes(b []byte) error {
<span id="L31" class="ln">    31</span>		return nil
<span id="L32" class="ln">    32</span>	}
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>	<span class="comment">// Creates a new Quibit Frame reading for sending from a Serializeable Object.</span>
<span id="L35" class="ln">    35</span>	func MakeFrame(command, t uint8, payload Serializer) *quibit.Frame {
<span id="L36" class="ln">    36</span>		frame := new(quibit.Frame)
<span id="L37" class="ln">    37</span>		frame.Configure(payload.GetBytes(), command, t)
<span id="L38" class="ln">    38</span>		frame.Peer = &#34;&#34;
<span id="L39" class="ln">    39</span>	
<span id="L40" class="ln">    40</span>		return frame
<span id="L41" class="ln">    41</span>	}
</pre><p><a href="/src/pkg/emp/objects/interface.go?m=text">View as plain text</a></p>

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

