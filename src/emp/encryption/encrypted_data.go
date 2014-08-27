<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/encryption/encrypted_data.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/encryption/encrypted_data.go</h1>




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
<span id="L12" class="ln">    12</span>	package encryption
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import &#34;errors&#34;
<span id="L15" class="ln">    15</span>	
<span id="L16" class="ln">    16</span>	type EncryptedMessage struct {
<span id="L17" class="ln">    17</span>		IV         [16]byte <span class="comment">// Initialization Vector for AES encryption</span>
<span id="L18" class="ln">    18</span>		PublicKey  [65]byte <span class="comment">// Random Public Key used for decryption</span>
<span id="L19" class="ln">    19</span>		CipherText []byte   <span class="comment">// CipherText, length is multiple of AES blocksize</span>
<span id="L20" class="ln">    20</span>		HMAC       [32]byte <span class="comment">// HMAC-SHA256, used to validate key before decryption</span>
<span id="L21" class="ln">    21</span>	}
<span id="L22" class="ln">    22</span>	
<span id="L23" class="ln">    23</span>	const (
<span id="L24" class="ln">    24</span>		ivLen     = 16
<span id="L25" class="ln">    25</span>		pubkeyLen = 65
<span id="L26" class="ln">    26</span>		hmacLen   = 32
<span id="L27" class="ln">    27</span>		minLen    = ivLen + pubkeyLen + hmacLen
<span id="L28" class="ln">    28</span>	)
<span id="L29" class="ln">    29</span>	
<span id="L30" class="ln">    30</span>	func (ret *EncryptedMessage) FromBytes(b []byte) error {
<span id="L31" class="ln">    31</span>		if len(b) &lt; minLen {
<span id="L32" class="ln">    32</span>			return errors.New(&#34;Bytes too short to create EncryptedMessage object.&#34;)
<span id="L33" class="ln">    33</span>		}
<span id="L34" class="ln">    34</span>		if ret == nil {
<span id="L35" class="ln">    35</span>			return errors.New(&#34;Can&#39;t fill nil object.&#34;)
<span id="L36" class="ln">    36</span>		}
<span id="L37" class="ln">    37</span>	
<span id="L38" class="ln">    38</span>		copy(ret.IV[:], b[:ivLen])
<span id="L39" class="ln">    39</span>		copy(ret.PublicKey[:], b[ivLen:ivLen+pubkeyLen])
<span id="L40" class="ln">    40</span>		ret.CipherText = append(ret.CipherText, b[ivLen+pubkeyLen:len(b)-hmacLen]...)
<span id="L41" class="ln">    41</span>		copy(ret.HMAC[:], b[len(b)-hmacLen:])
<span id="L42" class="ln">    42</span>	
<span id="L43" class="ln">    43</span>		return nil
<span id="L44" class="ln">    44</span>	}
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>	func (e *EncryptedMessage) GetBytes() []byte {
<span id="L47" class="ln">    47</span>		if e == nil {
<span id="L48" class="ln">    48</span>			return nil
<span id="L49" class="ln">    49</span>		}
<span id="L50" class="ln">    50</span>		ret := make([]byte, 0, 0)
<span id="L51" class="ln">    51</span>		ret = append(ret, e.IV[:]...)
<span id="L52" class="ln">    52</span>		ret = append(ret, e.PublicKey[:]...)
<span id="L53" class="ln">    53</span>		ret = append(ret, e.CipherText...)
<span id="L54" class="ln">    54</span>		ret = append(ret, e.HMAC[:]...)
<span id="L55" class="ln">    55</span>	
<span id="L56" class="ln">    56</span>		return ret
<span id="L57" class="ln">    57</span>	}
</pre><p><a href="/src/emp/encryption/encrypted_data.go?m=text">View as plain text</a></p>

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

