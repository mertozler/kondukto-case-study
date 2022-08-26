package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/dockerengine"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/git-clonner"
	"log"
	"net/http"
)

func NewScan(cloner *git_clonner.GitClonner, repo *repository.Repository) func(c *fiber.Ctx) error {
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

		err, scanData := dockerengine.NewScanResults(projectId)
		if err != nil {
			log.Println("Error while analyzing results ")
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		err = repo.SetScanData(scanData.ScanID, scanData.ScanData)
		if err != nil {
			log.Println("Error inserting scan data to database ")
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		return c.Status(fiber.StatusOK).SendString(projectId)
	}
}

func GetScan(repo *repository.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		scanId := c.AllParams()["scanid"]
		scanData, err := repo.GetScanData(scanId)
		if err != nil {
			log.Println("Error while getting scan data")
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		return c.Status(http.StatusOK).JSON(scanData)
		return nil
	}
}
