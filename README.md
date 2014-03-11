# Bellman #

Bellman is a Go library to create an HTTP service to broadcast data backed
in a key/value datastore under a fixed time update policy. Think of it as a
front-line cache for rapidly changing data.


In this diagram I describe the current architecture of bellman:

```
                                                             +-----------------
                                                             | HTTP(S) clients
                                                             |
             Update data in fixed intervals                  |
+----------------+                +------------------+<------+
|  object Store  |   (1s/1m/1h)   | bellman service  |<------+
|                +--------------->|                  |<------+
|                |<---------------+                  |<------+
+----------------+                |                  |<------+
                                  +------------------+<------+
                                                             |
                                                             |
                                                             +-----------------
```

Currently it takes the URL path and turns it into a colon separated key,
this is an example for a bellamn server listening to http://localhost:8080

Say you have a redis database with a JSON payload on the "foo:bar:baz" key,
to retrieve that key you make a request to http://localhost:8080/foo/bar/baz
and you get that JSON payload back.

After this, the object will be kept in the cache and will be updated every
second (or other interval if specified).

### ADDING OTHER DATA STORES ###
Other cache backends are easily added provided you
implement provided it implements the DataStore interface, the Config struct
is extended and the Cache object gets support for it.

### PERFORMANCE ###
Managed to serve 1kb object at 45.000 HTTP requests per second in a dual core
Intel i5-2557M CPU at 1.70GHz and 4GB of RAM.

HTTPS reduces this mark by an order of magnitude, reaching barely 5-7k rps.

### OTHER DATA PROVIDERS ###
You get a cache request for the key "foo:bar:baz", which for now only points
to a Redis server. Other cache backends are easily added provided you
implement provided it implements the DataStore interface, the Config struct
is extended and the Cache object gets support for it.

### INSTALL ###
```
go get github.com/aruiz/bellman
```

### SERVER USAGE ###
```
$GOCODE/bin/bellman -help
```

## TODO ##
- Discard objects from the internal cache that haven't been requested for a while
- Make the separator for the key more generic
- Write tests for the provided funcionality
- Add other key/value store providers
