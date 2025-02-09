package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pinger/internal/domain"
)

type BackendRepository struct {
	baseURL string
	client  *http.Client
}

func NewBackendRepository(baseURL string) *BackendRepository {
	return &BackendRepository{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (r *BackendRepository) UpdateContainer(containerInfo domain.ContainerInfo) error {
	url := fmt.Sprintf("%s/api/containers/%s", r.baseURL, containerInfo.ID)

	data, err := json.Marshal(containerInfo)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%d", resp.StatusCode)
	}

	return nil
}

func (r *BackendRepository) PostPing(pingRes domain.PingResult) error {
	url := fmt.Sprintf("%s/api/containers/%s/pings", r.baseURL, pingRes.ContainerID)

	data, err := json.Marshal(pingRes)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%d", resp.StatusCode)
	}

	return nil
}
