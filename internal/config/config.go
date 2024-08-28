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
	"github.com/spf13/viper"
	"log"
)

var Env *envConfig

type envConfig struct {
	Username                    string `mapstructure:"USERNAME"`
	Password                    string `mapstructure:"PASSWORD"`
	ServerPort                  string `mapstructure:"SERVER_PORT"`
	Repository                  string `mapstructure:"REPOSITORY"`
	ContainerRegistryHost       string `mapstructure:"CONTAINER_REGISTRY_HOST"`
	ContainerRegistryRepository string `mapstructure:"CONTAINER_REGISTRY_REPOSITORY"`
	ContainerRegistryUsername   string `mapstructure:"CONTAINER_REGISTRY_USERNAME"`
	ContainerRegistryPassword   string `mapstructure:"CONTAINER_REGISTRY_PASSWORD"`
	AllowAnonymousGetRequests   bool   `mapstructure:"ALLOW_ANONYMOUS_GET_REQUESTS"`
}

func Init() {
	Env = load()
}

func load() *envConfig {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(".env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Println("Failed to read config from .env file. ", err)
	}

	return &envConfig{
		Username:                    v.GetString("USERNAME"),
		Password:                    v.GetString("PASSWORD"),
		ServerPort:                  v.GetString("SERVER_PORT"),
		Repository:                  v.GetString("REPOSITORY"),
		ContainerRegistryHost:       v.GetString("CONTAINER_REGISTRY_HOST"),
		ContainerRegistryRepository: v.GetString("CONTAINER_REGISTRY_REPOSITORY"),
		ContainerRegistryUsername:   v.GetString("CONTAINER_REGISTRY_USERNAME"),
		ContainerRegistryPassword:   v.GetString("CONTAINER_REGISTRY_PASSWORD"),
		AllowAnonymousGetRequests:   v.GetBool("ALLOW_ANONYMOUS_GET_REQUESTS"),
	}

}
