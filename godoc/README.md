<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>README.md - The Go Programming Language</title>

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
  <h1>Text file README.md</h1>




<div id="nav"></div>


<pre><span id="L1" class="ln">     1</span>	EMP
<span id="L2" class="ln">     2</span>	=====
<span id="L3" class="ln">     3</span>	EMP is a fully encrypted, distributed messaging service designed with speed in mind.
<span id="L4" class="ln">     4</span>	Originally based off of BitMessage, EMP makes modifications to the API to include
<span id="L5" class="ln">     5</span>	both Read Receipts that Purge the network of read messages, and an extra identification field
<span id="L6" class="ln">     6</span>	to prevent clients from having to decrypt every single incoming message.
<span id="L7" class="ln">     7</span>	
<span id="L8" class="ln">     8</span>	Submodules
<span id="L9" class="ln">     9</span>	----------
<span id="L10" class="ln">    10</span>	
<span id="L11" class="ln">    11</span>	This repository contains submodules.  You will need to run:
<span id="L12" class="ln">    12</span>	```
<span id="L13" class="ln">    13</span>	git submodule init
<span id="L14" class="ln">    14</span>	git submodule update
<span id="L15" class="ln">    15</span>	```
<span id="L16" class="ln">    16</span>	
<span id="L17" class="ln">    17</span>	Required Tools
<span id="L18" class="ln">    18</span>	---------
<span id="L19" class="ln">    19</span>	In order to compile and run this software, you will need:
<span id="L20" class="ln">    20</span>	
<span id="L21" class="ln">    21</span>	* The [Go Compiler (gc)](http://golang.org/doc/install)
<span id="L22" class="ln">    22</span>	* For Downloading Dependencies: [Git](http://git-scm.com/book/en/Getting-Started-Installing-Git)
<span id="L23" class="ln">    23</span>	* For Downloading Dependencies: [Mercurial](http://mercurial.selenic.com/wiki/Download)
<span id="L24" class="ln">    24</span>	
<span id="L25" class="ln">    25</span>	Building and Launching
<span id="L26" class="ln">    26</span>	---------
<span id="L27" class="ln">    27</span>	
<span id="L28" class="ln">    28</span>	* `make build` will install the daemon to ./bin/emp
<span id="L29" class="ln">    29</span>	* `make start` will set up the config directory at ~/.config/emp/, then build and run the daemon, outputting to the log file at ~/.config/emp/log/log_&lt;date&gt;
<span id="L30" class="ln">    30</span>	* `make stop` will stop any existing emp daemon
<span id="L31" class="ln">    31</span>	* `make clean` will remove all build packages and log files
<span id="L32" class="ln">    32</span>	* `make clobber` will also remove all the dependency sources
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>	**Running as root user is NOT recommended!**
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>	Configuration
<span id="L37" class="ln">    37</span>	---------
<span id="L38" class="ln">    38</span>	All configuration is found in `~/.config/emp/msg.conf`, which is installed automatically with `make start`. An example is found in `./script/msg.conf.example`. The example should be good for most users, but if you plan on running a &#34;backbone&#34; node, make sure to add your external IP to msg.conf in order to have it circulated around the network.
<span id="L39" class="ln">    39</span>	
<span id="L40" class="ln">    40</span>	Debian/Ubuntu Installation
<span id="L41" class="ln">    41</span>	---------
<span id="L42" class="ln">    42</span>	* Add the APT repository with `add-apt-repository &#39;deb http://emp.jar.st/repos/apt/debian unstable main&#39;`
<span id="L43" class="ln">    43</span>	* Download and install the JARST GPG Key with `wget -O key http://emp.jar.st/repos/apt/debian/conf/jarst.gpg.key &amp;&amp; sudo apt-key add key; rm -f key`
<span id="L44" class="ln">    44</span>	* Update the APT Database: `sudo apt-get update`
<span id="L45" class="ln">    45</span>	* You can now install EMP with `sudo apt-get install emp`
<span id="L46" class="ln">    46</span>	
<span id="L47" class="ln">    47</span>	**Note:**
<span id="L48" class="ln">    48</span>	Configuration of the Debian installation will be stored in `/usr/share/emp/` instead of the home directory.
</pre><p><a href="/README.md?m=text">View as plain text</a></p>

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

