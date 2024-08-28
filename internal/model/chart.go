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
package model

import (
	"errors"
	"regexp"
	"strings"
)

type Chart struct {
	Name    string
	Version string
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
}

type ChartManifest struct {
	Layers []Layer `json:"layers"`
}

func ChartFromString(fullChart string) (*Chart, error) {
	if len(strings.TrimSpace(fullChart)) < 1 {
		return nil, errors.New("provided chart cannot be empty string")
	}

	name, err := extractName(fullChart)
	if err != nil {
		return nil, err
	}

	version, err := extractVersion(fullChart)
	if err != nil {
		return nil, err
	}

	return &Chart{Name: name, Version: version}, nil
}

func extractName(chart string) (string, error) {
	expr := regexp.MustCompile("^(.*?)(-\\d+(\\.\\d+)*.+\\.tgz)$")
	matches := expr.FindStringSubmatch(chart)

	if len(matches) < 1 {
		return "", failedToParseError(chart)
	}

	return matches[1], nil
}

func extractVersion(chart string) (string, error) {
	expr := regexp.MustCompile("(\\d+(\\.\\d+)*.*)\\.tgz$")
	matches := expr.FindStringSubmatch(chart)

	if len(matches) < 1 {
		return "", failedToParseError(chart)
	}
	return strings.ReplaceAll(matches[1], "+", "_"), nil
}

func failedToParseError(chart string) error {
	return errors.New("failed to parse name or version from the provided chart: " + chart)
}
