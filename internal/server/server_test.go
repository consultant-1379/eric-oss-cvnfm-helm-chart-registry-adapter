/*
 * COPYRIGHT Ericsson 2024
 *
 *
 *
 * The copyright to the computer program(s) herein is the property of
 *
 * Ericsson Inc. The programs may be used and/or copied only with written
 *
 * permission from Ericsson Inc. or in accordance with the terms and
 *
 * conditions stipulated in the agreement/contract under which the
 *
 * program(s) have been supplied.
 */
package server

import (
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/auth"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/config"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockRegistry struct {
}

type mockErrorRegistry struct {
}

func (mockRegistry mockRegistry) GetChartManifestBlob(chart *model.Chart) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(chart.Name + "-" + chart.Version)),
	}, nil
}

func (mockRegistry mockErrorRegistry) GetChartManifestBlob(chart *model.Chart) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusInternalServerError}, nil
}

func setupRouter(registry Registry) *gin.Engine {
	server := New()
	server.registry = registry

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(auth.BasicAuth())
	createRoutes(router, server)
	return router
}

func TestServerSetup(t *testing.T) {
	initEnvVars(t)
	r := setupRouter(&mockRegistry{})
	w := httptest.NewRecorder()

	// Request to non-existing route to check if server is properly set up
	req, _ := http.NewRequest(http.MethodGet, "/nonexistent", nil)
	req.SetBasicAuth(config.Env.Username, config.Env.Password)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetChartPositive(t *testing.T) {
	initEnvVars(t)

	r := setupRouter(&mockRegistry{})
	w := httptest.NewRecorder()

	// Make a GET request to /charts/:chartPathVar
	req, _ := http.NewRequest(http.MethodGet, "/charts/chart-1.2.1.tgz", nil)
	req.SetBasicAuth(config.Env.Username, config.Env.Password)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "chart-1.2.1", w.Body.String()) //mocked response, real behavior is different
}

func TestGetChartInvalidChart(t *testing.T) {
	initEnvVars(t)

	r := setupRouter(&mockRegistry{})
	w := httptest.NewRecorder()

	// Make a GET request to /charts/:chartPathVar
	req, _ := http.NewRequest(http.MethodGet, "/charts/chart", nil)
	req.SetBasicAuth(config.Env.Username, config.Env.Password)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetChartRegistryError(t *testing.T) {
	initEnvVars(t)
	r := setupRouter(&mockErrorRegistry{})
	w := httptest.NewRecorder()

	// Make a GET request to /charts/:chartPathVar
	req, _ := http.NewRequest(http.MethodGet, "/charts/chart-1.2.1.tgz", nil)
	req.SetBasicAuth(config.Env.Username, config.Env.Password)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func initEnvVars(t *testing.T) {
	t.Setenv("USERNAME", "user")
	t.Setenv("PASSWORD", "pass")
	t.Setenv("REGISTRY", "/onboarding")
	config.Init()
}
