package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"SalesTracker/internal/config"
	"SalesTracker/internal/handler"
	"SalesTracker/internal/repository"
	"SalesTracker/internal/service"

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
	h := handler.New(srv)

	router := h.InitRoutes(cfg.Env)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Info(fmt.Sprintf("HTTP-сервер запущен на %s", cfg.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Критическая ошибка HTTP-сервера", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-stop
	log.Info("Получен сигнал завершения", "signal", sig.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Принудительная остановка сервера по таймауту", "error", err)
	} else {
		log.Info("HTTP-сервер успешно остановлен")
	}

	log.Info("Приложение SalesTracker успешно завершило работу")
}
