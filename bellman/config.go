package bellman

import "flag"

type Config struct {
  host string
  port string
  ssl  string
  ssl_cert string
  ssl_key  string
  cache_backend string
  cache_interval string
  redis_host string
  redis_port string
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
  cfg.host = flag.Lookup("host").Value.String()
  cfg.port = flag.Lookup("port").Value.String()
  cfg.ssl      = flag.Lookup("ssl").Value.String()
  cfg.ssl_cert = flag.Lookup("ssl_cert").Value.String()
  cfg.ssl_key  = flag.Lookup("ssl_key").Value.String()
  cfg.cache_backend  = flag.Lookup("cache_backend").Value.String()
  cfg.cache_interval = flag.Lookup("cache_interval").Value.String()
  cfg.redis_host = flag.Lookup("redis_host").Value.String()
  cfg.redis_port = flag.Lookup("redis_port").Value.String()

  return cfg
}
