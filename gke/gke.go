// Copyright 2019, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gke // import "contrib.go.opencensus.io/resource/gke"

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"go.opencensus.io/resource"
	"go.opencensus.io/resource/resourcekeys"
	"os"
)

// Detect detects associated resources when running in GKE environment.
func Detect(ctx context.Context) (*resource.Resource, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return nil, nil
	}

	k8s := func(ctx context.Context) (*resource.Resource, error) {
		k8sRes := &resource.Resource{
			Type:   resourcekeys.K8SType,
			Labels: map[string]string{},
		}

		clusterName, err := metadata.InstanceAttributeValue("instance/attributes/cluster-name")
		logError(err)
		if clusterName != "" {
			k8sRes.Labels[resourcekeys.K8SKeyClusterName] = clusterName
		}

		k8sRes.Labels[resourcekeys.K8SKeyNamespaceName] = os.Getenv("NAMESPACE")
		k8sRes.Labels[resourcekeys.K8SKeyPodName] = os.Getenv("HOSTNAME")
		return k8sRes, nil
	}

	container := func(ctx context.Context) (*resource.Resource, error) {
		containerRes := &resource.Resource{
			Type:   resourcekeys.ContainerType,
			Labels: map[string]string{},
		}
		containerRes.Labels[resourcekeys.ContainerKeyName] = os.Getenv("CONTAINER_NAME")
		return containerRes, nil
	}

	host := func(ctx context.Context) (*resource.Resource, error) {
		hostRes := &resource.Resource{
			Type:   resourcekeys.HostType,
			Labels: map[string]string{},
		}
		instanceID, err := metadata.InstanceID()
		logError(err)
		if instanceID != "" {
			hostRes.Labels[resourcekeys.HostKeyID] = instanceID
		}
		return hostRes, nil
	}

	cloud := func(ctx context.Context) (*resource.Resource, error) {
		cloudRes := &resource.Resource{
			Type:   resourcekeys.CloudType,
			Labels: map[string]string{},
		}

		cloudRes.Labels[resourcekeys.CloudKeyProvider] = resourcekeys.CloudProviderGCP
		projectID, err := metadata.ProjectID()
		logError(err)
		if projectID != "" {
			cloudRes.Labels[resourcekeys.CloudKeyAccountID] = projectID
		}

		zone, err := metadata.Zone()
		logError(err)
		if zone != "" {
			cloudRes.Labels[resourcekeys.CloudKeyZone] = zone
		}

		cloudRes.Labels[resourcekeys.CloudKeyRegion] = ""
		return cloudRes, nil
	}

	return resource.MultiDetector(k8s, container, host, cloud)(ctx)
}

// logError logs error only if the error is present and it is not 'not defined'
func logError(err error) {
	if err != nil {
		if !strings.Contains(err.Error(), "not defined") {
			log.Printf("Error retrieving gcp metadata: %v", err)
		}
	}
}
