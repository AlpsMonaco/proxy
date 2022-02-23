# proxy
`This package contains multiple proxy/forward tools.`
`Uses unsafe package to reduce memory allocating and convert memory to protocol struct in net transfer part,which avoids memory-copy that will improve speed.`

## Progress

### Forward

- net.Conn interface forward complete. Usable
- TCP forward complete. Usable

### http

- unfinished

### socks5

- Socks5 client complete. Usable.

- Socks5 server complete. Usable.

### proxy  
a c/s architecture proxy tool over tcp.  
work properly.  
look `proxy\cmd\vpn\main.go` for usage.  

香港阿里云服务器测试通过，youtube流畅看4k，速度取决于服务器实际带宽，因为协议是自己新写的加上go的协程tcp为io非阻塞式，目前比较稳定。  
