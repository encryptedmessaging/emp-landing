<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/db/db_test.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/db/db_test.go</h1>




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
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;fmt&#34;
<span id="L16" class="ln">    16</span>		&#34;os/exec&#34;
<span id="L17" class="ln">    17</span>		&#34;emp/objects&#34;
<span id="L18" class="ln">    18</span>		&#34;testing&#34;
<span id="L19" class="ln">    19</span>		&#34;time&#34;
<span id="L20" class="ln">    20</span>	)
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	func TestDatabase(t *testing.T) {
<span id="L23" class="ln">    23</span>		<span class="comment">// Start Logger</span>
<span id="L24" class="ln">    24</span>		log := make(chan string, 100)
<span id="L25" class="ln">    25</span>		go func() {
<span id="L26" class="ln">    26</span>			for {
<span id="L27" class="ln">    27</span>				log_stmt := &lt;-log
<span id="L28" class="ln">    28</span>				fmt.Println(log_stmt)
<span id="L29" class="ln">    29</span>			}
<span id="L30" class="ln">    30</span>		}()
<span id="L31" class="ln">    31</span>	
<span id="L32" class="ln">    32</span>		err := Initialize(log, &#34;testdb.db&#34;)
<span id="L33" class="ln">    33</span>		if dbConn == nil || hashList == nil {
<span id="L34" class="ln">    34</span>			fmt.Println(&#34;ERROR! ERROR! WTF!!!&#34;)
<span id="L35" class="ln">    35</span>		}
<span id="L36" class="ln">    36</span>	
<span id="L37" class="ln">    37</span>		if err != nil {
<span id="L38" class="ln">    38</span>			fmt.Printf(&#34;ERROR: %s\n&#34;, err)
<span id="L39" class="ln">    39</span>			t.FailNow()
<span id="L40" class="ln">    40</span>		}
<span id="L41" class="ln">    41</span>	
<span id="L42" class="ln">    42</span>		txid := []byte{&#39;a&#39;, &#39;b&#39;, &#39;c&#39;, &#39;d&#39;, &#39;e&#39;, &#39;f&#39;, &#39;g&#39;, &#39;h&#39;, &#39;i&#39;, &#39;j&#39;, &#39;k&#39;, &#39;l&#39;, &#39;m&#39;, &#39;n&#39;, &#39;o&#39;, &#39;p&#39;}
<span id="L43" class="ln">    43</span>		purgeHash := objects.MakeHash(txid)
<span id="L44" class="ln">    44</span>		pubHash := objects.MakeHash([]byte{&#39;e&#39;, &#39;f&#39;, &#39;g&#39;, &#39;h&#39;})
<span id="L45" class="ln">    45</span>	
<span id="L46" class="ln">    46</span>		if Contains(purgeHash) != NOTFOUND {
<span id="L47" class="ln">    47</span>			fmt.Println(&#34;Purge Hash already in list...&#34;)
<span id="L48" class="ln">    48</span>			t.FailNow()
<span id="L49" class="ln">    49</span>		}
<span id="L50" class="ln">    50</span>		if Contains(pubHash) != NOTFOUND {
<span id="L51" class="ln">    51</span>			fmt.Println(&#34;Pubkey Hash already in list...&#34;)
<span id="L52" class="ln">    52</span>			t.FailNow()
<span id="L53" class="ln">    53</span>		}
<span id="L54" class="ln">    54</span>	
<span id="L55" class="ln">    55</span>		pub := new(objects.EncryptedPubkey)
<span id="L56" class="ln">    56</span>		pub.AddrHash = pubHash
<span id="L57" class="ln">    57</span>		pub.IV = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
<span id="L58" class="ln">    58</span>		pub.Payload = []byte{&#39;a&#39;, &#39;b&#39;, &#39;c&#39;, &#39;d&#39;}
<span id="L59" class="ln">    59</span>	
<span id="L60" class="ln">    60</span>		err = AddPubkey(log, *pub)
<span id="L61" class="ln">    61</span>		if err != nil {
<span id="L62" class="ln">    62</span>			fmt.Printf(&#34;ERROR: %s\n&#34;, err)
<span id="L63" class="ln">    63</span>			t.FailNow()
<span id="L64" class="ln">    64</span>		}
<span id="L65" class="ln">    65</span>	
<span id="L66" class="ln">    66</span>		if Contains(pubHash) != PUBKEY {
<span id="L67" class="ln">    67</span>			fmt.Println(&#34;Pubkey not in hash list&#34;)
<span id="L68" class="ln">    68</span>			time.Sleep(time.Millisecond)
<span id="L69" class="ln">    69</span>			t.FailNow()
<span id="L70" class="ln">    70</span>		}
<span id="L71" class="ln">    71</span>	
<span id="L72" class="ln">    72</span>		RemoveHash(log, pubHash)
<span id="L73" class="ln">    73</span>	
<span id="L74" class="ln">    74</span>		if Contains(pubHash) != NOTFOUND {
<span id="L75" class="ln">    75</span>			fmt.Println(&#34;Pubkey stuck in hash list&#34;)
<span id="L76" class="ln">    76</span>			t.FailNow()
<span id="L77" class="ln">    77</span>		}
<span id="L78" class="ln">    78</span>	
<span id="L79" class="ln">    79</span>		purge := new(objects.Purge)
<span id="L80" class="ln">    80</span>		copy(purge.Txid[:], txid)
<span id="L81" class="ln">    81</span>	
<span id="L82" class="ln">    82</span>		err = AddPurge(log, *purge)
<span id="L83" class="ln">    83</span>		if err != nil {
<span id="L84" class="ln">    84</span>			fmt.Printf(&#34;ERROR: %s\n&#34;, err)
<span id="L85" class="ln">    85</span>			t.FailNow()
<span id="L86" class="ln">    86</span>		}
<span id="L87" class="ln">    87</span>	
<span id="L88" class="ln">    88</span>		if Contains(purgeHash) != PURGE {
<span id="L89" class="ln">    89</span>			fmt.Println(&#34;Purge not in hash list&#34;)
<span id="L90" class="ln">    90</span>			t.FailNow()
<span id="L91" class="ln">    91</span>		}
<span id="L92" class="ln">    92</span>	
<span id="L93" class="ln">    93</span>		RemoveHash(log, purgeHash)
<span id="L94" class="ln">    94</span>	
<span id="L95" class="ln">    95</span>		if Contains(purgeHash) != NOTFOUND {
<span id="L96" class="ln">    96</span>			fmt.Println(&#34;Purge stuck in hash list&#34;)
<span id="L97" class="ln">    97</span>			t.FailNow()
<span id="L98" class="ln">    98</span>		}
<span id="L99" class="ln">    99</span>	
<span id="L100" class="ln">   100</span>		Cleanup()
<span id="L101" class="ln">   101</span>	
<span id="L102" class="ln">   102</span>		<span class="comment">// Remove DB</span>
<span id="L103" class="ln">   103</span>		err = exec.Command(&#34;rm&#34;, &#34;testdb.db&#34;).Run()
<span id="L104" class="ln">   104</span>	
<span id="L105" class="ln">   105</span>	}
</pre><p><a href="/src/emp/db/db_test.go?m=text">View as plain text</a></p>

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

