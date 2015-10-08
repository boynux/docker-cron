package docker

import (
  "bytes"
  "log"
  dockerApi "github.com/fsouza/go-dockerclient"
)

type Docker struct {
  client *dockerApi.Client;
  containers []*dockerApi.Container;
}

func NewDocker(endpoint string) (*Docker, error) {
  docker := new(Docker)

  docker.client, _ = dockerApi.NewClient(endpoint)

  return docker, nil
}

func (c *Docker) Run(name string, image string, cmd []string) string {
  var id string 

  container, err := c.createContainer(name, image, cmd)

  if nil == err {
    err = c.client.StartContainer(container.ID, container.HostConfig)
    c.containers = append(c.containers, container)

    id = container.ID
  } else {
    log.Println("Can not create container ... ", err)
  }

  return id
}

func (c *Docker) Read(id string, stream *bytes.Buffer) {
  err := c.client.AttachToContainer(dockerApi.AttachToContainerOptions{
      Container:    id,
      OutputStream: stream,
      Logs:         true,
      Stdout:       true,
      Stderr:       true,
  })

  if err != nil {
      log.Fatal(err)
  }
}

func (c *Docker) Wait(id string) int {
  res, _ := c.client.WaitContainer(id)

  return res
}

func (c *Docker) Stop(id string) {
  _ = c.client.StopContainer(id, 5)
}

func (c *Docker) Remove(id string, force bool) error {
  return c.client.RemoveContainer(dockerApi.RemoveContainerOptions{ ID: id, Force: force})
}

func (c *Docker) createContainer(name string, image string, cmd []string) (container *dockerApi.Container, err error) {
  container, err = c.client.CreateContainer(
    dockerApi.CreateContainerOptions{
      Name: name, Config: &dockerApi.Config{Image: image, Cmd: cmd, AttachStdout: true},
      HostConfig: &dockerApi.HostConfig{},
  })

  return container, err
}
