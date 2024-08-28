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
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/model"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func handleGetChartsManifestBlob(registry Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		chartParam := c.Param(chartPathVar)

		chart, err := model.ChartFromString(chartParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		registryResponse, err := registry.GetChartManifestBlob(chart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}
		defer registryResponse.Body.Close()

		c.Status(registryResponse.StatusCode)
		_, err = io.Copy(c.Writer, registryResponse.Body)
		if err != nil {
			log.Printf("Failed to transfer response from client: %v\n", err)
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}
	}
}
