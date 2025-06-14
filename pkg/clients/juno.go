package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// JunoConfig represents the configuration for a Juno node
type JunoConfig struct {
	Network     string   `yaml:"network"`
	Port        string   `yaml:"port"`
	UseSnapshot bool     `yaml:"use_snapshot"`
	DataDir     string   `yaml:"data_dir"`
	Environment []string `yaml:"environment"`
}

// DefaultJunoConfig returns a default configuration for Juno
func DefaultJunoConfig() *JunoConfig {
	return &JunoConfig{
		Network:     "mainnet",
		Port:        "5050",
		UseSnapshot: true,
		DataDir:     "/data",
		Environment: []string{
			"JUNO_NETWORK=mainnet",
			"JUNO_HTTP_PORT=5050",
			"JUNO_HTTP_HOST=0.0.0.0",
		},
	}
}

// JunoClient represents a client for interacting with a local Nethermind Juno node
type JunoClient struct {
	dockerClient *client.Client
	containerID  string
	config       *JunoConfig
}

// NewJunoClient creates a new Juno client instance
func NewJunoClient(config *JunoConfig) (*JunoClient, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	if config == nil {
		config = DefaultJunoConfig()
	}

	return &JunoClient{
		dockerClient: dockerClient,
		config:       config,
	}, nil
}

// StartNode starts a local Nethermind Juno node in a Docker container
func (c *JunoClient) StartNode(ctx context.Context) error {
	// Pull the Nethermind Juno image
	_, err := c.dockerClient.ImagePull(ctx, "nethermind/juno:latest", types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull juno image: %w", err)
	}

	// Configure container
	containerConfig := &container.Config{
		Image: "nethermind/juno:latest",
		ExposedPorts: nat.PortSet{
			"5050/tcp": struct{}{},
		},
		Env: c.config.Environment,
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5050/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: c.config.Port,
				},
			},
		},
		Binds: []string{
			fmt.Sprintf("%s:/data", c.config.DataDir),
		},
	}

	// Create and start the container
	resp, err := c.dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		"juno-node",
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	c.containerID = resp.ID

	err = c.dockerClient.ContainerStart(ctx, c.containerID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for the node to be ready
	time.Sleep(5 * time.Second)

	return nil
}

// StopNode stops the running Juno node
func (c *JunoClient) StopNode(ctx context.Context) error {
	if c.containerID == "" {
		return nil
	}

	timeout := 10
	err := c.dockerClient.ContainerStop(ctx, c.containerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	err = c.dockerClient.ContainerRemove(ctx, c.containerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// GetNodeStatus returns the status of the Juno node
func (c *JunoClient) GetNodeStatus(ctx context.Context) (string, error) {
	if c.containerID == "" {
		return "not running", nil
	}

	container, err := c.dockerClient.ContainerInspect(ctx, c.containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	return container.State.Status, nil
}

// Close cleans up resources
func (c *JunoClient) Close() error {
	if c.dockerClient != nil {
		return c.dockerClient.Close()
	}
	return nil
}
