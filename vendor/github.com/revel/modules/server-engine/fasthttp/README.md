# FastHTTP
This module is the for [FastHTTP](https://github.com/valyala/fasthttp) server engine.
It does not support WebSockets.

### app.conf
- **server.engine** You must set this to `fasthttp` in order to use this server engine

### Other Notes
All features from that the go HTTP engine is supported by this server engine.
Memory usage is decreased by avoiding memory allocations. This should 
increase overall runtime performance and throughput. 
