package main

import (
//  "net"
//  "net/http/fcgi"
  "net/http"
  "runtime"
  "github.com/aruiz/bellman/bellman"
)

func main () {
  runtime.GOMAXPROCS(runtime.NumCPU() * 2)

  cfg := bellman.NewConfig()

  cache, cerr := bellman.CreateCache(cfg)
  if cerr != nil {
    print (cerr.Error())
    return
  }

  mh := bellman.NewHttpHandler(cache)

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
  url := cfg.Host + ":" + cfg.Port
  if cfg.Ssl != "true" {
    print ("starting bellman at address http://" +url)
    err = http.ListenAndServe(url, mh)
  } else {
    print ("starting bellman at address https://" +url)
    err = http.ListenAndServeTLS(url, cfg.SslCert, cfg.SslKey, mh)
  }
  if err != nil {
    print(err.Error())
  }
}
