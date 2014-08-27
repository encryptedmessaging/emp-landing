<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>script/start.sh - The Go Programming Language</title>

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
  <h1>Text file script/start.sh</h1>




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
<span id="L22" class="ln">    22</span>	TMPGOPATH=$GOPATH
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	# Check for go
<span id="L25" class="ln">    25</span>	echo &#34;Checking for go command...&#34;
<span id="L26" class="ln">    26</span>	if ! which go &gt; /dev/null; then
<span id="L27" class="ln">    27</span>	  echo &#34;Go command not found, please install it.&#34;
<span id="L28" class="ln">    28</span>	  exit -1
<span id="L29" class="ln">    29</span>	fi
<span id="L30" class="ln">    30</span>	
<span id="L31" class="ln">    31</span>	# Setup environment variables
<span id="L32" class="ln">    32</span>	DIR=&#34;$( cd &#34;$( dirname &#34;${BASH_SOURCE[0]}&#34; )&#34; &amp;&amp; pwd )&#34;
<span id="L33" class="ln">    33</span>	export GOPATH=$DIR/..
<span id="L34" class="ln">    34</span>	
<span id="L35" class="ln">    35</span>	# Make and fill config directory
<span id="L36" class="ln">    36</span>	echo &#34;Checking config directory...&#34;
<span id="L37" class="ln">    37</span>	mkdir -p ~/.config/emp
<span id="L38" class="ln">    38</span>	mkdir -p ~/.config/emp/log
<span id="L39" class="ln">    39</span>	touch ~/.config/emp/known_nodes.dat
<span id="L40" class="ln">    40</span>	if [ ! -f ~/.config/emp/msg.conf ]; then
<span id="L41" class="ln">    41</span>	  cp &#34;$DIR/msg.conf.example&#34; ~/.config/emp/msg.conf
<span id="L42" class="ln">    42</span>	fi
<span id="L43" class="ln">    43</span>	rm -rf ~/.config/emp/client
<span id="L44" class="ln">    44</span>	cp -r &#34;$DIR/../client&#34; ~/.config/emp/
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>	# Kill existing process
<span id="L47" class="ln">    47</span>	if [ -f ~/.config/emp/pid ];
<span id="L48" class="ln">    48</span>	then
<span id="L49" class="ln">    49</span>	  echo &#39;Killing existing process...&#39;
<span id="L50" class="ln">    50</span>	  kill -15 `cat ~/.config/emp/pid`
<span id="L51" class="ln">    51</span>	  rm -f ~/.config/emp/pid
<span id="L52" class="ln">    52</span>	fi
<span id="L53" class="ln">    53</span>	
<span id="L54" class="ln">    54</span>	# Get Dependencies
<span id="L55" class="ln">    55</span>	echo &#34;Installing dependencies...&#34;
<span id="L56" class="ln">    56</span>	go get code.google.com/p/go.crypto/ripemd160
<span id="L57" class="ln">    57</span>	go get github.com/BurntSushi/toml
<span id="L58" class="ln">    58</span>	go get github.com/gorilla/rpc
<span id="L59" class="ln">    59</span>	go get github.com/mxk/go-sqlite/sqlite3
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>	# Install and go!
<span id="L62" class="ln">    62</span>	echo &#34;Building and running...&#34;
<span id="L63" class="ln">    63</span>	if `go install emp`; then
<span id="L64" class="ln">    64</span>		&#34;$GOPATH/bin/emp&#34; &#34;$HOME/.config/emp/&#34; &gt; ~/.config/emp/log/log_`date +%s` &amp;
<span id="L65" class="ln">    65</span>		echo $! &gt; ~/.config/emp/pid
<span id="L66" class="ln">    66</span>	
<span id="L67" class="ln">    67</span>		# Get Ports
<span id="L68" class="ln">    68</span>		PORTS=$(sed -n &#39;s/.*port *= *\([^ ]*.*\)/\1/p&#39; &lt; ~/.config/emp/msg.conf)
<span id="L69" class="ln">    69</span>		IFS=&#39; &#39; read -a array &lt;&lt;&lt; $PORTS
<span id="L70" class="ln">    70</span>	
<span id="L71" class="ln">    71</span>		# Final Output
<span id="L72" class="ln">    72</span>		echo &#34;Started EMP client on local port ${array[0]}.&#34;
<span id="L73" class="ln">    73</span>		echo &#34;Access local client at: http://localhost:${array[1]}&#34;
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>		export GOPATH=$TMPGOPATH
<span id="L76" class="ln">    76</span>		exit 0
<span id="L77" class="ln">    77</span>	
<span id="L78" class="ln">    78</span>	else echo &#34;Build Failed, could not start client.&#34;
<span id="L79" class="ln">    79</span>	fi
<span id="L80" class="ln">    80</span>	
<span id="L81" class="ln">    81</span>	export GOPATH=$TMPGOPATH
<span id="L82" class="ln">    82</span>	exit -1
</pre><p><a href="/script/start.sh?m=text">View as plain text</a></p>

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

