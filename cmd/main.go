package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/config"
	"github.com/mertozler/internal/dockerengine"
	"github.com/mertozler/internal/handler"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/git-clonner"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	config, err := config.LoadConfig("./configs")
	if err != nil {
		logrus.Fatalf("Error while getting configs", err)
	}
	repo := repository.NewRepository(config.Redis)
	cloner := git_clonner.NewGitClonner()
	dockerEngine, err := dockerengine.NewDockerEngine()
	if err != nil {
		logrus.Errorf("Error while creating docker engine: %v", err)
	}

	scanHandler := handler.NewScanHandler(cloner, repo, dockerEngine)

	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/newscan", scanHandler.NewScan())
	api.Get("/scan/:scanid", scanHandler.GetScan())

	err = app.Listen(config.Host.Port)
	if err != nil {
		logrus.Errorf("Error while initializing fiber", err)
	}
}
