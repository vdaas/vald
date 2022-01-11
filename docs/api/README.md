# Vald API Overview

Vald provides 6 kinds of API for handling vector with Vald cluster.<br>
Vald provide 2 API interface: gRPC and REST API.
Using **gRPC** is preferred for better performance.

The APIs overview tables is here:

|    Service     | Description                                                                                         | API NAMES                                                                                                                                     | LINK                                 |
| :------------: | :-------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------- |
| Insert Service | Insert new vector(s) into the Vald Agent Pods                                                       | [Insert](../api/insert.md#insert-rpc)<br>[StreamInsert](../api/insert.md#streaminsert-rpc)<br>[MultiInsert](../api/insert.md#multiinsert-rpc) | [Vald Insert APIs](../api/insert.md) |
| Update Service | Update the exists vector(s) in the Vald Agent Pods                                                  | [Update](../api/update.md#update-rpc)<br>[StreamUpdate](../api/update.md#streamupdate-rpc)<br>[MultiUpdate](../api/update.md#multiupdate-rpc) | [Vald Update APIs](../api/update.md) |
| Upsert Service | Update the exists vector(s) in the Vald Agent Pods or Insert new vector(s) into the Vald Agent Pods | [Upsert](../api/upsert.md#upsert-rpc)<br>[StreamUpsert](../api/upsert.md#streamupsert-rpc)<br>[MultiUpdate](../api/upsert.md#multiupsert-rpc) | [Vald Upsert APIs](../api/upsert.md) |
