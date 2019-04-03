<!-- TOC -->

- [说明](#说明)

<!-- /TOC -->


<a id="markdown-说明" name="说明"></a>
# 说明

有些服务只提供了tcp + jsonrpc接口,用末尾的"\n"来表示tcp流->包的逻辑. 的确有很多服务是用"\n"的方式来表示包的末尾,比如redis,memcached. 但是这样的一个明显的问题就是只能使用netcat进行发包,收包. 好在redis,memcached协议比较简单.但如果涉及到复杂的协议,有大量的不同类型的请求和应答,并且使用了jsonrpc作为序列化和反序列化的基础设施. 那么这样做测试起来就会很麻烦. 类似postman这样的只支持http收发包的工具也无法使用.

本工具是给tcp + jsonrpc接口提供了一个proxy

1. 将http + jsonrpc 转换成tcp + jsonrpc
2. 提供配置的自动登录(tcp + jsonrpc), 每次连接打开时将配置好的数据自动发送 (有些不必要的前置握手信息)
