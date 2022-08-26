package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/config"
	"github.com/mertozler/internal/handler"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/git-clonner"
	"log"
)

func main() {
	config, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepository(config.Redis)

	app := fiber.New()

	cloner := git_clonner.NewGitClonner()

	api := app.Group("/api/v1")
	api.Post("/newscan", handler.NewScan(cloner, repo))
	api.Get("/scan/:scanid", handler.GetScan(repo))

	app.Listen(":8080")
}
