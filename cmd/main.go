package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/config"
	"github.com/mertozler/internal/dockerengine"
	"github.com/mertozler/internal/handler"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/gitcloner"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	config, err := config.LoadConfig("./configs")
	if err != nil {
		logrus.Fatal("Error while getting configs", err)
	}

	scanHandler, err := initializingHandler(&config)
	if err != nil {
		logrus.Errorf("Error while initializing handler: %v", err)
	}

	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/newscan", scanHandler.NewScan())
	api.Get("/scan/:scanid", scanHandler.GetScan())

	err = app.Listen(config.Host.Port)
	if err != nil {
		logrus.Error("Error while initializing fiber", err)
	}
}

func initializingHandler(config *config.Config) (*handler.ScanHandler, error) {
	repo := repository.NewRepository(config.Redis)
	cloner := gitcloner.NewGitClonner()
	dockerEngine, err := dockerengine.NewDockerEngine(config.DockerEngine)
	if err != nil {
		return nil, err
	}

	scanHandler := handler.NewScanHandler(cloner, repo, dockerEngine)
	return scanHandler, nil
}
