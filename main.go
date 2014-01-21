package main

import (
  "net/http"
  "strings"
  "sync"
  "time"
  "runtime"
  "github.com/garyburd/redigo/redis"
)

type DataStore interface {
  GetStateForSession (session string) (string, error)
}

type RedisDataStore struct {
  conn        *redis.Conn
}

func (rds RedisDataStore) Connect (url string) error {
  c, err := redis.Dial("tcp", url)
  if err != nil {
    return err
  }

  rds.conn = &c
  return nil
}

func (rds RedisDataStore) GetStateForSession (session string) (string, error) {
  if rds.conn == nil {
    return "", nil
  }

  rep, err := (*rds.conn).Do ("GET", "sessions:" + session + ":state")
  payload , err := redis.String(rep, err)
  return payload, err
}

type Session struct {
  payload string
  lock    sync.RWMutex
}


//This handler implements the web service
type MainHandler struct {
  /* Make a map of payloads and locks for each session being served */
  sessions     map[string]*Session
  store        DataStore
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
  path := strings.Split(r.URL.Path, "/")
  if len (path) != 3 {

    return
  }
  sessionName := path[2]
  if r.Method == "GET" {
    session := mh.sessions[sessionName]

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
        payload , err := mh.store.GetStateForSession (sessionName)
        if err != nil {
          delete (mh.sessions, sessionName)
          return
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

  rds := RedisDataStore{}
  err := rds.Connect (":6379")
  if err != nil {
    return
  }

  mh := MainHandler{make(map[string]*Session), &rds}

  mh.AddSession("live")
  mh.StartSessions()

  http.ListenAndServe(":8080", mh)
}
