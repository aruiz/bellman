package main

import (
  "time"
  "sync"
  "errors"
)

type Cache struct {
  ds DataStore
  cache map[string]string
  lock sync.RWMutex
  running   bool
  stop chan bool
  ret  chan bool
}

//TODO: set interval and configuration as argument
func CreateCache (cfg *Config) (*Cache, error) {
  if cfg.cache_backend == "redis" {
    rds := RedisDataStore{}
    err := rds.Connect(cfg.redis_host + ":" + cfg.redis_port)
    if err != nil {
      return nil, err
    }
    cache := Cache{&rds,
                   make(map[string]string),
                   sync.RWMutex{},
                   true,
                   make(chan bool, 1),
                   make(chan bool, 1),
                 }

    go UpdateCache(&cache)
    return &cache, nil
  }
  return nil, errors.New("CreateCache: No valid cache_backend option was provided")
}

func UpdateCache (c *Cache) {
  dur, _ := time.ParseDuration("1s")
  for {
    timeout := time.After(dur)
    c.lock.RLock ()
    for key := range c.cache {
      c.lock.RUnlock ()

      payload, err := c.ds.GetObject(key)
      if err != nil {
        print (err.Error ())
        return
      }

      c.lock.Lock()
      c.cache[key] = payload
      c.lock.Unlock()

      c.lock.RLock ()
    }
    c.lock.RUnlock ()

    select {
    case foo := <-c.stop:
      c.ret <- foo
      return
    default:
    }

    //Wait for a second to pass
    <-timeout
  }
}

func (self *Cache) Close () {
  self.stop <- true
  self.running = <-self.ret
  self.ds.Close()
}

func (self *Cache) GetObject (key string) (string, error) {

  self.lock.RLock()
  value, ok := self.cache[key]
  self.lock.RUnlock()

  if !ok {
    v, err := self.ds.GetObject(key)
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
