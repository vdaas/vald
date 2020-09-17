# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [apis/proto/discoverer/discoverer.proto](#apis/proto/discoverer/discoverer.proto)
    - [Discoverer](#discoverer.Discoverer)
  
- [apis/proto/meta/meta.proto](#apis/proto/meta/meta.proto)
    - [Meta](#meta_manager.Meta)
  
- [apis/proto/agent/core/agent.proto](#apis/proto/agent/core/agent.proto)
    - [Agent](#core.Agent)
  
- [apis/proto/agent/sidecar/sidecar.proto](#apis/proto/agent/sidecar/sidecar.proto)
    - [Sidecar](#sidecar.Sidecar)
  
- [apis/proto/manager/traffic/traffic_manager.proto](#apis/proto/manager/traffic/traffic_manager.proto)
- [apis/proto/manager/compressor/compressor.proto](#apis/proto/manager/compressor/compressor.proto)
    - [Backup](#compressor.Backup)
  
- [apis/proto/manager/index/index_manager.proto](#apis/proto/manager/index/index_manager.proto)
    - [Index](#index_manager.Index)
  
- [apis/proto/manager/replication/agent/replication_manager.proto](#apis/proto/manager/replication/agent/replication_manager.proto)
    - [Replication](#replication_manager.Replication)
  
- [apis/proto/manager/replication/controller/replication_manager.proto](#apis/proto/manager/replication/controller/replication_manager.proto)
    - [ReplicationController](#replication_manager.ReplicationController)
  
- [apis/proto/manager/backup/backup_manager.proto](#apis/proto/manager/backup/backup_manager.proto)
    - [Backup](#backup_manager.Backup)
  
- [apis/proto/filter/egress/egress_filter.proto](#apis/proto/filter/egress/egress_filter.proto)
    - [EgressFilter](#egress_filter.EgressFilter)
  
- [apis/proto/filter/ingress/ingress_filter.proto](#apis/proto/filter/ingress/ingress_filter.proto)
    - [IngressFilter](#ingress_filter.IngressFilter)
  
- [apis/proto/errors/errors.proto](#apis/proto/errors/errors.proto)
    - [Errors](#errors.Errors)
    - [Errors.RPC](#errors.Errors.RPC)
  
- [apis/proto/gateway/vald/vald.proto](#apis/proto/gateway/vald/vald.proto)
    - [Vald](#vald.Vald)
  
- [apis/proto/gateway/filter/filter.proto](#apis/proto/gateway/filter/filter.proto)
    - [Filter](#filter.Filter)
  
- [Scalar Value Types](#scalar-value-types)



<a name="apis/proto/discoverer/discoverer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/discoverer/discoverer.proto


 

 

 


<a name="discoverer.Discoverer"></a>

### Discoverer


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Pods | [.payload.Discoverer.Request](#payload.Discoverer.Request) | [.payload.Info.Pods](#payload.Info.Pods) |  |
| Nodes | [.payload.Discoverer.Request](#payload.Discoverer.Request) | [.payload.Info.Nodes](#payload.Info.Nodes) |  |

 



<a name="apis/proto/meta/meta.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/meta/meta.proto


 

 

 


<a name="meta_manager.Meta"></a>

### Meta


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetMeta | [.payload.Meta.Key](#payload.Meta.Key) | [.payload.Meta.Val](#payload.Meta.Val) |  |
| GetMetas | [.payload.Meta.Keys](#payload.Meta.Keys) | [.payload.Meta.Vals](#payload.Meta.Vals) |  |
| GetMetaInverse | [.payload.Meta.Val](#payload.Meta.Val) | [.payload.Meta.Key](#payload.Meta.Key) |  |
| GetMetasInverse | [.payload.Meta.Vals](#payload.Meta.Vals) | [.payload.Meta.Keys](#payload.Meta.Keys) |  |
| SetMeta | [.payload.Meta.KeyVal](#payload.Meta.KeyVal) | [.payload.Empty](#payload.Empty) |  |
| SetMetas | [.payload.Meta.KeyVals](#payload.Meta.KeyVals) | [.payload.Empty](#payload.Empty) |  |
| DeleteMeta | [.payload.Meta.Key](#payload.Meta.Key) | [.payload.Meta.Val](#payload.Meta.Val) |  |
| DeleteMetas | [.payload.Meta.Keys](#payload.Meta.Keys) | [.payload.Meta.Vals](#payload.Meta.Vals) |  |
| DeleteMetaInverse | [.payload.Meta.Val](#payload.Meta.Val) | [.payload.Meta.Key](#payload.Meta.Key) |  |
| DeleteMetasInverse | [.payload.Meta.Vals](#payload.Meta.Vals) | [.payload.Meta.Keys](#payload.Meta.Keys) |  |

 



<a name="apis/proto/agent/core/agent.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/agent/core/agent.proto


 

 

 


<a name="core.Agent"></a>

### Agent


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateIndex | [.payload.Control.CreateIndexRequest](#payload.Control.CreateIndexRequest) | [.payload.Empty](#payload.Empty) |  |
| SaveIndex | [.payload.Empty](#payload.Empty) | [.payload.Empty](#payload.Empty) |  |
| CreateAndSaveIndex | [.payload.Control.CreateIndexRequest](#payload.Control.CreateIndexRequest) | [.payload.Empty](#payload.Empty) |  |
| IndexInfo | [.payload.Empty](#payload.Empty) | [.payload.Info.Index.Count](#payload.Info.Index.Count) |  |

 



<a name="apis/proto/agent/sidecar/sidecar.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/agent/sidecar/sidecar.proto


 

 

 


<a name="sidecar.Sidecar"></a>

### Sidecar


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|

 



<a name="apis/proto/manager/traffic/traffic_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/manager/traffic/traffic_manager.proto


 

 

 

 



<a name="apis/proto/manager/compressor/compressor.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/manager/compressor/compressor.proto


 

 

 


<a name="compressor.Backup"></a>

### Backup


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetVector | [.payload.Backup.GetVector.Request](#payload.Backup.GetVector.Request) | [.payload.Backup.MetaVector](#payload.Backup.MetaVector) |  |
| Locations | [.payload.Backup.Locations.Request](#payload.Backup.Locations.Request) | [.payload.Info.IPs](#payload.Info.IPs) |  |
| Register | [.payload.Backup.MetaVector](#payload.Backup.MetaVector) | [.payload.Empty](#payload.Empty) |  |
| RegisterMulti | [.payload.Backup.MetaVectors](#payload.Backup.MetaVectors) | [.payload.Empty](#payload.Empty) |  |
| Remove | [.payload.Backup.Remove.Request](#payload.Backup.Remove.Request) | [.payload.Empty](#payload.Empty) |  |
| RemoveMulti | [.payload.Backup.Remove.RequestMulti](#payload.Backup.Remove.RequestMulti) | [.payload.Empty](#payload.Empty) |  |
| RegisterIPs | [.payload.Backup.IP.Register.Request](#payload.Backup.IP.Register.Request) | [.payload.Empty](#payload.Empty) |  |
| RemoveIPs | [.payload.Backup.IP.Remove.Request](#payload.Backup.IP.Remove.Request) | [.payload.Empty](#payload.Empty) |  |

 



<a name="apis/proto/manager/index/index_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/manager/index/index_manager.proto


 

 

 


<a name="index_manager.Index"></a>

### Index


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| IndexInfo | [.payload.Empty](#payload.Empty) | [.payload.Info.Index.Count](#payload.Info.Index.Count) |  |

 



<a name="apis/proto/manager/replication/agent/replication_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/manager/replication/agent/replication_manager.proto


 

 

 


<a name="replication_manager.Replication"></a>

### Replication


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Recover | [.payload.Replication.Recovery](#payload.Replication.Recovery) | [.payload.Empty](#payload.Empty) |  |
| Rebalance | [.payload.Replication.Rebalance](#payload.Replication.Rebalance) | [.payload.Empty](#payload.Empty) |  |
| AgentInfo | [.payload.Empty](#payload.Empty) | [.payload.Replication.Agents](#payload.Replication.Agents) |  |

 



<a name="apis/proto/manager/replication/controller/replication_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/manager/replication/controller/replication_manager.proto


 

 

 


<a name="replication_manager.ReplicationController"></a>

### ReplicationController


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ReplicationInfo | [.payload.Empty](#payload.Empty) | [.payload.Replication.Agents](#payload.Replication.Agents) |  |

 



<a name="apis/proto/manager/backup/backup_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/manager/backup/backup_manager.proto


 

 

 


<a name="backup_manager.Backup"></a>

### Backup


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetVector | [.payload.Backup.GetVector.Request](#payload.Backup.GetVector.Request) | [.payload.Backup.Compressed.MetaVector](#payload.Backup.Compressed.MetaVector) |  |
| Locations | [.payload.Backup.Locations.Request](#payload.Backup.Locations.Request) | [.payload.Info.IPs](#payload.Info.IPs) |  |
| Register | [.payload.Backup.Compressed.MetaVector](#payload.Backup.Compressed.MetaVector) | [.payload.Empty](#payload.Empty) |  |
| RegisterMulti | [.payload.Backup.Compressed.MetaVectors](#payload.Backup.Compressed.MetaVectors) | [.payload.Empty](#payload.Empty) |  |
| Remove | [.payload.Backup.Remove.Request](#payload.Backup.Remove.Request) | [.payload.Empty](#payload.Empty) |  |
| RemoveMulti | [.payload.Backup.Remove.RequestMulti](#payload.Backup.Remove.RequestMulti) | [.payload.Empty](#payload.Empty) |  |
| RegisterIPs | [.payload.Backup.IP.Register.Request](#payload.Backup.IP.Register.Request) | [.payload.Empty](#payload.Empty) |  |
| RemoveIPs | [.payload.Backup.IP.Remove.Request](#payload.Backup.IP.Remove.Request) | [.payload.Empty](#payload.Empty) |  |

 



<a name="apis/proto/filter/egress/egress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/filter/egress/egress_filter.proto


 

 

 


<a name="egress_filter.EgressFilter"></a>

### EgressFilter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Filter | [.payload.Object.Distance](#payload.Object.Distance) | [.payload.Object.Distance](#payload.Object.Distance) |  |
| StreamFilter | [.payload.Object.Distance](#payload.Object.Distance) stream | [.payload.Object.Distance](#payload.Object.Distance) stream |  |

 



<a name="apis/proto/filter/ingress/ingress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/filter/ingress/ingress_filter.proto


 

 

 


<a name="ingress_filter.IngressFilter"></a>

### IngressFilter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GenVector | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Vector](#payload.Object.Vector) |  |
| StreamGenVector | [.payload.Object.Blob](#payload.Object.Blob) stream | [.payload.Object.Vector](#payload.Object.Vector) stream |  |
| FilterVector | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Object.Vector](#payload.Object.Vector) |  |
| StreamFilterVector | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Object.Vector](#payload.Object.Vector) stream |  |

 



<a name="apis/proto/errors/errors.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/errors/errors.proto



<a name="errors.Errors"></a>

### Errors







<a name="errors.Errors.RPC"></a>

### Errors.RPC



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  |  |
| msg | [string](#string) |  |  |
| details | [string](#string) | repeated |  |
| instance | [string](#string) |  |  |
| status | [int64](#int64) |  |  |
| error | [string](#string) |  |  |
| roots | [Errors.RPC](#errors.Errors.RPC) | repeated |  |





 

 

 

 



<a name="apis/proto/gateway/vald/vald.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/gateway/vald/vald.proto


 

 

 


<a name="vald.Vald"></a>

### Vald


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exists | [.payload.Object.ID](#payload.Object.ID) | [.payload.Object.ID](#payload.Object.ID) |  |
| Search | [.payload.Search.Request](#payload.Search.Request) | [.payload.Search.Response](#payload.Search.Response) |  |
| SearchByID | [.payload.Search.IDRequest](#payload.Search.IDRequest) | [.payload.Search.Response](#payload.Search.Response) |  |
| StreamSearch | [.payload.Search.Request](#payload.Search.Request) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| StreamSearchByID | [.payload.Search.IDRequest](#payload.Search.IDRequest) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| Insert | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamInsert | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiInsert | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Object.Locations](#payload.Object.Locations) |  |
| Update | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamUpdate | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiUpdate | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Object.Locations](#payload.Object.Locations) |  |
| Upsert | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamUpsert | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiUpsert | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Object.Locations](#payload.Object.Locations) |  |
| Remove | [.payload.Object.ID](#payload.Object.ID) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamRemove | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiRemove | [.payload.Object.IDs](#payload.Object.IDs) | [.payload.Object.Locations](#payload.Object.Locations) |  |
| GetObject | [.payload.Object.ID](#payload.Object.ID) | [.payload.Object.Vector](#payload.Object.Vector) |  |
| StreamGetObject | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Object.Vector](#payload.Object.Vector) stream |  |

 



<a name="apis/proto/gateway/filter/filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/gateway/filter/filter.proto


 

 

 


<a name="filter.Filter"></a>

### Filter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SearchObject | [.payload.Search.ObjectRequest](#payload.Search.ObjectRequest) | [.payload.Search.Response](#payload.Search.Response) |  |
| StreamSearchObject | [.payload.Search.ObjectRequest](#payload.Search.ObjectRequest) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| InsertObject | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamInsertObject | [.payload.Object.Blob](#payload.Object.Blob) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiInsertObject | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Locations](#payload.Object.Locations) |  |
| UpdateObject | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamUpdateObject | [.payload.Object.Blob](#payload.Object.Blob) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiUpdateObject | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Locations](#payload.Object.Locations) |  |
| UpsertObject | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Location](#payload.Object.Location) |  |
| StreamUpsertObject | [.payload.Object.Blob](#payload.Object.Blob) stream | [.payload.Object.Location](#payload.Object.Location) stream |  |
| MultiUpsertObject | [.payload.Object.Blob](#payload.Object.Blob) | [.payload.Object.Locations](#payload.Object.Locations) |  |

 



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

