# Protocol Documentation

<a name="top"></a>

## Table of Contents

- [v1/payload/payload.proto](#v1_payload_payload-proto)
  - [Control](#payload-v1-Control)
  - [Control.CreateIndexRequest](#payload-v1-Control-CreateIndexRequest)
  - [Discoverer](#payload-v1-Discoverer)
  - [Discoverer.Request](#payload-v1-Discoverer-Request)
  - [Empty](#payload-v1-Empty)
  - [Filter](#payload-v1-Filter)
  - [Filter.Config](#payload-v1-Filter-Config)
  - [Filter.Target](#payload-v1-Filter-Target)
  - [Flush](#payload-v1-Flush)
  - [Flush.Request](#payload-v1-Flush-Request)
  - [Info](#payload-v1-Info)
  - [Info.Annotations](#payload-v1-Info-Annotations)
  - [Info.Annotations.AnnotationsEntry](#payload-v1-Info-Annotations-AnnotationsEntry)
  - [Info.CPU](#payload-v1-Info-CPU)
  - [Info.IPs](#payload-v1-Info-IPs)
  - [Info.Index](#payload-v1-Info-Index)
  - [Info.Index.Count](#payload-v1-Info-Index-Count)
  - [Info.Index.Detail](#payload-v1-Info-Index-Detail)
  - [Info.Index.Detail.CountsEntry](#payload-v1-Info-Index-Detail-CountsEntry)
  - [Info.Index.Property](#payload-v1-Info-Index-Property)
  - [Info.Index.PropertyDetail](#payload-v1-Info-Index-PropertyDetail)
  - [Info.Index.PropertyDetail.DetailsEntry](#payload-v1-Info-Index-PropertyDetail-DetailsEntry)
  - [Info.Index.Statistics](#payload-v1-Info-Index-Statistics)
  - [Info.Index.StatisticsDetail](#payload-v1-Info-Index-StatisticsDetail)
  - [Info.Index.StatisticsDetail.DetailsEntry](#payload-v1-Info-Index-StatisticsDetail-DetailsEntry)
  - [Info.Index.UUID](#payload-v1-Info-Index-UUID)
  - [Info.Index.UUID.Committed](#payload-v1-Info-Index-UUID-Committed)
  - [Info.Index.UUID.Uncommitted](#payload-v1-Info-Index-UUID-Uncommitted)
  - [Info.Labels](#payload-v1-Info-Labels)
  - [Info.Labels.LabelsEntry](#payload-v1-Info-Labels-LabelsEntry)
  - [Info.Memory](#payload-v1-Info-Memory)
  - [Info.Node](#payload-v1-Info-Node)
  - [Info.Nodes](#payload-v1-Info-Nodes)
  - [Info.Pod](#payload-v1-Info-Pod)
  - [Info.Pods](#payload-v1-Info-Pods)
  - [Info.Service](#payload-v1-Info-Service)
  - [Info.ServicePort](#payload-v1-Info-ServicePort)
  - [Info.Services](#payload-v1-Info-Services)
  - [Insert](#payload-v1-Insert)
  - [Insert.Config](#payload-v1-Insert-Config)
  - [Insert.MultiObjectRequest](#payload-v1-Insert-MultiObjectRequest)
  - [Insert.MultiRequest](#payload-v1-Insert-MultiRequest)
  - [Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest)
  - [Insert.Request](#payload-v1-Insert-Request)
  - [Meta](#payload-v1-Meta)
  - [Meta.Key](#payload-v1-Meta-Key)
  - [Meta.KeyValue](#payload-v1-Meta-KeyValue)
  - [Meta.Value](#payload-v1-Meta-Value)
  - [Mirror](#payload-v1-Mirror)
  - [Mirror.Target](#payload-v1-Mirror-Target)
  - [Mirror.Targets](#payload-v1-Mirror-Targets)
  - [Object](#payload-v1-Object)
  - [Object.Blob](#payload-v1-Object-Blob)
  - [Object.Distance](#payload-v1-Object-Distance)
  - [Object.ID](#payload-v1-Object-ID)
  - [Object.IDs](#payload-v1-Object-IDs)
  - [Object.List](#payload-v1-Object-List)
  - [Object.List.Request](#payload-v1-Object-List-Request)
  - [Object.List.Response](#payload-v1-Object-List-Response)
  - [Object.Location](#payload-v1-Object-Location)
  - [Object.Locations](#payload-v1-Object-Locations)
  - [Object.ReshapeVector](#payload-v1-Object-ReshapeVector)
  - [Object.StreamBlob](#payload-v1-Object-StreamBlob)
  - [Object.StreamDistance](#payload-v1-Object-StreamDistance)
  - [Object.StreamLocation](#payload-v1-Object-StreamLocation)
  - [Object.StreamVector](#payload-v1-Object-StreamVector)
  - [Object.Timestamp](#payload-v1-Object-Timestamp)
  - [Object.TimestampRequest](#payload-v1-Object-TimestampRequest)
  - [Object.Vector](#payload-v1-Object-Vector)
  - [Object.VectorRequest](#payload-v1-Object-VectorRequest)
  - [Object.Vectors](#payload-v1-Object-Vectors)
  - [Remove](#payload-v1-Remove)
  - [Remove.Config](#payload-v1-Remove-Config)
  - [Remove.MultiRequest](#payload-v1-Remove-MultiRequest)
  - [Remove.Request](#payload-v1-Remove-Request)
  - [Remove.Timestamp](#payload-v1-Remove-Timestamp)
  - [Remove.TimestampRequest](#payload-v1-Remove-TimestampRequest)
  - [Search](#payload-v1-Search)
  - [Search.Config](#payload-v1-Search-Config)
  - [Search.IDRequest](#payload-v1-Search-IDRequest)
  - [Search.MultiIDRequest](#payload-v1-Search-MultiIDRequest)
  - [Search.MultiObjectRequest](#payload-v1-Search-MultiObjectRequest)
  - [Search.MultiRequest](#payload-v1-Search-MultiRequest)
  - [Search.ObjectRequest](#payload-v1-Search-ObjectRequest)
  - [Search.Request](#payload-v1-Search-Request)
  - [Search.Response](#payload-v1-Search-Response)
  - [Search.Responses](#payload-v1-Search-Responses)
  - [Search.StreamResponse](#payload-v1-Search-StreamResponse)
  - [Update](#payload-v1-Update)
  - [Update.Config](#payload-v1-Update-Config)
  - [Update.MultiObjectRequest](#payload-v1-Update-MultiObjectRequest)
  - [Update.MultiRequest](#payload-v1-Update-MultiRequest)
  - [Update.ObjectRequest](#payload-v1-Update-ObjectRequest)
  - [Update.Request](#payload-v1-Update-Request)
  - [Update.TimestampRequest](#payload-v1-Update-TimestampRequest)
  - [Upsert](#payload-v1-Upsert)
  - [Upsert.Config](#payload-v1-Upsert-Config)
  - [Upsert.MultiObjectRequest](#payload-v1-Upsert-MultiObjectRequest)
  - [Upsert.MultiRequest](#payload-v1-Upsert-MultiRequest)
  - [Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest)
  - [Upsert.Request](#payload-v1-Upsert-Request)
  - [Remove.Timestamp.Operator](#payload-v1-Remove-Timestamp-Operator)
  - [Search.AggregationAlgorithm](#payload-v1-Search-AggregationAlgorithm)

- [v1/agent/core/agent.proto](#v1_agent_core_agent-proto)
  - [Agent](#core-v1-Agent)
- [v1/agent/sidecar/sidecar.proto](#v1_agent_sidecar_sidecar-proto)
  - [Sidecar](#sidecar-v1-Sidecar)
- [v1/discoverer/discoverer.proto](#v1_discoverer_discoverer-proto)
  - [Discoverer](#discoverer-v1-Discoverer)
- [v1/filter/egress/egress_filter.proto](#v1_filter_egress_egress_filter-proto)
  - [Filter](#filter-egress-v1-Filter)
- [v1/filter/ingress/ingress_filter.proto](#v1_filter_ingress_ingress_filter-proto)
  - [Filter](#filter-ingress-v1-Filter)
- [v1/meta/meta.proto](#v1_meta_meta-proto)
  - [Meta](#meta-v1-Meta)
- [v1/mirror/mirror.proto](#v1_mirror_mirror-proto)
  - [Mirror](#mirror-v1-Mirror)
- [v1/rpc/errdetails/error_details.proto](#v1_rpc_errdetails_error_details-proto)
  - [BadRequest](#rpc-v1-BadRequest)
  - [BadRequest.FieldViolation](#rpc-v1-BadRequest-FieldViolation)
  - [DebugInfo](#rpc-v1-DebugInfo)
  - [ErrorInfo](#rpc-v1-ErrorInfo)
  - [ErrorInfo.MetadataEntry](#rpc-v1-ErrorInfo-MetadataEntry)
  - [Help](#rpc-v1-Help)
  - [Help.Link](#rpc-v1-Help-Link)
  - [LocalizedMessage](#rpc-v1-LocalizedMessage)
  - [PreconditionFailure](#rpc-v1-PreconditionFailure)
  - [PreconditionFailure.Violation](#rpc-v1-PreconditionFailure-Violation)
  - [QuotaFailure](#rpc-v1-QuotaFailure)
  - [QuotaFailure.Violation](#rpc-v1-QuotaFailure-Violation)
  - [RequestInfo](#rpc-v1-RequestInfo)
  - [ResourceInfo](#rpc-v1-ResourceInfo)
  - [RetryInfo](#rpc-v1-RetryInfo)
- [v1/vald/filter.proto](#v1_vald_filter-proto)
  - [Filter](#vald-v1-Filter)
- [v1/vald/flush.proto](#v1_vald_flush-proto)
  - [Flush](#vald-v1-Flush)
- [v1/vald/index.proto](#v1_vald_index-proto)
  - [Index](#vald-v1-Index)
- [v1/vald/insert.proto](#v1_vald_insert-proto)
  - [Insert](#vald-v1-Insert)
- [v1/vald/object.proto](#v1_vald_object-proto)
  - [Object](#vald-v1-Object)
- [v1/vald/remove.proto](#v1_vald_remove-proto)
  - [Remove](#vald-v1-Remove)
- [v1/vald/search.proto](#v1_vald_search-proto)
  - [Search](#vald-v1-Search)
- [v1/vald/update.proto](#v1_vald_update-proto)
  - [Update](#vald-v1-Update)
- [v1/vald/upsert.proto](#v1_vald_upsert-proto)
  - [Upsert](#vald-v1-Upsert)
- [Scalar Value Types](#scalar-value-types)

<a name="v1_payload_payload-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/payload/payload.proto

<a name="payload-v1-Control"></a>

### Control

Control related messages.

<a name="payload-v1-Control-CreateIndexRequest"></a>

### Control.CreateIndexRequest

Represent the create index request.

| Field     | Type              | Label | Description                                  |
| --------- | ----------------- | ----- | -------------------------------------------- |
| pool_size | [uint32](#uint32) |       | The pool size of the create index operation. |

<a name="payload-v1-Discoverer"></a>

### Discoverer

Discoverer related messages.

<a name="payload-v1-Discoverer-Request"></a>

### Discoverer.Request

Represent the dicoverer request.

| Field     | Type              | Label | Description                      |
| --------- | ----------------- | ----- | -------------------------------- |
| name      | [string](#string) |       | The agent name to be discovered. |
| namespace | [string](#string) |       | The namespace to be discovered.  |
| node      | [string](#string) |       | The node to be discovered.       |

<a name="payload-v1-Empty"></a>

### Empty

Represent an empty message.

<a name="payload-v1-Filter"></a>

### Filter

Filter related messages.

<a name="payload-v1-Filter-Config"></a>

### Filter.Config

Represent filter configuration.

| Field   | Type                                       | Label    | Description                                |
| ------- | ------------------------------------------ | -------- | ------------------------------------------ |
| targets | [Filter.Target](#payload-v1-Filter-Target) | repeated | Represent the filter target configuration. |

<a name="payload-v1-Filter-Target"></a>

### Filter.Target

Represent the target filter server.

| Field | Type              | Label | Description          |
| ----- | ----------------- | ----- | -------------------- |
| host  | [string](#string) |       | The target hostname. |
| port  | [uint32](#uint32) |       | The target port.     |

<a name="payload-v1-Flush"></a>

### Flush

Flush related messages.

<a name="payload-v1-Flush-Request"></a>

### Flush.Request

<a name="payload-v1-Info"></a>

### Info

Info related messages.

<a name="payload-v1-Info-Annotations"></a>

### Info.Annotations

Represent the kubernetes annotations.

| Field       | Type                                                                               | Label    | Description |
| ----------- | ---------------------------------------------------------------------------------- | -------- | ----------- |
| annotations | [Info.Annotations.AnnotationsEntry](#payload-v1-Info-Annotations-AnnotationsEntry) | repeated |             |

<a name="payload-v1-Info-Annotations-AnnotationsEntry"></a>

### Info.Annotations.AnnotationsEntry

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| key   | [string](#string) |       |             |
| value | [string](#string) |       |             |

<a name="payload-v1-Info-CPU"></a>

### Info.CPU

Represent the CPU information message.

| Field   | Type              | Label | Description                 |
| ------- | ----------------- | ----- | --------------------------- |
| limit   | [double](#double) |       | The CPU resource limit.     |
| request | [double](#double) |       | The CPU resource requested. |
| usage   | [double](#double) |       | The CPU usage.              |

<a name="payload-v1-Info-IPs"></a>

### Info.IPs

Represent the multiple IP message.

| Field | Type              | Label    | Description |
| ----- | ----------------- | -------- | ----------- |
| ip    | [string](#string) | repeated |             |

<a name="payload-v1-Info-Index"></a>

### Info.Index

Represent the index information messages.

<a name="payload-v1-Info-Index-Count"></a>

### Info.Index.Count

Represent the index count message.

| Field       | Type              | Label | Description                  |
| ----------- | ----------------- | ----- | ---------------------------- |
| stored      | [uint32](#uint32) |       | The stored index count.      |
| uncommitted | [uint32](#uint32) |       | The uncommitted index count. |
| indexing    | [bool](#bool)     |       | The indexing index count.    |
| saving      | [bool](#bool)     |       | The saving index count.      |

<a name="payload-v1-Info-Index-Detail"></a>

### Info.Index.Detail

Represent the index count for each Agents message.

| Field       | Type                                                                       | Label    | Description                        |
| ----------- | -------------------------------------------------------------------------- | -------- | ---------------------------------- |
| counts      | [Info.Index.Detail.CountsEntry](#payload-v1-Info-Index-Detail-CountsEntry) | repeated | count infos for each agents        |
| replica     | [uint32](#uint32)                                                          |          | index replica of vald cluster      |
| live_agents | [uint32](#uint32)                                                          |          | live agent replica of vald cluster |

<a name="payload-v1-Info-Index-Detail-CountsEntry"></a>

### Info.Index.Detail.CountsEntry

| Field | Type                                             | Label | Description |
| ----- | ------------------------------------------------ | ----- | ----------- |
| key   | [string](#string)                                |       |             |
| value | [Info.Index.Count](#payload-v1-Info-Index-Count) |       |             |

<a name="payload-v1-Info-Index-Property"></a>

### Info.Index.Property

Represents index Property

| Field                              | Type              | Label | Description |
| ---------------------------------- | ----------------- | ----- | ----------- |
| dimension                          | [int32](#int32)   |       |             |
| thread_pool_size                   | [int32](#int32)   |       |             |
| object_type                        | [string](#string) |       |             |
| distance_type                      | [string](#string) |       |             |
| index_type                         | [string](#string) |       |             |
| database_type                      | [string](#string) |       |             |
| object_alignment                   | [string](#string) |       |             |
| path_adjustment_interval           | [int32](#int32)   |       |             |
| graph_shared_memory_size           | [int32](#int32)   |       |             |
| tree_shared_memory_size            | [int32](#int32)   |       |             |
| object_shared_memory_size          | [int32](#int32)   |       |             |
| prefetch_offset                    | [int32](#int32)   |       |             |
| prefetch_size                      | [int32](#int32)   |       |             |
| accuracy_table                     | [string](#string) |       |             |
| search_type                        | [string](#string) |       |             |
| max_magnitude                      | [float](#float)   |       |             |
| n_of_neighbors_for_insertion_order | [int32](#int32)   |       |             |
| epsilon_for_insertion_order        | [float](#float)   |       |             |
| refinement_object_type             | [string](#string) |       |             |
| truncation_threshold               | [int32](#int32)   |       |             |
| edge_size_for_creation             | [int32](#int32)   |       |             |
| edge_size_for_search               | [int32](#int32)   |       |             |
| edge_size_limit_for_creation       | [int32](#int32)   |       |             |
| insertion_radius_coefficient       | [double](#double) |       |             |
| seed_size                          | [int32](#int32)   |       |             |
| seed_type                          | [string](#string) |       |             |
| truncation_thread_pool_size        | [int32](#int32)   |       |             |
| batch_size_for_creation            | [int32](#int32)   |       |             |
| graph_type                         | [string](#string) |       |             |
| dynamic_edge_size_base             | [int32](#int32)   |       |             |
| dynamic_edge_size_rate             | [int32](#int32)   |       |             |
| build_time_limit                   | [float](#float)   |       |             |
| outgoing_edge                      | [int32](#int32)   |       |             |
| incoming_edge                      | [int32](#int32)   |       |             |

<a name="payload-v1-Info-Index-PropertyDetail"></a>

### Info.Index.PropertyDetail

Represents index Properties for each Agents

| Field   | Type                                                                                         | Label    | Description |
| ------- | -------------------------------------------------------------------------------------------- | -------- | ----------- |
| details | [Info.Index.PropertyDetail.DetailsEntry](#payload-v1-Info-Index-PropertyDetail-DetailsEntry) | repeated |             |

<a name="payload-v1-Info-Index-PropertyDetail-DetailsEntry"></a>

### Info.Index.PropertyDetail.DetailsEntry

| Field | Type                                                   | Label | Description |
| ----- | ------------------------------------------------------ | ----- | ----------- |
| key   | [string](#string)                                      |       |             |
| value | [Info.Index.Property](#payload-v1-Info-Index-Property) |       |             |

<a name="payload-v1-Info-Index-Statistics"></a>

### Info.Index.Statistics

Represents index Statistics

| Field                                | Type              | Label    | Description |
| ------------------------------------ | ----------------- | -------- | ----------- |
| valid                                | [bool](#bool)     |          |             |
| median_indegree                      | [int32](#int32)   |          |             |
| median_outdegree                     | [int32](#int32)   |          |             |
| max_number_of_indegree               | [uint64](#uint64) |          |             |
| max_number_of_outdegree              | [uint64](#uint64) |          |             |
| min_number_of_indegree               | [uint64](#uint64) |          |             |
| min_number_of_outdegree              | [uint64](#uint64) |          |             |
| mode_indegree                        | [uint64](#uint64) |          |             |
| mode_outdegree                       | [uint64](#uint64) |          |             |
| nodes_skipped_for_10_edges           | [uint64](#uint64) |          |             |
| nodes_skipped_for_indegree_distance  | [uint64](#uint64) |          |             |
| number_of_edges                      | [uint64](#uint64) |          |             |
| number_of_indexed_objects            | [uint64](#uint64) |          |             |
| number_of_nodes                      | [uint64](#uint64) |          |             |
| number_of_nodes_without_edges        | [uint64](#uint64) |          |             |
| number_of_nodes_without_indegree     | [uint64](#uint64) |          |             |
| number_of_objects                    | [uint64](#uint64) |          |             |
| number_of_removed_objects            | [uint64](#uint64) |          |             |
| size_of_object_repository            | [uint64](#uint64) |          |             |
| size_of_refinement_object_repository | [uint64](#uint64) |          |             |
| variance_of_indegree                 | [double](#double) |          |             |
| variance_of_outdegree                | [double](#double) |          |             |
| mean_edge_length                     | [double](#double) |          |             |
| mean_edge_length_for_10_edges        | [double](#double) |          |             |
| mean_indegree_distance_for_10_edges  | [double](#double) |          |             |
| mean_number_of_edges_per_node        | [double](#double) |          |             |
| c1_indegree                          | [double](#double) |          |             |
| c5_indegree                          | [double](#double) |          |             |
| c95_outdegree                        | [double](#double) |          |             |
| c99_outdegree                        | [double](#double) |          |             |
| indegree_count                       | [int64](#int64)   | repeated |             |
| outdegree_histogram                  | [uint64](#uint64) | repeated |             |
| indegree_histogram                   | [uint64](#uint64) | repeated |             |

<a name="payload-v1-Info-Index-StatisticsDetail"></a>

### Info.Index.StatisticsDetail

Represents index Statistics for each Agents

| Field   | Type                                                                                             | Label    | Description                 |
| ------- | ------------------------------------------------------------------------------------------------ | -------- | --------------------------- |
| details | [Info.Index.StatisticsDetail.DetailsEntry](#payload-v1-Info-Index-StatisticsDetail-DetailsEntry) | repeated | count infos for each agents |

<a name="payload-v1-Info-Index-StatisticsDetail-DetailsEntry"></a>

### Info.Index.StatisticsDetail.DetailsEntry

| Field | Type                                                       | Label | Description |
| ----- | ---------------------------------------------------------- | ----- | ----------- |
| key   | [string](#string)                                          |       |             |
| value | [Info.Index.Statistics](#payload-v1-Info-Index-Statistics) |       |             |

<a name="payload-v1-Info-Index-UUID"></a>

### Info.Index.UUID

Represent the UUID message.

<a name="payload-v1-Info-Index-UUID-Committed"></a>

### Info.Index.UUID.Committed

The committed UUID.

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| uuid  | [string](#string) |       |             |

<a name="payload-v1-Info-Index-UUID-Uncommitted"></a>

### Info.Index.UUID.Uncommitted

The uncommitted UUID.

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| uuid  | [string](#string) |       |             |

<a name="payload-v1-Info-Labels"></a>

### Info.Labels

Represent the kubernetes labels.

| Field  | Type                                                           | Label    | Description |
| ------ | -------------------------------------------------------------- | -------- | ----------- |
| labels | [Info.Labels.LabelsEntry](#payload-v1-Info-Labels-LabelsEntry) | repeated |             |

<a name="payload-v1-Info-Labels-LabelsEntry"></a>

### Info.Labels.LabelsEntry

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| key   | [string](#string) |       |             |
| value | [string](#string) |       |             |

<a name="payload-v1-Info-Memory"></a>

### Info.Memory

Represent the memory information message.

| Field   | Type              | Label | Description           |
| ------- | ----------------- | ----- | --------------------- |
| limit   | [double](#double) |       | The memory limit.     |
| request | [double](#double) |       | The memory requested. |
| usage   | [double](#double) |       | The memory usage.     |

<a name="payload-v1-Info-Node"></a>

### Info.Node

Represent the node information message.

| Field         | Type                                   | Label | Description                          |
| ------------- | -------------------------------------- | ----- | ------------------------------------ |
| name          | [string](#string)                      |       | The name of the node.                |
| internal_addr | [string](#string)                      |       | The internal IP address of the node. |
| external_addr | [string](#string)                      |       | The external IP address of the node. |
| cpu           | [Info.CPU](#payload-v1-Info-CPU)       |       | The CPU information of the node.     |
| memory        | [Info.Memory](#payload-v1-Info-Memory) |       | The memory information of the node.  |
| Pods          | [Info.Pods](#payload-v1-Info-Pods)     |       | The pod information of the node.     |

<a name="payload-v1-Info-Nodes"></a>

### Info.Nodes

Represent the multiple node information message.

| Field | Type                               | Label    | Description                    |
| ----- | ---------------------------------- | -------- | ------------------------------ |
| nodes | [Info.Node](#payload-v1-Info-Node) | repeated | The multiple node information. |

<a name="payload-v1-Info-Pod"></a>

### Info.Pod

Represent the pod information message.

| Field     | Type                                   | Label | Description                           |
| --------- | -------------------------------------- | ----- | ------------------------------------- |
| app_name  | [string](#string)                      |       | The app name of the pod on the label. |
| name      | [string](#string)                      |       | The name of the pod.                  |
| namespace | [string](#string)                      |       | The namespace of the pod.             |
| ip        | [string](#string)                      |       | The IP of the pod.                    |
| cpu       | [Info.CPU](#payload-v1-Info-CPU)       |       | The CPU information of the pod.       |
| memory    | [Info.Memory](#payload-v1-Info-Memory) |       | The memory information of the pod.    |
| node      | [Info.Node](#payload-v1-Info-Node)     |       | The node information of the pod.      |

<a name="payload-v1-Info-Pods"></a>

### Info.Pods

Represent the multiple pod information message.

| Field | Type                             | Label    | Description                   |
| ----- | -------------------------------- | -------- | ----------------------------- |
| pods  | [Info.Pod](#payload-v1-Info-Pod) | repeated | The multiple pod information. |

<a name="payload-v1-Info-Service"></a>

### Info.Service

Represent the service information message.

| Field       | Type                                             | Label    | Description                     |
| ----------- | ------------------------------------------------ | -------- | ------------------------------- |
| name        | [string](#string)                                |          | The name of the svc.            |
| cluster_ip  | [string](#string)                                |          | The cluster ip of the svc.      |
| cluster_ips | [string](#string)                                | repeated | The cluster ips of the svc.     |
| ports       | [Info.ServicePort](#payload-v1-Info-ServicePort) | repeated | The port of the svc.            |
| labels      | [Info.Labels](#payload-v1-Info-Labels)           |          | The labels of the service.      |
| annotations | [Info.Annotations](#payload-v1-Info-Annotations) |          | The annotations of the service. |

<a name="payload-v1-Info-ServicePort"></a>

### Info.ServicePort

Represets the service port information message.

| Field | Type              | Label | Description           |
| ----- | ----------------- | ----- | --------------------- |
| name  | [string](#string) |       | The name of the port. |
| port  | [int32](#int32)   |       | The port number       |

<a name="payload-v1-Info-Services"></a>

### Info.Services

Represent the multiple service information message.

| Field    | Type                                     | Label    | Description                       |
| -------- | ---------------------------------------- | -------- | --------------------------------- |
| services | [Info.Service](#payload-v1-Info-Service) | repeated | The multiple service information. |

<a name="payload-v1-Insert"></a>

### Insert

Insert related messages.

<a name="payload-v1-Insert-Config"></a>

### Insert.Config

Represent insert configurations.

| Field                   | Type                                       | Label | Description                                         |
| ----------------------- | ------------------------------------------ | ----- | --------------------------------------------------- |
| skip_strict_exist_check | [bool](#bool)                              |       | A flag to skip exist check during insert operation. |
| filters                 | [Filter.Config](#payload-v1-Filter-Config) |       | Filter configurations.                              |
| timestamp               | [int64](#int64)                            |       | Insert timestamp.                                   |

<a name="payload-v1-Insert-MultiObjectRequest"></a>

### Insert.MultiObjectRequest

Represent the multiple insert by binary object request.

| Field    | Type                                                     | Label    | Description                                  |
| -------- | -------------------------------------------------------- | -------- | -------------------------------------------- |
| requests | [Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest) | repeated | Represent multiple insert by object content. |

<a name="payload-v1-Insert-MultiRequest"></a>

### Insert.MultiRequest

Represent the multiple insert request.

| Field    | Type                                         | Label    | Description                                |
| -------- | -------------------------------------------- | -------- | ------------------------------------------ |
| requests | [Insert.Request](#payload-v1-Insert-Request) | repeated | Represent multiple insert request content. |

<a name="payload-v1-Insert-ObjectRequest"></a>

### Insert.ObjectRequest

Represent the insert by binary object request.

| Field      | Type                                       | Label | Description                              |
| ---------- | ------------------------------------------ | ----- | ---------------------------------------- |
| object     | [Object.Blob](#payload-v1-Object-Blob)     |       | The binary object to be inserted.        |
| config     | [Insert.Config](#payload-v1-Insert-Config) |       | The configuration of the insert request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |       | Filter configurations.                   |

<a name="payload-v1-Insert-Request"></a>

### Insert.Request

Represent the insert request.

| Field  | Type                                       | Label | Description                              |
| ------ | ------------------------------------------ | ----- | ---------------------------------------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |       | The vector to be inserted.               |
| config | [Insert.Config](#payload-v1-Insert-Config) |       | The configuration of the insert request. |

<a name="payload-v1-Meta"></a>

### Meta

<a name="payload-v1-Meta-Key"></a>

### Meta.Key

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| key   | [string](#string) |       |             |

<a name="payload-v1-Meta-KeyValue"></a>

### Meta.KeyValue

| Field | Type                                 | Label | Description |
| ----- | ------------------------------------ | ----- | ----------- |
| key   | [Meta.Key](#payload-v1-Meta-Key)     |       |             |
| value | [Meta.Value](#payload-v1-Meta-Value) |       |             |

<a name="payload-v1-Meta-Value"></a>

### Meta.Value

| Field | Type                                        | Label | Description |
| ----- | ------------------------------------------- | ----- | ----------- |
| value | [google.protobuf.Any](#google-protobuf-Any) |       |             |

<a name="payload-v1-Mirror"></a>

### Mirror

Mirror related messages.

<a name="payload-v1-Mirror-Target"></a>

### Mirror.Target

Represent server information.

| Field | Type              | Label | Description          |
| ----- | ----------------- | ----- | -------------------- |
| host  | [string](#string) |       | The target hostname. |
| port  | [uint32](#uint32) |       | The target port.     |

<a name="payload-v1-Mirror-Targets"></a>

### Mirror.Targets

Represent the multiple Target message.

| Field   | Type                                       | Label    | Description                      |
| ------- | ------------------------------------------ | -------- | -------------------------------- |
| targets | [Mirror.Target](#payload-v1-Mirror-Target) | repeated | The multiple target information. |

<a name="payload-v1-Object"></a>

### Object

Common messages.

<a name="payload-v1-Object-Blob"></a>

### Object.Blob

Represent the binary object.

| Field  | Type              | Label | Description        |
| ------ | ----------------- | ----- | ------------------ |
| id     | [string](#string) |       | The object ID.     |
| object | [bytes](#bytes)   |       | The binary object. |

<a name="payload-v1-Object-Distance"></a>

### Object.Distance

Represent the ID and distance pair.

| Field    | Type              | Label | Description    |
| -------- | ----------------- | ----- | -------------- |
| id       | [string](#string) |       | The vector ID. |
| distance | [float](#float)   |       | The distance.  |

<a name="payload-v1-Object-ID"></a>

### Object.ID

Represent the vector ID.

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| id    | [string](#string) |       |             |

<a name="payload-v1-Object-IDs"></a>

### Object.IDs

Represent multiple vector IDs.

| Field | Type              | Label    | Description |
| ----- | ----------------- | -------- | ----------- |
| ids   | [string](#string) | repeated |             |

<a name="payload-v1-Object-List"></a>

### Object.List

Represent the list object vector stream request and response.

<a name="payload-v1-Object-List-Request"></a>

### Object.List.Request

<a name="payload-v1-Object-List-Response"></a>

### Object.List.Response

| Field  | Type                                       | Label | Description           |
| ------ | ------------------------------------------ | ----- | --------------------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |       | The vector            |
| status | [google.rpc.Status](#google-rpc-Status)    |       | The RPC error status. |

<a name="payload-v1-Object-Location"></a>

### Object.Location

Represent the vector location.

| Field | Type              | Label    | Description               |
| ----- | ----------------- | -------- | ------------------------- |
| name  | [string](#string) |          | The name of the location. |
| uuid  | [string](#string) |          | The UUID of the vector.   |
| ips   | [string](#string) | repeated | The IP list.              |

<a name="payload-v1-Object-Locations"></a>

### Object.Locations

Represent multiple vector locations.

| Field     | Type                                           | Label    | Description |
| --------- | ---------------------------------------------- | -------- | ----------- |
| locations | [Object.Location](#payload-v1-Object-Location) | repeated |             |

<a name="payload-v1-Object-ReshapeVector"></a>

### Object.ReshapeVector

Represent reshape vector.

| Field  | Type            | Label    | Description        |
| ------ | --------------- | -------- | ------------------ |
| object | [bytes](#bytes) |          | The binary object. |
| shape  | [int32](#int32) | repeated | The new shape.     |

<a name="payload-v1-Object-StreamBlob"></a>

### Object.StreamBlob

Represent stream response of binary objects.

| Field  | Type                                    | Label | Description           |
| ------ | --------------------------------------- | ----- | --------------------- |
| blob   | [Object.Blob](#payload-v1-Object-Blob)  |       | The binary object.    |
| status | [google.rpc.Status](#google-rpc-Status) |       | The RPC error status. |

<a name="payload-v1-Object-StreamDistance"></a>

### Object.StreamDistance

Represent stream response of distances.

| Field    | Type                                           | Label | Description           |
| -------- | ---------------------------------------------- | ----- | --------------------- |
| distance | [Object.Distance](#payload-v1-Object-Distance) |       | The distance.         |
| status   | [google.rpc.Status](#google-rpc-Status)        |       | The RPC error status. |

<a name="payload-v1-Object-StreamLocation"></a>

### Object.StreamLocation

Represent the stream response of the vector location.

| Field    | Type                                           | Label | Description           |
| -------- | ---------------------------------------------- | ----- | --------------------- |
| location | [Object.Location](#payload-v1-Object-Location) |       | The vector location.  |
| status   | [google.rpc.Status](#google-rpc-Status)        |       | The RPC error status. |

<a name="payload-v1-Object-StreamVector"></a>

### Object.StreamVector

Represent stream response of the vector.

| Field  | Type                                       | Label | Description           |
| ------ | ------------------------------------------ | ----- | --------------------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |       | The vector.           |
| status | [google.rpc.Status](#google-rpc-Status)    |       | The RPC error status. |

<a name="payload-v1-Object-Timestamp"></a>

### Object.Timestamp

Represent a vector meta data.

| Field     | Type              | Label | Description                                     |
| --------- | ----------------- | ----- | ----------------------------------------------- |
| id        | [string](#string) |       | The vector ID.                                  |
| timestamp | [int64](#int64)   |       | timestamp represents when this vector inserted. |

<a name="payload-v1-Object-TimestampRequest"></a>

### Object.TimestampRequest

Represent a request to fetch vector meta data.

| Field | Type                               | Label | Description                  |
| ----- | ---------------------------------- | ----- | ---------------------------- |
| id    | [Object.ID](#payload-v1-Object-ID) |       | The vector ID to be fetched. |

<a name="payload-v1-Object-Vector"></a>

### Object.Vector

Represent a vector.

| Field     | Type              | Label    | Description                                     |
| --------- | ----------------- | -------- | ----------------------------------------------- |
| id        | [string](#string) |          | The vector ID.                                  |
| vector    | [float](#float)   | repeated | The vector.                                     |
| timestamp | [int64](#int64)   |          | timestamp represents when this vector inserted. |

<a name="payload-v1-Object-VectorRequest"></a>

### Object.VectorRequest

Represent a request to fetch raw vector.

| Field   | Type                                       | Label | Description                  |
| ------- | ------------------------------------------ | ----- | ---------------------------- |
| id      | [Object.ID](#payload-v1-Object-ID)         |       | The vector ID to be fetched. |
| filters | [Filter.Config](#payload-v1-Filter-Config) |       | Filter configurations.       |

<a name="payload-v1-Object-Vectors"></a>

### Object.Vectors

Represent multiple vectors.

| Field   | Type                                       | Label    | Description |
| ------- | ------------------------------------------ | -------- | ----------- |
| vectors | [Object.Vector](#payload-v1-Object-Vector) | repeated |             |

<a name="payload-v1-Remove"></a>

### Remove

Remove related messages.

<a name="payload-v1-Remove-Config"></a>

### Remove.Config

Represent the remove configuration.

| Field                   | Type            | Label | Description                                         |
| ----------------------- | --------------- | ----- | --------------------------------------------------- |
| skip_strict_exist_check | [bool](#bool)   |       | A flag to skip exist check during upsert operation. |
| timestamp               | [int64](#int64) |       | Remove timestamp.                                   |

<a name="payload-v1-Remove-MultiRequest"></a>

### Remove.MultiRequest

Represent the multiple remove request.

| Field    | Type                                         | Label    | Description                                    |
| -------- | -------------------------------------------- | -------- | ---------------------------------------------- |
| requests | [Remove.Request](#payload-v1-Remove-Request) | repeated | Represent the multiple remove request content. |

<a name="payload-v1-Remove-Request"></a>

### Remove.Request

Represent the remove request.

| Field  | Type                                       | Label | Description                              |
| ------ | ------------------------------------------ | ----- | ---------------------------------------- |
| id     | [Object.ID](#payload-v1-Object-ID)         |       | The object ID to be removed.             |
| config | [Remove.Config](#payload-v1-Remove-Config) |       | The configuration of the remove request. |

<a name="payload-v1-Remove-Timestamp"></a>

### Remove.Timestamp

Represent the timestamp comparison.

| Field     | Type                                                               | Label | Description               |
| --------- | ------------------------------------------------------------------ | ----- | ------------------------- |
| timestamp | [int64](#int64)                                                    |       | The timestamp.            |
| operator  | [Remove.Timestamp.Operator](#payload-v1-Remove-Timestamp-Operator) |       | The conditional operator. |

<a name="payload-v1-Remove-TimestampRequest"></a>

### Remove.TimestampRequest

Represent the remove request based on timestamp.

| Field      | Type                                             | Label    | Description                                                                                |
| ---------- | ------------------------------------------------ | -------- | ------------------------------------------------------------------------------------------ |
| timestamps | [Remove.Timestamp](#payload-v1-Remove-Timestamp) | repeated | The timestamp comparison list. If more than one is specified, the `AND` search is applied. |

<a name="payload-v1-Search"></a>

### Search

Search related messages.

<a name="payload-v1-Search-Config"></a>

### Search.Config

Represent search configuration.

| Field                 | Type                                                                   | Label | Description                                  |
| --------------------- | ---------------------------------------------------------------------- | ----- | -------------------------------------------- |
| request_id            | [string](#string)                                                      |       | Unique request ID.                           |
| num                   | [uint32](#uint32)                                                      |       | Maximum number of result to be returned.     |
| radius                | [float](#float)                                                        |       | Search radius.                               |
| epsilon               | [float](#float)                                                        |       | Search coefficient.                          |
| timeout               | [int64](#int64)                                                        |       | Search timeout in nanoseconds.               |
| ingress_filters       | [Filter.Config](#payload-v1-Filter-Config)                             |       | Ingress filter configurations.               |
| egress_filters        | [Filter.Config](#payload-v1-Filter-Config)                             |       | Egress filter configurations.                |
| min_num               | [uint32](#uint32)                                                      |       | Minimum number of result to be returned.     |
| aggregation_algorithm | [Search.AggregationAlgorithm](#payload-v1-Search-AggregationAlgorithm) |       | Aggregation Algorithm                        |
| ratio                 | [google.protobuf.FloatValue](#google-protobuf-FloatValue)              |       | Search ratio for agent return result number. |
| nprobe                | [uint32](#uint32)                                                      |       | Search nprobe.                               |

<a name="payload-v1-Search-IDRequest"></a>

### Search.IDRequest

Represent a search by ID request.

| Field  | Type                                       | Label | Description                              |
| ------ | ------------------------------------------ | ----- | ---------------------------------------- |
| id     | [string](#string)                          |       | The vector ID to be searched.            |
| config | [Search.Config](#payload-v1-Search-Config) |       | The configuration of the search request. |

<a name="payload-v1-Search-MultiIDRequest"></a>

### Search.MultiIDRequest

Represent the multiple search by ID request.

| Field    | Type                                             | Label    | Description                                          |
| -------- | ------------------------------------------------ | -------- | ---------------------------------------------------- |
| requests | [Search.IDRequest](#payload-v1-Search-IDRequest) | repeated | Represent the multiple search by ID request content. |

<a name="payload-v1-Search-MultiObjectRequest"></a>

### Search.MultiObjectRequest

Represent the multiple search by binary object request.

| Field    | Type                                                     | Label    | Description                                                     |
| -------- | -------------------------------------------------------- | -------- | --------------------------------------------------------------- |
| requests | [Search.ObjectRequest](#payload-v1-Search-ObjectRequest) | repeated | Represent the multiple search by binary object request content. |

<a name="payload-v1-Search-MultiRequest"></a>

### Search.MultiRequest

Represent the multiple search request.

| Field    | Type                                         | Label    | Description                                    |
| -------- | -------------------------------------------- | -------- | ---------------------------------------------- |
| requests | [Search.Request](#payload-v1-Search-Request) | repeated | Represent the multiple search request content. |

<a name="payload-v1-Search-ObjectRequest"></a>

### Search.ObjectRequest

Represent a search by binary object request.

| Field      | Type                                       | Label | Description                              |
| ---------- | ------------------------------------------ | ----- | ---------------------------------------- |
| object     | [bytes](#bytes)                            |       | The binary object to be searched.        |
| config     | [Search.Config](#payload-v1-Search-Config) |       | The configuration of the search request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |       | Filter configuration.                    |

<a name="payload-v1-Search-Request"></a>

### Search.Request

Represent a search request.

| Field  | Type                                       | Label    | Description                              |
| ------ | ------------------------------------------ | -------- | ---------------------------------------- |
| vector | [float](#float)                            | repeated | The vector to be searched.               |
| config | [Search.Config](#payload-v1-Search-Config) |          | The configuration of the search request. |

<a name="payload-v1-Search-Response"></a>

### Search.Response

Represent a search response.

| Field      | Type                                           | Label    | Description            |
| ---------- | ---------------------------------------------- | -------- | ---------------------- |
| request_id | [string](#string)                              |          | The unique request ID. |
| results    | [Object.Distance](#payload-v1-Object-Distance) | repeated | Search results.        |

<a name="payload-v1-Search-Responses"></a>

### Search.Responses

Represent multiple search responses.

| Field     | Type                                           | Label    | Description                                     |
| --------- | ---------------------------------------------- | -------- | ----------------------------------------------- |
| responses | [Search.Response](#payload-v1-Search-Response) | repeated | Represent the multiple search response content. |

<a name="payload-v1-Search-StreamResponse"></a>

### Search.StreamResponse

Represent stream search response.

| Field    | Type                                           | Label | Description                    |
| -------- | ---------------------------------------------- | ----- | ------------------------------ |
| response | [Search.Response](#payload-v1-Search-Response) |       | Represent the search response. |
| status   | [google.rpc.Status](#google-rpc-Status)        |       | The RPC error status.          |

<a name="payload-v1-Update"></a>

### Update

Update related messages

<a name="payload-v1-Update-Config"></a>

### Update.Config

Represent the update configuration.

| Field                   | Type                                       | Label | Description                                                                                      |
| ----------------------- | ------------------------------------------ | ----- | ------------------------------------------------------------------------------------------------ |
| skip_strict_exist_check | [bool](#bool)                              |       | A flag to skip exist check during update operation.                                              |
| filters                 | [Filter.Config](#payload-v1-Filter-Config) |       | Filter configuration.                                                                            |
| timestamp               | [int64](#int64)                            |       | Update timestamp.                                                                                |
| disable_balanced_update | [bool](#bool)                              |       | A flag to disable balanced update (split remove -&gt; insert operation) during update operation. |

<a name="payload-v1-Update-MultiObjectRequest"></a>

### Update.MultiObjectRequest

Represent the multiple update binary object request.

| Field    | Type                                                     | Label    | Description                                           |
| -------- | -------------------------------------------------------- | -------- | ----------------------------------------------------- |
| requests | [Update.ObjectRequest](#payload-v1-Update-ObjectRequest) | repeated | Represent the multiple update object request content. |

<a name="payload-v1-Update-MultiRequest"></a>

### Update.MultiRequest

Represent the multiple update request.

| Field    | Type                                         | Label    | Description                                    |
| -------- | -------------------------------------------- | -------- | ---------------------------------------------- |
| requests | [Update.Request](#payload-v1-Update-Request) | repeated | Represent the multiple update request content. |

<a name="payload-v1-Update-ObjectRequest"></a>

### Update.ObjectRequest

Represent the update binary object request.

| Field      | Type                                       | Label | Description                              |
| ---------- | ------------------------------------------ | ----- | ---------------------------------------- |
| object     | [Object.Blob](#payload-v1-Object-Blob)     |       | The binary object to be updated.         |
| config     | [Update.Config](#payload-v1-Update-Config) |       | The configuration of the update request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |       | Filter target.                           |

<a name="payload-v1-Update-Request"></a>

### Update.Request

Represent the update request.

| Field  | Type                                       | Label | Description                              |
| ------ | ------------------------------------------ | ----- | ---------------------------------------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |       | The vector to be updated.                |
| config | [Update.Config](#payload-v1-Update-Config) |       | The configuration of the update request. |

<a name="payload-v1-Update-TimestampRequest"></a>

### Update.TimestampRequest

Represent a vector meta data.

| Field     | Type              | Label | Description                                       |
| --------- | ----------------- | ----- | ------------------------------------------------- |
| id        | [string](#string) |       | The vector ID.                                    |
| timestamp | [int64](#int64)   |       | timestamp represents when this vector inserted.   |
| force     | [bool](#bool)     |       | force represents forcefully update the timestamp. |

<a name="payload-v1-Upsert"></a>

### Upsert

Upsert related messages.

<a name="payload-v1-Upsert-Config"></a>

### Upsert.Config

Represent the upsert configuration.

| Field                   | Type                                       | Label | Description                                                                                      |
| ----------------------- | ------------------------------------------ | ----- | ------------------------------------------------------------------------------------------------ |
| skip_strict_exist_check | [bool](#bool)                              |       | A flag to skip exist check during upsert operation.                                              |
| filters                 | [Filter.Config](#payload-v1-Filter-Config) |       | Filter configuration.                                                                            |
| timestamp               | [int64](#int64)                            |       | Upsert timestamp.                                                                                |
| disable_balanced_update | [bool](#bool)                              |       | A flag to disable balanced update (split remove -&gt; insert operation) during update operation. |

<a name="payload-v1-Upsert-MultiObjectRequest"></a>

### Upsert.MultiObjectRequest

Represent the multiple upsert binary object request.

| Field    | Type                                                     | Label    | Description                                           |
| -------- | -------------------------------------------------------- | -------- | ----------------------------------------------------- |
| requests | [Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest) | repeated | Represent the multiple upsert object request content. |

<a name="payload-v1-Upsert-MultiRequest"></a>

### Upsert.MultiRequest

Represent mthe ultiple upsert request.

| Field    | Type                                         | Label    | Description                                    |
| -------- | -------------------------------------------- | -------- | ---------------------------------------------- |
| requests | [Upsert.Request](#payload-v1-Upsert-Request) | repeated | Represent the multiple upsert request content. |

<a name="payload-v1-Upsert-ObjectRequest"></a>

### Upsert.ObjectRequest

Represent the upsert binary object request.

| Field      | Type                                       | Label | Description                              |
| ---------- | ------------------------------------------ | ----- | ---------------------------------------- |
| object     | [Object.Blob](#payload-v1-Object-Blob)     |       | The binary object to be upserted.        |
| config     | [Upsert.Config](#payload-v1-Upsert-Config) |       | The configuration of the upsert request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |       | Filter target.                           |

<a name="payload-v1-Upsert-Request"></a>

### Upsert.Request

Represent the upsert request.

| Field  | Type                                       | Label | Description                              |
| ------ | ------------------------------------------ | ----- | ---------------------------------------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |       | The vector to be upserted.               |
| config | [Upsert.Config](#payload-v1-Upsert-Config) |       | The configuration of the upsert request. |

<a name="payload-v1-Remove-Timestamp-Operator"></a>

### Remove.Timestamp.Operator

Operator is enum of each conditional operator.

| Name | Number | Description                                                                   |
| ---- | ------ | ----------------------------------------------------------------------------- |
| Eq   | 0      | The timestamp is equal to the specified value in the request.                 |
| Ne   | 1      | The timestamp is not equal to the specified value in the request.             |
| Ge   | 2      | The timestamp is greater than or equal to the specified value in the request. |
| Gt   | 3      | The timestamp is greater than the specified value in the request.             |
| Le   | 4      | The timestamp is less than or equal to the specified value in the request.    |
| Lt   | 5      | The timestamp is less than the specified value in the request.                |

<a name="payload-v1-Search-AggregationAlgorithm"></a>

### Search.AggregationAlgorithm

AggregationAlgorithm is enum of each aggregation algorithms

| Name            | Number | Description |
| --------------- | ------ | ----------- |
| Unknown         | 0      |             |
| ConcurrentQueue | 1      |             |
| SortSlice       | 2      |             |
| SortPoolSlice   | 3      |             |
| PairingHeap     | 4      |             |

<a name="v1_agent_core_agent-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/agent/core/agent.proto

<a name="core-v1-Agent"></a>

### Agent

Represent the agent service.

| Method Name        | Request Type                                                                     | Response Type                          | Description                                  |
| ------------------ | -------------------------------------------------------------------------------- | -------------------------------------- | -------------------------------------------- |
| CreateIndex        | [.payload.v1.Control.CreateIndexRequest](#payload-v1-Control-CreateIndexRequest) | [.payload.v1.Empty](#payload-v1-Empty) | Represent the creating index RPC.            |
| SaveIndex          | [.payload.v1.Empty](#payload-v1-Empty)                                           | [.payload.v1.Empty](#payload-v1-Empty) | Represent the saving index RPC.              |
| CreateAndSaveIndex | [.payload.v1.Control.CreateIndexRequest](#payload-v1-Control-CreateIndexRequest) | [.payload.v1.Empty](#payload-v1-Empty) | Represent the creating and saving index RPC. |

<a name="v1_agent_sidecar_sidecar-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/agent/sidecar/sidecar.proto

<a name="sidecar-v1-Sidecar"></a>

### Sidecar

Represent the agent sidecar service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ----------- |

<a name="v1_discoverer_discoverer-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/discoverer/discoverer.proto

<a name="discoverer-v1-Discoverer"></a>

### Discoverer

Represent the discoverer service.

| Method Name | Request Type                                                     | Response Type                                          | Description                                               |
| ----------- | ---------------------------------------------------------------- | ------------------------------------------------------ | --------------------------------------------------------- |
| Pods        | [.payload.v1.Discoverer.Request](#payload-v1-Discoverer-Request) | [.payload.v1.Info.Pods](#payload-v1-Info-Pods)         | Represent the RPC to get the agent pods information.      |
| Nodes       | [.payload.v1.Discoverer.Request](#payload-v1-Discoverer-Request) | [.payload.v1.Info.Nodes](#payload-v1-Info-Nodes)       | Represent the RPC to get the node information.            |
| Services    | [.payload.v1.Discoverer.Request](#payload-v1-Discoverer-Request) | [.payload.v1.Info.Services](#payload-v1-Info-Services) | Represent the RPC to get the readreplica svc information. |

<a name="v1_filter_egress_egress_filter-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/filter/egress/egress_filter.proto

<a name="filter-egress-v1-Filter"></a>

### Filter

Represent the egress filter service.

| Method Name    | Request Type                                               | Response Type                                              | Description                               |
| -------------- | ---------------------------------------------------------- | ---------------------------------------------------------- | ----------------------------------------- |
| FilterDistance | [.payload.v1.Object.Distance](#payload-v1-Object-Distance) | [.payload.v1.Object.Distance](#payload-v1-Object-Distance) | Represent the RPC to filter the distance. |
| FilterVector   | [.payload.v1.Object.Vector](#payload-v1-Object-Vector)     | [.payload.v1.Object.Vector](#payload-v1-Object-Vector)     | Represent the RPC to filter the vector.   |

<a name="v1_filter_ingress_ingress_filter-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/filter/ingress/ingress_filter.proto

<a name="filter-ingress-v1-Filter"></a>

### Filter

Represent the ingress filter service.

| Method Name  | Request Type                                           | Response Type                                          | Description                               |
| ------------ | ------------------------------------------------------ | ------------------------------------------------------ | ----------------------------------------- |
| GenVector    | [.payload.v1.Object.Blob](#payload-v1-Object-Blob)     | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | Represent the RPC to generate the vector. |
| FilterVector | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | Represent the RPC to filter the vector.   |

<a name="v1_meta_meta-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/meta/meta.proto

<a name="meta-v1-Meta"></a>

### Meta

| Method Name | Request Type                                           | Response Type                                    | Description |
| ----------- | ------------------------------------------------------ | ------------------------------------------------ | ----------- |
| Get         | [.payload.v1.Meta.Key](#payload-v1-Meta-Key)           | [.payload.v1.Meta.Value](#payload-v1-Meta-Value) |             |
| Set         | [.payload.v1.Meta.KeyValue](#payload-v1-Meta-KeyValue) | [.payload.v1.Empty](#payload-v1-Empty)           |             |
| Delete      | [.payload.v1.Meta.Key](#payload-v1-Meta-Key)           | [.payload.v1.Empty](#payload-v1-Empty)           |             |

<a name="v1_mirror_mirror-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/mirror/mirror.proto

<a name="mirror-v1-Mirror"></a>

### Mirror

Overview
Mirror Service is responsible for providing the `Register` interface for the Vald Mirror Gateway.

| Method Name | Request Type                                             | Response Type                                            | Description                                                                                        |
| ----------- | -------------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | -------- | --- |
| Register    | [.payload.v1.Mirror.Targets](#payload-v1-Mirror-Targets) | [.payload.v1.Mirror.Targets](#payload-v1-Mirror-Targets) | Overview Register RPC is the method to register other Vald Mirror Gateway targets. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 13  | INTERNAL |     |

<a name="v1_rpc_errdetails_error_details-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/rpc/errdetails/error_details.proto

<a name="rpc-v1-BadRequest"></a>

### BadRequest

Describes violations in a client request. This error type focuses on the
syntactic aspects of the request.

| Field            | Type                                                           | Label    | Description                                   |
| ---------------- | -------------------------------------------------------------- | -------- | --------------------------------------------- |
| field_violations | [BadRequest.FieldViolation](#rpc-v1-BadRequest-FieldViolation) | repeated | Describes all violations in a client request. |

<a name="rpc-v1-BadRequest-FieldViolation"></a>

### BadRequest.FieldViolation

A message type used to describe a single bad request field.

| Field | Type              | Label | Description                                                                                                                                        |
| ----- | ----------------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| field | [string](#string) |       | A path that leads to a field in the request body. The value will be a sequence of dot-separated identifiers that identify a protocol buffer field. |

Consider the following:

message CreateContactRequest { message EmailAddress { enum Type { TYPE_UNSPECIFIED = 0; HOME = 1; WORK = 2; }

optional string email = 1; repeated EmailType type = 2; }

string full_name = 1; repeated EmailAddress email_addresses = 2; }

In this example, in proto `field` could take one of the following values:

- `full_name` for a violation in the `full_name` value _ `email_addresses[1].email` for a violation in the `email` field of the first `email_addresses` message _ `email_addresses[3].type[2]` for a violation in the second `type` value in the third `email_addresses` message.

In JSON, the same values are represented as:

- `fullName` for a violation in the `fullName` value _ `emailAddresses[1].email` for a violation in the `email` field of the first `emailAddresses` message _ `emailAddresses[3].type[2]` for a violation in the second `type` value in the third `emailAddresses` message. |
  | description | [string](#string) | | A description of why the request element is bad. |

<a name="rpc-v1-DebugInfo"></a>

### DebugInfo

Describes additional debugging info.

| Field         | Type              | Label    | Description                                                  |
| ------------- | ----------------- | -------- | ------------------------------------------------------------ |
| stack_entries | [string](#string) | repeated | The stack trace entries indicating where the error occurred. |
| detail        | [string](#string) |          | Additional debugging information provided by the server.     |

<a name="rpc-v1-ErrorInfo"></a>

### ErrorInfo

Describes the cause of the error with structured details.

Example of an error when contacting the &#34;pubsub.googleapis.com&#34; API when it
is not enabled:

    { &#34;reason&#34;: &#34;API_DISABLED&#34;
      &#34;domain&#34;: &#34;googleapis.com&#34;
      &#34;metadata&#34;: {
        &#34;resource&#34;: &#34;projects/123&#34;,
        &#34;service&#34;: &#34;pubsub.googleapis.com&#34;
      }
    }

This response indicates that the pubsub.googleapis.com API is not enabled.

Example of an error that is returned when attempting to create a Spanner
instance in a region that is out of stock:

    { &#34;reason&#34;: &#34;STOCKOUT&#34;
      &#34;domain&#34;: &#34;spanner.googleapis.com&#34;,
      &#34;metadata&#34;: {
        &#34;availableRegions&#34;: &#34;us-central1,us-east2&#34;
      }
    }

| Field    | Type                                                       | Label    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| -------- | ---------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| reason   | [string](#string)                                          |          | The reason of the error. This is a constant value that identifies the proximate cause of the error. Error reasons are unique within a particular domain of errors. This should be at most 63 characters and match a regular expression of `[A-Z][A-Z0-9_]&#43;[A-Z0-9]`, which represents UPPER_SNAKE_CASE.                                                                                                                                 |
| domain   | [string](#string)                                          |          | The logical grouping to which the &#34;reason&#34; belongs. The error domain is typically the registered service name of the tool or product that generates the error. Example: &#34;pubsub.googleapis.com&#34;. If the error is generated by some common infrastructure, the error domain must be a globally unique value that identifies the infrastructure. For Google API infrastructure, the error domain is &#34;googleapis.com&#34;. |
| metadata | [ErrorInfo.MetadataEntry](#rpc-v1-ErrorInfo-MetadataEntry) | repeated | Additional structured details about this error.                                                                                                                                                                                                                                                                                                                                                                                             |

Keys should match /[a-zA-Z0-9-_]/ and be limited to 64 characters in length. When identifying the current value of an exceeded limit, the units should be contained in the key, not the value. For example, rather than {&#34;instanceLimit&#34;: &#34;100/request&#34;}, should be returned as, {&#34;instanceLimitPerRequest&#34;: &#34;100&#34;}, if the client exceeds the number of instances that can be created in a single (batch) request. |

<a name="rpc-v1-ErrorInfo-MetadataEntry"></a>

### ErrorInfo.MetadataEntry

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| key   | [string](#string) |       |             |
| value | [string](#string) |       |             |

<a name="rpc-v1-Help"></a>

### Help

Provides links to documentation or for performing an out of band action.

For example, if a quota check failed with an error indicating the calling
project hasn&#39;t enabled the accessed service, this can contain a URL pointing
directly to the right place in the developer console to flip the bit.

| Field | Type                           | Label    | Description                                                              |
| ----- | ------------------------------ | -------- | ------------------------------------------------------------------------ |
| links | [Help.Link](#rpc-v1-Help-Link) | repeated | URL(s) pointing to additional information on handling the current error. |

<a name="rpc-v1-Help-Link"></a>

### Help.Link

Describes a URL link.

| Field       | Type              | Label | Description                     |
| ----------- | ----------------- | ----- | ------------------------------- |
| description | [string](#string) |       | Describes what the link offers. |
| url         | [string](#string) |       | The URL of the link.            |

<a name="rpc-v1-LocalizedMessage"></a>

### LocalizedMessage

Provides a localized error message that is safe to return to the user
which can be attached to an RPC error.

| Field   | Type              | Label | Description                                                                                                                                                          |
| ------- | ----------------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| locale  | [string](#string) |       | The locale used following the specification defined at https://www.rfc-editor.org/rfc/bcp/bcp47.txt. Examples are: &#34;en-US&#34;, &#34;fr-CH&#34;, &#34;es-MX&#34; |
| message | [string](#string) |       | The localized error message in the above locale.                                                                                                                     |

<a name="rpc-v1-PreconditionFailure"></a>

### PreconditionFailure

Describes what preconditions have failed.

For example, if an RPC failed because it required the Terms of Service to be
acknowledged, it could list the terms of service violation in the
PreconditionFailure message.

| Field      | Type                                                                   | Label    | Description                            |
| ---------- | ---------------------------------------------------------------------- | -------- | -------------------------------------- |
| violations | [PreconditionFailure.Violation](#rpc-v1-PreconditionFailure-Violation) | repeated | Describes all precondition violations. |

<a name="rpc-v1-PreconditionFailure-Violation"></a>

### PreconditionFailure.Violation

A message type used to describe a single precondition failure.

| Field       | Type              | Label | Description                                                                                                                                                                                                    |
| ----------- | ----------------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| type        | [string](#string) |       | The type of PreconditionFailure. We recommend using a service-specific enum type to define the supported precondition violation subjects. For example, &#34;TOS&#34; for &#34;Terms of Service violation&#34;. |
| subject     | [string](#string) |       | The subject, relative to the type, that failed. For example, &#34;google.com/cloud&#34; relative to the &#34;TOS&#34; type would indicate which terms of service is being referenced.                          |
| description | [string](#string) |       | A description of how the precondition failed. Developers can use this description to understand how to fix the failure.                                                                                        |

For example: &#34;Terms of service not accepted&#34;. |

<a name="rpc-v1-QuotaFailure"></a>

### QuotaFailure

Describes how a quota check failed.

For example if a daily limit was exceeded for the calling project,
a service could respond with a QuotaFailure detail containing the project
id and the description of the quota limit that was exceeded. If the
calling project hasn&#39;t enabled the service in the developer console, then
a service could respond with the project id and set `service_disabled`
to true.

Also see RetryInfo and Help types for other details about handling a
quota failure.

| Field      | Type                                                     | Label    | Description                     |
| ---------- | -------------------------------------------------------- | -------- | ------------------------------- |
| violations | [QuotaFailure.Violation](#rpc-v1-QuotaFailure-Violation) | repeated | Describes all quota violations. |

<a name="rpc-v1-QuotaFailure-Violation"></a>

### QuotaFailure.Violation

A message type used to describe a single quota violation. For example, a
daily quota or a custom quota that was exceeded.

| Field       | Type              | Label | Description                                                                                                                                                                                                                               |
| ----------- | ----------------- | ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| subject     | [string](#string) |       | The subject on which the quota check failed. For example, &#34;clientip:&lt;ip address of client&gt;&#34; or &#34;project:&lt;Google developer project id&gt;&#34;.                                                                       |
| description | [string](#string) |       | A description of how the quota check failed. Clients can use this description to find more about the quota configuration in the service&#39;s public documentation, or find the relevant quota limit to adjust through developer console. |

For example: &#34;Service disabled&#34; or &#34;Daily Limit for read operations exceeded&#34;. |

<a name="rpc-v1-RequestInfo"></a>

### RequestInfo

Contains metadata about the request that clients can attach when filing a bug
or providing other forms of feedback.

| Field        | Type              | Label | Description                                                                                                                                                |
| ------------ | ----------------- | ----- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| request_id   | [string](#string) |       | An opaque string that should only be interpreted by the service generating it. For example, it can be used to identify requests in the service&#39;s logs. |
| serving_data | [string](#string) |       | Any data that was used to serve this request. For example, an encrypted stack trace that can be sent back to the service provider for debugging.           |

<a name="rpc-v1-ResourceInfo"></a>

### ResourceInfo

Describes the resource that is being accessed.

| Field         | Type              | Label | Description                                                                                                                                                                                                                                      |
| ------------- | ----------------- | ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| resource_type | [string](#string) |       | A name for the type of resource being accessed, e.g. &#34;sql table&#34;, &#34;cloud storage bucket&#34;, &#34;file&#34;, &#34;Google calendar&#34;; or the type URL of the resource: e.g. &#34;type.googleapis.com/google.pubsub.v1.Topic&#34;. |
| resource_name | [string](#string) |       | The name of the resource being accessed. For example, a shared calendar name: &#34;example.com_4fghdhgsrgh@group.calendar.google.com&#34;, if the current error is [google.rpc.Code.PERMISSION_DENIED][google.rpc.Code.PERMISSION_DENIED].       |
| owner         | [string](#string) |       | The owner of the resource (optional). For example, &#34;user:&lt;owner email&gt;&#34; or &#34;project:&lt;Google developer project id&gt;&#34;.                                                                                                  |
| description   | [string](#string) |       | Describes what error is encountered when accessing this resource. For example, updating a cloud project may require the `writer` permission on the developer console project.                                                                    |

<a name="rpc-v1-RetryInfo"></a>

### RetryInfo

Describes when the clients can retry a failed request. Clients could ignore
the recommendation here or retry when this information is missing from error
responses.

It&#39;s always recommended that clients should use exponential backoff when
retrying.

Clients should wait until `retry_delay` amount of time has passed since
receiving the error response before retrying. If retrying requests also
fail, clients should use an exponential backoff scheme to gradually increase
the delay between retries based on `retry_delay`, until either a maximum
number of retries have been reached or a maximum retry delay cap has been
reached.

| Field       | Type                                                  | Label | Description                                                               |
| ----------- | ----------------------------------------------------- | ----- | ------------------------------------------------------------------------- |
| retry_delay | [google.protobuf.Duration](#google-protobuf-Duration) |       | Clients should wait at least this long between retrying the same request. |

<a name="v1_vald_filter-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/filter.proto

<a name="vald-v1-Filter"></a>

### Filter

Overview
Filter Server is responsible for providing insert, update, upsert and search interface for `Vald Filter Gateway`.

Vald Filter Gateway forward user request to user-defined ingress/egress filter components allowing user to run custom logic.

| Method Name        | Request Type                                                                   | Response Type                                                                 | Description                                                                                                                                                                                                                                                                                                                                                                                                             |
| ------------------ | ------------------------------------------------------------------------------ | ----------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | -------------- | --- | --- | -------- | --- |
| SearchObject       | [.payload.v1.Search.ObjectRequest](#payload-v1-Search-ObjectRequest)           | [.payload.v1.Search.Response](#payload-v1-Search-Response)                    | Overview SearchObject RPC is the method to search object(s) similar to request object. --- Status Code                                                                                                                                                                                                                                                                                                                  | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 6   | ALREADY_EXISTS |     | 13  | INTERNAL |     |
| MultiSearchObject  | [.payload.v1.Search.MultiObjectRequest](#payload-v1-Search-MultiObjectRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses)                  | Overview StreamSearchObject RPC is the method to search vectors with multi queries(objects) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server. Each Search request and response are independent. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 6   | ALREADY_EXISTS |     | 13  | INTERNAL |     |
| StreamSearchObject | [.payload.v1.Search.ObjectRequest](#payload-v1-Search-ObjectRequest) stream    | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | Overview MultiSearchObject RPC is the method to search objects with multiple objects in **1** request.                                                                                                                                                                                                                                                                                                                  |

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| InsertObject | [.payload.v1.Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview InsertObject RPC is the method to insert object through Vald Filter Gateway. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| StreamInsertObject | [.payload.v1.Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview StreamInsertObject RPC is the method to add new multiple object using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).

By using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server. Each Insert request and response are independent. It&#39;s the recommended method to insert a large number of objects. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| MultiInsertObject | [.payload.v1.Insert.MultiObjectRequest](#payload-v1-Insert-MultiObjectRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiInsertObject RPC is the method to add multiple new objects in **1** request. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| UpdateObject | [.payload.v1.Update.ObjectRequest](#payload-v1-Update-ObjectRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview UpdateObject RPC is the method to update a single vector. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| StreamUpdateObject | [.payload.v1.Update.ObjectRequest](#payload-v1-Update-ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview StreamUpdateObject RPC is the method to update multiple objects using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; By using the bidirectional streaming RPC, the update request can be communicated in any order between client and server. Each Update request and response are independent. It&#39;s the recommended method to update the large amount of objects. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| MultiUpdateObject | [.payload.v1.Update.MultiObjectRequest](#payload-v1-Update-MultiObjectRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiUpdateObject is the method to update multiple objects in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has the message size limitation.&lt;br&gt; Please be careful that the size of the request exceed the limit. &lt;/div&gt; --- Status Code

| 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| UpsertObject | [.payload.v1.Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview UpsertObject RPC is the method to update a single object and add a new single object. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| StreamUpsertObject | [.payload.v1.Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview UpsertObject RPC is the method to update a single object and add a new single object. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |
| MultiUpsertObject | [.payload.v1.Upsert.MultiObjectRequest](#payload-v1-Upsert-MultiObjectRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiUpsertObject is the method to update existing multiple objects and add new multiple objects in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 13 | INTERNAL | |

<a name="v1_vald_flush-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/flush.proto

<a name="vald-v1-Flush"></a>

### Flush

Overview
Flush Service is responsible for removing all vectors that are indexed and uncommitted in the `vald-agent`.

| Method Name | Request Type                                           | Response Type                                                | Description                                                             |
| ----------- | ------------------------------------------------------ | ------------------------------------------------------------ | ----------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | -------- | --- |
| Flush       | [.payload.v1.Flush.Request](#payload-v1-Flush-Request) | [.payload.v1.Info.Index.Count](#payload-v1-Info-Index-Count) | Overview Flush RPC is the method to remove all vectors. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 13  | INTERNAL |     |

<a name="v1_vald_index-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/index.proto

<a name="vald-v1-Index"></a>

### Index

Overview
Represent the index manager service.

| Method Name           | Request Type                           | Response Type                                                                      | Description                                                              |
| --------------------- | -------------------------------------- | ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| IndexInfo             | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.Count](#payload-v1-Info-Index-Count)                       | Overview Represent the RPC to get the index information.                 |
| IndexDetail           | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.Detail](#payload-v1-Info-Index-Detail)                     | Overview Represent the RPC to get the index information for each agents. |
| IndexStatistics       | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.Statistics](#payload-v1-Info-Index-Statistics)             | Overview Represent the RPC to get the index statistics.                  |
| IndexStatisticsDetail | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.StatisticsDetail](#payload-v1-Info-Index-StatisticsDetail) | Overview Represent the RPC to get the index statistics for each agents.  |
| IndexProperty         | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.PropertyDetail](#payload-v1-Info-Index-PropertyDetail)     | Overview Represent the RPC to get the index property.                    |

<a name="v1_vald_insert-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/insert.proto

<a name="vald-v1-Insert"></a>

### Insert

Overview
Insert Service is responsible for inserting new vectors into the `vald-agent`.

| Method Name | Request Type                                             | Response Type                                              | Description                                                                   |
| ----------- | -------------------------------------------------------- | ---------------------------------------------------------- | ----------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | -------- | ---------------------------------------------------------------------------------------------------- |
| Insert      | [.payload.v1.Insert.Request](#payload-v1-Insert-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview Insert RPC is the method to add a new single vector. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 13  | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`. |

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | ALREADY_EXISTS | Request ID is already inserted. | Change request ID. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamInsert | [.payload.v1.Insert.Request](#payload-v1-Insert-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server. Each Insert request and response are independent. It&#39;s the recommended method to insert a large number of vectors. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | ALREADY_EXISTS | Request ID is already inserted. | Change request ID. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiInsert | [.payload.v1.Insert.MultiRequest](#payload-v1-Insert-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiInsert RPC is the method to add multiple new vectors in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 6 | ALREADY_EXISTS | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | ALREADY_EXISTS | Request ID is already inserted. | Change request ID. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |

<a name="v1_vald_object-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/object.proto

<a name="vald-v1-Object"></a>

### Object

Overview
Object Service is responsible for getting inserted vectors and checking whether vectors are inserted into the `vald-agent`.

| Method Name | Request Type                                   | Response Type                                  | Description                                                                                          |
| ----------- | ---------------------------------------------- | ---------------------------------------------- | ---------------------------------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | -------- | ---------------------------------------------------------------------------------------------------- |
| Exists      | [.payload.v1.Object.ID](#payload-v1-Object-ID) | [.payload.v1.Object.ID](#payload-v1-Object-ID) | Overview Exists RPC is the method to check that a vector exists in the `vald-agent`. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 13  | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`. |

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| GetObject | [.payload.v1.Object.VectorRequest](#payload-v1-Object-VectorRequest) | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | Overview GetObject RPC is the method to get the metadata of a vector inserted into the `vald-agent`. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamGetObject | [.payload.v1.Object.VectorRequest](#payload-v1-Object-VectorRequest) stream | [.payload.v1.Object.StreamVector](#payload-v1-Object-StreamVector) stream | Overview StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server. Each Upsert request and response are independent. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamListObject | [.payload.v1.Object.List.Request](#payload-v1-Object-List-Request) | [.payload.v1.Object.List.Response](#payload-v1-Object-List-Response) stream | Overview A method to get all the vectors with server streaming --- Status Code TODO --- Troubleshooting TODO |
| GetTimestamp | [.payload.v1.Object.TimestampRequest](#payload-v1-Object-TimestampRequest) | [.payload.v1.Object.Timestamp](#payload-v1-Object-Timestamp) | Overview Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process --- Status Code TODO --- Troubleshooting TODO |

<a name="v1_vald_remove-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/remove.proto

<a name="vald-v1-Remove"></a>

### Remove

Overview
Remove Service is responsible for removing vectors indexed in the `vald-agent`.

| Method Name | Request Type                                             | Response Type                                              | Description                                                                  |
| ----------- | -------------------------------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | ------- | --- | --- | -------- | ---------------------------------------------------------------------------------------------------- |
| Remove      | [.payload.v1.Remove.Request](#payload-v1-Remove-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview Remove RPC is the method to remove a single vector. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 10  | ABORTED |     | 13  | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`. |

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| RemoveByTimestamp | [.payload.v1.Remove.TimestampRequest](#payload-v1-Remove-TimestampRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview RemoveByTimestamp RPC is the method to remove vectors based on timestamp.

&lt;div class=&#34;notice&#34;&gt; In the TimestampRequest message, the &#39;timestamps&#39; field is repeated, allowing the inclusion of multiple Timestamp.&lt;br&gt; When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.&lt;br&gt; This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | No vectors in the system match the specified timestamp conditions. | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. |
| StreamRemove | [.payload.v1.Remove.Request](#payload-v1-Remove-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview A method to remove multiple indexed vectors by bidirectional streaming.

StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server. Each Remove request and response are independent. It&#39;s the recommended method to remove a large number of vectors. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiRemove | [.payload.v1.Remove.MultiRequest](#payload-v1-Remove-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiRemove is the method to remove multiple vectors in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |

<a name="v1_vald_search-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/search.proto

<a name="vald-v1-Search"></a>

### Search

Overview
Search Service is responsible for searching vectors similar to the user request vector from `vald-agent`.

| Method Name | Request Type                                             | Response Type                                              | Description                                                                                          |
| ----------- | -------------------------------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | ------- | --- | --- | -------- | ---------------------------------------------------------------------------------------------------- |
| Search      | [.payload.v1.Search.Request](#payload-v1-Search-Request) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | Overview Search RPC is the method to search vector(s) similar to the request vector. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 10  | ABORTED |     | 13  | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`. |

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Search result is empty or insufficient to request result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| SearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | Overview SearchByID RPC is the method to search similar vectors using a user-defined vector ID.&lt;br&gt; The vector with the same requested ID should be indexed into the `vald-agent` before searching. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamSearch | [.payload.v1.Search.Request](#payload-v1-Search-Request) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | Overview StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server. Each Search request and response are independent. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Search result is empty or insufficient to request result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamSearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | Overview StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server. Each SearchByID request and response are independent. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiSearch | [.payload.v1.Search.MultiRequest](#payload-v1-Search-MultiRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | Overview MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Search result is empty or insufficient to request result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiSearchByID | [.payload.v1.Search.MultiIDRequest](#payload-v1-Search-MultiIDRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | Overview MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| LinearSearch | [.payload.v1.Search.Request](#payload-v1-Search-Request) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | Overview LinearSearch RPC is the method to linear search vector(s) similar to the request vector. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Search result is empty or insufficient to request result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| LinearSearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | Overview LinearSearchByID RPC is the method to linear search similar vectors using a user-defined vector ID.&lt;br&gt; The vector with the same requested ID should be indexed into the `vald-agent` before searching. You will get a `NOT_FOUND` error if the vector isn&#39;t stored. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamLinearSearch | [.payload.v1.Search.Request](#payload-v1-Search-Request) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | Overview StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server. Each LinearSearch request and response are independent. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Search result is empty or insufficient to request result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamLinearSearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | Overview StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server. Each LinearSearchByID request and response are independent. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiLinearSearch | [.payload.v1.Search.MultiRequest](#payload-v1-Search-MultiRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | Overview MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Search result is empty or insufficient to request result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiLinearSearchByID | [.payload.v1.Search.MultiIDRequest](#payload-v1-Search-MultiIDRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | Overview MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; // --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Requested vector&#39;s ID is empty, or some request payload is invalid. | Check request payload and fix request payload. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |

<a name="v1_vald_update-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/update.proto

<a name="vald-v1-Update"></a>

### Update

Overview
Update Service updates to new vector from inserted vector in the `vald-agent` components.

| Method Name | Request Type                                             | Response Type                                              | Description                                                                  |
| ----------- | -------------------------------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | -------------- | --- | --- | ------- | --- | --- | -------- | ---------------------------------------------------------------------------------------------------- |
| Update      | [.payload.v1.Update.Request](#payload-v1-Update-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview Update RPC is the method to update a single vector. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 6   | ALREADY_EXISTS |     | 10  | ABORTED |     | 13  | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`. |

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | ALREADY_EXISTS | Request pair of ID and vector is already inserted. | Change request ID. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamUpdate | [.payload.v1.Update.Request](#payload-v1-Update-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview StreamUpdate RPC is the method to update multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server. Each Update request and response are independent. It&#39;s the recommended method to update the large amount of vectors. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 6 | ALREADY_EXISTS | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | ALREADY_EXISTS | Request pair of ID and vector is already inserted. | Change request ID. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiUpdate | [.payload.v1.Update.MultiRequest](#payload-v1-Update-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiUpdate is the method to update multiple vectors in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 6 | ALREADY_EXISTS | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | NOT_FOUND | Requested ID is NOT inserted. | Send a request with an ID that is already inserted. | | ALREADY_EXISTS | Request pair of ID and vector is already inserted. | Change request ID. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| UpdateTimestamp | [.payload.v1.Update.TimestampRequest](#payload-v1-Update-TimestampRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview A method to update timestamp an indexed vector. --- Status Code TODO --- Troubleshooting TODO |

<a name="v1_vald_upsert-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## v1/vald/upsert.proto

<a name="vald-v1-Upsert"></a>

### Upsert

Overview
Upsert Service is responsible for updating existing vectors in the `vald-agent` or inserting new vectors into the `vald-agent` if the vector does not exist.

| Method Name | Request Type                                             | Response Type                                              | Description                                                                                                                                               |
| ----------- | -------------------------------------------------------- | ---------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | --- | --- | --- | --- | --------- | --- | --- | ---------------- | --- | --- | ----------------- | --- | --- | --------- | --- | --- | -------------- | --- | --- | ------- | --- | --- | -------- | ---------------------------------------------------------------------------------------------------- |
| Upsert      | [.payload.v1.Upsert.Request](#payload-v1-Upsert-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | Overview Upsert RPC is the method to update the inserted vector to a new single vector or add a new single vector if not inserted before. --- Status Code | 0   | OK  |     | 1   | CANCELLED |     | 3   | INVALID_ARGUMENT |     | 4   | DEADLINE_EXCEEDED |     | 5   | NOT_FOUND |     | 6   | ALREADY_EXISTS |     | 10  | ABORTED |     | 13  | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`. |

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | ALREADY_EXISTS | Requested pair of ID and vector is already inserted | Change request payload or nothing to do if update is unnecessary. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| StreamUpsert | [.payload.v1.Upsert.Request](#payload-v1-Upsert-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Overview StreamUpsert RPC is the method to update multiple existing vectors or add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).&lt;br&gt; Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server. Each Upsert request and response are independent. It’s the recommended method to upsert a large number of vectors. --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 6 | ALREADY_EXISTS | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | ALREADY_EXISTS | Requested pair of ID and vector is already inserted | Change request payload or nothing to do if update is unnecessary. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |
| MultiUpsert | [.payload.v1.Upsert.MultiRequest](#payload-v1-Upsert-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | Overview MultiUpsert is the method to update existing multiple vectors and add new multiple vectors in **1** request.

&lt;div class=&#34;notice&#34;&gt; gRPC has a message size limitation.&lt;br&gt; Please be careful that the size of the request exceeds the limit. &lt;/div&gt; --- Status Code | 0 | OK | | 1 | CANCELLED | | 3 | INVALID_ARGUMENT | | 4 | DEADLINE_EXCEEDED | | 5 | NOT_FOUND | | 6 | ALREADY_EXISTS | | 10 | ABORTED | | 13 | INTERNAL | --- Troubleshooting The request process may not be completed when the response code is NOT `0 (OK)`.

Here are some common reasons and how to resolve each error.

| name | common reason | how to resolve | | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | | CANCELLED | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed. | | INVALID_ARGUMENT | The Dimension of the request vector is NOT the same as Vald Agent&#39;s config, the requested vector&#39;s ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config. | | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side. | Check the gRPC timeout setting on both the client and server sides and fix it if needed. | | ALREADY_EXISTS | Requested pair of ID and vector is already inserted | Change request payload or nothing to do if update is unnecessary. | | INTERNAL | Target Vald cluster or network route has some critical error. | Check target Vald cluster first and check network route including ingress as second. | |

## Scalar Value Types

| .proto Type                    | Notes                                                                                                                                           | C++    | Java       | Python      | Go      | C#         | PHP            | Ruby                           |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ---------- | ----------- | ------- | ---------- | -------------- | ------------------------------ |
| <a name="double" /> double     |                                                                                                                                                 | double | double     | float       | float64 | double     | float          | Float                          |
| <a name="float" /> float       |                                                                                                                                                 | float  | float      | float       | float32 | float      | float          | Float                          |
| <a name="int32" /> int32       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="int64" /> int64       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="uint32" /> uint32     | Uses variable-length encoding.                                                                                                                  | uint32 | int        | int/long    | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64     | Uses variable-length encoding.                                                                                                                  | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32     | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s.                            | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64     | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s.                            | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="fixed32" /> fixed32   | Always four bytes. More efficient than uint32 if values are often greater than 2^28.                                                            | uint32 | int        | int         | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64   | Always eight bytes. More efficient than uint64 if values are often greater than 2^56.                                                           | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum                         |
| <a name="sfixed32" /> sfixed32 | Always four bytes.                                                                                                                              | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes.                                                                                                                             | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="bool" /> bool         |                                                                                                                                                 | bool   | boolean    | boolean     | bool    | bool       | boolean        | TrueClass/FalseClass           |
| <a name="string" /> string     | A string must always contain UTF-8 encoded or 7-bit ASCII text.                                                                                 | string | String     | str/unicode | string  | string     | string         | String (UTF-8)                 |
| <a name="bytes" /> bytes       | May contain any arbitrary sequence of bytes.                                                                                                    | string | ByteString | str         | []byte  | ByteString | string         | String (ASCII-8BIT)            |
