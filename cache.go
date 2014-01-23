package main

import (
  "time"
  "sync"
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
func CreateCache () (*Cache, error) {
  // Use Redis connection
  //TODO: Make type of datastore configurable
  rds := RedisDataStore{}
  err := rds.Connect(":6379")
  if err != nil {
    return &Cache{}, err
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

func UpdateCache (c *Cache) {
  dur, _ := time.ParseDuration("1s")
  for {
    timeout := time.After(dur)
    c.lock.RLock ()
    for key := range c.cache {
      c.lock.RUnlock ()

      payload, err := c.ds.GetStateForSession(key)
      if err != nil {
        //FIXME: Log error
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
