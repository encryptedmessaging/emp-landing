<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/encryption/symmetric.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/encryption/symmetric.go</h1>




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
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;crypto/aes&#34;
<span id="L16" class="ln">    16</span>		&#34;crypto/cipher&#34;
<span id="L17" class="ln">    17</span>		&#34;crypto/rand&#34;
<span id="L18" class="ln">    18</span>		&#34;fmt&#34;
<span id="L19" class="ln">    19</span>	)
<span id="L20" class="ln">    20</span>	
<span id="L21" class="ln">    21</span>	<span class="comment">// Symmetrically encryptes plainText using AES-256 and the given key. Returns the IV and ciphertext. </span>
<span id="L22" class="ln">    22</span>	<span class="comment">//</span>
<span id="L23" class="ln">    23</span>	<span class="comment">// Special Case: If the key is length 25-bytes (an address), it is padded to 32-bytes with 0x00.</span>
<span id="L24" class="ln">    24</span>	func SymmetricEncrypt(key []byte, plainText string) ([aes.BlockSize]byte, []byte, error) {
<span id="L25" class="ln">    25</span>	
<span id="L26" class="ln">    26</span>		<span class="comment">// Make Initialization Vector</span>
<span id="L27" class="ln">    27</span>		var IV [aes.BlockSize]byte
<span id="L28" class="ln">    28</span>		n, err := rand.Reader.Read(IV[:])
<span id="L29" class="ln">    29</span>		if err != nil || n != 16 {
<span id="L30" class="ln">    30</span>			return IV, nil, err
<span id="L31" class="ln">    31</span>		}
<span id="L32" class="ln">    32</span>	
<span id="L33" class="ln">    33</span>		<span class="comment">// Pad Plaintext</span>
<span id="L34" class="ln">    34</span>		plainBytes := []byte(plainText)
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>		pad_len := aes.BlockSize - (len(plainBytes) % aes.BlockSize)
<span id="L37" class="ln">    37</span>	
<span id="L38" class="ln">    38</span>		padding := make([]byte, pad_len, pad_len)
<span id="L39" class="ln">    39</span>		plainBytes = append(plainBytes, padding...)
<span id="L40" class="ln">    40</span>	
<span id="L41" class="ln">    41</span>		<span class="comment">// Generate AES Cipher</span>
<span id="L42" class="ln">    42</span>	
<span id="L43" class="ln">    43</span>		if len(key) == 25 {
<span id="L44" class="ln">    44</span>			key = append(key, make([]byte, 7, 7)...)
<span id="L45" class="ln">    45</span>		}
<span id="L46" class="ln">    46</span>	
<span id="L47" class="ln">    47</span>		block, err := aes.NewCipher(key)
<span id="L48" class="ln">    48</span>		if err != nil {
<span id="L49" class="ln">    49</span>			fmt.Println(&#34;Uh Oh: &#34;, err)
<span id="L50" class="ln">    50</span>			return IV, nil, err
<span id="L51" class="ln">    51</span>		}
<span id="L52" class="ln">    52</span>		mode := cipher.NewCBCEncrypter(block, IV[:])
<span id="L53" class="ln">    53</span>	
<span id="L54" class="ln">    54</span>		<span class="comment">// Do encryption</span>
<span id="L55" class="ln">    55</span>		cipherText := make([]byte, len(plainBytes), len(plainBytes))
<span id="L56" class="ln">    56</span>		mode.CryptBlocks(cipherText, plainBytes)
<span id="L57" class="ln">    57</span>	
<span id="L58" class="ln">    58</span>		return IV, cipherText, nil
<span id="L59" class="ln">    59</span>	}
<span id="L60" class="ln">    60</span>	
<span id="L61" class="ln">    61</span>	
<span id="L62" class="ln">    62</span>	<span class="comment">// Decrypted cipherText encrypted with AES-256 using the given IV and key.</span>
<span id="L63" class="ln">    63</span>	func SymmetricDecrypt(IV [aes.BlockSize]byte, key, cipherText []byte) []byte {
<span id="L64" class="ln">    64</span>		plainText := make([]byte, len(cipherText), len(cipherText))
<span id="L65" class="ln">    65</span>	
<span id="L66" class="ln">    66</span>		if len(key) == 25 {
<span id="L67" class="ln">    67</span>			key = append(key, make([]byte, 7, 7)...)
<span id="L68" class="ln">    68</span>		}
<span id="L69" class="ln">    69</span>	
<span id="L70" class="ln">    70</span>		<span class="comment">// Generate AES Cipher</span>
<span id="L71" class="ln">    71</span>		block, _ := aes.NewCipher(key)
<span id="L72" class="ln">    72</span>		mode := cipher.NewCBCDecrypter(block, IV[:])
<span id="L73" class="ln">    73</span>	
<span id="L74" class="ln">    74</span>		<span class="comment">// Do decryption</span>
<span id="L75" class="ln">    75</span>		plainText = make([]byte, len(cipherText), len(cipherText))
<span id="L76" class="ln">    76</span>		mode.CryptBlocks(plainText, cipherText[:])
<span id="L77" class="ln">    77</span>	
<span id="L78" class="ln">    78</span>		return plainText
<span id="L79" class="ln">    79</span>	}
</pre><p><a href="/src/emp/encryption/symmetric.go?m=text">View as plain text</a></p>

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

