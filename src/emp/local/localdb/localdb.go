<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/local/localdb/localdb.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/local/localdb/localdb.go</h1>




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
<span id="L12" class="ln">    12</span>	<span class="comment">// Package localdb provdes a local SQLite3 Database for the EMPLocal client.</span>
<span id="L13" class="ln">    13</span>	package localdb
<span id="L14" class="ln">    14</span>	
<span id="L15" class="ln">    15</span>	import (
<span id="L16" class="ln">    16</span>		&#34;emp/objects&#34;
<span id="L17" class="ln">    17</span>		&#34;fmt&#34;
<span id="L18" class="ln">    18</span>		&#34;github.com/mxk/go-sqlite/sqlite3&#34;
<span id="L19" class="ln">    19</span>		&#34;sync&#34;
<span id="L20" class="ln">    20</span>	)
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	<span class="comment">// Database Connection</span>
<span id="L23" class="ln">    23</span>	var LocalDB *sqlite3.Conn
<span id="L24" class="ln">    24</span>	var localMutex *sync.Mutex
<span id="L25" class="ln">    25</span>	
<span id="L26" class="ln">    26</span>	<span class="comment">// Initialize database with mutexes from file.</span>
<span id="L27" class="ln">    27</span>	func Initialize(log chan string, dbFile string) error {
<span id="L28" class="ln">    28</span>		var err error
<span id="L29" class="ln">    29</span>		if LocalDB != nil {
<span id="L30" class="ln">    30</span>			return nil
<span id="L31" class="ln">    31</span>		}
<span id="L32" class="ln">    32</span>	
<span id="L33" class="ln">    33</span>		localMutex = new(sync.Mutex)
<span id="L34" class="ln">    34</span>	
<span id="L35" class="ln">    35</span>		<span class="comment">// Create Database Connection</span>
<span id="L36" class="ln">    36</span>		LocalDB, err = sqlite3.Open(dbFile)
<span id="L37" class="ln">    37</span>		if err != nil || LocalDB == nil {
<span id="L38" class="ln">    38</span>			log &lt;- fmt.Sprintf(&#34;Error opening sqlite database at %s... %s&#34;, dbFile, err)
<span id="L39" class="ln">    39</span>			LocalDB = nil
<span id="L40" class="ln">    40</span>			return err
<span id="L41" class="ln">    41</span>		}
<span id="L42" class="ln">    42</span>	
<span id="L43" class="ln">    43</span>		<span class="comment">// Create Database Schema</span>
<span id="L44" class="ln">    44</span>	
<span id="L45" class="ln">    45</span>		err = LocalDB.Exec(&#34;CREATE TABLE IF NOT EXISTS addressbook (hash BLOB NOT NULL UNIQUE, address BLOB NOT NULL UNIQUE, registered INTEGER NOT NULL, pubkey BLOB, privkey BLOB, label TEXT, subscribed INTEGER NOT NULL, PRIMARY KEY (hash) ON CONFLICT REPLACE)&#34;)
<span id="L46" class="ln">    46</span>		if err != nil {
<span id="L47" class="ln">    47</span>			log &lt;- fmt.Sprintf(&#34;Error setting up addressbook schema... %s&#34;, err)
<span id="L48" class="ln">    48</span>			LocalDB = nil
<span id="L49" class="ln">    49</span>			return err
<span id="L50" class="ln">    50</span>		}
<span id="L51" class="ln">    51</span>	
<span id="L52" class="ln">    52</span>		<span class="comment">// Migration, Ignore error</span>
<span id="L53" class="ln">    53</span>		LocalDB.Exec(&#34;ALTER TABLE addressbook ADD COLUMN subscribed INTEGER NOT NULL DEFAULT 0&#34;)
<span id="L54" class="ln">    54</span>		LocalDB.Exec(&#34;ALTER TABLE addressbook ADD COLUMN encprivkey BLOB&#34;)
<span id="L55" class="ln">    55</span>	
<span id="L56" class="ln">    56</span>		err = LocalDB.Exec(&#34;CREATE TABLE IF NOT EXISTS msg (txid_hash BLOB NOT NULL, recipient BLOB, timestamp INTEGER, box INTEGER, encrypted BLOB, decrypted BLOB, purged INTEGER, sender BLOB, PRIMARY KEY (txid_hash) ON CONFLICT REPLACE)&#34;)
<span id="L57" class="ln">    57</span>		if err != nil {
<span id="L58" class="ln">    58</span>			log &lt;- fmt.Sprintf(&#34;Error setting up msg schema... %s&#34;, err)
<span id="L59" class="ln">    59</span>			LocalDB = nil
<span id="L60" class="ln">    60</span>			return err
<span id="L61" class="ln">    61</span>		}
<span id="L62" class="ln">    62</span>	
<span id="L63" class="ln">    63</span>		if hashList == nil {
<span id="L64" class="ln">    64</span>			hashList = make(map[string]int)
<span id="L65" class="ln">    65</span>			return populateHashes()
<span id="L66" class="ln">    66</span>		}
<span id="L67" class="ln">    67</span>	
<span id="L68" class="ln">    68</span>		if LocalDB == nil || hashList == nil {
<span id="L69" class="ln">    69</span>			fmt.Println(&#34;ERROR! ERROR! WTF!!! SHOULD BE INITIALIZED!&#34;)
<span id="L70" class="ln">    70</span>		}
<span id="L71" class="ln">    71</span>	
<span id="L72" class="ln">    72</span>		return nil
<span id="L73" class="ln">    73</span>	}
<span id="L74" class="ln">    74</span>	
<span id="L75" class="ln">    75</span>	func populateHashes() error {
<span id="L76" class="ln">    76</span>		for s, err := LocalDB.Query(&#34;SELECT hash FROM addressbook&#34;); err == nil; err = s.Next() {
<span id="L77" class="ln">    77</span>			var hash []byte
<span id="L78" class="ln">    78</span>			s.Scan(&amp;hash) <span class="comment">// Assigns 1st column to rowid, the rest to row</span>
<span id="L79" class="ln">    79</span>			hashList[string(hash)] = ADDRESS
<span id="L80" class="ln">    80</span>		}
<span id="L81" class="ln">    81</span>	
<span id="L82" class="ln">    82</span>		for s, err := LocalDB.Query(&#34;SELECT txid_hash, box FROM msg&#34;); err == nil; err = s.Next() {
<span id="L83" class="ln">    83</span>			var hash []byte
<span id="L84" class="ln">    84</span>			var box int
<span id="L85" class="ln">    85</span>			s.Scan(&amp;hash, &amp;box) <span class="comment">// Assigns 1st column to rowid, the rest to row</span>
<span id="L86" class="ln">    86</span>			hashList[string(hash)] = box
<span id="L87" class="ln">    87</span>		}
<span id="L88" class="ln">    88</span>	
<span id="L89" class="ln">    89</span>		return nil
<span id="L90" class="ln">    90</span>	}
<span id="L91" class="ln">    91</span>	
<span id="L92" class="ln">    92</span>	<span class="comment">// Close and Cleanup database.</span>
<span id="L93" class="ln">    93</span>	func Cleanup() {
<span id="L94" class="ln">    94</span>		LocalDB.Close()
<span id="L95" class="ln">    95</span>		LocalDB = nil
<span id="L96" class="ln">    96</span>		hashList = nil
<span id="L97" class="ln">    97</span>	}
<span id="L98" class="ln">    98</span>	
<span id="L99" class="ln">    99</span>	<span class="comment">// Hash Types</span>
<span id="L100" class="ln">   100</span>	const (
<span id="L101" class="ln">   101</span>		INBOX    = iota <span class="comment">// Incoming Messages</span>
<span id="L102" class="ln">   102</span>		OUTBOX   = iota <span class="comment">// Outgoing, Unsent Messages</span>
<span id="L103" class="ln">   103</span>		SENDBOX  = iota <span class="comment">// Outgoing, Sent Messages</span>
<span id="L104" class="ln">   104</span>		ADDRESS  = iota <span class="comment">// EMP Addresses</span>
<span id="L105" class="ln">   105</span>		NOTFOUND = iota <span class="comment">// Not Found in DB</span>
<span id="L106" class="ln">   106</span>	)
<span id="L107" class="ln">   107</span>	
<span id="L108" class="ln">   108</span>	<span class="comment">// Hash List</span>
<span id="L109" class="ln">   109</span>	var hashList map[string]int
<span id="L110" class="ln">   110</span>	
<span id="L111" class="ln">   111</span>	<span class="comment">// Add to global Hash List.</span>
<span id="L112" class="ln">   112</span>	func Add(hashObj objects.Hash, hashType int) {
<span id="L113" class="ln">   113</span>		hash := string(hashObj.GetBytes())
<span id="L114" class="ln">   114</span>		if hashList != nil {
<span id="L115" class="ln">   115</span>			hashList[hash] = hashType
<span id="L116" class="ln">   116</span>		}
<span id="L117" class="ln">   117</span>	}
<span id="L118" class="ln">   118</span>	
<span id="L119" class="ln">   119</span>	<span class="comment">// Delete hash from Hash List</span>
<span id="L120" class="ln">   120</span>	func Del(hashObj objects.Hash) {
<span id="L121" class="ln">   121</span>		hash := string(hashObj.GetBytes())
<span id="L122" class="ln">   122</span>		if hashList != nil {
<span id="L123" class="ln">   123</span>			delete(hashList, hash)
<span id="L124" class="ln">   124</span>		}
<span id="L125" class="ln">   125</span>	}
<span id="L126" class="ln">   126</span>	
<span id="L127" class="ln">   127</span>	<span class="comment">// Get type of object in Hash List</span>
<span id="L128" class="ln">   128</span>	func Contains(hashObj objects.Hash) int {
<span id="L129" class="ln">   129</span>		hash := string(hashObj.GetBytes())
<span id="L130" class="ln">   130</span>		if hashList != nil {
<span id="L131" class="ln">   131</span>			hashType, ok := hashList[hash]
<span id="L132" class="ln">   132</span>			if ok {
<span id="L133" class="ln">   133</span>				return hashType
<span id="L134" class="ln">   134</span>			} else {
<span id="L135" class="ln">   135</span>				return NOTFOUND
<span id="L136" class="ln">   136</span>			}
<span id="L137" class="ln">   137</span>		}
<span id="L138" class="ln">   138</span>		return NOTFOUND
<span id="L139" class="ln">   139</span>	}
</pre><p><a href="/src/emp/local/localdb/localdb.go?m=text">View as plain text</a></p>

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

