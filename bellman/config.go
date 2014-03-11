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

import "flag"

type Config struct {
  Host string
  Port string
  Ssl  string
  SslCert string
  SslKey  string
  CacheBackend string
  CacheInterval string
  RedisHost string
  RedisPort string
}

func NewConfig () (*Config) {
  flag.String("host", "localhost", "default host to listen from")
  flag.String("port", "8080", "default port to listen from")
  flag.String("ssl", "false", "whether we serve data using HTTPS")
  flag.String("ssl_cert", "./cert.pem", "path to SSL certificate")
  flag.String("ssl_key",  "./cert.key", "path to SSL key")
  flag.String("cache_backend",  "redis", "cache backend to use")
  flag.String("cache_interval", "1s",    "interval used to update the cache keys")
  flag.String("redis_host", "localhost", "redis server host")
  flag.String("redis_port", "6379", "redis server port")

  flag.Parse()

  cfg := new(Config)
  cfg.Host = flag.Lookup("host").Value.String()
  cfg.Port = flag.Lookup("port").Value.String()
  cfg.Ssl      = flag.Lookup("ssl").Value.String()
  cfg.SslCert = flag.Lookup("ssl_cert").Value.String()
  cfg.SslKey  = flag.Lookup("ssl_key").Value.String()
  cfg.CacheBackend  = flag.Lookup("cache_backend").Value.String()
  cfg.CacheInterval = flag.Lookup("cache_interval").Value.String()
  cfg.RedisHost = flag.Lookup("redis_host").Value.String()
  cfg.RedisPort = flag.Lookup("redis_port").Value.String()

  return cfg
}
