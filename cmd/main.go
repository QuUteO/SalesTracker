package main

import (
	"SalesTracker/internal/config"
	"SalesTracker/internal/repository"
	"SalesTracker/internal/service"
	"context"
	"fmt"
	"os"
	"time"

	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
	"github.com/wb-go/wbf/logger"
)

func main() {
	cfg, err := config.New("./config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.InitLogger(
		logger.ZapEngine,
		"SalesTracker",
		cfg.Env,
		logger.WithLevel(logger.InfoLevel),
		logger.WithRotation("logs/app.log", 100, 5, 30),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации логгера: %v\n", err)
		os.Exit(1)
	}

	log.Info("Запуск приложения SalesTracker...")

	pg, err := pgxdriver.New(
		cfg.DSN,
		log,
		pgxdriver.MaxPoolSize(50),
		pgxdriver.MaxConnAttempts(5),
		pgxdriver.BaseRetryDelay(100*time.Millisecond),
	)
	if err != nil {
		log.Error("Не удалось подключиться к PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	ctx := context.Background()
	if err := pg.Ping(ctx); err != nil {
		log.Error("PostgreSQL недоступен", "error", err)
		os.Exit(1)
	}
	log.Info("PostgreSQL успешно подключен")

	repo := repository.New(pg)
	srv := service.New(repo)
}
