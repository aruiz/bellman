package main

import (
  "sync"
  "github.com/garyburd/redigo/redis"
)

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

  return nil
}

func (rds RedisDataStore) GetStateForSession (session string) (string, error) {
  if rds.conn == nil {
    return "", nil
  }

  rds.mut.Lock()
  rep, err := (*rds.conn).Do ("GET", "sessions:" + session + ":state")
  rds.mut.Unlock()

  payload , err := redis.String(rep, err)
  return payload, err
}

func (rds RedisDataStore) Close () {
  rds.mut.Lock()
  (*rds.conn).Close()
  rds.mut.Unlock()
}
