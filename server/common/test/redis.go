package test

import (
	"GoYin/server/common/consts"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"testing"
)

func RunRedisInDocker(db int, t *testing.T) (cleanUpFunc func(), cli *redis.Client, err error) {
	c, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return func() {}, nil, err
	}

	ctx := context.Background()

	resp, err := c.ContainerCreate(ctx,
		&container.Config{
			Image: consts.RedisImage,
			ExposedPorts: nat.PortSet{
				consts.RedisContainerPort: {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				consts.RedisContainerPort: []nat.PortBinding{
					{
						HostIP:   consts.RedisContainerIP,
						HostPort: consts.RedisPort,
					},
				},
			},
		}, nil, nil, "")

	if err != nil {
		return func() {}, nil, err
	}

	containerID := resp.ID
	cleanUpFunc = func() {
		err = c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			t.Error("remove test docker failed", err)
		}
	}

	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return cleanUpFunc, nil, err
	}

	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		return cleanUpFunc, nil, err
	}
	hostPort := inspRes.NetworkSettings.Ports[consts.RedisContainerPort][0]

	return cleanUpFunc, redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", consts.RedisContainerIP, hostPort.HostPort),
		DB:   db,
	}), nil
}
