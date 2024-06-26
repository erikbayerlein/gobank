package main

import (
	"log"
	"projects/gobank/storage"
)

func main() {
  store, err := storage.NewPostgresStore()
  if err != nil {
    log.Fatal(err)
  }

  if err := store.Init(); err != nil {
    log.Fatal(err)
  }

  server := NewAPIServer(":3000", store)
  server.Run()
}
