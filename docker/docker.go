package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func GenAuthString(userName string, password string) (string, error) {
	authConfig := types.AuthConfig{
		Username: userName,
		Password: password,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	return authStr, nil
}

func RunContainer(userName string,
	password string,
	imageName string,
	containerName string,
	port string,
	networkName string,
	shouldWait bool,
) (container.ContainerCreateCreatedBody, error) {

	authStr, err := GenAuthString(userName, password)

	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	reader, err := cli.ImagePull(ctx, imageName,
		types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}
	// TODO: writing to log
	io.Copy(os.Stdout, reader)

	config := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}

	networkingConfig := networktypes.NetworkingConfig{
		EndpointsConfig: map[string]*networktypes.EndpointSettings{
			networkName: {},
		},
	}

	resp, err := cli.ContainerCreate(
		ctx, config, hostConfig, &networkingConfig, nil, containerName)
	if err != nil {
		return resp, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return resp, err
	}

	// TODO: need debug this to see if this fits cadence
	if shouldWait {
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				return container.ContainerCreateCreatedBody{}, err
			}
		case <-statusCh:
		}
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return resp, err
	}

	// TODO: writing to log
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return resp, nil
}

func StopContainer(
	containerName string,
) error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	return cli.ContainerStop(ctx, containerName, nil)
}
