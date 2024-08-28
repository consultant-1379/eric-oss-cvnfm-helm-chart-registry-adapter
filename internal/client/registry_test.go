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
package client

import (
	"encoding/json"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/config"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/model"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateChart(t *testing.T) {
	testCases := []struct {
		name          string
		chart         *model.Chart
		expectedError error
	}{
		{
			name:          "ValidChart",
			chart:         &model.Chart{Name: "name", Version: "version"},
			expectedError: nil,
		},
		{
			name:          "ChartIsNil",
			chart:         nil,
			expectedError: errors.New("chart cannot be nil"),
		},
		{
			name:          "ChartNameIsEmpty",
			chart:         &model.Chart{Name: "", Version: "1.0.0"},
			expectedError: errors.New("chart should contain name and version"),
		},
		{
			name:          "ChartVerionIsEmpty",
			chart:         &model.Chart{Name: "chart", Version: ""},
			expectedError: errors.New("chart should contain name and version"),
		},
		{
			name:          "ChartNameAndVersionAreEmpty",
			chart:         &model.Chart{Name: "", Version: ""},
			expectedError: errors.New("chart should contain name and version"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateChart(tc.chart)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestBuildChartManifestUrl(t *testing.T) {
	initEnvVars(t, "http://helm.registry.se")

	chart := model.Chart{Name: "test-chart", Version: "1.2.1-SNAPSHOT"}
	expectedUrl := "http://helm.registry.se/charts/onboarded/test-chart/manifests/1.2.1-SNAPSHOT"

	assert.Equal(t, expectedUrl, buildChartManifestUrl(&chart))
}

func TestBuildChartBlobUrl(t *testing.T) {
	initEnvVars(t, "http://helm.registry.se")

	chart := model.Chart{Name: "test-chart", Version: "1.2.1-SNAPSHOT"}
	expectedUrl := "http://helm.registry.se/charts/onboarded/test-chart/blobs/sha256:610b2bacfeed"

	assert.Equal(t, expectedUrl, buildChartBlobUrl(&chart, "sha256:610b2bacfeed"))
}

func TestGetChartSuccessCase(t *testing.T) {
	registry := NewRegistry()
	expectedManifest := model.ChartManifest{Layers: []model.Layer{
		{MediaType: "invalid", Digest: "sha256:invalid"},
		{MediaType: layerMediaType, Digest: "sha256:610b2bacfeed"},
	}}
	expectedBlobResponse := "raw data"
	chart := model.Chart{Name: "test-chart", Version: "1.2.1-SNAPSHOT"}

	//starts a server under mockServer.URL address running on localhost
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "test-chart/manifests/1.2.1-SNAPSHOT") {
			if r.Header.Get("Authorization") != "Basic dXNlcjpwYXNz" {
				w.WriteHeader(http.StatusUnauthorized)
			}
			resp, _ := json.Marshal(expectedManifest)
			_, _ = fmt.Fprintln(w, string(resp))
		} else if strings.Contains(r.URL.Path, "test-chart/blobs/sha256:610b2bacfeed") {
			w.Header().Add("Authorization", r.Header.Get("Authorization"))
			_, _ = fmt.Fprintln(w, expectedBlobResponse)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()
	initEnvVars(t, mockServer.URL)

	resp, err := registry.GetChartManifestBlob(&chart)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBodyByteArray, _ := io.ReadAll(resp.Body)
	assert.Equal(t, expectedBlobResponse, strings.TrimSpace(string(respBodyByteArray)))
	assert.Equal(t, "Basic dXNlcjpwYXNz", resp.Header.Get("Authorization"))
}

func TestGetChartFailResponse(t *testing.T) {
	registry := NewRegistry()
	chart := model.Chart{Name: "test-chart", Version: "1.2.1-SNAPSHOT"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()
	initEnvVars(t, mockServer.URL)

	resp, err := registry.GetChartManifestBlob(&chart)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func initEnvVars(t *testing.T, serverUrl string) {
	t.Setenv("CONTAINER_REGISTRY_HOST", serverUrl)
	t.Setenv("CONTAINER_REGISTRY_REPOSITORY", "charts/onboarded")
	t.Setenv("CONTAINER_REGISTRY_USERNAME", "user")
	t.Setenv("CONTAINER_REGISTRY_PASSWORD", "pass")
	config.Init()
	assert.NotNil(t, config.Env.ContainerRegistryHost)
}
