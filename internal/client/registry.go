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
	"crypto/tls"
	"encoding/json"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/config"
	"eric-oss-cvnfm-helm-chart-registry-adapter/internal/model"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const layerMediaType = "application/vnd.cncf.helm.chart.content.v1.tar+gzip"
const ociMediaType = "application/vnd.oci.image.manifest.v1+json"

type Registry struct {
	client *http.Client
}

func NewRegistry() *Registry {
	//currently client has TLS disabled on premise. Should be changed further
	return &Registry{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 10,
		},
	}
}

func (registry Registry) GetChartManifestBlob(chart *model.Chart) (*http.Response, error) {
	if err := validateChart(chart); err != nil {
		return nil, err
	}

	//get chart's manifest
	resp, err := doGetRequest(registry.client, buildChartManifestUrl(chart))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, nil
	}

	//fetch digest from the response
	digest, err := parseDigest(resp)
	if err != nil {
		return nil, err
	}

	//return chart blob as data stream
	return doGetRequest(registry.client, buildChartBlobUrl(chart, digest))
}

func parseDigest(response *http.Response) (string, error) {
	manifest := model.ChartManifest{}
	manifestJson, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(manifestJson, &manifest)
	if err != nil {
		return "", err
	}

	for _, layer := range manifest.Layers {
		if layer.MediaType == layerMediaType {
			return layer.Digest, nil
		}
	}

	return "", errors.New("failed to parse digest from the manifest: " + string(manifestJson))
}

func validateChart(chart *model.Chart) error {
	if chart == nil {
		return errors.New("chart cannot be nil")
	}

	if len(chart.Name) == 0 || len(chart.Version) == 0 {
		return errors.New("chart should contain name and version")
	}
	return nil
}

func doGetRequest(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.Env.ContainerRegistryUsername, config.Env.ContainerRegistryPassword)
	req.Header.Add("Accept", ociMediaType)

	return client.Do(req)
}

func buildChartManifestUrl(chart *model.Chart) string {
	return fmt.Sprintf("%s/v2/%s/%s/manifests/%s",
		config.Env.ContainerRegistryHost,
		config.Env.ContainerRegistryRepository,
		chart.Name,
		chart.Version)
}

func buildChartBlobUrl(chart *model.Chart, digest string) string {
	return fmt.Sprintf("%s/v2/%s/%s/blobs/%s",
		config.Env.ContainerRegistryHost,
		config.Env.ContainerRegistryRepository,
		chart.Name,
		digest)
}
