package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"time"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	c := context.Background()
	cc, err := cli.ContainerCreate(c, &container.Config{
		Image: "mongo",
		ExposedPorts: nat.PortSet{
			"27017/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"27017/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	err2 := cli.ContainerStart(c, cc.ID, types.ContainerStartOptions{})
	if err2 != nil {
		panic(err2)
		return
	}
	time.Sleep(10 * time.Second)

	err3 := cli.ContainerRemove(c, cc.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err3 != nil {
		return
	}
}
