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
// CreateObjectStoreBucketProperties struct for CreateObjectStoreBucketProperties
type CreateObjectStoreBucketProperties struct {
	Alibaba *CreateAlibabaObjectStoreBucketProperties `json:"alibaba,omitempty"`
	Amazon *CreateAmazonObjectStoreBucketProperties `json:"amazon,omitempty"`
	Azure *CreateAzureObjectStoreBucketProperties `json:"azure,omitempty"`
	Google *CreateGoogleObjectStoreBucketProperties `json:"google,omitempty"`
	Oracle *CreateOracleObjectStoreBucketProperties `json:"oracle,omitempty"`
}
