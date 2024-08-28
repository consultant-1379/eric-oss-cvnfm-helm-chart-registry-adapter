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
package auth

import (
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	cases := []struct {
		name                      string
		method                    string
		user                      string
		password                  string
		allowAnonymousGetRequests bool
		expectedStatus            int
	}{
		{"ValidCredentialsGet", http.MethodGet, "user", "pass", false, http.StatusOK},
		{"ValidCredentialsPost", http.MethodPost, "user", "pass", false, http.StatusOK},
		{"InvalidCredentialsGet", http.MethodGet, "wrongUser", "wrongPass", false, http.StatusUnauthorized},
		{"InvalidCredentialsPost", http.MethodPost, "wrongUser", "wrongPass", false, http.StatusUnauthorized},
		{"InvalidCredentialsAllowAnonymousGet", http.MethodGet, "wrongUser", "wrongPass", true, http.StatusOK},
		{"ValidCredentialsAllowAnonymousPost", http.MethodPost, "user", "pass", true, http.StatusOK},
		{"InvalidCredentialsAllowAnonymousPost", http.MethodPost, "wrongUser", "wrongPass", true, http.StatusUnauthorized},
	}

	r := gin.Default()
	r.Use(BasicAuth())
	r.POST("/", func(c *gin.Context) {})
	r.GET("/", func(c *gin.Context) {})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initEnvVars(t, strconv.FormatBool(tc.allowAnonymousGetRequests))
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tc.method, "/", nil)
			req.SetBasicAuth(tc.user, tc.password)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, tc.name)
		})
	}
}

func initEnvVars(t *testing.T, allowAnonymousGetRequests string) {
	t.Setenv("USERNAME", "user")
	t.Setenv("PASSWORD", "pass")
	t.Setenv("ALLOW_ANONYMOUS_GET_REQUESTS", allowAnonymousGetRequests)
	config.Init()
}
