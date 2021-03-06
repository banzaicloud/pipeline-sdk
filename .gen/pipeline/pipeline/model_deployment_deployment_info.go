/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments. 
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pipeline
// DeploymentDeploymentInfo struct for DeploymentDeploymentInfo
type DeploymentDeploymentInfo struct {
	Chart string `json:"chart,omitempty"`
	ChartName string `json:"chartName,omitempty"`
	ChartVersion string `json:"chartVersion,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Description string `json:"description,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	ReleaseName string `json:"releaseName,omitempty"`
	TargetClusters []DeploymentTargetClusterStatus `json:"targetClusters,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	ValueOverrides map[string]interface{} `json:"valueOverrides,omitempty"`
	Values map[string]interface{} `json:"values,omitempty"`
	Version int32 `json:"version,omitempty"`
}
