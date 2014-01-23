package main

import (
//  "time"
  "sync"
  "github.com/garyburd/redigo/redis"
)


type DataStore interface {
  GetStateForSession (session string) (string, error)
}


type RedisDataStore struct {
  mut         *sync.Mutex
  conn        *redis.Conn
}

func (rds *RedisDataStore) Connect (url string) error {
  c, err := redis.Dial("tcp", url)
  if err != nil {
    return err
  }

  rds.mut = new(sync.Mutex)
  rds.conn = &c

  print (rds.conn)

  return nil
}

func (rds RedisDataStore) GetStateForSession (session string) (string, error) {
  print (rds.conn)
  if rds.conn == nil {
    return "", nil
  }

  rds.mut.Lock()
  rep, err := (*rds.conn).Do ("GET", "sessions:" + session + ":state")
  rds.mut.Unlock()

  payload , err := redis.String(rep, err)
  return payload, err
}

type Cache struct {
  ds DataStore
  cache map[string]string
  lock *sync.RWMutex

}

func CreateCache () (Cache, error) {

  //go func () {
  //} ()

  rds := RedisDataStore{}
  err := rds.Connect(":6379")
  if err != nil {
    return Cache{}, err
  }

  cache := Cache{&rds,
                 make(map[string]string),
                 new(sync.RWMutex)}

  print (rds.conn)
  return cache, nil
}


func (self *Cache) Close () {
}

func (self *Cache) SetPayload (key string, payload string) error {
  return nil
}

func (self *Cache) GetPayload (key string) (string, error) {

  self.lock.RLock()
  value, ok := self.cache[key]
  self.lock.RUnlock()

  if !ok {
    v, err := self.ds.GetStateForSession(key)
    if err != nil {
      return "", err
    }

    value = v

    self.lock.Lock ()
    self.cache[key] = value
    self.lock.Unlock ()
  }

  return value, nil
}
