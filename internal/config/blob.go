//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package config providers configuration type and load configuration logic
package config

import "strings"

type BlobStorageType uint8

const (
	S3 BlobStorageType = 1 + iota
)

func (bst BlobStorageType) String() string {
	switch bst {
	case S3:
		return "s3"
	}
	return "unknown"
}

func AtoBST(bst string) BlobStorageType {
	switch strings.ToLower(bst) {
	case S3.String():
		return S3
	}
	return 0
}

type Blob struct {
	// StorageType represents blob storaget type
	StorageType string `json:"storage_type" yaml:"storage_type"`

	// Bucket represents bucket name
	Bucket string `json:"bucket" yaml:"bucket"`

	// S3 represents S3 config
	S3 *S3Config `json:"s3" yaml:"s3"`

	// CloudStrage represents CloudStrage config
	CloudStrage *CloudStrage `json:"cloud_strage" yaml:"cloud_strage"`
}

type S3Config struct {
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	Region          string `json:"region" yaml:"region"`
	AccessKey       string `json:"access_key" yaml:"access_key"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
	Token           string `json:"token" yaml:"token"`

	MaxRetries                 int  `json:"max_retries" yaml:"max_retries"`
	ForcePathStyle             bool `json:"force_path_style" yaml:"force_path_style"`
	UseAccelerate              bool `json:"use_accelerate" yaml:"use_accelerate"`
	UseARNRegion               bool `json:"use_arn_region" yaml:"use_arn_region"`
	UseDualStack               bool `json:"use_dual_stack" yaml:"use_dual_stack"`
	EnableSSL                  bool `json:"enable_ssl" yaml:"enable_ssl"`
	EnableParamValidation      bool `json:"enable_param_validation" yaml:"enable_param_validation"`
	Enable100Continue          bool `json:"enable_100_continue" yaml:"enable_100_continue"`
	EnableContentMD5Validation bool `json:"enable_content_md5_validation" yaml:"enable_content_md5_validation"`
	EnableEndpointDiscovery    bool `json:"enable_endpoint_discovery" yaml:"enable_endpoint_discovery"`
	EnableEndpointHostPrefix   bool `json:"enable_endpoint_host_prefix" yaml:"enable_endpoint_host_prefix"`

	MaxPartSize string `json:"max_part_size" yaml:"max_part_size"`
}

type CloudStrage struct {
	URL string `json:"url" yaml:"url"`

	Reader struct {
	} `json:"reader" yaml:"reader"`

	Writer struct {
		BufferSize         int    `json:"buffer_size" yaml:"buffer_size"`
		CacheControl       string `json:"cache_control" yaml:"cache_control"`
		ContentDisposition string `json:"content_disposition" yaml:"content_disposition"`
		ContentEncoding    string `json:"content_encoding" yaml:"content_encoding"`
		ContentLanguage    string `json:"content_language" yaml:"content_language"`
		ContentType        string `json:"content_type" yaml:"content_type"`
	} `json:"writer" yaml:"writer"`
}

func (b *Blob) Bind() *Blob {
	b.StorageType = GetActualValue(b.StorageType)
	b.Bucket = GetActualValue(b.Bucket)

	if b.S3 != nil {
		b.S3 = b.S3.Bind()
	} else {
		b.S3 = new(S3Config)
	}

	if b.CloudStrage != nil {
		b.CloudStrage = b.CloudStrage.Bind()
	} else {
		b.CloudStrage = new(CloudStrage)
	}

	return b
}

func (s *S3Config) Bind() *S3Config {
	s.Endpoint = GetActualValue(s.Endpoint)
	s.Region = GetActualValue(s.Region)
	s.AccessKey = GetActualValue(s.AccessKey)
	s.SecretAccessKey = GetActualValue(s.SecretAccessKey)
	s.Token = GetActualValue(s.Token)
	s.MaxPartSize = GetActualValue(s.MaxPartSize)

	return s
}

func (c *CloudStrage) Bind() *CloudStrage {
	c.Writer.CacheControl = GetActualValue(c.Writer.CacheControl)
	c.Writer.ContentDisposition = GetActualValue(c.Writer.ContentDisposition)
	c.Writer.ContentEncoding = GetActualValue(c.Writer.ContentEncoding)
	c.Writer.ContentLanguage = GetActualValue(c.Writer.ContentLanguage)
	c.Writer.ContentType = GetActualValue(c.Writer.ContentType)

	return c
}
