<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

  <title>src/emp/local/localapi/rpcfunc.go - The Go Programming Language</title>

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
  <h1>Source file src/emp/local/localapi/rpcfunc.go</h1>




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
<span id="L12" class="ln">    12</span>	package localapi
<span id="L13" class="ln">    13</span>	
<span id="L14" class="ln">    14</span>	import (
<span id="L15" class="ln">    15</span>		&#34;emp/encryption&#34;
<span id="L16" class="ln">    16</span>		&#34;emp/local/localdb&#34;
<span id="L17" class="ln">    17</span>		&#34;emp/objects&#34;
<span id="L18" class="ln">    18</span>		&#34;errors&#34;
<span id="L19" class="ln">    19</span>		&#34;fmt&#34;
<span id="L20" class="ln">    20</span>		&#34;net/http&#34;
<span id="L21" class="ln">    21</span>		&#34;quibit&#34;
<span id="L22" class="ln">    22</span>	)
<span id="L23" class="ln">    23</span>	
<span id="L24" class="ln">    24</span>	var logChan chan string
<span id="L25" class="ln">    25</span>	
<span id="L26" class="ln">    26</span>	func (service *EMPService) ForgetAddress(r *http.Request, args *string, reply *NilParam) error {
<span id="L27" class="ln">    27</span>		if !basicAuth(service.Config, r) {
<span id="L28" class="ln">    28</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L29" class="ln">    29</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L30" class="ln">    30</span>		}
<span id="L31" class="ln">    31</span>	
<span id="L32" class="ln">    32</span>		address := encryption.StringToAddress(*args)
<span id="L33" class="ln">    33</span>		if len(address) != 25 {
<span id="L34" class="ln">    34</span>			return errors.New(fmt.Sprintf(&#34;Invalid Address: %s&#34;, address))
<span id="L35" class="ln">    35</span>		}
<span id="L36" class="ln">    36</span>	
<span id="L37" class="ln">    37</span>		addrHash := objects.MakeHash(address)
<span id="L38" class="ln">    38</span>	
<span id="L39" class="ln">    39</span>		return localdb.DeleteAddress(&amp;addrHash)
<span id="L40" class="ln">    40</span>	}
<span id="L41" class="ln">    41</span>	
<span id="L42" class="ln">    42</span>	func (service *EMPService) ConnectionStatus(r *http.Request, args *NilParam, reply *int) error {
<span id="L43" class="ln">    43</span>		if !basicAuth(service.Config, r) {
<span id="L44" class="ln">    44</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L45" class="ln">    45</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L46" class="ln">    46</span>		}
<span id="L47" class="ln">    47</span>		
<span id="L48" class="ln">    48</span>		*reply = quibit.Status()
<span id="L49" class="ln">    49</span>		return nil
<span id="L50" class="ln">    50</span>	}
<span id="L51" class="ln">    51</span>	
<span id="L52" class="ln">    52</span>	func (service *EMPService) GetLabel(r *http.Request, args *string, reply *string) error {
<span id="L53" class="ln">    53</span>		if !basicAuth(service.Config, r) {
<span id="L54" class="ln">    54</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L55" class="ln">    55</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L56" class="ln">    56</span>		}
<span id="L57" class="ln">    57</span>	
<span id="L58" class="ln">    58</span>		var err error
<span id="L59" class="ln">    59</span>	
<span id="L60" class="ln">    60</span>		address := encryption.StringToAddress(*args)
<span id="L61" class="ln">    61</span>		if len(address) != 25 {
<span id="L62" class="ln">    62</span>			return errors.New(fmt.Sprintf(&#34;Invalid Address: %s&#34;, address))
<span id="L63" class="ln">    63</span>		}
<span id="L64" class="ln">    64</span>	
<span id="L65" class="ln">    65</span>		addrHash := objects.MakeHash(address)
<span id="L66" class="ln">    66</span>	
<span id="L67" class="ln">    67</span>		detail, err := localdb.GetAddressDetail(addrHash)
<span id="L68" class="ln">    68</span>		if err != nil {
<span id="L69" class="ln">    69</span>			return err
<span id="L70" class="ln">    70</span>		}
<span id="L71" class="ln">    71</span>	
<span id="L72" class="ln">    72</span>		if len(detail.Label) &gt; 0 {
<span id="L73" class="ln">    73</span>			*reply = detail.Label
<span id="L74" class="ln">    74</span>		} else {
<span id="L75" class="ln">    75</span>			*reply = *args
<span id="L76" class="ln">    76</span>		}
<span id="L77" class="ln">    77</span>		return nil
<span id="L78" class="ln">    78</span>	}
<span id="L79" class="ln">    79</span>	
<span id="L80" class="ln">    80</span>	func (service *EMPService) CreateAddress(r *http.Request, args *NilParam, reply *objects.AddressDetail) error {
<span id="L81" class="ln">    81</span>		if !basicAuth(service.Config, r) {
<span id="L82" class="ln">    82</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L83" class="ln">    83</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L84" class="ln">    84</span>		}
<span id="L85" class="ln">    85</span>	
<span id="L86" class="ln">    86</span>		<span class="comment">// Create Address</span>
<span id="L87" class="ln">    87</span>	
<span id="L88" class="ln">    88</span>		priv, x, y := encryption.CreateKey(service.Config.Log)
<span id="L89" class="ln">    89</span>		reply.Privkey = priv
<span id="L90" class="ln">    90</span>		if x == nil {
<span id="L91" class="ln">    91</span>			return errors.New(&#34;Key Pair Generation Error&#34;)
<span id="L92" class="ln">    92</span>		}
<span id="L93" class="ln">    93</span>	
<span id="L94" class="ln">    94</span>		reply.Pubkey = encryption.MarshalPubkey(x, y)
<span id="L95" class="ln">    95</span>	
<span id="L96" class="ln">    96</span>		reply.IsRegistered = true
<span id="L97" class="ln">    97</span>	
<span id="L98" class="ln">    98</span>		reply.Address = encryption.GetAddress(service.Config.Log, x, y)
<span id="L99" class="ln">    99</span>	
<span id="L100" class="ln">   100</span>		if reply.Address == nil {
<span id="L101" class="ln">   101</span>			return errors.New(&#34;Could not create address, function returned nil.&#34;)
<span id="L102" class="ln">   102</span>		}
<span id="L103" class="ln">   103</span>	
<span id="L104" class="ln">   104</span>		reply.String = encryption.AddressToString(reply.Address)
<span id="L105" class="ln">   105</span>	
<span id="L106" class="ln">   106</span>		<span class="comment">// Add Address to Database</span>
<span id="L107" class="ln">   107</span>		err := localdb.AddUpdateAddress(reply)
<span id="L108" class="ln">   108</span>		if err != nil {
<span id="L109" class="ln">   109</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Error Adding Address: &#34;, err)
<span id="L110" class="ln">   110</span>			return err
<span id="L111" class="ln">   111</span>		}
<span id="L112" class="ln">   112</span>	
<span id="L113" class="ln">   113</span>		<span class="comment">// Send Pubkey to Network</span>
<span id="L114" class="ln">   114</span>		encPub := new(objects.EncryptedPubkey)
<span id="L115" class="ln">   115</span>	
<span id="L116" class="ln">   116</span>		encPub.AddrHash = objects.MakeHash(reply.Address)
<span id="L117" class="ln">   117</span>	
<span id="L118" class="ln">   118</span>		encPub.IV, encPub.Payload, err = encryption.SymmetricEncrypt(reply.Address, string(reply.Pubkey))
<span id="L119" class="ln">   119</span>		if err != nil {
<span id="L120" class="ln">   120</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Error Encrypting Pubkey: &#34;, err)
<span id="L121" class="ln">   121</span>			return nil
<span id="L122" class="ln">   122</span>		}
<span id="L123" class="ln">   123</span>	
<span id="L124" class="ln">   124</span>		<span class="comment">// Record Pubkey for Network</span>
<span id="L125" class="ln">   125</span>		service.Config.RecvQueue &lt;- *objects.MakeFrame(objects.PUBKEY, objects.BROADCAST, encPub)
<span id="L126" class="ln">   126</span>		return nil
<span id="L127" class="ln">   127</span>	}
<span id="L128" class="ln">   128</span>	
<span id="L129" class="ln">   129</span>	func (service *EMPService) GetAddress(r *http.Request, args *string, reply *objects.AddressDetail) error {
<span id="L130" class="ln">   130</span>	
<span id="L131" class="ln">   131</span>		if !basicAuth(service.Config, r) {
<span id="L132" class="ln">   132</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L133" class="ln">   133</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L134" class="ln">   134</span>		}
<span id="L135" class="ln">   135</span>	
<span id="L136" class="ln">   136</span>		var err error
<span id="L137" class="ln">   137</span>	
<span id="L138" class="ln">   138</span>		address := encryption.StringToAddress(*args)
<span id="L139" class="ln">   139</span>		if len(address) != 25 {
<span id="L140" class="ln">   140</span>			return errors.New(fmt.Sprintf(&#34;Invalid Address: %s&#34;, address))
<span id="L141" class="ln">   141</span>		}
<span id="L142" class="ln">   142</span>	
<span id="L143" class="ln">   143</span>		addrHash := objects.MakeHash(address)
<span id="L144" class="ln">   144</span>	
<span id="L145" class="ln">   145</span>		detail, err := localdb.GetAddressDetail(addrHash)
<span id="L146" class="ln">   146</span>		if err != nil {
<span id="L147" class="ln">   147</span>			return err
<span id="L148" class="ln">   148</span>		}
<span id="L149" class="ln">   149</span>	
<span id="L150" class="ln">   150</span>		<span class="comment">// Check for pubkey</span>
<span id="L151" class="ln">   151</span>		if len(detail.Pubkey) == 0 {
<span id="L152" class="ln">   152</span>			detail.Pubkey = checkPubkey(service.Config, objects.MakeHash(detail.Address))
<span id="L153" class="ln">   153</span>		}
<span id="L154" class="ln">   154</span>	
<span id="L155" class="ln">   155</span>		*reply = *detail
<span id="L156" class="ln">   156</span>	
<span id="L157" class="ln">   157</span>		return nil
<span id="L158" class="ln">   158</span>	}
<span id="L159" class="ln">   159</span>	
<span id="L160" class="ln">   160</span>	func (service *EMPService) AddUpdateAddress(r *http.Request, args *objects.AddressDetail, reply *NilParam) error {
<span id="L161" class="ln">   161</span>		if !basicAuth(service.Config, r) {
<span id="L162" class="ln">   162</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L163" class="ln">   163</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L164" class="ln">   164</span>		}
<span id="L165" class="ln">   165</span>	
<span id="L166" class="ln">   166</span>		err := localdb.AddUpdateAddress(args)
<span id="L167" class="ln">   167</span>		if err != nil {
<span id="L168" class="ln">   168</span>			return err
<span id="L169" class="ln">   169</span>		}
<span id="L170" class="ln">   170</span>	
<span id="L171" class="ln">   171</span>		checkPubkey(service.Config, objects.MakeHash(args.Address))
<span id="L172" class="ln">   172</span>	
<span id="L173" class="ln">   173</span>		return nil
<span id="L174" class="ln">   174</span>	}
<span id="L175" class="ln">   175</span>	
<span id="L176" class="ln">   176</span>	func (service *EMPService) ListAddresses(r *http.Request, args *bool, reply *([][2]string)) error {
<span id="L177" class="ln">   177</span>		if !basicAuth(service.Config, r) {
<span id="L178" class="ln">   178</span>			service.Config.Log &lt;- fmt.Sprintf(&#34;Unauthorized RPC Request from: %s&#34;, r.RemoteAddr)
<span id="L179" class="ln">   179</span>			return errors.New(&#34;Unauthorized&#34;)
<span id="L180" class="ln">   180</span>		}
<span id="L181" class="ln">   181</span>	
<span id="L182" class="ln">   182</span>		strs := localdb.ListAddresses(*args)
<span id="L183" class="ln">   183</span>		*reply = strs
<span id="L184" class="ln">   184</span>		return nil
<span id="L185" class="ln">   185</span>	}
</pre><p><a href="/src/emp/local/localapi/rpcfunc.go?m=text">View as plain text</a></p>

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

