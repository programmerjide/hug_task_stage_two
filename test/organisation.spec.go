package test

import (
	"bytes"
	"encoding/json"
	"hng_stage_two_task/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"hng_stage_two_task/internal/api"
	"hng_stage_two_task/internal/dto"
)

func TestCreateOrganisation(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	reqBody := dto.CreateOrganisationRequest{
		Name:        "New Organisation",
		Description: "This is a new organisation.",
	}

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organisations", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_jwt_token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

func TestGetUserOrganisations(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	req := httptest.NewRequest("GET", "/api/organisations", nil)
	req.Header.Set("Authorization", "Bearer valid_jwt_token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

func TestGetOrganisationByID(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	req := httptest.NewRequest("GET", "/api/organisations/org123", nil)
	req.Header.Set("Authorization", "Bearer valid_jwt_token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

// setupTestEnvironment initializes the test environment and returns AppConfig for testing
func setupTestEnvironment(t *testing.T) config.AppConfig {
	cfg, err := config.SetupEnv()
	if err != nil {
		t.Fatalf("config file is not loaded properly %v\n", err)
	}
	return cfg
}
