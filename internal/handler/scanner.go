package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/dockerengine"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/git-clonner"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ScanHandler struct {
	cloner       *git_clonner.GitClonner
	repo         *repository.Repository
	dockerengine *dockerengine.Docker
}

func NewScanHandler(cloner *git_clonner.GitClonner, repo *repository.Repository, dockerengine *dockerengine.Docker) *ScanHandler {
	return &ScanHandler{
		cloner:       cloner,
		repo:         repo,
		dockerengine: dockerengine,
	}
}

func (s *ScanHandler) NewScan() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logrus.Info("NewScan request receiver")
		var url models.NewScanRequest

		if err := c.BodyParser(&url); err != nil {
			logrus.Errorf("Error while parsing body", err)
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		projectId, err := s.cloner.CloneRepo(url.Url)
		if err != nil {
			logrus.Errorf("Error while cloning repo", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		err, scanData := s.dockerengine.NewScanResults(projectId)
		if err != nil {
			logrus.Errorf("Error while analyzing scan results", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		err = s.repo.SetScanData(scanData.ScanID, scanData.ScanData)
		if err != nil {
			logrus.Errorf("Error while inserting scan result to database", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		response := checkSeverity(scanData)
		logrus.Info("NewScan request ended")
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func (s *ScanHandler) GetScan() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		scanId := c.AllParams()["scanid"]
		logrus.Info("GetScan request received for scan_id: %v", scanId)
		scanData, err := s.repo.GetScanData(scanId)
		if err != nil {
			logrus.Errorf("Error getting scan data from database", err)
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		logrus.Info("GetScan request ended for scan_id: %v", scanId)
		return c.Status(http.StatusOK).JSON(scanData)
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

	var response = models.ScanResponse{
		ScanId:         scanData.ScanID,
		SecurityIssues: SecurityData,
		SecurityStatus: securityStatusCheck(severityCounter),
	}
	return response
}

func securityStatusCheck(severityCounter int) string {
	if severityCounter >= 1 {
		return "Unsecure"
	}
	return "Secure"
}
