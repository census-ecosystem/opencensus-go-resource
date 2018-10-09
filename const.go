package resource

// Constants for Kubernetes resources.
const (
	K8STypeContainer = "k8s.io/container"

	K8SKeyNamespaceName = "k8s.io/namespace/name"
	K8SKeyPodName       = "k8s.io/pod/name"
	K8SKeyContainerName = "k8s.io/container/name"
)

// Constants for AWS resources.
const (
	AWSTypeEC2Instance = "aws.com/ec2/instance"

	AWSKeyEC2AccountID  = "aws.com/ec2/account_id"
	AWSKeyEC2Region     = "aws.com/ec2/region"
	AWSKeyEC2InstanceID = "aws.com/ec2/instance_id"
)

// Constants for GCP resources.
const (
	GCPTypeGCEInstance = "cloud.google.com/gce/instance"

	// ProjectID of the GCE VM. This is not the project ID of the used client credentials.
	GCPKeyGCEProjectID  = "cloud.google.com/gce/project_id"
	GCPKeyGCEZone       = "cloud.google.com/gce/zone"
	GCPKeyGCEInstanceID = "cloud.google.com/gce/instance_id"
	// Key for the node attribute that's automatically set with the GKE cluster name.
	GCPKeyGCEClusterName = "cloud.google.com/gce/attributes/cluster_name"
)
