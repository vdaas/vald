//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package config provides configuration types and logic for loading and binding configuration values.
// This file includes detailed Bind methods for all configuration types with extensive comments.
package config

type OperationType string

const (
	OpSearch           OperationType = "search"
	OpSearchByID       OperationType = "search_by_id"
	OpLinearSearch     OperationType = "linear_search"
	OpLinearSearchByID OperationType = "linear_search_by_id"

	OpInsert            OperationType = "insert"
	OpUpdate            OperationType = "update"
	OpUpsert            OperationType = "upsert"
	OpRemove            OperationType = "remove"
	OpRemoveByTimestamp OperationType = "remove_by_timestamp"

	OpObject     OperationType = "object"
	OpListObject OperationType = "list_object"
	OpTimestamp  OperationType = "timestamp"
	OpExists     OperationType = "exists"

	OpIndexInfo             OperationType = "index_info"
	OpIndexDetail           OperationType = "index_detail"
	OpIndexStatistics       OperationType = "index_statistics"
	OpIndexStatisticsDetail OperationType = "index_statistics_detail"
	OpIndexProperty         OperationType = "index_property"
	OpFlush                 OperationType = "flush"

	OpKubernetes OperationType = "kubernetes"
	OpClient     OperationType = "client"
	OpWait       OperationType = "wait"
)

type StatusCode string

type StatusCodes []StatusCode

const (
	StatusCodeOK                 StatusCode = "ok"
	StatusCodeCanceled           StatusCode = "canceled"
	StatusCodeUnknown            StatusCode = "unknown"
	StatusCodeInvalidArgument    StatusCode = "invalidargument"
	StatusCodeDeadlineExceeded   StatusCode = "deadlineexceeded"
	StatusCodeNotFound           StatusCode = "notfound"
	StatusCodeAlreadyExists      StatusCode = "alreadyexists"
	StatusCodePermissionDenied   StatusCode = "permissiondenied"
	StatusCodeResourceExhausted  StatusCode = "resourceexhausted"
	StatusCodeFailedPrecondition StatusCode = "failedprecondition"
	StatusCodeAborted            StatusCode = "aborted"
	StatusCodeOutOfRange         StatusCode = "outofrange"
	StatusCodeUnimplemented      StatusCode = "unimplemented"
	StatusCodeInternal           StatusCode = "internal"
	StatusCodeUnavailable        StatusCode = "unavailable"
	StatusCodeDataLoss           StatusCode = "dataloss"
	StatusCodeUnauthenticated    StatusCode = "unauthenticated"
)

type OperationMode string

const (
	OperationUnary    OperationMode = "unary"
	OperationStream   OperationMode = "stream"
	OperationMultiple OperationMode = "multiple"
	OperationOther    OperationMode = "other"
)

type KubernetesAction string

const (
	KubernetesActionRollout KubernetesAction = "rollout"
	KubernetesActionDelete  KubernetesAction = "delete"
	KubernetesActionGet     KubernetesAction = "get"
	KubernetesActionExec    KubernetesAction = "exec"
	KubernetesActionApply   KubernetesAction = "apply"
	KubernetesActionCreate  KubernetesAction = "create"
	KubernetesActionPatch   KubernetesAction = "patch"
	KubernetesActionScale   KubernetesAction = "scale"
)

type KubernetesKind string

const (
	ConfigMap   KubernetesKind = "configmap"
	CronJob     KubernetesKind = "cronjob"
	DaemonSet   KubernetesKind = "daemonset"
	Deployment  KubernetesKind = "deployment"
	Job         KubernetesKind = "job"
	Pod         KubernetesKind = "pod"
	Secret      KubernetesKind = "secret"
	Service     KubernetesKind = "service"
	StatefulSet KubernetesKind = "statefulset"
)
