//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import "github.com/vdaas/vald/internal/strings"

// BlobStorageType represents blob storage type.
type BlobStorageType uint8

const (
	// S3 represents s3 storage type.
	S3 BlobStorageType = 1 + iota
	CloudStorage
)

// String returns blob storage type.
func (bst BlobStorageType) String() string {
	switch bst {
	case S3:
		return "s3"
	case CloudStorage:
		return "cloud_storage"
	}
	return "unknown"
}

// AtoBST returns BlobStorageType converted from string.
func AtoBST(bst string) BlobStorageType {
	switch strings.ToLower(bst) {
	case S3.String():
		return S3
	case CloudStorage.String():
		return CloudStorage
	}
	return 0
}

// Blob represents Blob configuration.
type Blob struct {
	// S3 represents S3 configuration.
	S3 *S3Config `json:"s3" yaml:"s3"`
	// CloudStorage represents CloudStorage configuration.
	CloudStorage *CloudStorageConfig `json:"cloud_storage" yaml:"cloud_storage"`
	// StorageType represents the storage type.
	StorageType string `json:"storage_type" yaml:"storage_type"`
	// Bucket represents the bucket name.
	Bucket string `json:"bucket" yaml:"bucket"`
}

// S3Config represents S3Config configuration.
type S3Config struct {
	// Endpoint represents the S3 endpoint.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// Region represents the S3 region.
	Region string `json:"region" yaml:"region"`
	// AccessKey represents the access key.
	AccessKey string `json:"access_key" yaml:"access_key"`
	// SecretAccessKey represents the secret access key.
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
	// Token represents the session token.
	Token string `json:"token" yaml:"token"`
	// MaxChunkSize represents the maximum chunk size.
	MaxChunkSize string `json:"max_chunk_size" yaml:"max_chunk_size"`
	// MaxPartSize represents the maximum part size.
	MaxPartSize string `json:"max_part_size" yaml:"max_part_size"`
	// MaxRetries represents the maximum retry count.
	MaxRetries int `json:"max_retries" yaml:"max_retries"`
	// UseARNRegion enables using ARN region.
	UseARNRegion bool `json:"use_arn_region" yaml:"use_arn_region"`
	// UseDualStack enables using dual stack.
	UseDualStack bool `json:"use_dual_stack" yaml:"use_dual_stack"`
	// EnableSSL enables SSL.
	EnableSSL bool `json:"enable_ssl" yaml:"enable_ssl"`
	// EnableParamValidation enables parameter validation.
	EnableParamValidation bool `json:"enable_param_validation" yaml:"enable_param_validation"`
	// Enable100Continue enables 100-continue.
	Enable100Continue bool `json:"enable_100_continue" yaml:"enable_100_continue"`
	// EnableContentMD5Validation enables content MD5 validation.
	EnableContentMD5Validation bool `json:"enable_content_md5_validation" yaml:"enable_content_md5_validation"`
	// EnableEndpointDiscovery enables endpoint discovery.
	EnableEndpointDiscovery bool `json:"enable_endpoint_discovery" yaml:"enable_endpoint_discovery"`
	// EnableEndpointHostPrefix enables endpoint host prefix.
	EnableEndpointHostPrefix bool `json:"enable_endpoint_host_prefix" yaml:"enable_endpoint_host_prefix"`
	// UseAccelerate enables using accelerate.
	UseAccelerate bool `json:"use_accelerate" yaml:"use_accelerate"`
	// ForcePathStyle enables forcing path style.
	ForcePathStyle bool `json:"force_path_style" yaml:"force_path_style"`
}

// CloudStorageConfig represents CloudStorage configuration.
type CloudStorageConfig struct {
	// Client represents the CloudStorage client configuration.
	Client *CloudStorageClient `json:"client" yaml:"client"`
	// URL represents the storage URL.
	URL string `json:"url" yaml:"url"`
	// WriteCacheControl represents the cache control header.
	WriteCacheControl string `json:"write_cache_control" yaml:"write_cache_control"`
	// WriteContentDisposition represents the content disposition header.
	WriteContentDisposition string `json:"write_content_disposition" yaml:"write_content_disposition"`
	// WriteContentEncoding represents the content encoding header.
	WriteContentEncoding string `json:"write_content_encoding" yaml:"write_content_encoding"`
	// WriteContentLanguage represents the content language header.
	WriteContentLanguage string `json:"write_content_language" yaml:"write_content_language"`
	// WriteContentType represents the content type header.
	WriteContentType string `json:"write_content_type" yaml:"write_content_type"`
	// WriteBufferSize represents the write buffer size.
	WriteBufferSize int `json:"write_buffer_size" yaml:"write_buffer_size"`
}

// CloudStorageClient represents CloudStorage client configuration.
type CloudStorageClient struct {
	// CredentialsFilePath represents the path to the credentials file.
	CredentialsFilePath string `json:"credentials_file_path" yaml:"credentials_file_path"`
	// CredentialsJSON represents the credentials JSON string.
	CredentialsJSON string `json:"credentials_json" yaml:"credentials_json"`
}

// Bind binds the actual data from the CloudStorageClient receiver fields.
func (csc *CloudStorageClient) Bind() *CloudStorageClient {
	csc.CredentialsFilePath = GetActualValue(csc.CredentialsFilePath)
	csc.CredentialsJSON = GetActualValue(csc.CredentialsJSON)
	return csc
}

// Bind binds the actual data from the Blob receiver field.
func (b *Blob) Bind() *Blob {
	b.StorageType = GetActualValue(b.StorageType)
	b.Bucket = GetActualValue(b.Bucket)

	if b.S3 == nil {
		b.S3 = new(S3Config)
	}
	b.S3.Bind()

	if b.CloudStorage == nil {
		b.CloudStorage = new(CloudStorageConfig)
	}
	b.CloudStorage.Bind()

	return b
}

// Bind binds the actual data from the S3Config receiver field.
func (s *S3Config) Bind() *S3Config {
	s.Endpoint = GetActualValue(s.Endpoint)
	s.Region = GetActualValue(s.Region)
	s.AccessKey = GetActualValue(s.AccessKey)
	s.SecretAccessKey = GetActualValue(s.SecretAccessKey)
	s.Token = GetActualValue(s.Token)
	s.MaxPartSize = GetActualValue(s.MaxPartSize)
	s.MaxChunkSize = GetActualValue(s.MaxChunkSize)

	return s
}

func (c *CloudStorageConfig) Bind() *CloudStorageConfig {
	c.URL = GetActualValue(c.URL)

	if c.Client == nil {
		c.Client = new(CloudStorageClient)
	}
	c.Client.Bind() // Call the new Bind method for CloudStorageClient

	c.WriteCacheControl = GetActualValue(c.WriteCacheControl)
	c.WriteContentDisposition = GetActualValue(c.WriteContentDisposition)
	c.WriteContentEncoding = GetActualValue(c.WriteContentEncoding)
	c.WriteContentLanguage = GetActualValue(c.WriteContentLanguage)
	c.WriteContentType = GetActualValue(c.WriteContentType)

	return c
}
