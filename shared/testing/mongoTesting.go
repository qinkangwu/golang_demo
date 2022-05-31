package mongotesting

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"testing"
)

const (
	image         = "mongo"
	containerPort = "27017/tcp"
)

func RunWithMongoInDocker(m *testing.M, mongoUri *string) int {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	c := context.Background()
	cc, err := cli.ContainerCreate(c, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			containerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "")
	defer func(cli *client.Client, ctx context.Context, containerID string, options types.ContainerRemoveOptions) {
		err := cli.ContainerRemove(ctx, containerID, options)
		if err != nil {
			panic(err)
		}
	}(cli, c, cc.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}
	err2 := cli.ContainerStart(c, cc.ID, types.ContainerStartOptions{})
	if err2 != nil {
		panic(err2)
	}
	containerInspect, err4 := cli.ContainerInspect(c, cc.ID)
	if err4 != nil {
		panic(err4)
	}
	hostPort := containerInspect.NetworkSettings.Ports[containerPort][0]
	*mongoUri = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}
