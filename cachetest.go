package main

import "time"

func main () {
  cache, err := CreateCache()
  if err != nil {
    print (err.Error())
    return
  }

  p, err :=  cache.GetPayload("live")
  if err != nil {
    print (err.Error())
    return
  }

  for {
    dur, _ := time.ParseDuration("1s")
    timeout := time.After(dur)

    for i := 0; i < 30000; i++ {
      go func (c *Cache, p string) {
        payload, _ := cache.GetPayload("live")
        if payload != p {
          print ("corruption")
          print (payload)
          return
        }
      } (&cache, p)
    }

    print ("iteration\n")
    <-timeout
  }
}
