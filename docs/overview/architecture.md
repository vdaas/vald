# Vald Architecture <!-- omit in toc -->

This document describes the high-level architecture design of Vald and explains each component in Vald.

## Table of Contents <!-- omit in toc -->

- [Overview](#overview)
- [Data Flow](#data-flow)
  - [Insert](#insert)
  - [Search](#search)

## Overview

Vald uses a cloud-native architecture focusing on [Kubernetes](https://kubernetes.io/).
Some components in Vald use Kubernetes API to control the behavior of distributed vector indexes.
Before reading this document, you need to have some understanding of the basic idea of cloud-native architecture and Kubernetes.

The below image is Vald's architecture.

<img src="../../design/Vald Future Architecture Overview.svg" />

We will explain this image in the following section.

## Data Flow

This section describes the data flow inside Vald and how Vald's vector indexes are stored.
This is the most important part for the users to understand Vald.

### Insert

<img src="../../assets/docs/insert_flow.svg" />

When the user inserts data into Vald:

1. Vald Ingress receives the request from the user. The request includes the vector and the vector ID.
2. Vald Ingress will forward the request to the Vald Filter Gateway to pre-process the request data.
3. Vald Filter Gateway will forward the request to the user-defined Vald Ingress Filter. After the Vald Ingress Filter received the request, it will perform the pre-processing logic defined by the user, for example, padding the vector to match the vector dimension in Vald.
4. After the request is processed by the user-defined Vald Ingress Filter, the result will return to the Vald Filter Gateway.
5. Vald Filter Gateway will forward the processed data to the Vald Meta Gateway. Vald Meta Gateway will generate the UUID for each vector for internal use and the UUID will be mapped to the vector ID from the user's request. The reason of using UUID instead of vector ID is because the vector ID may be too long and it may increase the memory usage in Vald Agent.
6. Vald Meta Gateway will forward the request with the UUID to the Vald Backup Gateway, which will process the backup logic in 14-16 to prevent the data lost in Vald.
7. Vald Backup Gateway will forward the request to Vald LB Gateway. Vald LB Gateway will determine which Vald Agent(s) to process the request based on the resource usage of the nodes and pods, and the number of vector replicas.
8. Vald LB Gateway will forward the UUID and the vector data to the selected Vald Agents in parallel. Vald Agent will insert the vector and UUID in an on-memory vector queue. A vector queue will be committed to an ANN graph index by a `CreateIndex` instruction executed by the Vald Index Manager.
9. If Vald Agent successfully inserts the request data, it will return success to the Vald LB Gateway.
10. After Vald LB Gateway receives success from the selected Vald Agents, it will respond the IP addresses of all selected Vald Agents to the Vald Backup Gateway.
11. Vald Backup Gateway returns success to Vald Meta Gateway.
12. Vald Meta Gateway will forward the UUID(s) and vector ID(s) to the Vald Meta.
13. Vald Meta will store the UUID(s) and vector ID(s) that were successfully processed by the Vald Agent(s) to the persistent layer such as Redis, Cassandra, MySQL, etc.
14. Vald Backup Gateway will asynchronously send all the inserted the data (vector(s), vector ID(s), UUID(s) and IP address(es)) to the Vald Compressor. Vald Compressor will compress the vector data asynchronously to reduce the size of the vector data.
15. Vald Compressor will forward the data (compressed vector(s), vector ID(s), UUID(s) and IP address(es)) to the Vald Backup Manager.
16. Vald Backup Manager will store all of the data to the persistent layer such as MySQL, Cassandra, etc., to prevent the data lost in Vald.
17. Vald Meta Gateway will return success to the Vald Filter Gateway.
18. Vald Filter Gateway will return success to the Vald Ingress.

### Search

<img src="../../assets/docs/search_flow.svg" />

When the user searches a vector from Vald:

1. Vald Ingress receives a search request from the user. Vald provides 2 searching interfaces to the user, the user can search by vector or search by the vector ID.
2. Vald Ingress will forward the request to the Vald Filter Gateway to pre-process the request data.
3. Vald Filter Gateway will forward the request to the user-defined Vald Ingress Filter. After the Vald Ingress Filter received the request, it will perform the pre-processing logic defined by the user, for example, padding the vector to match the vector dimension in Vald.
4. After the request is processed by the user-defined Vald Ingress Filter, the result will return to the Vald Filter Gateway.
5. Vald Filter Gateway will forward the request to the Vald Meta Gateway. Vald Meta Gateway is used to resolve the internal used UUID to the user inserted vector ID in step 10-11.
6. Vald Meta Gateway will forward the request to the Vald LB Gateway. Vald LB Gateway will preform the post-processing of the result in step 9 after the Vald Agent(s) return in step 8.
7. Vald LB Gateway will forward the request to all Vald Agents in parallel. Each Vald Agent will search the _k_ nearest neighbor vectors in an on memory graph index.
8. Vald Agent returns the searching result to the Vald LB Gateway. The searching result includes the UUID, the vector distance, and the vector. The number of the result will be the same as requested.
9. Vald LB Gateway will aggregate all searching results from all Vald Agents, rank the result by the vector distance, and return the ranked result to the Vald Meta Gateway.
10. Vald Meta Gateway will forward the searching result to the Vald Meta to resolve the user-defined vector IDs from the internal used UUIDs.
11. Vald Meta will perform a search for the Vector IDs based on the internal used UUIDs.
12. Vald Meta will return the Vector IDs to the Vald Meta Gateway.
13. Vald Meta Gateway will combine the vectors and the vector IDs from the searching result and return to the Vald Filter Gateway.
14. Vald Filter Gateway will forward the request to the user-defined Vald Egress Filter to filter the final result. For example exclude the specific type of the result from the vector ID.
15. Vald Egress Filter will return the filtered result to the Vald Filter Gateway.
16. Vald Filter Gateway will return the final result to the Vald Ingress.

<!-- ### Update -->

<!-- ### Delete -->
