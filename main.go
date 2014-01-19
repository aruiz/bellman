package main

import (
  "net/http"
  "strings"
  "sync"
  "time"
  "runtime"
)

type MainHandler struct {
  payload     *string
  payloadLock *sync.RWMutex
}

func (mh MainHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  switch {
    case strings.HasPrefix(path, "/session"):
        mh.HandleSessionRequest (w, r)
  }
}

func (mh MainHandler) HandleSessionRequest (w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET"{
    mh.payloadLock.RLock()
    w.Write([]byte(*mh.payload))
    mh.payloadLock.RUnlock()
  }
}

func main () {
  runtime.GOMAXPROCS(runtime.NumCPU())
  mh := MainHandler{payload: new(string), payloadLock: new(sync.RWMutex)}

  go func (mh MainHandler) {
    dur, err := time.ParseDuration("1s")
    if err != nil {
      return;
    }

    for {
      time.Sleep(dur)
      mh.payloadLock.Lock()
      *mh.payload = "{1: [1, 1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]}"
      mh.payloadLock.Unlock()
    }
  } (mh)

  http.ListenAndServe(":8080", mh)
}
