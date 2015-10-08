package main

import (
  "time"
  "fmt"
  "log"
  "net/http"
  "github.com/pborman/uuid"
  "github.com/boynux/docker-cron/docker-helper"
)

const (
  Duration = 5
)

func handler(w http.ResponseWriter, r *http.Request) {
  log.Println("Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

  endpoint := "unix:///var/run/docker.sock"
  client, _ := docker.NewDocker(endpoint)
  ticker := time.NewTicker(Duration * time.Second)

  // Launch HTTP server for callbacks
  go func () {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
  }()

  q := make(chan string)

  // Docker garbage collector
  go func(c chan string) {
    for id := range q {
      client.Wait(id)
      client.Remove(id, false)
    }
  }(q)

  // set a timer job to spawn new Docker container
  for _ = range ticker.C {
    log.Println("Timer's ticking ... ", time.Now())

    for i := 0; i < 10; i++ {
      go func(c chan string) {
        c <- client.Run(uuid.New(), "tutum/curl:latest", []string{"curl", "http://172.17.42.1:8080/" + "Hello/No/" + string(i)})
      }(q)
    }
  }
}

