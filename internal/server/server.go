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
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/client"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/config"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const chartsPath = "/charts"
const chartPathVar = "chart"

type Registry interface {
	GetChartManifestBlob(chart *model.Chart) (*http.Response, error)
}

type Server struct {
	registry Registry
}

func New() *Server {
	return &Server{
		registry: client.NewRegistry(),
	}
}

func (server *Server) Run() {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(auth.BasicAuth())
	createRoutes(router, server)

	err := router.Run(":" + config.Env.ServerPort)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func createRoutes(router *gin.Engine, server *Server) {
	charts := router.Group(config.Env.Repository + chartsPath)
	charts.GET("/:"+chartPathVar, handleGetChartsManifestBlob(server.registry))
}
