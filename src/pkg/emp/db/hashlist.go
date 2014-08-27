<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/pkg/emp/db/hashlist.go - The Go Programming Language</title>

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
  <h1>Source file src/pkg/emp/db/hashlist.go</h1>




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
<span id="L12" class="ln">    12</span>	package db
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import &#34;emp/objects&#34;
<span id="L15" class="ln">    15</span>	
<span id="L16" class="ln">    16</span>	const (
<span id="L17" class="ln">    17</span>		PUBKEY   = iota <span class="comment">// Encrypted Public Key</span>
<span id="L18" class="ln">    18</span>		PURGE    = iota <span class="comment">// Purged Message</span>
<span id="L19" class="ln">    19</span>		MSG      = iota <span class="comment">// Basic Message</span>
<span id="L20" class="ln">    20</span>		PUBKEYRQ = iota <span class="comment">// Public Key Request</span>
<span id="L21" class="ln">    21</span>		PUB      = iota <span class="comment">// Published Message</span>
<span id="L22" class="ln">    22</span>		NOTFOUND = iota <span class="comment">// Object not in hash list.</span>
<span id="L23" class="ln">    23</span>	)
<span id="L24" class="ln">    24</span>	
<span id="L25" class="ln">    25</span>	<span class="comment">// Hash List</span>
<span id="L26" class="ln">    26</span>	var hashList map[string]int
<span id="L27" class="ln">    27</span>	
<span id="L28" class="ln">    28</span>	<span class="comment">// Add an object to the hash list with a given type.</span>
<span id="L29" class="ln">    29</span>	func Add(hashObj objects.Hash, hashType int) {
<span id="L30" class="ln">    30</span>		hash := string(hashObj.GetBytes())
<span id="L31" class="ln">    31</span>		if hashList != nil {
<span id="L32" class="ln">    32</span>			hashList[hash] = hashType
<span id="L33" class="ln">    33</span>		}
<span id="L34" class="ln">    34</span>	}
<span id="L35" class="ln">    35</span>	
<span id="L36" class="ln">    36</span>	<span class="comment">// Remove an object from the hash list.</span>
<span id="L37" class="ln">    37</span>	func Delete(hashObj objects.Hash) {
<span id="L38" class="ln">    38</span>		hash := string(hashObj.GetBytes())
<span id="L39" class="ln">    39</span>		if hashList != nil {
<span id="L40" class="ln">    40</span>			delete(hashList, hash)
<span id="L41" class="ln">    41</span>		}
<span id="L42" class="ln">    42</span>	}
<span id="L43" class="ln">    43</span>	
<span id="L44" class="ln">    44</span>	<span class="comment">// Return the type the item in the hash list (see constants).</span>
<span id="L45" class="ln">    45</span>	func Contains(hashObj objects.Hash) int {
<span id="L46" class="ln">    46</span>		hash := string(hashObj.GetBytes())
<span id="L47" class="ln">    47</span>		if hashList != nil {
<span id="L48" class="ln">    48</span>			hashType, ok := hashList[hash]
<span id="L49" class="ln">    49</span>			if ok {
<span id="L50" class="ln">    50</span>				return hashType
<span id="L51" class="ln">    51</span>			} else {
<span id="L52" class="ln">    52</span>				return NOTFOUND
<span id="L53" class="ln">    53</span>			}
<span id="L54" class="ln">    54</span>		}
<span id="L55" class="ln">    55</span>		return NOTFOUND
<span id="L56" class="ln">    56</span>	}
<span id="L57" class="ln">    57</span>	
<span id="L58" class="ln">    58</span>	<span class="comment">// List of all hashes in the hash list.</span>
<span id="L59" class="ln">    59</span>	func ObjList() *objects.Obj {
<span id="L60" class="ln">    60</span>		if hashList == nil {
<span id="L61" class="ln">    61</span>			return nil
<span id="L62" class="ln">    62</span>		}
<span id="L63" class="ln">    63</span>	
<span id="L64" class="ln">    64</span>		ret := new(objects.Obj)
<span id="L65" class="ln">    65</span>		ret.HashList = make([]objects.Hash, 0, 0)
<span id="L66" class="ln">    66</span>	
<span id="L67" class="ln">    67</span>		hash := new(objects.Hash)
<span id="L68" class="ln">    68</span>	
<span id="L69" class="ln">    69</span>		for key, _ := range hashList {
<span id="L70" class="ln">    70</span>			hash.FromBytes([]byte(key))
<span id="L71" class="ln">    71</span>			ret.HashList = append(ret.HashList, *hash)
<span id="L72" class="ln">    72</span>		}
<span id="L73" class="ln">    73</span>		return ret
<span id="L74" class="ln">    74</span>	}
</pre><p><a href="/src/pkg/emp/db/hashlist.go?m=text">View as plain text</a></p>

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

