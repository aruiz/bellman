package bellman

import  (
  "net/http"
  "strings"
  "strconv"
)

//This handler implements the web service
type MainHandler struct {
  cache *Cache
}

func (mh MainHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  if len (path) < 2 {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  //TODO: Make separator configurable? Through a handler?
  key := strings.Join(path[1:len(path)-1], ":")

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

