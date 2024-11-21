package main

import (
	"github.com/JoseBeteta/surfe/app"
	user_application "github.com/JoseBeteta/surfe/app/application"
	"github.com/JoseBeteta/surfe/app/infrastructure/common/configx"
	http2 "github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	action_infrastructure "github.com/JoseBeteta/surfe/app/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

var (
	hash, version string
)

func main() {
	config, err := app.LoadConfig(configx.NewLoader())

	// panic if there is an error
	if err != nil {
		panic(err)
	}

	logger := createLogger(hash, version, config.ServiceName, "test", "server")

	// Setup and start the server
	server := setupServer(*config, *logger)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("error listening server", slog.String("error", err.Error()))
		panic(err)
	}
}

func setupServer(cfg app.Config, log slog.Logger) *http.Server {
	r := gin.New()
	httpMapper := http2.NewHttpMapper(log)
	httpMapper.Initialize(r)

	usersFile := os.Getenv("USERS_FILE")
	actionsFile := os.Getenv("ACTIONS_FILE")

	userReadRepository := action_infrastructure.NewUserJSONRepository(usersFile)
	actionReadRepository := action_infrastructure.NewActionJSONRepository(actionsFile)

	http2.RegisterHomeHandler(r)

	userHandler := user_application.NewUserHandler(
		userReadRepository,
		log,
		httpMapper,
	)

	actionHandler := user_application.NewActionHandler(
		actionReadRepository,
		log,
		httpMapper,
	)

	userHandler.Initialize(r)
	actionHandler.Initialize(r)

	return http2.NewServer(cfg.Server, r)
}

// Custom function to create a logger with default fields
func createLogger(hash, version, serviceName, env, component string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	return logger.With(
		slog.String("hash", hash),
		slog.String("version", version),
		slog.String("service", serviceName),
		slog.String("environment", env),
		slog.String("component", component),
	)
}
