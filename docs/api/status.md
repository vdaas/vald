# Response Status Code

This page describes each status code from API response.

## Status Codes

This table shows the main status code and its name using Vald.
Below sections describes the meaning of each code and the reason why API returns.

| code | name                                      |
| :--: | :---------------------------------------- |
|  0   | [OK](#OK)                                 |
|  3   | [INVALID_ARGUMENT](#INVALID_ARGUMENT)     |
|  4   | [DEADLINE_EXCEEDED](#DEADLINE_EXCEEDED)   |
|  5   | [NOT_FOUND](#NOT_FOUND)                   |
|  6   | [ALREADY_EXISTS](#ALREADY_EXISTS)         |
|  8   | [RESOURCE_EXHAUSTED](#RESOURCE_EXHAUSTED) |
|  13  | [INTERNAL](#INTERNAL)                     |
|  14  | [UNAVAILABLE](#UNAVAILABLE)               |

## OK

`OK` means complete process with success.

Services that return this code are all services.

## INVALID_ARGUMENT

`INVALID_ARGUMENT` means the something wrong in the request configuration.

Services that return status are all services.
If you get this code, please verify your request is correct.

## DEADLINE_EXCEEDED

`DEADLINE_EXCEEDED` returns when the process ends due to timeout.

Services that returns status are:
- [Object Service(only Exists RPC)](../api/object.md#Exists-RPC)
- [Insert Service](../api/insert.md)
- [Remove Service](../api/insert.md)
- [Search Service](../api/search.md)

The timeout configuration is on the Vald cluster side.
If it appears constantly, you need to review the cluster settings.
However, it appears only when using search service, you can overwrite timeout configuration by lengthening the time setting in the search config.

## NOT_FOUND

`NOT_FOUND` appears when there is no result corresponding to the request.

The example cases are:

    - No search result when using SearchById api
    - No index data corresponding to the request vector when using Update/Exists/GetObject api

Services that returns status are:
- [Object Service](../api/object.md)
- [Update Service](../api/update.md)
- [Remove Service](../api/insert.md)
- [Search Service](../api/search.md)

There are two reasons.
One is there is no index data in Vald Agent components or index process is running in.
When Vald Agent component runs the index process (createIndex/saveIndex), the any process won't run and it will return with no result.

The other one, which occurs using search / update / remove service, is the request query vector or id is wrong.
Especially, both of update service and remove service requires the ID of vector Vald Agent component already indexed.

## ALREADY_EXISTS

`ALREADY_EXISTS` means that Vald Agent component already index the vector same as the query vector when set `skip_strict_exist_check` as `true` in request config.

Services that returns status are:
- [Insert Service](../api/insert.md)
- [Update Service](../api/update.md)

The way to avoid it, you have to change the query vector with `skip_strict_exist_check` is `true` or set `skip_strict_exist_check` as `false` instead of change the query vector.

## RESOURCE_EXHAUSTED

`RESOURCE_EXHAUSTED` means the gRPC message size is bigger than limit (default is 4MB).

Services that returns status are all services.
The most case in the Vald is the query vector is too large in other word setting vector dimension size it too large.

## INTERNAL

`INTERNAL` appears when some wrong happens in the Vald cluster.

Services that returns status are all services.
If you get it, please verify the state of Vald cluster.

## UNAVAILABLE

`UNAVAILABLE` means the gRPC message cannot reach to the Vald cluster.

You need to verify whether the Vald cluster is running and host and port is correct.
