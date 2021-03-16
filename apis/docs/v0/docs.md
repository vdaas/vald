# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [apis/proto/filter/egress/egress_filter.proto](#apis/proto/filter/egress/egress_filter.proto)
    - [EgressFilter](#egress_filter.EgressFilter)
  
- [apis/proto/gateway/vald/vald.proto](#apis/proto/gateway/vald/vald.proto)
    - [Vald](#vald.Vald)
  
- [apis/proto/payload/payload.proto](#apis/proto/payload/payload.proto)
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
    - [Object.Blob](#payload.Object.Blob)
    - [Object.Distance](#payload.Object.Distance)
    - [Object.ID](#payload.Object.ID)
    - [Object.IDs](#payload.Object.IDs)
    - [Object.Location](#payload.Object.Location)
    - [Object.Locations](#payload.Object.Locations)
    - [Object.Vector](#payload.Object.Vector)
    - [Object.Vectors](#payload.Object.Vectors)
    - [Replication](#payload.Replication)
    - [Replication.Agents](#payload.Replication.Agents)
    - [Replication.Rebalance](#payload.Replication.Rebalance)
    - [Replication.Recovery](#payload.Replication.Recovery)
    - [Search](#payload.Search)
    - [Search.Config](#payload.Search.Config)
    - [Search.IDRequest](#payload.Search.IDRequest)
    - [Search.MultiIDRequest](#payload.Search.MultiIDRequest)
    - [Search.MultiRequest](#payload.Search.MultiRequest)
    - [Search.ObjectRequest](#payload.Search.ObjectRequest)
    - [Search.Request](#payload.Search.Request)
    - [Search.Response](#payload.Search.Response)
    - [Search.Responses](#payload.Search.Responses)
  
- [Scalar Value Types](#scalar-value-types)



<a name="apis/proto/filter/egress/egress_filter.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/filter/egress/egress_filter.proto


 

 

 


<a name="egress_filter.EgressFilter"></a>

### EgressFilter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Filter | [.payload.Search.Response](#payload.Search.Response) | [.payload.Search.Response](#payload.Search.Response) |  |
| StreamFilter | [.payload.Object.Distance](#payload.Object.Distance) stream | [.payload.Object.Distance](#payload.Object.Distance) stream |  |

 



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

 



<a name="apis/proto/payload/payload.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## apis/proto/payload/payload.proto



<a name="payload.Backup"></a>

### Backup







<a name="payload.Backup.Compressed"></a>

### Backup.Compressed







<a name="payload.Backup.Compressed.MetaVector"></a>

### Backup.Compressed.MetaVector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
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







<a name="payload.Object.Blob"></a>

### Object.Blob



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| object | [bytes](#bytes) |  |  |






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






<a name="payload.Object.Location"></a>

### Object.Location



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| ips | [string](#string) | repeated |  |






<a name="payload.Object.Locations"></a>

### Object.Locations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| locations | [Object.Location](#payload.Object.Location) | repeated |  |






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
| request_id | [string](#string) |  |  |
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






<a name="payload.Search.MultiIDRequest"></a>

### Search.MultiIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.IDRequest](#payload.Search.IDRequest) | repeated |  |






<a name="payload.Search.MultiRequest"></a>

### Search.MultiRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requests | [Search.Request](#payload.Search.Request) | repeated |  |






<a name="payload.Search.ObjectRequest"></a>

### Search.ObjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| object | [bytes](#bytes) |  |  |
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
| request_id | [string](#string) |  |  |
| results | [Object.Distance](#payload.Object.Distance) | repeated |  |






<a name="payload.Search.Responses"></a>

### Search.Responses



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| responses | [Search.Response](#payload.Search.Response) | repeated |  |





 

 

 

 



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

