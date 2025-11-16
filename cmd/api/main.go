package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/de4et/avito-test/internal/adapters/postgres"
	"github.com/de4et/avito-test/internal/server"
	"github.com/de4et/avito-test/internal/service"
	logger "github.com/de4et/avito-test/pkg"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.SetupLog("")

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic("DB_PORT MUST be integer")
	}

	pgClient := postgres.MustGetPostgresqlClient(postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_DATABASE"),
	})

	userRepository := postgres.NewPostgresqlUserRepository(pgClient)
	teamRepository := postgres.NewPostgresqlTeamRepository(pgClient)
	prRepository := postgres.NewPostgresqlPullRequestRepository(pgClient)

	transactor := postgres.NewPostgresqlTransactor(pgClient)
	teamService := service.NewTeamService(teamRepository, transactor)
	userService := service.NewUserService(userRepository, prRepository, transactor)
	prService := service.NewPullRequestService(prRepository, userRepository, teamRepository, transactor)

	routes := server.RegisterRoutes(teamService, userService, prService)
	server := server.NewServer(routes)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	slog.Debug("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Debug("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown with error")
	}

	slog.Debug("Server exiting")

	done <- true
}
