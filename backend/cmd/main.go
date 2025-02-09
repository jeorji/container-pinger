package main

import (
	"context"
	"log"
    "os"

	"backend/internal/infrastructure/in"
	"backend/internal/infrastructure/out"
	"backend/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GET    /api/containers             // Список контейнеров со статусами
// PUT    /api/containers/:id         // Создать / Обновить информацию о контейнере
// POST   /api/containers/:id/pings   // Добавить запись о пинге для конкретного контейнера

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer pool.Close()

	pingRepo := repository.NewPingPGRepo(pool)
	containerRepo := repository.NewContainerPGRepo(pool)

	containerStatsUC := usecase.NewContainerStatsUseCase(containerRepo, pingRepo)
	containerUpdateUC := usecase.NewContainerUpdateUseCase(containerRepo)
	pingUC := usecase.NewPingUseCase(pingRepo)

	handler := restapi.NewHandler(containerStatsUC, containerUpdateUC, pingUC)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	restapi.SetupRoutes(r, handler)

	r.Run(":8080")
}
