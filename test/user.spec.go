package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"hng_stage_two_task/internal/api"
)

func TestGetUserByID(t *testing.T) {
	cfg := setupTestEnvironment(t)
	app := api.StartServer(cfg)
	defer app.Shutdown()

	req := httptest.NewRequest("GET", "/api/users/user123", nil)
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
