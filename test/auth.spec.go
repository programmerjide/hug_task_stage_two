package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"hng_stage_two_task/internal/api"
	"hng_stage_two_task/internal/dto"
)

func TestRegisterUser(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	reqBody := dto.UserSignupRequestDto{
		UserLoginDto: dto.UserLoginDto{
			Email:    "john.doe@example.com",
			Password: "password123",
		},
		Phone:     "1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"].(map[string]interface{})["accessToken"])
}

func TestLoginUser(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	reqBody := dto.UserLoginDto{
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"].(map[string]interface{})["accessToken"])
}

func TestRegisterUserValidationErrors(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	// Test missing fields
	reqBody := dto.UserSignupRequestDto{
		UserLoginDto: dto.UserLoginDto{
			Email:    "",
			Password: "",
		},
		Phone:     "",
		FirstName: "",
		LastName:  "",
	}

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "Bad request", response["status"])
	assert.NotNil(t, response["errors"])
}

func TestRegisterDuplicateUser(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	// Register the first user
	reqBody := dto.UserSignupRequestDto{
		UserLoginDto: dto.UserLoginDto{
			Email:    "john.doe@example.com",
			Password: "password123",
		},
		Phone:     "1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Try to register the same user again
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "Bad request", response["status"])
	assert.NotNil(t, response["errors"])
}
