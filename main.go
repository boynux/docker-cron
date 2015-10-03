package main

import (
  "time"
  "fmt"
  "bytes"
  "github.com/boynux/docker-cron/docker-helper"
)

const (
  Duration = 5
)

func main() {

  endpoint := "unix:///var/run/docker.sock"
  client, _ := docker.NewDocker(endpoint)
  ticker := time.NewTicker(Duration * time.Second)

  for _ = range ticker.C {
    fmt.Println("Timer's ticking ... ", time.Now())
    go func() {
      var stream bytes.Buffer

      id := client.Run("test", "busybox", []string{"echo", "-n", "Hello world!"})

      client.Wait(id)
      client.Read(id, &stream)

      fmt.Println(stream.String())

      client.Remove(id, false)
    }()
  }
}
