package main

import (
//  "net"
//  "net/http/fcgi"
  "net/http"
  "runtime"
)

func main () {
  runtime.GOMAXPROCS(runtime.NumCPU() * 2)

  cfg := NewConfig()

  cache, cerr := CreateCache(cfg)
  if cerr != nil {
    print (cerr.Error())
    return
  }

  mh := MainHandler{cache}

  /*
  // FCFI Server
  unix, uerr := net.Listen("unix", "/tmp/sock.foo")
  if uerr != nil {
    print (uerr.Error())
    return
  }
  err = fcgi.Serve (unix, mh)
  */

  var err error
  url := cfg.host + ":" + cfg.port
  if cfg.ssl != "true" {
    err = http.ListenAndServe(url, mh)
  } else {
    err = http.ListenAndServeTLS(url, cfg.ssl_cert, cfg.ssl_key, mh)
  }
  if err != nil {
    print(err.Error())
  }
}
