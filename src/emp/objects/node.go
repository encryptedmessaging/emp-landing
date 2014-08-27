<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/objects/node.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/objects/node.go</h1>




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
<span id="L18" class="ln">    18</span>		&#34;net&#34;
<span id="L19" class="ln">    19</span>		&#34;strconv&#34;
<span id="L20" class="ln">    20</span>		&#34;time&#34;
<span id="L21" class="ln">    21</span>	)
<span id="L22" class="ln">    22</span>	
<span id="L23" class="ln">    23</span>	const (
<span id="L24" class="ln">    24</span>		nodeLen = 26
<span id="L25" class="ln">    25</span>	)
<span id="L26" class="ln">    26</span>	
<span id="L27" class="ln">    27</span>	type Node struct {
<span id="L28" class="ln">    28</span>		IP       net.IP    <span class="comment">// Public IPv6 or IPv4 Address</span>
<span id="L29" class="ln">    29</span>		Port     uint16    <span class="comment">// Port on which TCP Server is running</span>
<span id="L30" class="ln">    30</span>		LastSeen time.Time <span class="comment">// Time of last connection to Node.</span>
<span id="L31" class="ln">    31</span>		Attempts uint8     <span class="comment">// Number of reconnection attempt. Currently, node is forgotten after 3 failed attempts.</span>
<span id="L32" class="ln">    32</span>	}
<span id="L33" class="ln">    33</span>	
<span id="L34" class="ln">    34</span>	type NodeList struct {
<span id="L35" class="ln">    35</span>		Nodes map[string]Node
<span id="L36" class="ln">    36</span>	}
<span id="L37" class="ln">    37</span>	
<span id="L38" class="ln">    38</span>	func (n *Node) FromString(hostPort string) error {
<span id="L39" class="ln">    39</span>		if n == nil {
<span id="L40" class="ln">    40</span>			return errors.New(&#34;Can&#39;t fill nil object.&#34;)
<span id="L41" class="ln">    41</span>		}
<span id="L42" class="ln">    42</span>	
<span id="L43" class="ln">    43</span>		ip, port, err := net.SplitHostPort(hostPort)
<span id="L44" class="ln">    44</span>		if err != nil {
<span id="L45" class="ln">    45</span>			return nil
<span id="L46" class="ln">    46</span>		}
<span id="L47" class="ln">    47</span>	
<span id="L48" class="ln">    48</span>		n.IP = net.ParseIP(ip)
<span id="L49" class="ln">    49</span>		prt, err := strconv.Atoi(port)
<span id="L50" class="ln">    50</span>		if err != nil {
<span id="L51" class="ln">    51</span>			return err
<span id="L52" class="ln">    52</span>		}
<span id="L53" class="ln">    53</span>		n.Port = uint16(prt)
<span id="L54" class="ln">    54</span>		n.LastSeen = time.Now().Round(time.Second)
<span id="L55" class="ln">    55</span>		n.Attempts = 0
<span id="L56" class="ln">    56</span>		return nil
<span id="L57" class="ln">    57</span>	}
<span id="L58" class="ln">    58</span>	
<span id="L59" class="ln">    59</span>	func (n *Node) String() string {
<span id="L60" class="ln">    60</span>		if n == nil {
<span id="L61" class="ln">    61</span>			return &#34;&#34;
<span id="L62" class="ln">    62</span>		}
<span id="L63" class="ln">    63</span>		return net.JoinHostPort(n.IP.String(), strconv.Itoa(int(n.Port)))
<span id="L64" class="ln">    64</span>	}
<span id="L65" class="ln">    65</span>	
<span id="L66" class="ln">    66</span>	func (n *NodeList) GetBytes() []byte {
<span id="L67" class="ln">    67</span>		if n == nil {
<span id="L68" class="ln">    68</span>			return nil
<span id="L69" class="ln">    69</span>		}
<span id="L70" class="ln">    70</span>		if n.Nodes == nil {
<span id="L71" class="ln">    71</span>			return nil
<span id="L72" class="ln">    72</span>		}
<span id="L73" class="ln">    73</span>	
<span id="L74" class="ln">    74</span>		ret := make([]byte, 0, nodeLen*len(n.Nodes))
<span id="L75" class="ln">    75</span>	
<span id="L76" class="ln">    76</span>		for _, node := range n.Nodes {
<span id="L77" class="ln">    77</span>			nBytes := make([]byte, nodeLen, nodeLen)
<span id="L78" class="ln">    78</span>			copy(nBytes, []byte(node.IP))
<span id="L79" class="ln">    79</span>			binary.BigEndian.PutUint16(nBytes[16:18], node.Port)
<span id="L80" class="ln">    80</span>			binary.BigEndian.PutUint64(nBytes[18:26], uint64(node.LastSeen.Unix()))
<span id="L81" class="ln">    81</span>			ret = append(ret, nBytes...)
<span id="L82" class="ln">    82</span>		}
<span id="L83" class="ln">    83</span>	
<span id="L84" class="ln">    84</span>		return ret
<span id="L85" class="ln">    85</span>	}
<span id="L86" class="ln">    86</span>	
<span id="L87" class="ln">    87</span>	func (n *NodeList) FromBytes(data []byte) error {
<span id="L88" class="ln">    88</span>		if len(data)%nodeLen != 0 {
<span id="L89" class="ln">    89</span>			return errors.New(&#34;Incorrect length for a Node List.&#34;)
<span id="L90" class="ln">    90</span>		}
<span id="L91" class="ln">    91</span>		if n == nil {
<span id="L92" class="ln">    92</span>			return errors.New(&#34;Can&#39;t configure nil Node List&#34;)
<span id="L93" class="ln">    93</span>		}
<span id="L94" class="ln">    94</span>		if n.Nodes == nil {
<span id="L95" class="ln">    95</span>			n.Nodes = make(map[string]Node)
<span id="L96" class="ln">    96</span>		}
<span id="L97" class="ln">    97</span>	
<span id="L98" class="ln">    98</span>		for i := 0; i &lt; len(data); i += nodeLen {
<span id="L99" class="ln">    99</span>			b := bytes.NewBuffer(data[i : i+nodeLen])
<span id="L100" class="ln">   100</span>			node := new(Node)
<span id="L101" class="ln">   101</span>			node.IP = net.IP(b.Next(16))
<span id="L102" class="ln">   102</span>			node.Port = binary.BigEndian.Uint16(b.Next(2))
<span id="L103" class="ln">   103</span>			node.LastSeen = time.Unix(int64(binary.BigEndian.Uint64(b.Next(8))), 0)
<span id="L104" class="ln">   104</span>			node.Attempts = 0
<span id="L105" class="ln">   105</span>			n.Nodes[node.String()] = *node
<span id="L106" class="ln">   106</span>		}
<span id="L107" class="ln">   107</span>	
<span id="L108" class="ln">   108</span>		return nil
<span id="L109" class="ln">   109</span>	}
</pre><p><a href="/src/emp/objects/node.go?m=text">View as plain text</a></p>

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

