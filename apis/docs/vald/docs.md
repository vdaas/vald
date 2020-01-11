# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [vald.proto](#vald.proto)
  
  
  
    - [Vald](#vald.Vald)
  

- [Scalar Value Types](#scalar-value-types)



<a name="vald.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vald.proto


 

 

 


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
| Remove | [.payload.Object.ID](#payload.Object.ID) | [.payload.Empty](#payload.Empty) |  |
| StreamRemove | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Empty](#payload.Empty) stream |  |
| MultiRemove | [.payload.Object.IDs](#payload.Object.IDs) | [.payload.Empty](#payload.Empty) |  |
| GetObject | [.payload.Object.ID](#payload.Object.ID) | [.payload.Backup.MetaVector](#payload.Backup.MetaVector) |  |
| StreamGetObject | [.payload.Object.ID](#payload.Object.ID) stream | [.payload.Backup.MetaVector](#payload.Backup.MetaVector) stream |  |

 



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

