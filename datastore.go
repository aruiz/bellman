package main

type DataStore interface {
  GetObject (session string) (string, error)
  Close ()
}
