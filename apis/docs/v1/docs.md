# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [apis/proto/v1/vald/upsert.proto](#apis/proto/v1/vald/upsert.proto)
    - [Upsert](#vald.v1.Upsert)
  
- [apis/proto/v1/vald/filter.proto](#apis/proto/v1/vald/filter.proto)
    - [Filter](#vald.v1.Filter)
  
- [apis/proto/v1/vald/object.proto](#apis/proto/v1/vald/object.proto)
    - [Object](#vald.v1.Object)
  
- [apis/proto/v1/vald/insert.proto](#apis/proto/v1/vald/insert.proto)
    - [Insert](#vald.v1.Insert)
  
- [apis/proto/v1/vald/search.proto](#apis/proto/v1/vald/search.proto)
    - [Search](#vald.v1.Search)
  
- [apis/proto/v1/vald/remove.proto](#apis/proto/v1/vald/remove.proto)
    - [Remove](#vald.v1.Remove)
  
- [apis/proto/v1/vald/update.proto](#apis/proto/v1/vald/update.proto)
    - [Update](#vald.v1.Update)
  
- [apis/proto/v1/discoverer/discoverer.proto](#apis/proto/v1/discoverer/discoverer.proto)
    - [Discoverer](#discoverer.v1.Discoverer)
  
- [apis/proto/v1/payload/payload.proto](#apis/proto/v1/payload/payload.proto)
    - [Backup](#payload.v1.Backup)
    - [Backup.Compressed](#payload.v1.Backup.Compressed)
    - [Backup.Compressed.Vector](#payload.v1.Backup.Compressed.Vector)
    - [Backup.Compressed.Vectors](#payload.v1.Backup.Compressed.Vectors)
    - [Backup.GetVector](#payload.v1.Backup.GetVector)
    - [Backup.GetVector.Owner](#payload.v1.Backup.GetVector.Owner)
    - [Backup.GetVector.Request](#payload.v1.Backup.GetVector.Request)
    - [Backup.IP](#payload.v1.Backup.IP)
    - [Backup.IP.Register](#payload.v1.Backup.IP.Register)
    - [Backup.IP.Register.Request](#payload.v1.Backup.IP.Register.Request)
    - [Backup.IP.Remove](#payload.v1.Backup.IP.Remove)
    - [Backup.IP.Remove.Request](#payload.v1.Backup.IP.Remove.Request)
    - [Backup.Locations](#payload.v1.Backup.Locations)
    - [Backup.Locations.Request](#payload.v1.Backup.Locations.Request)
    - [Backup.Remove](#payload.v1.Backup.Remove)
    - [Backup.Remove.Request](#payload.v1.Backup.Remove.Request)
    - [Backup.Remove.RequestMulti](#payload.v1.Backup.Remove.RequestMulti)
    - [Backup.Vector](#payload.v1.Backup.Vector)
    - [Backup.Vectors](#payload.v1.Backup.Vectors)
    - [Control](#payload.v1.Control)
    - [Control.CreateIndexRequest](#payload.v1.Control.CreateIndexRequest)
    - [Discoverer](#payload.v1.Discoverer)
    - [Discoverer.Request](#payload.v1.Discoverer.Request)
    - [Empty](#payload.v1.Empty)
    - [Filter](#payload.v1.Filter)
    - [Filter.Config](#payload.v1.Filter.Config)
    - [Filter.Target](#payload.v1.Filter.Target)
    - [Info](#payload.v1.Info)
    - [Info.CPU](#payload.v1.Info.CPU)
    - [Info.IPs](#payload.v1.Info.IPs)
    - [Info.Index](#payload.v1.Info.Index)
    - [Info.Index.Count](#payload.v1.Info.Index.Count)
    - [Info.Index.UUID](#payload.v1.Info.Index.UUID)
    - [Info.Index.UUID.Committed](#payload.v1.Info.Index.UUID.Committed)
    - [Info.Index.UUID.Uncommitted](#payload.v1.Info.Index.UUID.Uncommitted)
    - [Info.Memory](#payload.v1.Info.Memory)
    - [Info.Node](#payload.v1.Info.Node)
    - [Info.Nodes](#payload.v1.Info.Nodes)
    - [Info.Pod](#payload.v1.Info.Pod)
    - [Info.Pods](#payload.v1.Info.Pods)
    - [Insert](#payload.v1.Insert)
    - [Insert.Config](#payload.v1.Insert.Config)
    - [Insert.MultiObjectRequest](#payload.v1.Insert.MultiObjectRequest)
    - [Insert.MultiRequest](#payload.v1.Insert.MultiRequest)
    - [Insert.ObjectRequest](#payload.v1.Insert.ObjectRequest)
    - [Insert.Request](#payload.v1.Insert.Request)
    - [Meta](#payload.v1.Meta)
    - [Meta.Key](#payload.v1.Meta.Key)
    - [Meta.KeyVal](#payload.v1.Meta.KeyVal)
    - [Meta.KeyVals](#payload.v1.Meta.KeyVals)
    - [Meta.Keys](#payload.v1.Meta.Keys)
    - [Meta.Val](#payload.v1.Meta.Val)
    - [Meta.Vals](#payload.v1.Meta.Vals)
    - [Object](#payload.v1.Object)
    - [Object.Blob](#payload.v1.Object.Blob)
    - [Object.Distance](#payload.v1.Object.Distance)
    - [Object.ID](#payload.v1.Object.ID)
    - [Object.IDs](#payload.v1.Object.IDs)
    - [Object.Location](#payload.v1.Object.Location)
    - [Object.Locations](#payload.v1.Object.Locations)
    - [Object.StreamBlob](#payload.v1.Object.StreamBlob)
    - [Object.StreamDistance](#payload.v1.Object.StreamDistance)
    - [Object.StreamLocation](#payload.v1.Object.StreamLocation)
    - [Object.StreamVector](#payload.v1.Object.StreamVector)
    - [Object.Vector](#payload.v1.Object.Vector)
    - [Object.VectorRequest](#payload.v1.Object.VectorRequest)
    - [Object.Vectors](#payload.v1.Object.Vectors)
    - [Remove](#payload.v1.Remove)
    - [Remove.Config](#payload.v1.Remove.Config)
    - [Remove.MultiRequest](#payload.v1.Remove.MultiRequest)
    - [Remove.Request](#payload.v1.Remove.Request)
    - [Replication](#payload.v1.Replication)
    - [Replication.Agents](#payload.v1.Replication.Agents)
    - [Replication.Rebalance](#payload.v1.Replication.Rebalance)
    - [Replication.Recovery](#payload.v1.Replication.Recovery)
    - [Search](#payload.v1.Search)
    - [Search.Config](#payload.v1.Search.Config)
    - [Search.IDRequest](#payload.v1.Search.IDRequest)
    - [Search.MultiIDRequest](#payload.v1.Search.MultiIDRequest)
    - [Search.MultiObjectRequest](#payload.v1.Search.MultiObjectRequest)
    - [Search.MultiRequest](#payload.v1.Search.MultiRequest)
    - [Search.ObjectRequest](#payload.v1.Search.ObjectRequest)
    - [Search.Request](#payload.v1.Search.Request)
    - [Search.Response](#payload.v1.Search.Response)
    - [Search.Responses](#payload.v1.Search.Responses)
    - [Search.StreamResponse](#payload.v1.Search.StreamResponse)
    - [Update](#payload.v1.Update)
    - [Update.Config](#payload.v1.Update.Config)
    - [Update.MultiObjectRequest](#payload.v1.Update.MultiObjectRequest)
    - [Update.MultiRequest](#payload.v1.Update.MultiRequest)
    - [Update.ObjectRequest](#payload.v1.Update.ObjectRequest)
    - [Update.Request](#payload.v1.Update.Request)
    - [Upsert](#payload.v1.Upsert)
    - [Upsert.Config](#payload.v1.Upsert.Config)
    - [Upsert.MultiObjectRequest](#payload.v1.Upsert.MultiObjectRequest)
    - [Upsert.MultiRequest](#payload.v1.Upsert.MultiRequest)
    - [Upsert.ObjectRequest](#payload.v1.Upsert.ObjectRequest)
    - [Upsert.Request](#payload.v1.Upsert.Request)
  
- [apis/proto/v1/agent/core/agent.proto](#apis/proto/v1/agent/core/agent.proto)
    - [Agent](#core.v1.Agent)
  
- [apis/proto/v1/agent/sidecar/sidecar.proto](#apis/proto/v1/agent/sidecar/sidecar.proto)
    - [Sidecar](#sidecar.v1.Sidecar)
  
- [apis/proto/v1/manager/index/index_manager.proto](#apis/proto/v1/manager/index/index_manager.proto)
    - [Index](#manager.index.v1.Index)
  
- [apis/proto/v1/filter/egress/egress_filter.proto](#apis/proto/v1/filter/egress/egress_filter.proto)
    - [Filter](#filter.egress.v1.Filter)
  
- [apis/proto/v1/filter/ingress/ingress_filter.proto](#apis/proto/v1/filter/ingress/ingress_filter.proto)
    - [Filter](#filter.ingress.v1.Filter)
  
- [Scalar Value Types](#scalar-value-types)



<a name="apis/proto/v1/vald/upsert.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/upsert.proto


 

 

 


<a name="vald.v1.Upsert"></a>

### Upsert
Upsert service provides ways to insert/update vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Upsert | [.payload.v1.Upsert.Request](#payload.v1.Upsert.Request) | [.payload.v1.Object.Location](#payload.v1.Object.Location) | A method to insert/update a vector. |
| StreamUpsert | [.payload.v1.Upsert.Request](#payload.v1.Upsert.Request) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream | A method to insert/update multiple vectors by bidirectional streaming. |
| MultiUpsert | [.payload.v1.Upsert.MultiRequest](#payload.v1.Upsert.MultiRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) | A method to insert/update multiple vectors in a single request. |

 



<a name="apis/proto/v1/vald/filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/filter.proto


 

 

 


<a name="vald.v1.Filter"></a>

### Filter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SearchObject | [.payload.v1.Search.ObjectRequest](#payload.v1.Search.ObjectRequest) | [.payload.v1.Search.Response](#payload.v1.Search.Response) |  |
| MultiSearchObject | [.payload.v1.Search.MultiObjectRequest](#payload.v1.Search.MultiObjectRequest) | [.payload.v1.Search.Responses](#payload.v1.Search.Responses) |  |
| StreamSearchObject | [.payload.v1.Search.ObjectRequest](#payload.v1.Search.ObjectRequest) stream | [.payload.v1.Search.StreamResponse](#payload.v1.Search.StreamResponse) stream |  |
| InsertObject | [.payload.v1.Insert.ObjectRequest](#payload.v1.Insert.ObjectRequest) | [.payload.v1.Object.Location](#payload.v1.Object.Location) |  |
| StreamInsertObject | [.payload.v1.Insert.ObjectRequest](#payload.v1.Insert.ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream |  |
| MultiInsertObject | [.payload.v1.Insert.MultiObjectRequest](#payload.v1.Insert.MultiObjectRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) |  |
| UpdateObject | [.payload.v1.Update.ObjectRequest](#payload.v1.Update.ObjectRequest) | [.payload.v1.Object.Location](#payload.v1.Object.Location) |  |
| StreamUpdateObject | [.payload.v1.Update.ObjectRequest](#payload.v1.Update.ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream |  |
| MultiUpdateObject | [.payload.v1.Update.MultiObjectRequest](#payload.v1.Update.MultiObjectRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) |  |
| UpsertObject | [.payload.v1.Upsert.ObjectRequest](#payload.v1.Upsert.ObjectRequest) | [.payload.v1.Object.Location](#payload.v1.Object.Location) |  |
| StreamUpsertObject | [.payload.v1.Upsert.ObjectRequest](#payload.v1.Upsert.ObjectRequest) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream |  |
| MultiUpsertObject | [.payload.v1.Upsert.MultiObjectRequest](#payload.v1.Upsert.MultiObjectRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) |  |

 



<a name="apis/proto/v1/vald/object.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/object.proto


 

 

 


<a name="vald.v1.Object"></a>

### Object
Object service provides ways to fetch indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exists | [.payload.v1.Object.ID](#payload.v1.Object.ID) | [.payload.v1.Object.ID](#payload.v1.Object.ID) | A method to check whether a specified ID is indexed or not. |
| GetObject | [.payload.v1.Object.VectorRequest](#payload.v1.Object.VectorRequest) | [.payload.v1.Object.Vector](#payload.v1.Object.Vector) | A method to fetch a vector. |
| StreamGetObject | [.payload.v1.Object.VectorRequest](#payload.v1.Object.VectorRequest) stream | [.payload.v1.Object.StreamVector](#payload.v1.Object.StreamVector) stream | A method to fetch vectors by bidirectional streaming. |

 



<a name="apis/proto/v1/vald/insert.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/insert.proto


 

 

 


<a name="vald.v1.Insert"></a>

### Insert
Insert service provides ways to add new vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Insert | [.payload.v1.Insert.Request](#payload.v1.Insert.Request) | [.payload.v1.Object.Location](#payload.v1.Object.Location) | A method to add a new single vector. |
| StreamInsert | [.payload.v1.Insert.Request](#payload.v1.Insert.Request) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream | A method to add new multiple vectors by bidirectional streaming. |
| MultiInsert | [.payload.v1.Insert.MultiRequest](#payload.v1.Insert.MultiRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) | A method to add new multiple vectors in a single request. |

 



<a name="apis/proto/v1/vald/search.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/search.proto


 

 

 


<a name="vald.v1.Search"></a>

### Search
Search service provides ways to search indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Search | [.payload.v1.Search.Request](#payload.v1.Search.Request) | [.payload.v1.Search.Response](#payload.v1.Search.Response) | A method to search indexed vectors by a raw vector. |
| SearchByID | [.payload.v1.Search.IDRequest](#payload.v1.Search.IDRequest) | [.payload.v1.Search.Response](#payload.v1.Search.Response) | A method to search indexed vectors by ID. |
| StreamSearch | [.payload.v1.Search.Request](#payload.v1.Search.Request) stream | [.payload.v1.Search.StreamResponse](#payload.v1.Search.StreamResponse) stream | A method to search indexed vectors by multiple vectors. |
| StreamSearchByID | [.payload.v1.Search.IDRequest](#payload.v1.Search.IDRequest) stream | [.payload.v1.Search.StreamResponse](#payload.v1.Search.StreamResponse) stream | A method to search indexed vectors by multiple IDs. |
| MultiSearch | [.payload.v1.Search.MultiRequest](#payload.v1.Search.MultiRequest) | [.payload.v1.Search.Responses](#payload.v1.Search.Responses) | A method to search indexed vectors by multiple vectors in a single request. |
| MultiSearchByID | [.payload.v1.Search.MultiIDRequest](#payload.v1.Search.MultiIDRequest) | [.payload.v1.Search.Responses](#payload.v1.Search.Responses) | A method to search indexed vectors by multiple IDs in a single request. |

 



<a name="apis/proto/v1/vald/remove.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/remove.proto


 

 

 


<a name="vald.v1.Remove"></a>

### Remove
Remove service provides ways to remove indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Remove | [.payload.v1.Remove.Request](#payload.v1.Remove.Request) | [.payload.v1.Object.Location](#payload.v1.Object.Location) | A method to remove an indexed vector. |
| StreamRemove | [.payload.v1.Remove.Request](#payload.v1.Remove.Request) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream | A method to remove multiple indexed vectors by bidirectional streaming. |
| MultiRemove | [.payload.v1.Remove.MultiRequest](#payload.v1.Remove.MultiRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) | A method to remove multiple indexed vectors in a single request. |

 



<a name="apis/proto/v1/vald/update.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/vald/update.proto


 

 

 


<a name="vald.v1.Update"></a>

### Update
Update service provides ways to update indexed vectors.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Update | [.payload.v1.Update.Request](#payload.v1.Update.Request) | [.payload.v1.Object.Location](#payload.v1.Object.Location) | A method to update an indexed vector. |
| StreamUpdate | [.payload.v1.Update.Request](#payload.v1.Update.Request) stream | [.payload.v1.Object.StreamLocation](#payload.v1.Object.StreamLocation) stream | A method to update multiple indexed vectors by bidirectional streaming. |
| MultiUpdate | [.payload.v1.Update.MultiRequest](#payload.v1.Update.MultiRequest) | [.payload.v1.Object.Locations](#payload.v1.Object.Locations) | A method to update multiple indexed vectors in a single request. |

 



<a name="apis/proto/v1/discoverer/discoverer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/discoverer/discoverer.proto


 

 

 


<a name="discoverer.v1.Discoverer"></a>

### Discoverer


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Pods | [.payload.v1.Discoverer.Request](#payload.v1.Discoverer.Request) | [.payload.v1.Info.Pods](#payload.v1.Info.Pods) |  |
| Nodes | [.payload.v1.Discoverer.Request](#payload.v1.Discoverer.Request) | [.payload.v1.Info.Nodes](#payload.v1.Info.Nodes) |  |

 



<a name="apis/proto/v1/payload/payload.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/payload/payload.proto



<a name="payload.v1.Backup"></a>

### Backup







<a name="payload.v1.Backup.Compressed"></a>

### Backup.Compressed







<a name="payload.v1.Backup.Compressed.Vector"></a>

### Backup.Compressed.Vector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| vector | [bytes](#bytes) |  |  |
| ips | [string](#string) | repeated |  |






<a name="payload.v1.Backup.Compressed.Vectors"></a>

### Backup.Compressed.Vectors



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Backup.Compressed.Vector](#payload.v1.Backup.Compressed.Vector) | repeated |  |






<a name="payload.v1.Backup.GetVector"></a>

### Backup.GetVector







<a name="payload.v1.Backup.GetVector.Owner"></a>

### Backup.GetVector.Owner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) |  |  |






<a name="payload.v1.Backup.GetVector.Request"></a>

### Backup.GetVector.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.v1.Backup.IP"></a>

### Backup.IP







<a name="payload.v1.Backup.IP.Register"></a>

### Backup.IP.Register







<a name="payload.v1.Backup.IP.Register.Request"></a>

### Backup.IP.Register.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| ips | [string](#string) | repeated |  |






<a name="payload.v1.Backup.IP.Remove"></a>

### Backup.IP.Remove







<a name="payload.v1.Backup.IP.Remove.Request"></a>

### Backup.IP.Remove.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ips | [string](#string) | repeated |  |






<a name="payload.v1.Backup.Locations"></a>

### Backup.Locations







<a name="payload.v1.Backup.Locations.Request"></a>

### Backup.Locations.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.v1.Backup.Remove"></a>

### Backup.Remove







<a name="payload.v1.Backup.Remove.Request"></a>

### Backup.Remove.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.v1.Backup.Remove.RequestMulti"></a>

### Backup.Remove.RequestMulti



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuids | [string](#string) | repeated |  |






<a name="payload.v1.Backup.Vector"></a>

### Backup.Vector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| vector | [float](#float) | repeated |  |
| ips | [string](#string) | repeated |  |






<a name="payload.v1.Backup.Vectors"></a>

### Backup.Vectors



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Backup.Vector](#payload.v1.Backup.Vector) | repeated |  |






<a name="payload.v1.Control"></a>

### Control







<a name="payload.v1.Control.CreateIndexRequest"></a>

### Control.CreateIndexRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_size | [uint32](#uint32) |  |  |






<a name="payload.v1.Discoverer"></a>

### Discoverer







<a name="payload.v1.Discoverer.Request"></a>

### Discoverer.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| namespace | [string](#string) |  |  |
| node | [string](#string) |  |  |






<a name="payload.v1.Empty"></a>

### Empty







<a name="payload.v1.Filter"></a>

### Filter
Filter related messages.






<a name="payload.v1.Filter.Config"></a>

### Filter.Config
Represents filter server configurations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| targets | [Filter.Target](#payload.v1.Filter.Target) | repeated |  |






<a name="payload.v1.Filter.Target"></a>

### Filter.Target
Represents filter server.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | hostname. |
| port | [uint32](#uint32) |  | port. |






<a name="payload.v1.Info"></a>

### Info







<a name="payload.v1.Info.CPU"></a>

### Info.CPU



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [double](#double) |  |  |
| request | [double](#double) |  |  |
| usage | [double](#double) |  |  |






<a name="payload.v1.Info.IPs"></a>

### Info.IPs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) | repeated |  |






<a name="payload.v1.Info.Index"></a>

### Info.Index







<a name="payload.v1.Info.Index.Count"></a>

### Info.Index.Count



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stored | [uint32](#uint32) |  |  |
| uncommitted | [uint32](#uint32) |  |  |
| indexing | [bool](#bool) |  |  |
| saving | [bool](#bool) |  |  |






<a name="payload.v1.Info.Index.UUID"></a>

### Info.Index.UUID







<a name="payload.v1.Info.Index.UUID.Committed"></a>

### Info.Index.UUID.Committed



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.v1.Info.Index.UUID.Uncommitted"></a>

### Info.Index.UUID.Uncommitted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.v1.Info.Memory"></a>

### Info.Memory



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [double](#double) |  |  |
| request | [double](#double) |  |  |
| usage | [double](#double) |  |  |






<a name="payload.v1.Info.Node"></a>

### Info.Node



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| internal_addr | [string](#string) |  |  |
| external_addr | [string](#string) |  |  |
| cpu | [Info.CPU](#payload.v1.Info.CPU) |  |  |
| memory | [Info.Memory](#payload.v1.Info.Memory) |  |  |
| Pods | [Info.Pods](#payload.v1.Info.Pods) |  |  |






<a name="payload.v1.Info.Nodes"></a>

### Info.Nodes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nodes | [Info.Node](#payload.v1.Info.Node) | repeated |  |






<a name="payload.v1.Info.Pod"></a>

### Info.Pod



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app_name | [string](#string) |  |  |
| name | [string](#string) |  |  |
| namespace | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| cpu | [Info.CPU](#payload.v1.Info.CPU) |  |  |
| memory | [Info.Memory](#payload.v1.Info.Memory) |  |  |
| node | [Info.Node](#payload.v1.Info.Node) |  |  |






<a name="payload.v1.Info.Pods"></a>

### Info.Pods



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pods | [Info.Pod](#payload.v1.Info.Pod) | repeated |  |






<a name="payload.v1.Insert"></a>

### Insert
Insert related messages.






<a name="payload.v1.Insert.Config"></a>

### Insert.Config
Represents insert configurations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  | If it is enabled, exist checking will be skipped during insert operation. |
| filters | [Filter.Config](#payload.v1.Filter.Config) |  | filter configurations. |
| timestamp | [int64](#int64) |  | timestamp. |






<a name="payload.v1.Insert.MultiObjectRequest"></a>

### Insert.MultiObjectRequest
Represents multiple insert request by binary object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Insert.ObjectRequest](#payload.v1.Insert.ObjectRequest) | repeated |  |






<a name="payload.v1.Insert.MultiRequest"></a>

### Insert.MultiRequest
Represents multiple insert requests.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Insert.Request](#payload.v1.Insert.Request) | repeated |  |






<a name="payload.v1.Insert.ObjectRequest"></a>

### Insert.ObjectRequest
Represents an insert request by binary object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object.Blob](#payload.v1.Object.Blob) |  | binary object. |
| config | [Insert.Config](#payload.v1.Insert.Config) |  | insert configurations. |
| vectorizer | [Filter.Target](#payload.v1.Filter.Target) |  | filter configurations. |






<a name="payload.v1.Insert.Request"></a>

### Insert.Request
Represents an insert request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload.v1.Object.Vector) |  | vector. |
| config | [Insert.Config](#payload.v1.Insert.Config) |  | insert configurations. |






<a name="payload.v1.Meta"></a>

### Meta







<a name="payload.v1.Meta.Key"></a>

### Meta.Key



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |






<a name="payload.v1.Meta.KeyVal"></a>

### Meta.KeyVal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| val | [string](#string) |  |  |






<a name="payload.v1.Meta.KeyVals"></a>

### Meta.KeyVals



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kvs | [Meta.KeyVal](#payload.v1.Meta.KeyVal) | repeated |  |






<a name="payload.v1.Meta.Keys"></a>

### Meta.Keys



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| keys | [string](#string) | repeated |  |






<a name="payload.v1.Meta.Val"></a>

### Meta.Val



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [string](#string) |  |  |






<a name="payload.v1.Meta.Vals"></a>

### Meta.Vals



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [string](#string) | repeated |  |






<a name="payload.v1.Object"></a>

### Object
Common messages.






<a name="payload.v1.Object.Blob"></a>

### Object.Blob
Represents binary object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID. |
| object | [bytes](#bytes) |  | binary object. |






<a name="payload.v1.Object.Distance"></a>

### Object.Distance
Represents ID and distance pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID. |
| distance | [float](#float) |  | distance. |






<a name="payload.v1.Object.ID"></a>

### Object.ID
Represents ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="payload.v1.Object.IDs"></a>

### Object.IDs
Represents IDs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ids | [string](#string) | repeated |  |






<a name="payload.v1.Object.Location"></a>

### Object.Location
Represents a vector location.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the location. |
| uuid | [string](#string) |  | UUID. |
| ips | [string](#string) | repeated | IP list. |






<a name="payload.v1.Object.Locations"></a>

### Object.Locations
Represents multiple vector locations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| locations | [Object.Location](#payload.v1.Object.Location) | repeated |  |






<a name="payload.v1.Object.StreamBlob"></a>

### Object.StreamBlob
Represents stream response of binary objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob | [Object.Blob](#payload.v1.Object.Blob) |  | binary object. |
| status | [google.rpc.Status](#google.rpc.Status) |  | error status. |






<a name="payload.v1.Object.StreamDistance"></a>

### Object.StreamDistance
Represents stream response of distances.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| distance | [Object.Distance](#payload.v1.Object.Distance) |  | distance. |
| status | [google.rpc.Status](#google.rpc.Status) |  | error status. |






<a name="payload.v1.Object.StreamLocation"></a>

### Object.StreamLocation
Represents stream response of locations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| location | [Object.Location](#payload.v1.Object.Location) |  | location. |
| status | [google.rpc.Status](#google.rpc.Status) |  | error status. |






<a name="payload.v1.Object.StreamVector"></a>

### Object.StreamVector
Represents stream response of vectors.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload.v1.Object.Vector) |  | vectors. |
| status | [google.rpc.Status](#google.rpc.Status) |  | error status. |






<a name="payload.v1.Object.Vector"></a>

### Object.Vector
Represents a vector.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID. |
| vector | [float](#float) | repeated | vector. |






<a name="payload.v1.Object.VectorRequest"></a>

### Object.VectorRequest
Represents a request to fetch raw vector.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [Object.ID](#payload.v1.Object.ID) |  | ID. |
| filters | [Filter.Config](#payload.v1.Filter.Config) |  | filter configurations. |






<a name="payload.v1.Object.Vectors"></a>

### Object.Vectors
Represents multiple vectors.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Object.Vector](#payload.v1.Object.Vector) | repeated |  |






<a name="payload.v1.Remove"></a>

### Remove







<a name="payload.v1.Remove.Config"></a>

### Remove.Config



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  |  |
| timestamp | [int64](#int64) |  |  |






<a name="payload.v1.Remove.MultiRequest"></a>

### Remove.MultiRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Remove.Request](#payload.v1.Remove.Request) | repeated |  |






<a name="payload.v1.Remove.Request"></a>

### Remove.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [Object.ID](#payload.v1.Object.ID) |  |  |
| config | [Remove.Config](#payload.v1.Remove.Config) |  |  |






<a name="payload.v1.Replication"></a>

### Replication







<a name="payload.v1.Replication.Agents"></a>

### Replication.Agents



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agents | [string](#string) | repeated |  |
| removed_agents | [string](#string) | repeated |  |
| replicating_agent | [string](#string) | repeated |  |






<a name="payload.v1.Replication.Rebalance"></a>

### Replication.Rebalance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| high_usage_agents | [string](#string) | repeated |  |
| low_usage_agents | [string](#string) | repeated |  |






<a name="payload.v1.Replication.Recovery"></a>

### Replication.Recovery



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deleted_agents | [string](#string) | repeated |  |






<a name="payload.v1.Search"></a>

### Search
Search related messages.






<a name="payload.v1.Search.Config"></a>

### Search.Config
Represents search configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  | unique request ID. |
| num | [uint32](#uint32) |  | number of results. |
| radius | [float](#float) |  | search radius. |
| epsilon | [float](#float) |  | search coefficient. |
| timeout | [int64](#int64) |  | timeout in nanoseconds. |
| ingress_filters | [Filter.Config](#payload.v1.Filter.Config) |  | ingress filter configurations. |
| egress_filters | [Filter.Config](#payload.v1.Filter.Config) |  | egress filter configurations. |






<a name="payload.v1.Search.IDRequest"></a>

### Search.IDRequest
Represents a search request by ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | vector ID. |
| config | [Search.Config](#payload.v1.Search.Config) |  | search configuration. |






<a name="payload.v1.Search.MultiIDRequest"></a>

### Search.MultiIDRequest
Represents multiple search requests by IDs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.IDRequest](#payload.v1.Search.IDRequest) | repeated |  |






<a name="payload.v1.Search.MultiObjectRequest"></a>

### Search.MultiObjectRequest
Represents multiple search requests by binary objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.ObjectRequest](#payload.v1.Search.ObjectRequest) | repeated |  |






<a name="payload.v1.Search.MultiRequest"></a>

### Search.MultiRequest
Represents multiple search requests.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.Request](#payload.v1.Search.Request) | repeated |  |






<a name="payload.v1.Search.ObjectRequest"></a>

### Search.ObjectRequest
Represents a search request by binary object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [bytes](#bytes) |  | binary object. |
| config | [Search.Config](#payload.v1.Search.Config) |  | search configuration. |
| vectorizer | [Filter.Target](#payload.v1.Filter.Target) |  | filter configuration. |






<a name="payload.v1.Search.Request"></a>

### Search.Request
Represents a search request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [float](#float) | repeated | vector. |
| config | [Search.Config](#payload.v1.Search.Config) |  | search configuration. |






<a name="payload.v1.Search.Response"></a>

### Search.Response
Represents a search response.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  | unique request ID. |
| results | [Object.Distance](#payload.v1.Object.Distance) | repeated | search results. |






<a name="payload.v1.Search.Responses"></a>

### Search.Responses
Represents multiple search responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| responses | [Search.Response](#payload.v1.Search.Response) | repeated |  |






<a name="payload.v1.Search.StreamResponse"></a>

### Search.StreamResponse
Represents stream response.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response | [Search.Response](#payload.v1.Search.Response) |  | search response. |
| status | [google.rpc.Status](#google.rpc.Status) |  | error status. |






<a name="payload.v1.Update"></a>

### Update







<a name="payload.v1.Update.Config"></a>

### Update.Config



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  |  |
| filters | [Filter.Config](#payload.v1.Filter.Config) |  |  |
| timestamp | [int64](#int64) |  |  |






<a name="payload.v1.Update.MultiObjectRequest"></a>

### Update.MultiObjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Update.ObjectRequest](#payload.v1.Update.ObjectRequest) | repeated |  |






<a name="payload.v1.Update.MultiRequest"></a>

### Update.MultiRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Update.Request](#payload.v1.Update.Request) | repeated |  |






<a name="payload.v1.Update.ObjectRequest"></a>

### Update.ObjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object.Blob](#payload.v1.Object.Blob) |  |  |
| config | [Update.Config](#payload.v1.Update.Config) |  |  |
| vectorizer | [Filter.Target](#payload.v1.Filter.Target) |  |  |






<a name="payload.v1.Update.Request"></a>

### Update.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload.v1.Object.Vector) |  |  |
| config | [Update.Config](#payload.v1.Update.Config) |  |  |






<a name="payload.v1.Upsert"></a>

### Upsert







<a name="payload.v1.Upsert.Config"></a>

### Upsert.Config



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip_strict_exist_check | [bool](#bool) |  |  |
| filters | [Filter.Config](#payload.v1.Filter.Config) |  |  |
| timestamp | [int64](#int64) |  |  |






<a name="payload.v1.Upsert.MultiObjectRequest"></a>

### Upsert.MultiObjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Upsert.ObjectRequest](#payload.v1.Upsert.ObjectRequest) | repeated |  |






<a name="payload.v1.Upsert.MultiRequest"></a>

### Upsert.MultiRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Upsert.Request](#payload.v1.Upsert.Request) | repeated |  |






<a name="payload.v1.Upsert.ObjectRequest"></a>

### Upsert.ObjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [Object.Blob](#payload.v1.Object.Blob) |  |  |
| config | [Upsert.Config](#payload.v1.Upsert.Config) |  |  |
| vectorizer | [Filter.Target](#payload.v1.Filter.Target) |  |  |






<a name="payload.v1.Upsert.Request"></a>

### Upsert.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [Object.Vector](#payload.v1.Object.Vector) |  |  |
| config | [Upsert.Config](#payload.v1.Upsert.Config) |  |  |





 

 

 

 



<a name="apis/proto/v1/agent/core/agent.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/agent/core/agent.proto


 

 

 


<a name="core.v1.Agent"></a>

### Agent


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateIndex | [.payload.v1.Control.CreateIndexRequest](#payload.v1.Control.CreateIndexRequest) | [.payload.v1.Empty](#payload.v1.Empty) |  |
| SaveIndex | [.payload.v1.Empty](#payload.v1.Empty) | [.payload.v1.Empty](#payload.v1.Empty) |  |
| CreateAndSaveIndex | [.payload.v1.Control.CreateIndexRequest](#payload.v1.Control.CreateIndexRequest) | [.payload.v1.Empty](#payload.v1.Empty) |  |
| IndexInfo | [.payload.v1.Empty](#payload.v1.Empty) | [.payload.v1.Info.Index.Count](#payload.v1.Info.Index.Count) |  |

 



<a name="apis/proto/v1/agent/sidecar/sidecar.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/agent/sidecar/sidecar.proto


 

 

 


<a name="sidecar.v1.Sidecar"></a>

### Sidecar


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|

 



<a name="apis/proto/v1/manager/index/index_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/manager/index/index_manager.proto


 

 

 


<a name="manager.index.v1.Index"></a>

### Index


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| IndexInfo | [.payload.v1.Empty](#payload.v1.Empty) | [.payload.v1.Info.Index.Count](#payload.v1.Info.Index.Count) |  |

 



<a name="apis/proto/v1/filter/egress/egress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/filter/egress/egress_filter.proto


 

 

 


<a name="filter.egress.v1.Filter"></a>

### Filter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| FilterDistance | [.payload.v1.Object.Distance](#payload.v1.Object.Distance) | [.payload.v1.Object.Distance](#payload.v1.Object.Distance) |  |
| FilterVector | [.payload.v1.Object.Vector](#payload.v1.Object.Vector) | [.payload.v1.Object.Vector](#payload.v1.Object.Vector) |  |

 



<a name="apis/proto/v1/filter/ingress/ingress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/v1/filter/ingress/ingress_filter.proto


 

 

 


<a name="filter.ingress.v1.Filter"></a>

### Filter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GenVector | [.payload.v1.Object.Blob](#payload.v1.Object.Blob) | [.payload.v1.Object.Vector](#payload.v1.Object.Vector) |  |
| FilterVector | [.payload.v1.Object.Vector](#payload.v1.Object.Vector) | [.payload.v1.Object.Vector](#payload.v1.Object.Vector) |  |

 



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

