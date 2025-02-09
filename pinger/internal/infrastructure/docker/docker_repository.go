package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"pinger/internal/domain"
)

type DockerRepository struct {
	cli *client.Client
}

func NewDockerRepository(cli *client.Client) *DockerRepository {
	return &DockerRepository{cli: cli}
}

func (r *DockerRepository) ListContainers(all bool) ([]domain.ContainerInfo, error) {
	containers, err := r.cli.ContainerList(context.Background(), container.ListOptions{
		All: all,
	})
	if err != nil {
		return nil, err
	}

	var results []domain.ContainerInfo
	for _, c := range containers {
		ci := domain.ContainerInfo{
			ID:        c.ID,
			Name:      trimSlash(c.Names),
			Image:     c.Image,
			State:     c.State,
			Status:    c.Status,
			CreatedAt: time.Unix(c.Created, 0),
		}
		results = append(results, ci)
	}
	return results, nil
}

func (r *DockerRepository) InspectContainer(containerID string) (domain.ContainerInfo, []string, error) {
	f := filters.NewArgs()
	f.Add("id", containerID)

	containers, err := r.cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return domain.ContainerInfo{}, nil, fmt.Errorf("%w", err)
	}

	if len(containers) == 0 {
		return domain.ContainerInfo{}, nil, fmt.Errorf("%s not found", containerID)
	}

	c := containers[0]

	ci := domain.ContainerInfo{
		ID:        c.ID,
		Name:      trimSlash(c.Names),
		Image:     c.Image,
		State:     c.State,
		Status:    c.Status,
		CreatedAt: time.Unix(c.Created, 0),
	}

	var ips []string
	if c.NetworkSettings != nil && c.NetworkSettings.Networks != nil {
		for _, netSettings := range c.NetworkSettings.Networks {
			if netSettings.IPAddress != "" {
				ips = append(ips, netSettings.IPAddress)
			}
		}
	}

	return ci, ips, nil
}

func trimSlash(names []string) string {
	if len(names) == 0 {
		return ""
	}
	name := names[0]
	if len(name) > 0 && name[0] == '/' {
		return name[1:]
	}
	return name
}

func trimLeadingSlash(name string) string {
	if len(name) > 0 && name[0] == '/' {
		return name[1:]
	}
	return name
}

func parseTime(dockerTime string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, dockerTime)
	if err != nil {
		return time.Now()
	}
	return t
}
