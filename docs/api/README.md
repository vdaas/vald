# Vald API Overview

Vald provides 6 kinds of API for handling vector with Vald cluster.<br>
There are 2 choices using gRPC or REST API.
Vald recommends to use **gRPC** in the sense of high performance.

The APIs overview tables is here:

|    Service     | Description                                   | API NAMES                                                                                                                                                | LINK             |
| :------------: | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------- |
| Insert Service | Insert new vector(s) into the Vald Agent Pods | [Insert](../api/insert.md#insert-rpc)<br>[StreamInsert](../api/insert.md#streaminsert-rpcrecommended)<br>[MultiInsert](../api/insert.md#multiinsert-rpc) | [Vald Insert APIs](../api/insert.md) |
