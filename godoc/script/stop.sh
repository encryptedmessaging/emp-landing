<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>script/stop.sh - The Go Programming Language</title>

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
  <h1>Text file script/stop.sh</h1>




<div id="nav"></div>


<pre><span id="L1" class="ln">     1</span>	#!/bin/bash
<span id="L2" class="ln">     2</span>	: &#39;
<span id="L3" class="ln">     3</span>	    Copyright 2014 JARST, LLC
<span id="L4" class="ln">     4</span>	    
<span id="L5" class="ln">     5</span>	    This file is part of EMP.
<span id="L6" class="ln">     6</span>	
<span id="L7" class="ln">     7</span>	    EMP is free software: you can redistribute it and/or modify
<span id="L8" class="ln">     8</span>	    it under the terms of the GNU General Public License as published by
<span id="L9" class="ln">     9</span>	    the Free Software Foundation, either version 3 of the License, or
<span id="L10" class="ln">    10</span>	    (at your option) any later version.
<span id="L11" class="ln">    11</span>	
<span id="L12" class="ln">    12</span>	    EMP is distributed in the hope that it will be useful,
<span id="L13" class="ln">    13</span>	    but WITHOUT ANY WARRANTY; without even the implied warranty of
<span id="L14" class="ln">    14</span>	    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
<span id="L15" class="ln">    15</span>	    GNU General Public License for more details.
<span id="L16" class="ln">    16</span>	
<span id="L17" class="ln">    17</span>	    You should have received a copy of the GNU General Public License
<span id="L18" class="ln">    18</span>	    along with Foobar.  If not, see &lt;http://www.gnu.org/licenses/&gt;.
<span id="L19" class="ln">    19</span>	&#39;
<span id="L20" class="ln">    20</span>	
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	# Kill existing process
<span id="L23" class="ln">    23</span>	if [ -f ~/.config/emp/pid ];
<span id="L24" class="ln">    24</span>	then
<span id="L25" class="ln">    25</span>	  echo &#34;Killing emp server (pid=&#34;`cat ~/.config/emp/pid`&#34;)&#34;
<span id="L26" class="ln">    26</span>	  kill -2 `cat ~/.config/emp/pid`
<span id="L27" class="ln">    27</span>	  rm ~/.config/emp/pid
<span id="L28" class="ln">    28</span>	else
<span id="L29" class="ln">    29</span>	  echo &#34;Server not running, execute start.sh&#34;
<span id="L30" class="ln">    30</span>	fi
<span id="L31" class="ln">    31</span>	
</pre><p><a href="/script/stop.sh?m=text">View as plain text</a></p>

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

