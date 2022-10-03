# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [apis/proto/v1/agent/core/agent.proto](#apis_proto_v1_agent_core_agent-proto)
    - [Agent](#core-v1-Agent)
  
- [apis/proto/v1/agent/sidecar/sidecar.proto](#apis_proto_v1_agent_sidecar_sidecar-proto)
    - [Sidecar](#sidecar-v1-Sidecar)
  
- [apis/proto/v1/discoverer/discoverer.proto](#apis_proto_v1_discoverer_discoverer-proto)
    - [Discoverer](#discoverer-v1-Discoverer)
  
- [apis/proto/v1/filter/egress/egress_filter.proto](#apis_proto_v1_filter_egress_egress_filter-proto)
    - [Filter](#filter-egress-v1-Filter)
  
- [apis/proto/v1/filter/ingress/ingress_filter.proto](#apis_proto_v1_filter_ingress_ingress_filter-proto)
    - [Filter](#filter-ingress-v1-Filter)
  
- [apis/proto/v1/manager/index/index_manager.proto](#apis_proto_v1_manager_index_index_manager-proto)
    - [Index](#manager-index-v1-Index)
  
- [apis/proto/v1/payload/payload.proto](#apis_proto_v1_payload_payload-proto)
    - [Control](#payload-v1-Control)
    - [Control.CreateIndexRequest](#payload-v1-Control-CreateIndexRequest)
    - [Discoverer](#payload-v1-Discoverer)
    - [Discoverer.Request](#payload-v1-Discoverer-Request)
    - [Empty](#payload-v1-Empty)
    - [Filter](#payload-v1-Filter)
    - [Filter.Config](#payload-v1-Filter-Config)
    - [Filter.Target](#payload-v1-Filter-Target)
    - [Info](#payload-v1-Info)
    - [Info.CPU](#payload-v1-Info-CPU)
    - [Info.IPs](#payload-v1-Info-IPs)
    - [Info.Index](#payload-v1-Info-Index)
    - [Info.Index.Count](#payload-v1-Info-Index-Count)
    - [Info.Index.UUID](#payload-v1-Info-Index-UUID)
    - [Info.Index.UUID.Committed](#payload-v1-Info-Index-UUID-Committed)
    - [Info.Index.UUID.Uncommitted](#payload-v1-Info-Index-UUID-Uncommitted)
    - [Info.Memory](#payload-v1-Info-Memory)
    - [Info.Node](#payload-v1-Info-Node)
    - [Info.Nodes](#payload-v1-Info-Nodes)
    - [Info.Pod](#payload-v1-Info-Pod)
    - [Info.Pods](#payload-v1-Info-Pods)
    - [Insert](#payload-v1-Insert)
    - [Insert.Config](#payload-v1-Insert-Config)
    - [Insert.MultiObjectRequest](#payload-v1-Insert-MultiObjectRequest)
    - [Insert.MultiRequest](#payload-v1-Insert-MultiRequest)
    - [Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest)
    - [Insert.Request](#payload-v1-Insert-Request)
    - [Object](#payload-v1-Object)
    - [Object.Blob](#payload-v1-Object-Blob)
    - [Object.Distance](#payload-v1-Object-Distance)
    - [Object.ID](#payload-v1-Object-ID)
    - [Object.IDs](#payload-v1-Object-IDs)
    - [Object.Location](#payload-v1-Object-Location)
    - [Object.Locations](#payload-v1-Object-Locations)
    - [Object.ReshapeVector](#payload-v1-Object-ReshapeVector)
    - [Object.StreamBlob](#payload-v1-Object-StreamBlob)
    - [Object.StreamDistance](#payload-v1-Object-StreamDistance)
    - [Object.StreamLocation](#payload-v1-Object-StreamLocation)
    - [Object.StreamVector](#payload-v1-Object-StreamVector)
    - [Object.Vector](#payload-v1-Object-Vector)
    - [Object.VectorRequest](#payload-v1-Object-VectorRequest)
    - [Object.Vectors](#payload-v1-Object-Vectors)
    - [Remove](#payload-v1-Remove)
    - [Remove.Config](#payload-v1-Remove-Config)
    - [Remove.MultiRequest](#payload-v1-Remove-MultiRequest)
    - [Remove.Request](#payload-v1-Remove-Request)
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
    - [Upsert](#payload-v1-Upsert)
    - [Upsert.Config](#payload-v1-Upsert-Config)
    - [Upsert.MultiObjectRequest](#payload-v1-Upsert-MultiObjectRequest)
    - [Upsert.MultiRequest](#payload-v1-Upsert-MultiRequest)
    - [Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest)
    - [Upsert.Request](#payload-v1-Upsert-Request)
  
- [apis/proto/v1/vald/filter.proto](#apis_proto_v1_vald_filter-proto)
    - [Filter](#vald-v1-Filter)
  
- [apis/proto/v1/vald/insert.proto](#apis_proto_v1_vald_insert-proto)
    - [Insert](#vald-v1-Insert)
  
- [apis/proto/v1/vald/object.proto](#apis_proto_v1_vald_object-proto)
    - [Object](#vald-v1-Object)
  
- [apis/proto/v1/vald/remove.proto](#apis_proto_v1_vald_remove-proto)
    - [Remove](#vald-v1-Remove)
  
- [apis/proto/v1/vald/search.proto](#apis_proto_v1_vald_search-proto)
    - [Search](#vald-v1-Search)
  
- [apis/proto/v1/vald/update.proto](#apis_proto_v1_vald_update-proto)
    - [Update](#vald-v1-Update)
  
- [apis/proto/v1/vald/upsert.proto](#apis_proto_v1_vald_upsert-proto)
    - [Upsert](#vald-v1-Upsert)
  
- [Scalar Value Types](#scalar-value-types)



<a name="apis_proto_v1_agent_core_agent-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/agent/core/agent.proto


 

 

 


<a name="core-v1-Agent"></a>

### Agent
Represent the agent service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateIndex | [.payload.v1.Control.CreateIndexRequest](#payload-v1-Control-CreateIndexRequest) | [.payload.v1.Empty](#payload-v1-Empty) | Represent the create index RPC. |
| SaveIndex | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Empty](#payload-v1-Empty) | Represent the save index RPC. |
| CreateAndSaveIndex | [.payload.v1.Control.CreateIndexRequest](#payload-v1-Control-CreateIndexRequest) | [.payload.v1.Empty](#payload-v1-Empty) | Represent the create and save index RPC. |
| IndexInfo | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.Count](#payload-v1-Info-Index-Count) | Represent the RPC to get the agent index information. |

 



<a name="apis_proto_v1_agent_sidecar_sidecar-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/agent/sidecar/sidecar.proto


 

 

 


<a name="sidecar-v1-Sidecar"></a>

### Sidecar
Represent the agent sidecar service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|

 



<a name="apis_proto_v1_discoverer_discoverer-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/discoverer/discoverer.proto


 

 

 


<a name="discoverer-v1-Discoverer"></a>

### Discoverer
Represent the discoverer service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Pods | [.payload.v1.Discoverer.Request](#payload-v1-Discoverer-Request) | [.payload.v1.Info.Pods](#payload-v1-Info-Pods) | Represent the RPC to get the agent pods information. |
| Nodes | [.payload.v1.Discoverer.Request](#payload-v1-Discoverer-Request) | [.payload.v1.Info.Nodes](#payload-v1-Info-Nodes) | Represent the RPC to get the node information. |

 



<a name="apis_proto_v1_filter_egress_egress_filter-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/filter/egress/egress_filter.proto


 

 

 


<a name="filter-egress-v1-Filter"></a>

### Filter
Represent the egress filter service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| FilterDistance | [.payload.v1.Object.Distance](#payload-v1-Object-Distance) | [.payload.v1.Object.Distance](#payload-v1-Object-Distance) | Represent the RPC to filter the distance. |
| FilterVector | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | Represent the RPC to filter the vector. |

 



<a name="apis_proto_v1_filter_ingress_ingress_filter-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/filter/ingress/ingress_filter.proto


 

 

 


<a name="filter-ingress-v1-Filter"></a>

### Filter
Represent the ingress filter service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GenVector | [.payload.v1.Object.Blob](#payload-v1-Object-Blob) | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | Represent the RPC to generate the vector. |
| FilterVector | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | Represent the RPC to filter the vector. |

 



<a name="apis_proto_v1_manager_index_index_manager-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/manager/index/index_manager.proto


 

 

 


<a name="manager-index-v1-Index"></a>

### Index
Represent the index manager service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| IndexInfo | [.payload.v1.Empty](#payload-v1-Empty) | [.payload.v1.Info.Index.Count](#payload-v1-Info-Index-Count) | Represent the RPC to get the index information. |

 



<a name="apis_proto_v1_payload_payload-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/payload/payload.proto



<a name="payload-v1-Control"></a>

### Control
Control related messages.






<a name="payload-v1-Control-CreateIndexRequest"></a>

### Control.CreateIndexRequest
Represent the create index request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_size | [uint32](#uint32) |  | The pool size of the create index operation. |






<a name="payload-v1-Discoverer"></a>

### Discoverer
Discoverer related messages.






<a name="payload-v1-Discoverer-Request"></a>

### Discoverer.Request
Represent the dicoverer request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The agent name to be discover. |
| namespace | [string](#string) |  | The namespace to be discover. |
| node | [string](#string) |  | The node to be discover. |






<a name="payload-v1-Empty"></a>

### Empty
Represent an empty message.






<a name="payload-v1-Filter"></a>

### Filter
Filter related messages.






<a name="payload-v1-Filter-Config"></a>

### Filter.Config
Represent filter configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| targets | [Filter.Target](#payload-v1-Filter-Target) | repeated | Represent the filter target configuration. |






<a name="payload-v1-Filter-Target"></a>

### Filter.Target
Represent the target filter server.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | The target hostname. |
| port | [uint32](#uint32) |  | The target port. |






<a name="payload-v1-Info"></a>

### Info
Info related messages.






<a name="payload-v1-Info-CPU"></a>

### Info.CPU
Represent the CPU information message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [double](#double) |  | The CPU resource limit. |
| request | [double](#double) |  | The CPU resource requested. |
| usage | [double](#double) |  | The CPU usage. |






<a name="payload-v1-Info-IPs"></a>

### Info.IPs
Represent the multiple IP message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) | repeated |  |






<a name="payload-v1-Info-Index"></a>

### Info.Index
Represent the index information messages.






<a name="payload-v1-Info-Index-Count"></a>

### Info.Index.Count
Represent the index count message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stored | [uint32](#uint32) |  | The stored index count. |
| uncommitted | [uint32](#uint32) |  | The uncommitted index count. |
| indexing | [bool](#bool) |  | The indexing index count. |
| saving | [bool](#bool) |  | The saving index count. |






<a name="payload-v1-Info-Index-UUID"></a>

### Info.Index.UUID
Represent the UUID message.






<a name="payload-v1-Info-Index-UUID-Committed"></a>

### Info.Index.UUID.Committed
The committed UUID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload-v1-Info-Index-UUID-Uncommitted"></a>

### Info.Index.UUID.Uncommitted
The uncommitted UUID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload-v1-Info-Memory"></a>

### Info.Memory
Represent the memory information message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [double](#double) |  | The memory limit. |
| request | [double](#double) |  | The memory requested. |
| usage | [double](#double) |  | The memory usage. |






<a name="payload-v1-Info-Node"></a>

### Info.Node
Represent the node information message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the node. |
| internal_addr | [string](#string) |  | The internal IP address of the node. |
| external_addr | [string](#string) |  | The external IP address of the node. |
| cpu | [Info.CPU](#payload-v1-Info-CPU) |  | The CPU information of the node. |
| memory | [Info.Memory](#payload-v1-Info-Memory) |  | The memory information of the node. |
| Pods | [Info.Pods](#payload-v1-Info-Pods) |  | The pod information of the node. |






<a name="payload-v1-Info-Nodes"></a>

### Info.Nodes
Represent the multiple node information message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nodes | [Info.Node](#payload-v1-Info-Node) | repeated | The multiple node information. |






<a name="payload-v1-Info-Pod"></a>

### Info.Pod
Represent the pod information message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app_name | [string](#string) |  | The app name of the pod on the label. |
| name | [string](#string) |  | The name of the pod. |
| namespace | [string](#string) |  | The namespace of the pod. |
| ip | [string](#string) |  | The IP of the pod. |
| cpu | [Info.CPU](#payload-v1-Info-CPU) |  | The CPU information of the pod. |
| memory | [Info.Memory](#payload-v1-Info-Memory) |  | The memory information of the pod. |
| node | [Info.Node](#payload-v1-Info-Node) |  | The node information of the pod. |






<a name="payload-v1-Info-Pods"></a>

### Info.Pods
Represent the multiple pod information message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pods | [Info.Pod](#payload-v1-Info-Pod) | repeated | The multiple pod information. |






<a name="payload-v1-Insert"></a>

### Insert
Insert related messages.






<a name="payload-v1-Insert-Config"></a>

### Insert.Config
Represent insert configurations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  | A flag to skip exist check during insert operation. |
| filters | [Filter.Config](#payload-v1-Filter-Config) |  | Filter configurations. |
| timestamp | [int64](#int64) |  | Insert timestamp. |






<a name="payload-v1-Insert-MultiObjectRequest"></a>

### Insert.MultiObjectRequest
Represent the multiple insert by binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest) | repeated | Represent multiple insert by object content. |






<a name="payload-v1-Insert-MultiRequest"></a>

### Insert.MultiRequest
Represent the multiple insert request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Insert.Request](#payload-v1-Insert-Request) | repeated | Represent multiple insert request content. |






<a name="payload-v1-Insert-ObjectRequest"></a>

### Insert.ObjectRequest
Represent the insert by binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object.Blob](#payload-v1-Object-Blob) |  | The binary object to be inserted. |
| config | [Insert.Config](#payload-v1-Insert-Config) |  | The configuration of the insert request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |  | Filter configurations. |






<a name="payload-v1-Insert-Request"></a>

### Insert.Request
Represent the insert request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |  | The vector to be inserted. |
| config | [Insert.Config](#payload-v1-Insert-Config) |  | The configuration of the insert request. |






<a name="payload-v1-Object"></a>

### Object
Common messages.






<a name="payload-v1-Object-Blob"></a>

### Object.Blob
Represent the binary object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The object ID. |
| object | [bytes](#bytes) |  | The binary object. |






<a name="payload-v1-Object-Distance"></a>

### Object.Distance
Represent the ID and distance pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The vector ID. |
| distance | [float](#float) |  | The distance. |






<a name="payload-v1-Object-ID"></a>

### Object.ID
Represent the vector ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="payload-v1-Object-IDs"></a>

### Object.IDs
Represent multiple vector IDs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ids | [string](#string) | repeated |  |






<a name="payload-v1-Object-Location"></a>

### Object.Location
Represent the vector location.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the location. |
| uuid | [string](#string) |  | The UUID of the vector. |
| ips | [string](#string) | repeated | The IP list. |






<a name="payload-v1-Object-Locations"></a>

### Object.Locations
Represent multiple vector locations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| locations | [Object.Location](#payload-v1-Object-Location) | repeated |  |






<a name="payload-v1-Object-ReshapeVector"></a>

### Object.ReshapeVector
Represent reshape vector.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [bytes](#bytes) |  | The binary object. |
| shape | [int32](#int32) | repeated | The new shape. |






<a name="payload-v1-Object-StreamBlob"></a>

### Object.StreamBlob
Represent stream response of binary objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob | [Object.Blob](#payload-v1-Object-Blob) |  | The binary object. |
| status | [google.rpc.Status](#google-rpc-Status) |  | The RPC error status. |






<a name="payload-v1-Object-StreamDistance"></a>

### Object.StreamDistance
Represent stream response of distances.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| distance | [Object.Distance](#payload-v1-Object-Distance) |  | The distance. |
| status | [google.rpc.Status](#google-rpc-Status) |  | The RPC error status. |






<a name="payload-v1-Object-StreamLocation"></a>

### Object.StreamLocation
Represent the stream response of the vector location.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| location | [Object.Location](#payload-v1-Object-Location) |  | The vector location. |
| status | [google.rpc.Status](#google-rpc-Status) |  | The RPC error status. |






<a name="payload-v1-Object-StreamVector"></a>

### Object.StreamVector
Represent stream response of the vector.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |  | The vector. |
| status | [google.rpc.Status](#google-rpc-Status) |  | The RPC error status. |






<a name="payload-v1-Object-Vector"></a>

### Object.Vector
Represent a vector.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The vector ID. |
| vector | [float](#float) | repeated | The vector. |






<a name="payload-v1-Object-VectorRequest"></a>

### Object.VectorRequest
Represent a request to fetch raw vector.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [Object.ID](#payload-v1-Object-ID) |  | The vector ID to be fetch. |
| filters | [Filter.Config](#payload-v1-Filter-Config) |  | Filter configurations. |






<a name="payload-v1-Object-Vectors"></a>

### Object.Vectors
Represent multiple vectors.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Object.Vector](#payload-v1-Object-Vector) | repeated |  |






<a name="payload-v1-Remove"></a>

### Remove
Remove related messages.






<a name="payload-v1-Remove-Config"></a>

### Remove.Config
Represent the remove configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  | A flag to skip exist check during upsert operation. |
| timestamp | [int64](#int64) |  | Remove timestamp. |






<a name="payload-v1-Remove-MultiRequest"></a>

### Remove.MultiRequest
Represent the multiple remove request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Remove.Request](#payload-v1-Remove-Request) | repeated | Represent the multiple remove request content. |






<a name="payload-v1-Remove-Request"></a>

### Remove.Request
Represent the remove request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [Object.ID](#payload-v1-Object-ID) |  | The object ID to be removed. |
| config | [Remove.Config](#payload-v1-Remove-Config) |  | The configuration of the remove request. |






<a name="payload-v1-Search"></a>

### Search
Search related messages.






<a name="payload-v1-Search-Config"></a>

### Search.Config
Represent search configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  | Unique request ID. |
| num | [uint32](#uint32) |  | Maximum number of result to be returned. |
| radius | [float](#float) |  | Search radius. |
| epsilon | [float](#float) |  | Search coefficient. |
| timeout | [int64](#int64) |  | Search timeout in nanoseconds. |
| ingress_filters | [Filter.Config](#payload-v1-Filter-Config) |  | Ingress filter configurations. |
| egress_filters | [Filter.Config](#payload-v1-Filter-Config) |  | Egress filter configurations. |
| min_num | [uint32](#uint32) |  | Minimum number of result to be returned. |






<a name="payload-v1-Search-IDRequest"></a>

### Search.IDRequest
Represent a search by ID request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The vector ID to be searched. |
| config | [Search.Config](#payload-v1-Search-Config) |  | The configuration of the search request. |






<a name="payload-v1-Search-MultiIDRequest"></a>

### Search.MultiIDRequest
Represent the multiple search by ID request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.IDRequest](#payload-v1-Search-IDRequest) | repeated | Represent the multiple search by ID request content. |






<a name="payload-v1-Search-MultiObjectRequest"></a>

### Search.MultiObjectRequest
Represent the multiple search by binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.ObjectRequest](#payload-v1-Search-ObjectRequest) | repeated | Represent the multiple search by binary object request content. |






<a name="payload-v1-Search-MultiRequest"></a>

### Search.MultiRequest
Represent the multiple search request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.Request](#payload-v1-Search-Request) | repeated | Represent the multiple search request content. |






<a name="payload-v1-Search-ObjectRequest"></a>

### Search.ObjectRequest
Represent a search by binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [bytes](#bytes) |  | The binary object to be searched. |
| config | [Search.Config](#payload-v1-Search-Config) |  | The configuration of the search request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |  | Filter configuration. |






<a name="payload-v1-Search-Request"></a>

### Search.Request
Represent a search request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [float](#float) | repeated | The vector to be searched. |
| config | [Search.Config](#payload-v1-Search-Config) |  | The configuration of the search request. |






<a name="payload-v1-Search-Response"></a>

### Search.Response
Represent a search response.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  | The unique request ID. |
| results | [Object.Distance](#payload-v1-Object-Distance) | repeated | Search results. |






<a name="payload-v1-Search-Responses"></a>

### Search.Responses
Represent multiple search responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| responses | [Search.Response](#payload-v1-Search-Response) | repeated | Represent the multiple search response content. |






<a name="payload-v1-Search-StreamResponse"></a>

### Search.StreamResponse
Represent stream search response.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response | [Search.Response](#payload-v1-Search-Response) |  | Represent the search response. |
| status | [google.rpc.Status](#google-rpc-Status) |  | The RPC error status. |






<a name="payload-v1-Update"></a>

### Update
Update related messages






<a name="payload-v1-Update-Config"></a>

### Update.Config
Represent the update configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  | A flag to skip exist check during update operation. |
| filters | [Filter.Config](#payload-v1-Filter-Config) |  | Filter configuration. |
| timestamp | [int64](#int64) |  | Update timestamp. |






<a name="payload-v1-Update-MultiObjectRequest"></a>

### Update.MultiObjectRequest
Represent the multiple update binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Update.ObjectRequest](#payload-v1-Update-ObjectRequest) | repeated | Represent the multiple update object request content. |






<a name="payload-v1-Update-MultiRequest"></a>

### Update.MultiRequest
Represent the multiple update request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Update.Request](#payload-v1-Update-Request) | repeated | Represent the multiple update request content. |






<a name="payload-v1-Update-ObjectRequest"></a>

### Update.ObjectRequest
Represent the update binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object.Blob](#payload-v1-Object-Blob) |  | The binary object to be updated. |
| config | [Update.Config](#payload-v1-Update-Config) |  | The configuration of the update request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |  | Filter target. |






<a name="payload-v1-Update-Request"></a>

### Update.Request
Represent the update request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |  | The vector to be updated. |
| config | [Update.Config](#payload-v1-Update-Config) |  | The configuration of the update request. |






<a name="payload-v1-Upsert"></a>

### Upsert
Upsert related messages.






<a name="payload-v1-Upsert-Config"></a>

### Upsert.Config
Represent the upsert configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  | A flag to skip exist check during upsert operation. |
| filters | [Filter.Config](#payload-v1-Filter-Config) |  | Filter configuration. |
| timestamp | [int64](#int64) |  | Upsert timestamp. |






<a name="payload-v1-Upsert-MultiObjectRequest"></a>

### Upsert.MultiObjectRequest
Represent the multiple upsert binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest) | repeated | Represent the multiple upsert object request content. |






<a name="payload-v1-Upsert-MultiRequest"></a>

### Upsert.MultiRequest
Represent mthe ultiple upsert request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Upsert.Request](#payload-v1-Upsert-Request) | repeated | Represent the multiple upsert request content. |






<a name="payload-v1-Upsert-ObjectRequest"></a>

### Upsert.ObjectRequest
Represent the upsert binary object request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object.Blob](#payload-v1-Object-Blob) |  | The binary object to be upserted. |
| config | [Upsert.Config](#payload-v1-Upsert-Config) |  | The configuration of the upsert request. |
| vectorizer | [Filter.Target](#payload-v1-Filter-Target) |  | Filter target. |






<a name="payload-v1-Upsert-Request"></a>

### Upsert.Request
Represent the upsert request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload-v1-Object-Vector) |  | The vector to be upserted. |
| config | [Upsert.Config](#payload-v1-Upsert-Config) |  | The configuration of the upsert request. |





 

 

 

 



<a name="apis_proto_v1_vald_filter-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/filter.proto


 

 

 


<a name="vald-v1-Filter"></a>

### Filter
Filter service provides ways to connect to Vald through filter.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SearchObject | [.payload.v1.Search.ObjectRequest](#payload-v1-Search-ObjectRequest) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | A method to search object. |
| MultiSearchObject | [.payload.v1.Search.MultiObjectRequest](#payload-v1-Search-MultiObjectRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | A method to search multiple objects. |
| StreamSearchObject | [.payload.v1.Search.ObjectRequest](#payload-v1-Search-ObjectRequest) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | A method to search object by bidirectional streaming. |
| InsertObject | [.payload.v1.Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method insert object. |
| StreamInsertObject | [.payload.v1.Insert.ObjectRequest](#payload-v1-Insert-ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | Represent the streaming RPC to insert object by bidirectional streaming. |
| MultiInsertObject | [.payload.v1.Insert.MultiObjectRequest](#payload-v1-Insert-MultiObjectRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to insert multiple objects. |
| UpdateObject | [.payload.v1.Update.ObjectRequest](#payload-v1-Update-ObjectRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method to update object. |
| StreamUpdateObject | [.payload.v1.Update.ObjectRequest](#payload-v1-Update-ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | A method to update object by bidirectional streaming. |
| MultiUpdateObject | [.payload.v1.Update.MultiObjectRequest](#payload-v1-Update-MultiObjectRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to update multiple objects. |
| UpsertObject | [.payload.v1.Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method to upsert object. |
| StreamUpsertObject | [.payload.v1.Upsert.ObjectRequest](#payload-v1-Upsert-ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | A method to upsert object by bidirectional streaming. |
| MultiUpsertObject | [.payload.v1.Upsert.MultiObjectRequest](#payload-v1-Upsert-MultiObjectRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to upsert multiple objects. |

 



<a name="apis_proto_v1_vald_insert-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/insert.proto


 

 

 


<a name="vald-v1-Insert"></a>

### Insert
Insert service provides ways to add new vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Insert | [.payload.v1.Insert.Request](#payload-v1-Insert-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method to add a new single vector. |
| StreamInsert | [.payload.v1.Insert.Request](#payload-v1-Insert-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | A method to add new multiple vectors by bidirectional streaming. |
| MultiInsert | [.payload.v1.Insert.MultiRequest](#payload-v1-Insert-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to add new multiple vectors in a single request. |

 



<a name="apis_proto_v1_vald_object-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/object.proto


 

 

 


<a name="vald-v1-Object"></a>

### Object
Object service provides ways to fetch indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exists | [.payload.v1.Object.ID](#payload-v1-Object-ID) | [.payload.v1.Object.ID](#payload-v1-Object-ID) | A method to check whether a specified ID is indexed or not. |
| GetObject | [.payload.v1.Object.VectorRequest](#payload-v1-Object-VectorRequest) | [.payload.v1.Object.Vector](#payload-v1-Object-Vector) | A method to fetch a vector. |
| StreamGetObject | [.payload.v1.Object.VectorRequest](#payload-v1-Object-VectorRequest) stream | [.payload.v1.Object.StreamVector](#payload-v1-Object-StreamVector) stream | A method to fetch vectors by bidirectional streaming. |

 



<a name="apis_proto_v1_vald_remove-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/remove.proto


 

 

 


<a name="vald-v1-Remove"></a>

### Remove
Remove service provides ways to remove indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Remove | [.payload.v1.Remove.Request](#payload-v1-Remove-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method to remove an indexed vector. |
| StreamRemove | [.payload.v1.Remove.Request](#payload-v1-Remove-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | A method to remove multiple indexed vectors by bidirectional streaming. |
| MultiRemove | [.payload.v1.Remove.MultiRequest](#payload-v1-Remove-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to remove multiple indexed vectors in a single request. |

 



<a name="apis_proto_v1_vald_search-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/search.proto


 

 

 


<a name="vald-v1-Search"></a>

### Search
Search service provides ways to search indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Search | [.payload.v1.Search.Request](#payload-v1-Search-Request) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | A method to search indexed vectors by a raw vector. |
| SearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | A method to search indexed vectors by ID. |
| StreamSearch | [.payload.v1.Search.Request](#payload-v1-Search-Request) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | A method to search indexed vectors by multiple vectors. |
| StreamSearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | A method to search indexed vectors by multiple IDs. |
| MultiSearch | [.payload.v1.Search.MultiRequest](#payload-v1-Search-MultiRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | A method to search indexed vectors by multiple vectors in a single request. |
| MultiSearchByID | [.payload.v1.Search.MultiIDRequest](#payload-v1-Search-MultiIDRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | A method to search indexed vectors by multiple IDs in a single request. |
| LinearSearch | [.payload.v1.Search.Request](#payload-v1-Search-Request) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | A method to linear search indexed vectors by a raw vector. |
| LinearSearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) | [.payload.v1.Search.Response](#payload-v1-Search-Response) | A method to linear search indexed vectors by ID. |
| StreamLinearSearch | [.payload.v1.Search.Request](#payload-v1-Search-Request) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | A method to linear search indexed vectors by multiple vectors. |
| StreamLinearSearchByID | [.payload.v1.Search.IDRequest](#payload-v1-Search-IDRequest) stream | [.payload.v1.Search.StreamResponse](#payload-v1-Search-StreamResponse) stream | A method to linear search indexed vectors by multiple IDs. |
| MultiLinearSearch | [.payload.v1.Search.MultiRequest](#payload-v1-Search-MultiRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | A method to linear search indexed vectors by multiple vectors in a single request. |
| MultiLinearSearchByID | [.payload.v1.Search.MultiIDRequest](#payload-v1-Search-MultiIDRequest) | [.payload.v1.Search.Responses](#payload-v1-Search-Responses) | A method to linear search indexed vectors by multiple IDs in a single request. |

 



<a name="apis_proto_v1_vald_update-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/update.proto


 

 

 


<a name="vald-v1-Update"></a>

### Update
Update service provides ways to update indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Update | [.payload.v1.Update.Request](#payload-v1-Update-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method to update an indexed vector. |
| StreamUpdate | [.payload.v1.Update.Request](#payload-v1-Update-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | A method to update multiple indexed vectors by bidirectional streaming. |
| MultiUpdate | [.payload.v1.Update.MultiRequest](#payload-v1-Update-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to update multiple indexed vectors in a single request. |

 



<a name="apis_proto_v1_vald_upsert-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/upsert.proto


 

 

 


<a name="vald-v1-Upsert"></a>

### Upsert
Upsert service provides ways to insert/update vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Upsert | [.payload.v1.Upsert.Request](#payload-v1-Upsert-Request) | [.payload.v1.Object.Location](#payload-v1-Object-Location) | A method to insert/update a vector. |
| StreamUpsert | [.payload.v1.Upsert.Request](#payload-v1-Upsert-Request) stream | [.payload.v1.Object.StreamLocation](#payload-v1-Object-StreamLocation) stream | A method to insert/update multiple vectors by bidirectional streaming. |
| MultiUpsert | [.payload.v1.Upsert.MultiRequest](#payload-v1-Upsert-MultiRequest) | [.payload.v1.Object.Locations](#payload-v1-Object-Locations) | A method to insert/update multiple vectors in a single request. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

