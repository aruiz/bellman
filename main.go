package main

import (
  "net/http"
  "strings"
  "sync"
  "time"
  "runtime"
  "github.com/garyburd/redigo/redis"
  "fmt"
)


type Session struct {
  payload string
  lock    sync.RWMutex
}

type MainHandler struct {
  /* Make a map of payloads and locks for each session being served */
  sessions     map[string]*Session
  rconn        *redis.Conn
}

func (mh MainHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  switch {
    case strings.HasPrefix(path, "/session"):
        mh.HandleSessionRequest (w, r)
  }
}

func (mh MainHandler) HandleSessionRequest (w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  if len (path) != 3 {
    return
  }
  sessionName := path[2]
  if r.Method == "GET" {
    session := mh.sessions[sessionName]
    fmt.Println (session.payload)

    session.lock.RLock()
    w.Write([]byte(session.payload))
    session.lock.RUnlock()
  }
  if r.Method == "POST" {
    mh.AddSession(sessionName)
  }
}

func (mh MainHandler) AddSession (session string) {
  sessions := mh.sessions
  _, ok := sessions[session]
  if ok {
    return
  }

  sessions[session] = new(Session)
}

func (mh MainHandler) StartSessions () {
  go func (mh MainHandler) {
    dur, err := time.ParseDuration("1s")
    if err != nil {
      return;
    }

    for {
      for sessionName, session := range mh.sessions {
        rep, err := (*mh.rconn).Do ("GET", "sessions:" + sessionName + ":state")
        payload , err := redis.String(rep, err)

        if err != nil {
          delete (mh.sessions, sessionName)
        }

        session.lock.Lock()
        session.payload = payload
        session.lock.Unlock()
      }
      time.Sleep(dur)
    }
  } (mh)
}

func main () {
  runtime.GOMAXPROCS(runtime.NumCPU())

  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    return;
  }

  mh := MainHandler{make(map[string]*Session), &c}
  mh.AddSession("live")
  mh.StartSessions()

  http.ListenAndServe(":8080", mh)
}
