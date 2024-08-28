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
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && config.Env.AllowAnonymousGetRequests {
			c.Next()
		} else {
			user, password, hasAuth := c.Request.BasicAuth()
			if hasAuth && user == config.Env.Username && password == config.Env.Password {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
				c.Abort()
			}
		}
	}
}
