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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChartFromStringPositiveCases(t *testing.T) {
	testCases := []struct {
		name      string
		fullChart string
		expected  *Chart
	}{
		{
			name:      "ValidChart",
			fullChart: "test-chart-1.2.3.tgz",
			expected:  &Chart{Name: "test-chart", Version: "1.2.3"},
		},
		{
			name:      "ValidChartSnapshotVersion",
			fullChart: "test-chart-1.2.3-SNAPSHOT.tgz",
			expected:  &Chart{Name: "test-chart", Version: "1.2.3-SNAPSHOT"},
		},
		{
			name:      "ValidChartPrereleaseVersion",
			fullChart: "pre-release-chart-1.0.0-0.3.7.tgz",
			expected:  &Chart{Name: "pre-release-chart", Version: "1.0.0-0.3.7"},
		},
		{
			name:      "ValidChartPrereleaseVersionMetadata",
			fullChart: "pre-release-chart-1.0.0-alpha+001.tgz",
			expected:  &Chart{Name: "pre-release-chart", Version: "1.0.0-alpha_001"},
		},
		{
			name:      "ValidChartOrderedVersion",
			fullChart: "ordered-1.0.0-beta.11.tgz",
			expected:  &Chart{Name: "ordered", Version: "1.0.0-beta.11"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ChartFromString(tc.fullChart)
			assert.Equal(t, tc.expected, actual)
			assert.Nil(t, err)
		})
	}
}

func TestGetChartFromStringNegativeCases(t *testing.T) {
	testCases := []struct {
		name      string
		fullChart string
		expected  string
	}{
		{
			name:      "EmptyInput",
			fullChart: "",
			expected:  "provided chart cannot be empty string",
		},
		{
			name:      "NoVersion",
			fullChart: "test-chart",
			expected:  "failed to parse name or version from the provided chart: test-chart",
		},
		{
			name:      "NoName",
			fullChart: "1.2.3-SNAPSHOT.tgz",
			expected:  "failed to parse name or version from the provided chart: 1.2.3-SNAPSHOT.tgz",
		},
		{
			name:      "NoExtension",
			fullChart: "test-chart-1.2.3-SNAPSHOT",
			expected:  "failed to parse name or version from the provided chart: test-chart-1.2.3-SNAPSHOT",
		},
		{
			name:      "WhiteSpaceInput",
			fullChart: "    ",
			expected:  "provided chart cannot be empty string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ChartFromString(tc.fullChart)
			assert.Nil(t, actual)
			assert.Equal(t, tc.expected, err.Error())
		})
	}
}
