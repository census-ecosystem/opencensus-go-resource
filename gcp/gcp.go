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

package gcp

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"go.opencensus.io/resource"
)

const (
	KeyProjectID   = "cloud.google.com/project_id"
	KeyZone        = "cloud.google.com/zone"
	KeyInstanceID  = "cloud.google.com/gce/instance_id"
	KeyClusterName = "cloud.google.com/gce/attributes/cluster_name"
)

func DetectGCEInstance(context.Context) (*resource.Resource, error) {
	if !metadata.OnGCE() {
		return nil, nil
	}
	res := &resource.Resource{
		Type: "cloud.google.com/gce/instance",
		Tags: map[string]string{},
	}
	instanceID, err := metadata.InstanceID()
	logError(err)
	if instanceID != "" {
		res.Tags[KeyInstanceID] = instanceID
	}

	projectID, err := metadata.ProjectID()
	logError(err)
	if projectID != "" {
		res.Tags[KeyProjectID] = projectID
	}

	zone, err := metadata.Zone()
	logError(err)
	if zone != "" {
		res.Tags[KeyZone] = zone
	}

	clusterName, err := metadata.InstanceAttributeValue("cluster-name")
	logError(err)
	if clusterName != "" {
		res.Tags[KeyClusterName] = strings.TrimSpace(clusterName)
	}

	return res, nil
}

// logError logs error only if the error is present and it is not 'not defined'
func logError(err error) {
	if err != nil {
		if !strings.Contains(err.Error(), "not defined") {
			log.Printf("Error retrieving gcp metadata: %v", err)
		}
	}
}
