package main

import (
  "net/http"
  "strings"
  "runtime"
)

/*-----------------------------------------*/

//This handler implements the web service
type MainHandler struct {
  cache *Cache
}

func (mh MainHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  switch {
    case strings.HasPrefix(path, "/session"):
        mh.HandleSessionRequest (w, r)
    default:
      w.WriteHeader(http.StatusNotFound)
  }
}

func (mh MainHandler) HandleSessionRequest (w http.ResponseWriter, r *http.Request) {
  //TODO: Write some payload in the response
  //TODO: Improve headers a bit
  path := strings.Split(r.URL.Path, "/")
  if len (path) != 3 {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  session := path[2]

  if r.Method == "GET" {
    payload, err := mh.cache.GetPayload(session)
    if err != nil {
      w.WriteHeader(http.StatusNotFound)
      return
    }

    w.Write([]byte(payload))
    return
  }

  w.WriteHeader(http.StatusNotFound)
}

func main () {
  runtime.GOMAXPROCS(runtime.NumCPU())
  cache, err := CreateCache()
  if err != nil {
    //TODO: Log error
    return
  }

  mh := MainHandler{cache}

  http.ListenAndServe(":8080", mh)
}
