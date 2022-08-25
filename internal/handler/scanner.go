package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/dockerengine"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/pkg/git-clonner"
	"log"
	"net/http"
)

func NewScan(cloner *git_clonner.GitClonner) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Println("NewScan request received")
		var url models.NewScanRequest

		if err := c.BodyParser(&url); err != nil {
			log.Println("Error while parsing URL request")
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		projectId, err := cloner.CloneRepo(url.Url)
		if err != nil {
			log.Println("Error while cloning repo ")
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		dockerengine.NewClient(projectId)

		return c.Status(fiber.StatusOK).SendString(projectId)
	}
}
