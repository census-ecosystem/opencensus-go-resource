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

package aws // import "contrib.go.opencensus.io/resource/aws"

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.opencensus.io/resource"
	"go.opencensus.io/resource/resourcekeys"
)

// Detect detects associated resources when running in AWS environment.
func Detect(ctx context.Context) (*resource.Resource, error) {
	c := ec2metadata.New(session.New())
	if !c.Available() {
		return nil, nil
	}
	doc, err := c.GetInstanceIdentityDocument()
	if err != nil {
		return nil, err
	}

	cloud := func(ctx context.Context) (*resource.Resource, error) {
		cloudRes := &resource.Resource{
			Type:   resourcekeys.CloudType,
			Labels: map[string]string{},
		}
		cloudRes.Labels[resourcekeys.CloudKeyProvider] = resourcekeys.CloudProviderAWS
		cloudRes.Labels[resourcekeys.CloudKeyRegion] = doc.Region
		cloudRes.Labels[resourcekeys.CloudKeyAccountID] = doc.AccountID
		return cloudRes, nil
	}

	host := func(ctx context.Context) (*resource.Resource, error) {
		hostRes := &resource.Resource{
			Type:   resourcekeys.HostType,
			Labels: map[string]string{},
		}
		hostRes.Labels[resourcekeys.HostKeyID] = doc.InstanceID
		return hostRes, nil
	}

	return resource.MultiDetector(cloud, host)(ctx)
}
