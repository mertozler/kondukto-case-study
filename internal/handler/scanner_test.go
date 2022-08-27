package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewScanHandler_Should_Success(t *testing.T) {
	// given
	requestUrl := "http://localhost:8080/api/v1/newscan"
	scanRequest := models.NewScanRequest{Url: "https://github.com/anxolerd/dvpwa"}
	mockCloner := mocks.NewCloner(t)
	mockCloner.On("CloneRepo", mock.Anything).Return("uuidexample", nil)
	mockRepo := mocks.NewRepo(t)
	mockRepo.On("SetScanData", mock.Anything, mock.Anything).Return(nil)
	mockEngine := mocks.NewEngine(t)
	mockEngine.On("NewScanResults", mock.Anything).Return(nil, models.ScanData{})
	handler := NewScanHandler(mockCloner, mockRepo, mockEngine)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/newscan", handler.NewScan())
	scanRequestJson, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestUrl, bytes.NewBuffer(scanRequestJson))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestNewScanHandler_Should_Return_Error_When_Clonning_Repo(t *testing.T) {
	// given
	requestUrl := "http://localhost:8080/api/v1/newscan"
	scanRequest := models.NewScanRequest{Url: "https://github.com/anxolerd/dvpwa"}
	mockCloner := mocks.NewCloner(t)
	mockCloner.On("CloneRepo", mock.Anything).Return("", errors.New("Error while cloning repository"))
	mockRepo := mocks.NewRepo(t)
	mockEngine := mocks.NewEngine(t)
	handler := NewScanHandler(mockCloner, mockRepo, mockEngine)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/newscan", handler.NewScan())
	scanRequestJson, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestUrl, bytes.NewBuffer(scanRequestJson))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestNewScanHandler_Should_Return_Error_When_Scanning_Results(t *testing.T) {
	// given
	requestUrl := "http://localhost:8080/api/v1/newscan"
	scanRequest := models.NewScanRequest{Url: "https://github.com/anxolerd/dvpwa"}
	mockCloner := mocks.NewCloner(t)
	mockCloner.On("CloneRepo", mock.Anything).Return("uuidexample", nil)
	mockRepo := mocks.NewRepo(t)
	mockEngine := mocks.NewEngine(t)
	mockEngine.On("NewScanResults", mock.Anything).Return(errors.New("Error while scanning results"), models.ScanData{})
	handler := NewScanHandler(mockCloner, mockRepo, mockEngine)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/newscan", handler.NewScan())
	scanRequestJson, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestUrl, bytes.NewBuffer(scanRequestJson))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestNewScanHandler_Should_Return_Error_When_Inserting_Database(t *testing.T) {
	// given
	requestUrl := "http://localhost:8080/api/v1/newscan"
	scanRequest := models.NewScanRequest{Url: "https://github.com/anxolerd/dvpwa"}
	mockCloner := mocks.NewCloner(t)
	mockCloner.On("CloneRepo", mock.Anything).Return("uuidexample", nil)
	mockRepo := mocks.NewRepo(t)
	mockRepo.On("SetScanData", mock.Anything, mock.Anything).Return(errors.New("Error while interting db"))
	mockEngine := mocks.NewEngine(t)
	mockEngine.On("NewScanResults", mock.Anything).Return(nil, models.ScanData{})
	handler := NewScanHandler(mockCloner, mockRepo, mockEngine)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/newscan", handler.NewScan())
	scanRequestJson, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestUrl, bytes.NewBuffer(scanRequestJson))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}
