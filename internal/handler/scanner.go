package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/dockerengine"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/gitcloner"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ScanHandler struct {
	cloner       gitcloner.Cloner
	repo         repository.Repo
	dockerengine dockerengine.Engine
}

func NewScanHandler(cloner gitcloner.Cloner, repo repository.Repo, dockerengine dockerengine.Engine) *ScanHandler {
	return &ScanHandler{
		cloner:       cloner,
		repo:         repo,
		dockerengine: dockerengine,
	}
}

func (s *ScanHandler) NewScan() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logrus.Info("NewScan request received")
		var url models.NewScanRequest

		if err := c.BodyParser(&url); err != nil {
			logrus.Error("Error while parsing body", err)
			return c.Status(http.StatusBadRequest).JSON(models.Response{
				Status:  fiber.StatusBadRequest,
				Message: "Error while parsing body",
				Data: &fiber.Map{
					"Error Message": err.Error(),
				},
			})
		}

		projectId, err := s.cloner.CloneRepo(url.Url)
		if err != nil {
			logrus.Error("Error while cloning repo", err)
			return c.Status(http.StatusInternalServerError).JSON(models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error while cloning repo",
				Data: &fiber.Map{
					"Error Message": err.Error(),
				},
			})
		}

		err, scanData := s.dockerengine.NewScanResults(projectId)
		if err != nil {
			logrus.Error("Error while analyzing scan results", err)
			return c.Status(http.StatusInternalServerError).JSON(models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error while analyzing scan results",
				Data: &fiber.Map{
					"Error Message": err.Error(),
				},
			})
		}

		err = s.repo.SetScanData(scanData.ScanID, scanData.ScanData)
		if err != nil {
			logrus.Error("Error while inserting scan result to database", err)
			return c.Status(http.StatusInternalServerError).JSON(models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error while inserting scan result to database",
				Data: &fiber.Map{
					"Error Message": err.Error(),
				},
			})
		}
		response := checkSeverity(scanData)
		logrus.Info("NewScan request ended")
		return c.Status(fiber.StatusOK).JSON(models.Response{
			Status:  fiber.StatusOK,
			Message: "New Scan request ended successfully",
			Data: &fiber.Map{
				"Scan Results": response,
			},
		})
	}
}

func (s *ScanHandler) GetScan() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		scanId := c.AllParams()["scanid"]
		logrus.Infof("GetScan request received for scan_id: %v", scanId)
		scanData, err := s.repo.GetScanData(scanId)
		if err != nil {
			logrus.Error("Error getting scan data from database", err)
			return c.Status(http.StatusBadRequest).JSON(models.Response{
				Status:  fiber.StatusBadRequest,
				Message: "Error getting scan data from database",
				Data: &fiber.Map{
					"Error Message": err.Error(),
				},
			})
		}
		logrus.Infof("GetScan request ended for scan_id: %v", scanId)
		return c.Status(http.StatusOK).JSON(models.Response{
			Status:  fiber.StatusOK,
			Message: "Get scan result from database successfully",
			Data: &fiber.Map{
				"Scan Results": scanData,
			},
		})
	}
}

func checkSeverity(scanData models.ScanData) models.ScanResponse {
	var severityCounter int
	var SecurityData []models.SecurityData
	for key, element := range scanData.ScanData.Metrics {
		if element.SEVERITYHIGH >= 1 {
			SecurityData = append(SecurityData, models.SecurityData{
				IssuePath:         key,
				HighSeverityCount: element.SEVERITYHIGH,
			})
			severityCounter += 1
		}
	}

	var scanResponse = models.ScanResponse{
		ScanId:         scanData.ScanID,
		SecurityIssues: SecurityData,
		SecurityStatus: securityStatusCheck(severityCounter),
	}
	return scanResponse
}

func securityStatusCheck(severityCounter int) string {
	if severityCounter >= 1 {
		return "Unsecure"
	}
	return "Secure"
}
