package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"

	"pinger/internal/infrastructure/backend"
	"pinger/internal/infrastructure/docker"
	"pinger/internal/infrastructure/ping"
	"pinger/internal/usecase"
)

func main() {
    backendURL := os.Getenv("BACKEND_URL")
    intervalStr := os.Getenv("PING_INTERVAL")

    intervalSec, err := strconv.Atoi(intervalStr)
    if err != nil {
        log.Fatalf("%v", err)
    }
    pingInterval := time.Duration(intervalSec) * time.Second

    dockerClient, err := client.NewClientWithOpts(
        client.FromEnv,
        client.WithAPIVersionNegotiation(),
    )
    if err != nil {
        log.Fatalf("%v", err)
    }

    dockerRepo := docker.NewDockerRepository(dockerClient)
    backendRepo := backend.NewBackendRepository(backendURL)
    pingService := ping.NewProBingService()

    containerSyncUC := usecase.NewContainerSyncUseCase(dockerRepo, backendRepo)
    pingContainersUC := usecase.NewPingContainersUseCase(dockerRepo, backendRepo, pingService)

    ctx := context.Background()
    ticker := time.NewTicker(pingInterval)
    defer ticker.Stop()

    for {
        if err := containerSyncUC.SyncContainers(ctx, true); err != nil {
            log.Printf("Error SyncContainers: %v\n", err)
        }

        if err := pingContainersUC.PingContainers(ctx, true); err != nil {
            log.Printf("Error PingContainers: %v\n", err)
        }

        <-ticker.C
    }
}
