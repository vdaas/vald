# Vald API Overview

Vald provides 6 kinds of API for handling vector with Vald cluster.<br>
Vald provide 2 API interface: gRPC and REST API.
Using **gRPC** is preferred for better performance.

The APIs overview tables is here:

|    Service     | Description                                   | API NAMES                                                                                                                                                | LINK             |
| :------------: | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------- |
| Insert Service | Insert new vector(s) into the Vald Agent Pods | [Insert](../api/insert.md#insert-rpc)<br>[StreamInsert](../api/insert.md#streaminsert-rpcrecommended)<br>[MultiInsert](../api/insert.md#multiinsert-rpc) | [Vald Insert APIs](../api/insert.md) |
