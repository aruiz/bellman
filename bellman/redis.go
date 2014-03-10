/*
 * Bellman Project
 * Copyright (c) 2014, Alberto Ruiz <aruiz@gnome.org>, All rights reserved.
 * 
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 3.0 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library.
 */

package bellman

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

func (rds RedisDataStore) GetObject (session string) (string, error) {
  if rds.conn == nil {
    return "", nil
  }

  rds.mut.Lock()
  rep, err := (*rds.conn).Do ("GET", session)
  rds.mut.Unlock()

  payload , err := redis.String(rep, err)
  return payload, err
}

func (rds RedisDataStore) Close () {
  rds.mut.Lock()
  (*rds.conn).Close()
  rds.mut.Unlock()
}
