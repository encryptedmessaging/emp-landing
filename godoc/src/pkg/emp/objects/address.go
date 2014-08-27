<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/objects/address.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/objects/address.go</h1>




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
<span id="L12" class="ln">    12</span>	package objects
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	type AddressDetail struct {
<span id="L15" class="ln">    15</span>		String       string `json:&#34;address&#34;`           <span class="comment">// String representation of address (33-35 characters)</span>
<span id="L16" class="ln">    16</span>		Address      []byte `json:&#34;address_bytes&#34;`     <span class="comment">// Byte representation of address (25 bytes)</span>
<span id="L17" class="ln">    17</span>		IsRegistered bool   `json:&#34;registered&#34;`        <span class="comment">// Whether EMPLocal is saving messages to this address.</span>
<span id="L18" class="ln">    18</span>		IsSubscribed bool   `json:&#34;subscribed&#34;`        <span class="comment">// Whether EMPLocal is saving Publications from this address.</span>
<span id="L19" class="ln">    19</span>		Pubkey       []byte `json:&#34;public_key&#34;`        <span class="comment">// Unencrypted 65-byte public key.</span>
<span id="L20" class="ln">    20</span>		Privkey      []byte `json:&#34;private_key&#34;`       <span class="comment">// Unencrypted 32-byte private key.</span>
<span id="L21" class="ln">    21</span>		EncPrivkey   []byte `json:&#34;encrypted_privkey&#34;` <span class="comment">// Encrypted private key (any length).</span>
<span id="L22" class="ln">    22</span>		Label        string `json:&#34;address_label&#34;`     <span class="comment">// Human-readable label for this address.</span>
<span id="L23" class="ln">    23</span>	}
</pre><p><a href="/src/pkg/emp/objects/address.go?m=text">View as plain text</a></p>

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

