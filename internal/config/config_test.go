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
package config

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestConfigEnvAssignment(t *testing.T) {
	expected := envConfig{
		Username:                    "user",
		Password:                    "pass",
		ServerPort:                  "8080",
		Repository:                  "onboarded/charts",
		ContainerRegistryHost:       "example.com",
		ContainerRegistryRepository: "/onboarded",
		ContainerRegistryUsername:   "ruser",
		ContainerRegistryPassword:   "rpass",
		AllowAnonymousGetRequests:   false,
	}

	t.Setenv("USERNAME", expected.Username)
	t.Setenv("PASSWORD", expected.Password)
	t.Setenv("SERVER_PORT", expected.ServerPort)
	t.Setenv("REPOSITORY", expected.Repository)
	t.Setenv("CONTAINER_REGISTRY_HOST", expected.ContainerRegistryHost)
	t.Setenv("CONTAINER_REGISTRY_REPOSITORY", expected.ContainerRegistryRepository)
	t.Setenv("CONTAINER_REGISTRY_USERNAME", expected.ContainerRegistryUsername)
	t.Setenv("CONTAINER_REGISTRY_PASSWORD", expected.ContainerRegistryPassword)
	t.Setenv("ALLOW_ANONYMOUS_GET_REQUESTS", strconv.FormatBool(expected.AllowAnonymousGetRequests))

	Init()

	assert.Equal(t, &expected, Env)
}
