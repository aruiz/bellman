package main

type DataStore interface {
  GetStateForSession (session string) (string, error)
  Close ()
}
