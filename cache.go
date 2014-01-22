package main

import (
//  "time"
//  "sync"
  "github.com/garyburd/redigo/redis"
)


type DataStore interface {
  GetStateForSession (session string) (string, error)
}

func (rds RedisDataStore) GetStateForSession (session string) (string, error) {
  if rds.conn == nil {
    return "", nil
  }

  rep, err := (*rds.conn).Do ("GET", "sessions:" + session + ":state")
  payload , err := redis.String(rep, err)

  return payload, err
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

type Cache struct {
  ds DataStore
  cache map[string]string
}

func CreateCache () (Cache, error) {
  conn, err := redis.Dial("tcp", ":6379")
  if err != nil {
    return Cache{}, err
  }

  ds := new(RedisDataStore)
  ds.conn = &conn

  go func (

  return Cache{ds}, nil
}

func (self *Cache) SetPayload (key string, payload string) error {
  return nil
}

func (self *Cache) GetPayload (key string) (string, error) {
  if 
  return  self.ds.GetStateForSession(key)
}
