// Copyright 2018, OpenCensus Authors
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

package gcp // import "contrib.go.opencensus.io/resource/gcp"

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"go.opencensus.io/resource"
	"go.opencensus.io/resource/resourcekeys"
)

// Detect detects associated resources when running in GCP environment.
func Detect(ctx context.Context) (*resource.Resource, error) {
	if !metadata.OnGCE() {
		return nil, nil
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

		name, err := metadata.InstanceName()
		logError(err)
		if instanceID != "" {
			hostRes.Labels[resourcekeys.HostKeyName] = name
		}

		hostname, err := metadata.Hostname()
		logError(err)
		if instanceID != "" {
			hostRes.Labels[resourcekeys.HostKeyHostName] = hostname
		}

		hostType, err := metadata.InstanceAttributeValue("instance/machine-type")
		logError(err)
		if instanceID != "" {
			hostRes.Labels[resourcekeys.HostKeyType] = hostType
		}

		return hostRes, nil
	}

	return resource.MultiDetector(cloud, host)(ctx)
}

// logError logs error only if the error is present and it is not 'not defined'
func logError(err error) {
	if err != nil {
		if !strings.Contains(err.Error(), "not defined") {
			log.Printf("Error retrieving gcp metadata: %v", err)
		}
	}
}
