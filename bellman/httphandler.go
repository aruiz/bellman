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

import  (
  "net/http"
  "strings"
  "strconv"
)

//This handler implements the web service
type HttpHandler struct {
  cache *Cache
}

func NewHttpHandler (cache *Cache) (*HttpHandler){
  hh := new(HttpHandler)
  hh.cache = cache
  return hh
}

func (mh HttpHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  if len (path) < 2 {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  //TODO: Make separator configurable? Through a handler?
  key := strings.Join(path[1:len(path)], ":")
  print(key+"\n")

  if r.Method != "GET" {
    w.WriteHeader(http.StatusNotFound)
  }
  payload, err := mh.cache.GetObject(key)
  if err != nil {
    w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
    w.WriteHeader(http.StatusNotFound)
    return
  }

  w.Write([]byte(payload))
  return
}

