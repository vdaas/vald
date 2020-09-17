# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [discoverer.proto](#discoverer.proto)
    - [Discoverer](#discoverer.Discoverer)
  
- [payload.proto](#payload.proto)
    - [Backup](#payload.Backup)
    - [Backup.Compressed](#payload.Backup.Compressed)
    - [Backup.Compressed.MetaVector](#payload.Backup.Compressed.MetaVector)
    - [Backup.Compressed.MetaVectors](#payload.Backup.Compressed.MetaVectors)
    - [Backup.GetVector](#payload.Backup.GetVector)
    - [Backup.GetVector.Owner](#payload.Backup.GetVector.Owner)
    - [Backup.GetVector.Request](#payload.Backup.GetVector.Request)
    - [Backup.IP](#payload.Backup.IP)
    - [Backup.IP.Register](#payload.Backup.IP.Register)
    - [Backup.IP.Register.Request](#payload.Backup.IP.Register.Request)
    - [Backup.IP.Remove](#payload.Backup.IP.Remove)
    - [Backup.IP.Remove.Request](#payload.Backup.IP.Remove.Request)
    - [Backup.Locations](#payload.Backup.Locations)
    - [Backup.Locations.Request](#payload.Backup.Locations.Request)
    - [Backup.MetaVector](#payload.Backup.MetaVector)
    - [Backup.MetaVectors](#payload.Backup.MetaVectors)
    - [Backup.Remove](#payload.Backup.Remove)
    - [Backup.Remove.Request](#payload.Backup.Remove.Request)
    - [Backup.Remove.RequestMulti](#payload.Backup.Remove.RequestMulti)
    - [Control](#payload.Control)
    - [Control.CreateIndexRequest](#payload.Control.CreateIndexRequest)
    - [Discoverer](#payload.Discoverer)
    - [Discoverer.Request](#payload.Discoverer.Request)
    - [Empty](#payload.Empty)
    - [Info](#payload.Info)
    - [Info.CPU](#payload.Info.CPU)
    - [Info.IPs](#payload.Info.IPs)
    - [Info.Index](#payload.Info.Index)
    - [Info.Index.Count](#payload.Info.Index.Count)
    - [Info.Index.UUID](#payload.Info.Index.UUID)
    - [Info.Index.UUID.Committed](#payload.Info.Index.UUID.Committed)
    - [Info.Index.UUID.Uncommitted](#payload.Info.Index.UUID.Uncommitted)
    - [Info.Memory](#payload.Info.Memory)
    - [Info.Node](#payload.Info.Node)
    - [Info.Nodes](#payload.Info.Nodes)
    - [Info.Pod](#payload.Info.Pod)
    - [Info.Pods](#payload.Info.Pods)
    - [Meta](#payload.Meta)
    - [Meta.Key](#payload.Meta.Key)
    - [Meta.KeyVal](#payload.Meta.KeyVal)
    - [Meta.KeyVals](#payload.Meta.KeyVals)
    - [Meta.Keys](#payload.Meta.Keys)
    - [Meta.Val](#payload.Meta.Val)
    - [Meta.Vals](#payload.Meta.Vals)
    - [Object](#payload.Object)
    - [Object.Distance](#payload.Object.Distance)
    - [Object.ID](#payload.Object.ID)
    - [Object.IDs](#payload.Object.IDs)
    - [Object.Vector](#payload.Object.Vector)
    - [Object.Vectors](#payload.Object.Vectors)
    - [Replication](#payload.Replication)
    - [Replication.Agents](#payload.Replication.Agents)
    - [Replication.Rebalance](#payload.Replication.Rebalance)
    - [Replication.Recovery](#payload.Replication.Recovery)
    - [Search](#payload.Search)
    - [Search.Config](#payload.Search.Config)
    - [Search.IDRequest](#payload.Search.IDRequest)
    - [Search.Request](#payload.Search.Request)
    - [Search.Response](#payload.Search.Response)
  
- [meta.proto](#meta.proto)
    - [Meta](#meta_manager.Meta)
  
- [core/agent.proto](#core/agent.proto)
    - [Agent](#core.Agent)
  
- [sidecar/sidecar.proto](#sidecar/sidecar.proto)
    - [Sidecar](#sidecar.Sidecar)
  
- [traffic/traffic_manager.proto](#traffic/traffic_manager.proto)
- [compressor/compressor.proto](#compressor/compressor.proto)
    - [Backup](#compressor.Backup)
  
- [index/index_manager.proto](#index/index_manager.proto)
    - [Index](#index_manager.Index)
  
- [replication/agent/replication_manager.proto](#replication/agent/replication_manager.proto)
    - [Replication](#replication_manager.Replication)
  
- [replication/controller/replication_manager.proto](#replication/controller/replication_manager.proto)
    - [ReplicationController](#replication_manager.ReplicationController)
  
- [backup/backup_manager.proto](#backup/backup_manager.proto)
    - [Backup](#backup_manager.Backup)
  
- [egress/egress_filter.proto](#egress/egress_filter.proto)
    - [EgressFilter](#egress_filter.EgressFilter)
  
- [ingress/ingress_filter.proto](#ingress/ingress_filter.proto)
    - [IngressFilter](#ingress_filter.IngressFilter)
  
- [errors.proto](#errors.proto)
    - [Errors](#errors.Errors)
    - [Errors.RPC](#errors.Errors.RPC)
  
- [vald/vald.proto](#vald/vald.proto)
    - [Vald](#vald.Vald)
  
- [Scalar Value Types](#scalar-value-types)



<a name="discoverer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## discoverer.proto


 

 

 


<a name="discoverer.Discoverer"></a>

### Discoverer


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Pods | [.payload.Discoverer.Request](#payload.Discoverer.Request) | [.payload.Info.Pods](#payload.Info.Pods) |  |
| Nodes | [.payload.Discoverer.Request](#payload.Discoverer.Request) | [.payload.Info.Nodes](#payload.Info.Nodes) |  |

 



<a name="payload.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## payload.proto



<a name="payload.Backup"></a>

### Backup







<a name="payload.Backup.Compressed"></a>

### Backup.Compressed







<a name="payload.Backup.Compressed.MetaVector"></a>

### Backup.Compressed.MetaVector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| meta | [string](#string) |  |  |
| vector | [bytes](#bytes) |  |  |
| ips | [string](#string) | repeated |  |






<a name="payload.Backup.Compressed.MetaVectors"></a>

### Backup.Compressed.MetaVectors



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Backup.Compressed.MetaVector](#payload.Backup.Compressed.MetaVector) | repeated |  |






<a name="payload.Backup.GetVector"></a>

### Backup.GetVector







<a name="payload.Backup.GetVector.Owner"></a>

### Backup.GetVector.Owner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) |  |  |






<a name="payload.Backup.GetVector.Request"></a>

### Backup.GetVector.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.Backup.IP"></a>

### Backup.IP







<a name="payload.Backup.IP.Register"></a>

### Backup.IP.Register







<a name="payload.Backup.IP.Register.Request"></a>

### Backup.IP.Register.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| ips | [string](#string) | repeated |  |






<a name="payload.Backup.IP.Remove"></a>

### Backup.IP.Remove







<a name="payload.Backup.IP.Remove.Request"></a>

### Backup.IP.Remove.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ips | [string](#string) | repeated |  |






<a name="payload.Backup.Locations"></a>

### Backup.Locations







<a name="payload.Backup.Locations.Request"></a>

### Backup.Locations.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.Backup.MetaVector"></a>

### Backup.MetaVector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| meta | [string](#string) |  |  |
| vector | [float](#float) | repeated |  |
| ips | [string](#string) | repeated |  |






<a name="payload.Backup.MetaVectors"></a>

### Backup.MetaVectors



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Backup.MetaVector](#payload.Backup.MetaVector) | repeated |  |






<a name="payload.Backup.Remove"></a>

### Backup.Remove







<a name="payload.Backup.Remove.Request"></a>

### Backup.Remove.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.Backup.Remove.RequestMulti"></a>

### Backup.Remove.RequestMulti



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuids | [string](#string) | repeated |  |






<a name="payload.Control"></a>

### Control







<a name="payload.Control.CreateIndexRequest"></a>

### Control.CreateIndexRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_size | [uint32](#uint32) |  |  |






<a name="payload.Discoverer"></a>

### Discoverer







<a name="payload.Discoverer.Request"></a>

### Discoverer.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| namespace | [string](#string) |  |  |
| node | [string](#string) |  |  |






<a name="payload.Empty"></a>

### Empty







<a name="payload.Info"></a>

### Info







<a name="payload.Info.CPU"></a>

### Info.CPU



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [double](#double) |  |  |
| request | [double](#double) |  |  |
| usage | [double](#double) |  |  |






<a name="payload.Info.IPs"></a>

### Info.IPs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) | repeated |  |






<a name="payload.Info.Index"></a>

### Info.Index







<a name="payload.Info.Index.Count"></a>

### Info.Index.Count



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stored | [uint32](#uint32) |  |  |
| uncommitted | [uint32](#uint32) |  |  |
| indexing | [bool](#bool) |  |  |






<a name="payload.Info.Index.UUID"></a>

### Info.Index.UUID







<a name="payload.Info.Index.UUID.Committed"></a>

### Info.Index.UUID.Committed



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.Info.Index.UUID.Uncommitted"></a>

### Info.Index.UUID.Uncommitted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="payload.Info.Memory"></a>

### Info.Memory



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [double](#double) |  |  |
| request | [double](#double) |  |  |
| usage | [double](#double) |  |  |






<a name="payload.Info.Node"></a>

### Info.Node



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| internal_addr | [string](#string) |  |  |
| external_addr | [string](#string) |  |  |
| cpu | [Info.CPU](#payload.Info.CPU) |  |  |
| memory | [Info.Memory](#payload.Info.Memory) |  |  |
| Pods | [Info.Pods](#payload.Info.Pods) |  |  |






<a name="payload.Info.Nodes"></a>

### Info.Nodes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nodes | [Info.Node](#payload.Info.Node) | repeated |  |






<a name="payload.Info.Pod"></a>

### Info.Pod



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app_name | [string](#string) |  |  |
| name | [string](#string) |  |  |
| namespace | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| cpu | [Info.CPU](#payload.Info.CPU) |  |  |
| memory | [Info.Memory](#payload.Info.Memory) |  |  |
| node | [Info.Node](#payload.Info.Node) |  |  |






<a name="payload.Info.Pods"></a>

### Info.Pods



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pods | [Info.Pod](#payload.Info.Pod) | repeated |  |






<a name="payload.Meta"></a>

### Meta







<a name="payload.Meta.Key"></a>

### Meta.Key



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |






<a name="payload.Meta.KeyVal"></a>

### Meta.KeyVal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| val | [string](#string) |  |  |






<a name="payload.Meta.KeyVals"></a>

### Meta.KeyVals



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kvs | [Meta.KeyVal](#payload.Meta.KeyVal) | repeated |  |






<a name="payload.Meta.Keys"></a>

### Meta.Keys



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| keys | [string](#string) | repeated |  |






<a name="payload.Meta.Val"></a>

### Meta.Val



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [string](#string) |  |  |






<a name="payload.Meta.Vals"></a>

### Meta.Vals



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [string](#string) | repeated |  |






<a name="payload.Object"></a>

### Object







<a name="payload.Object.Distance"></a>

### Object.Distance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| distance | [float](#float) |  |  |






<a name="payload.Object.ID"></a>

### Object.ID



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="payload.Object.IDs"></a>

### Object.IDs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ids | [string](#string) | repeated |  |






<a name="payload.Object.Vector"></a>

### Object.Vector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| vector | [float](#float) | repeated |  |






<a name="payload.Object.Vectors"></a>

### Object.Vectors



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Object.Vector](#payload.Object.Vector) | repeated |  |






<a name="payload.Replication"></a>

### Replication







<a name="payload.Replication.Agents"></a>

### Replication.Agents



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| agents | [string](#string) | repeated |  |
| removed_agents | [string](#string) | repeated |  |
| replicating_agent | [string](#string) | repeated |  |






<a name="payload.Replication.Rebalance"></a>

### Replication.Rebalance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| high_usage_agents | [string](#string) | repeated |  |
| low_usage_agents | [string](#string) | repeated |  |






<a name="payload.Replication.Recovery"></a>

### Replication.Recovery



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deleted_agents | [string](#string) | repeated |  |






<a name="payload.Search"></a>

### Search







<a name="payload.Search.Config"></a>

### Search.Config



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| num | [uint32](#uint32) |  |  |
| radius | [float](#float) |  |  |
| epsilon | [float](#float) |  |  |
| timeout | [int64](#int64) |  |  |






<a name="payload.Search.IDRequest"></a>

### Search.IDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| config | [Search.Config](#payload.Search.Config) |  |  |






<a name="payload.Search.Request"></a>

### Search.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vector | [float](#float) | repeated |  |
| config | [Search.Config](#payload.Search.Config) |  |  |






<a name="payload.Search.Response"></a>

### Search.Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| results | [Object.Distance](#payload.Object.Distance) | repeated |  |





 

 

 

 



<a name="meta.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## meta.proto


 

 

 


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

 



<a name="core/agent.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## core/agent.proto


 

 

 


<a name="core.Agent"></a>

### Agent


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exists | [.payload.Object.ID](#payload.Object.ID) | [.payload.Object.ID](#payload.Object.ID) |  |
| Search | [.payload.Search.Request](#payload.Search.Request) | [.payload.Search.Response](#payload.Search.Response) |  |
| SearchByID | [.payload.Search.IDRequest](#payload.Search.IDRequest) | [.payload.Search.Response](#payload.Search.Response) |  |
| StreamSearch | [.payload.Search.Request](#payload.Search.Request) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| StreamSearchByID | [.payload.Search.IDRequest](#payload.Search.IDRequest) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| Insert | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Empty](#payload.Empty) |  |
| StreamInsert | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiInsert | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Empty](#payload.Empty) |  |
| Update | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Empty](#payload.Empty) |  |
| StreamUpdate | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiUpdate | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Empty](#payload.Empty) |  |
| Remove | [.payload.Object.ID](#payload.Object.ID) | [.payload.Empty](#payload.Empty) |  |
| StreamRemove | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiRemove | [.payload.Object.IDs](#payload.Object.IDs) | [.payload.Empty](#payload.Empty) |  |
| GetObject | [.payload.Object.ID](#payload.Object.ID) | [.payload.Object.Vector](#payload.Object.Vector) |  |
| StreamGetObject | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Object.Vector](#payload.Object.Vector) stream |  |
| CreateIndex | [.payload.Control.CreateIndexRequest](#payload.Control.CreateIndexRequest) | [.payload.Empty](#payload.Empty) |  |
| SaveIndex | [.payload.Empty](#payload.Empty) | [.payload.Empty](#payload.Empty) |  |
| CreateAndSaveIndex | [.payload.Control.CreateIndexRequest](#payload.Control.CreateIndexRequest) | [.payload.Empty](#payload.Empty) |  |
| IndexInfo | [.payload.Empty](#payload.Empty) | [.payload.Info.Index.Count](#payload.Info.Index.Count) |  |

 



<a name="sidecar/sidecar.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## sidecar/sidecar.proto


 

 

 


<a name="sidecar.Sidecar"></a>

### Sidecar


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|

 



<a name="traffic/traffic_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## traffic/traffic_manager.proto


 

 

 

 



<a name="compressor/compressor.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## compressor/compressor.proto


 

 

 


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

 



<a name="index/index_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## index/index_manager.proto


 

 

 


<a name="index_manager.Index"></a>

### Index


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| IndexInfo | [.payload.Empty](#payload.Empty) | [.payload.Info.Index.Count](#payload.Info.Index.Count) |  |

 



<a name="replication/agent/replication_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## replication/agent/replication_manager.proto


 

 

 


<a name="replication_manager.Replication"></a>

### Replication


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Recover | [.payload.Replication.Recovery](#payload.Replication.Recovery) | [.payload.Empty](#payload.Empty) |  |
| Rebalance | [.payload.Replication.Rebalance](#payload.Replication.Rebalance) | [.payload.Empty](#payload.Empty) |  |
| AgentInfo | [.payload.Empty](#payload.Empty) | [.payload.Replication.Agents](#payload.Replication.Agents) |  |

 



<a name="replication/controller/replication_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## replication/controller/replication_manager.proto


 

 

 


<a name="replication_manager.ReplicationController"></a>

### ReplicationController


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ReplicationInfo | [.payload.Empty](#payload.Empty) | [.payload.Replication.Agents](#payload.Replication.Agents) |  |

 



<a name="backup/backup_manager.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## backup/backup_manager.proto


 

 

 


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

 



<a name="egress/egress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## egress/egress_filter.proto


 

 

 


<a name="egress_filter.EgressFilter"></a>

### EgressFilter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Filter | [.payload.Search.Response](#payload.Search.Response) | [.payload.Search.Response](#payload.Search.Response) |  |
| StreamFilter | [.payload.Object.Distance](#payload.Object.Distance) stream | [.payload.Object.Distance](#payload.Object.Distance) stream |  |

 



<a name="ingress/ingress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ingress/ingress_filter.proto


 

 

 


<a name="ingress_filter.IngressFilter"></a>

### IngressFilter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|

 



<a name="errors.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## errors.proto



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





 

 

 

 



<a name="vald/vald.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vald/vald.proto


 

 

 


<a name="vald.Vald"></a>

### Vald


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exists | [.payload.Object.ID](#payload.Object.ID) | [.payload.Object.ID](#payload.Object.ID) |  |
| Search | [.payload.Search.Request](#payload.Search.Request) | [.payload.Search.Response](#payload.Search.Response) |  |
| SearchByID | [.payload.Search.IDRequest](#payload.Search.IDRequest) | [.payload.Search.Response](#payload.Search.Response) |  |
| StreamSearch | [.payload.Search.Request](#payload.Search.Request) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| StreamSearchByID | [.payload.Search.IDRequest](#payload.Search.IDRequest) stream | [.payload.Search.Response](#payload.Search.Response) stream |  |
| Insert | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Empty](#payload.Empty) |  |
| StreamInsert | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiInsert | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Empty](#payload.Empty) |  |
| Update | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Empty](#payload.Empty) |  |
| StreamUpdate | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiUpdate | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Empty](#payload.Empty) |  |
| Upsert | [.payload.Object.Vector](#payload.Object.Vector) | [.payload.Empty](#payload.Empty) |  |
| StreamUpsert | [.payload.Object.Vector](#payload.Object.Vector) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiUpsert | [.payload.Object.Vectors](#payload.Object.Vectors) | [.payload.Empty](#payload.Empty) |  |
| Remove | [.payload.Object.ID](#payload.Object.ID) | [.payload.Empty](#payload.Empty) |  |
| StreamRemove | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiRemove | [.payload.Object.IDs](#payload.Object.IDs) | [.payload.Empty](#payload.Empty) |  |
| GetObject | [.payload.Object.ID](#payload.Object.ID) | [.payload.Backup.MetaVector](#payload.Backup.MetaVector) |  |
| StreamGetObject | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Backup.MetaVector](#payload.Backup.MetaVector) stream |  |

 



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

