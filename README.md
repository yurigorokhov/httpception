#HTTPCeption

Description
===========
HTTPCeption is a HTTP proxy that allows the interception of HTTP requests and responses for debugging purposes. During normal operation it acts as a transparent proxy, but once you connect to it over a web interface, you are able to analyze the network traffic step by step.

Compilation
===========
Install godep executable and make sure it is in your path

```
go get github.com/tools/godep
```

then build

```
make
```

Running
=======
```
Usage of ./httpception:
  -debug=":9999": Address to listen for debugging connection (ex: :9999)
  -listen="": Address to listen for new connections (ex: localhost:3333)
  -send="": Address to forward traffic to (ex: www.w3.org:80)
```

Example: forward traffic from localhost to http://www.w3.org/

```
./httpception -listen="localhost:3333" -send="www.w3.org:80"
```

Now in your browser navigate to

```
http://localhost:3333
```

and in a separate window navigate to the debugging interface
```
http://localhost:9999
```

You should see a trail of requests coming through:
![Screenshot](/images/screenshot.png)

TODO
====
- [ ] Support Host header rewriting
- [ ] Support modifying requests and responses in the debugger
- [ ] Support HTTPS
- [ ] Allow replaying of requests
- [ ] Add ability to save / load requests
