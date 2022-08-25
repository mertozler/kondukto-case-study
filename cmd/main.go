package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/handler"
	"github.com/mertozler/pkg/git-clonner"
)

func main() {
	app := fiber.New()

	cloner := git_clonner.NewGitClonner()

	api := app.Group("/api/v1")
	api.Post("/newscan", handler.NewScan(cloner))

	app.Listen(":8080")
}
