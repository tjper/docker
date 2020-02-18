package docker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var docker, _ = client.NewEnvClient()

// Run a docker image in the local docker environment. Pass RunOption(s) in
// order to configure the Run.
//
// The return values of this function included the exported ports for the
// started container, a function to stop the docker container, and an error
// indicating the status of the call.
func Run(image string, opts ...RunOption) (ports []string, stop func() error, err error) {
	cfg := &runCfg{
		imagePrefix: "docker.io/library/",
		out:         ioutil.Discard,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	ctx := context.Background()
	image = cfg.imagePrefix + image

	out, err := docker.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to Run image; err: %w", err)
	}
	if _, err := io.Copy(cfg.out, out); err != nil {
		return nil, nil, fmt.Errorf("failed to Run image; err: %w", err)
	}

	containerCfg := &container.Config{
		Image: image,
	}
	hostCfg := &container.HostConfig{
		AutoRemove:      true,
		PublishAllPorts: true,
	}
	resp, err := docker.ContainerCreate(ctx, containerCfg, hostCfg, nil, "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to Run image; err: %w", err)
	}
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, nil, fmt.Errorf("failed to Run image; err: %w", err)
	}
	info, err := docker.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to Run image; err: %w", err)
	}
	for port := range info.Config.ExposedPorts {
		ports = append(ports, string(port))
	}
	stop = func() error {
		ctx := context.Background()
		return wrap(docker.ContainerStop(ctx, resp.ID, nil), "failed to stop running image")
	}
	return
}

// runCfg is a configuration for the Run function.
type runCfg struct {
	imagePrefix string
	out         io.Writer
}

// RunOption is a function that is used to modify the runCfg object.
type RunOption func(*runCfg)

// WithOut modifies the Run function to output to w.
func WithOut(out io.Writer) RunOption {
	return func(cfg *runCfg) {
		cfg.out = out
	}
}

// WithImagePrefix modifies the Run function to utilize a new image prefix when
// pulling an image.
func WithImagePrefix(prefix string) RunOption {
	return func(cfg *runCfg) {
		cfg.imagePrefix = prefix
	}
}

func wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(msg+"; err: %w", err)
}
