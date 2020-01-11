# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [payload.proto](#payload.proto)
    - [Backup](#payload.Backup)
    - [Backup.Compressed](#payload.Backup.Compressed)
    - [Backup.Compressed.MetaVector](#payload.Backup.Compressed.MetaVector)
    - [Backup.Compressed.MetaVectors](#payload.Backup.Compressed.MetaVectors)
    - [Backup.GetVector](#payload.Backup.GetVector)
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
    - [Controll](#payload.Controll)
    - [Controll.CreateIndexRequest](#payload.Controll.CreateIndexRequest)
    - [Discoverer](#payload.Discoverer)
    - [Discoverer.Request](#payload.Discoverer.Request)
    - [Empty](#payload.Empty)
    - [Info](#payload.Info)
    - [Info.IPs](#payload.Info.IPs)
    - [Info.Index](#payload.Info.Index)
    - [Info.Server](#payload.Info.Server)
    - [Info.Servers](#payload.Info.Servers)
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
    - [Search](#payload.Search)
    - [Search.Config](#payload.Search.Config)
    - [Search.IDRequest](#payload.Search.IDRequest)
    - [Search.Request](#payload.Search.Request)
    - [Search.Response](#payload.Search.Response)
  
  
  
  

- [Scalar Value Types](#scalar-value-types)



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
| vector | [double](#double) | repeated |  |
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
| uuid | [string](#string) | repeated |  |






<a name="payload.Controll"></a>

### Controll







<a name="payload.Controll.CreateIndexRequest"></a>

### Controll.CreateIndexRequest



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
| node | [string](#string) |  |  |






<a name="payload.Empty"></a>

### Empty







<a name="payload.Info"></a>

### Info







<a name="payload.Info.IPs"></a>

### Info.IPs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) | repeated |  |






<a name="payload.Info.Index"></a>

### Info.Index



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [uint32](#uint32) |  |  |
| uncommitted_index | [uint32](#uint32) |  |  |
| uuids | [string](#string) | repeated |  |
| uncommitted_uuid | [string](#string) | repeated |  |






<a name="payload.Info.Server"></a>

### Info.Server



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| server | [Info.Server](#payload.Info.Server) |  |  |
| cpu | [double](#double) |  |  |
| mem | [double](#double) |  |  |






<a name="payload.Info.Servers"></a>

### Info.Servers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Servers | [Info.Server](#payload.Info.Server) | repeated |  |






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
| vector | [double](#double) | repeated |  |






<a name="payload.Object.Vectors"></a>

### Object.Vectors



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vectors | [Object.Vector](#payload.Object.Vector) | repeated |  |






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
| vector | [double](#double) | repeated |  |
| config | [Search.Config](#payload.Search.Config) |  |  |






<a name="payload.Search.Response"></a>

### Search.Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| results | [Object.Distance](#payload.Object.Distance) | repeated |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

