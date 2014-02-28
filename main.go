package main

import (
//  "net"
//  "net/http/fcgi"
  "net/http"
  "strings"
  "runtime"
  "strconv"
  "flag"
)

/*-----------------------------------------*/

//This handler implements the web service
type MainHandler struct {
  cache *Cache
}

func (mh MainHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  switch {
    case strings.HasPrefix(path, "/sessions"):
        mh.HandleSessionRequest (w, r)
    default:
      w.WriteHeader(http.StatusNotFound)
  }
}

func (mh MainHandler) HandleSessionRequest (w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  if len (path) != 3 {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  session := path[2]

  if r.Method == "GET" {
    payload, err := mh.cache.GetObject("sessions:" + session + ":state")
    if err != nil {
      w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
      w.WriteHeader(http.StatusNotFound)
      return
    }

    w.Write([]byte(payload))
    return
  }

  w.WriteHeader(http.StatusNotFound)
}

func main () {
  flag.String("config", "", "path to configuration file")
  flag.Parse()
  configFile := flag.Lookup("config").Value.String()

  if configFile == "" {
    //TODO: Find file in $HOME/.local/etc and /etc
  }

  //TODO: Load JSON file

  runtime.GOMAXPROCS(runtime.NumCPU() * 2)
  cache, cerr := CreateCache()
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

  err := http.ListenAndServe(":8080", mh)
  //err = http.ListenAndServeTLS(":8080", "vhost1.crt", "vhost1.key", mh)
  if err != nil {
    print(err.Error())
  }
}
