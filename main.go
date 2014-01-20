package main

import (
  "net/http"
  "strings"
  "sync"
  "time"
  "runtime"
  "github.com/garyburd/redigo/redis"
)

type MainHandler struct {
  payload     *string
  payloadLock *sync.RWMutex
  rconn       *redis.Conn
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

  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    return;
  }

  go func (mh MainHandler, conn *redis.Conn, session string) {
    dur, err := time.ParseDuration("1s")
    if err != nil {
      return;
    }

    for {
      rep, err := (*conn).Do ("GET", "sessions:" + session + ":state")
      payload , err := redis.String(rep, err)
      if err != nil {
        continue;
      }

      mh.payloadLock.Lock()
      *mh.payload = payload
      mh.payloadLock.Unlock()

      time.Sleep(dur)
    }
  } (mh, &c, "live")

  http.ListenAndServe(":8080", mh)
}
