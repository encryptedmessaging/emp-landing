<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/objects/message.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/objects/message.go</h1>




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
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;bytes&#34;
<span id="L16" class="ln">    16</span>		&#34;encoding/binary&#34;
<span id="L17" class="ln">    17</span>		&#34;errors&#34;
<span id="L18" class="ln">    18</span>		&#34;emp/encryption&#34;
<span id="L19" class="ln">    19</span>		&#34;time&#34;
<span id="L20" class="ln">    20</span>	)
<span id="L21" class="ln">    21</span>	
<span id="L22" class="ln">    22</span>	type Message struct {
<span id="L23" class="ln">    23</span>		AddrHash  Hash                        <span class="comment">// Hash of Recipient&#39;s Address</span>
<span id="L24" class="ln">    24</span>		TxidHash  Hash                        <span class="comment">// Hash of random identifier</span>
<span id="L25" class="ln">    25</span>		Timestamp time.Time                   <span class="comment">// Time that message was first broadcast</span>
<span id="L26" class="ln">    26</span>		Content   encryption.EncryptedMessage <span class="comment">// see package encryption</span>
<span id="L27" class="ln">    27</span>	}
<span id="L28" class="ln">    28</span>	
<span id="L29" class="ln">    29</span>	const (
<span id="L30" class="ln">    30</span>		msgLen = 2*hashLen + 8
<span id="L31" class="ln">    31</span>	)
<span id="L32" class="ln">    32</span>	
<span id="L33" class="ln">    33</span>	<span class="comment">// Message Command Types, See EMPv1 Specification</span>
<span id="L34" class="ln">    34</span>	const (
<span id="L35" class="ln">    35</span>		VERSION = iota
<span id="L36" class="ln">    36</span>		PEER    = iota
<span id="L37" class="ln">    37</span>		OBJ     = iota
<span id="L38" class="ln">    38</span>		GETOBJ  = iota
<span id="L39" class="ln">    39</span>	
<span id="L40" class="ln">    40</span>		PUBKEY_REQUEST = iota
<span id="L41" class="ln">    41</span>		PUBKEY         = iota
<span id="L42" class="ln">    42</span>		MSG            = iota
<span id="L43" class="ln">    43</span>		PURGE          = iota
<span id="L44" class="ln">    44</span>	
<span id="L45" class="ln">    45</span>		CHECKTXID      = iota
<span id="L46" class="ln">    46</span>		PUB            = iota
<span id="L47" class="ln">    47</span>	)
<span id="L48" class="ln">    48</span>	
<span id="L49" class="ln">    49</span>	<span class="comment">// Quibit Types, see package quibit.</span>
<span id="L50" class="ln">    50</span>	const (
<span id="L51" class="ln">    51</span>		BROADCAST = iota
<span id="L52" class="ln">    52</span>		REQUEST   = iota
<span id="L53" class="ln">    53</span>		REPLY     = iota
<span id="L54" class="ln">    54</span>	)
<span id="L55" class="ln">    55</span>	
<span id="L56" class="ln">    56</span>	func (m *Message) FromBytes(data []byte) error {
<span id="L57" class="ln">    57</span>		if len(data) &lt; msgLen {
<span id="L58" class="ln">    58</span>			return errors.New(&#34;Data too short to create message!&#34;)
<span id="L59" class="ln">    59</span>		}
<span id="L60" class="ln">    60</span>		if m == nil {
<span id="L61" class="ln">    61</span>			return errors.New(&#34;Can&#39;t fill nil Message object!&#34;)
<span id="L62" class="ln">    62</span>		}
<span id="L63" class="ln">    63</span>		buffer := bytes.NewBuffer(data)
<span id="L64" class="ln">    64</span>		m.AddrHash.FromBytes(buffer.Next(hashLen))
<span id="L65" class="ln">    65</span>		m.TxidHash.FromBytes(buffer.Next(hashLen))
<span id="L66" class="ln">    66</span>		m.Timestamp = time.Unix(int64(binary.BigEndian.Uint64(buffer.Next(8))), 0)
<span id="L67" class="ln">    67</span>		m.Content.FromBytes(buffer.Bytes())
<span id="L68" class="ln">    68</span>	
<span id="L69" class="ln">    69</span>		return nil
<span id="L70" class="ln">    70</span>	
<span id="L71" class="ln">    71</span>	}
<span id="L72" class="ln">    72</span>	
<span id="L73" class="ln">    73</span>	func (m *Message) GetBytes() []byte {
<span id="L74" class="ln">    74</span>		if m == nil {
<span id="L75" class="ln">    75</span>			return nil
<span id="L76" class="ln">    76</span>		}
<span id="L77" class="ln">    77</span>	
<span id="L78" class="ln">    78</span>		ret := make([]byte, 0, msgLen)
<span id="L79" class="ln">    79</span>		ret = append(m.AddrHash.GetBytes(), m.TxidHash.GetBytes()...)
<span id="L80" class="ln">    80</span>		time := make([]byte, 8, 8)
<span id="L81" class="ln">    81</span>		binary.BigEndian.PutUint64(time, uint64(m.Timestamp.Unix()))
<span id="L82" class="ln">    82</span>		ret = append(ret, time...)
<span id="L83" class="ln">    83</span>		ret = append(ret, m.Content.GetBytes()...)
<span id="L84" class="ln">    84</span>		return ret
<span id="L85" class="ln">    85</span>	}
</pre><p><a href="/src/emp/objects/message.go?m=text">View as plain text</a></p>

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

