# Vald Architecture <!-- omit in toc -->
This document describes the high-level architecture design of Vald and explains each component in Vald.
## Table of Contents <!-- omit in toc -->

- [Overview](#overview)
- [Data Flow](#data-flow)
  - [Insert](#insert)
  - [Search](#search)
- [Components](#components)
  - [Vald Filter](#vald-filter)
    - [Vald Ingress Filter](#vald-ingress-filter)
    - [Vald Egress Filter](#vald-egress-filter)
    - [Vald Filter Gateway](#vald-filter-gateway)
  - [Vald Metadata](#vald-metadata)
    - [Vald Meta Gateway](#vald-meta-gateway)
    - [Vald Meta](#vald-meta)
  - [Vald Backup](#vald-backup)
    - [Vald Compressor](#vald-compressor)
    - [Vald Backup Manager](#vald-backup-manager)
    - [Vald Backup Gateway](#vald-backup-gateway)
  - [Vald Load Balancing](#vald-load-balancing)
    - [Vald LB Gateway](#vald-lb-gateway)
    - [Agent Discoverer](#agent-discoverer)
  - [Vald Core Engine](#vald-core-engine)
    - [Vald Agent](#vald-agent)
    - [Vald Agent Scheduler](#vald-agent-scheduler)
    - [Vald Index Manager](#vald-index-manager)
  - [Vald Replication Manager](#vald-replication-manager)
    - [Vald Replication Manager Agent](#vald-replication-manager-agent)
    - [Vald Replication Manager Controller](#vald-replication-manager-controller)
  - [Kubernetes Components](#kubernetes-components)
    - [Kube-apiserver](#kube-apiserver)
    - [Custom Resources](#custom-resources)

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

## Components

### Vald Filter

Vald Filter is an optional functionality in Vald.
User can implement the custom filtering logic and integrate with Vald.

Vald Filter provides the following functionalities.

- Custom filter based on request query
- Custom filter for the searching result

#### Vald Ingress Filter

Vald Ingress Filter filters the incoming request before processing it.

Users can implement custom filtering logic such as changing the vectors or filtering based on user ID.

#### Vald Egress Filter

Vald Egress Filter filters the response before sending it to the user.

This component can reorder the searching result from multiple Vald Agents based on the user-defined ranking.

#### Vald Filter Gateway

Vald Filter Gateway forwards the request to Vald Ingress Filter before processing it and forwards the response to the Vald Egress Filter before returning the searching result to the user.

### Vald Metadata

In Vald, metadata consists of the vector data and the corresponding additional data to represent the set of the searching criteria and the result.

Vald Metadata includes the user inputted metadata(vector ID) and the vector, and the internal generated UUID.

#### Vald Meta Gateway

The main responsibility of the Vald Meta Gateway is to process the Vald metadata and to forward the information to Vald Backup Gateway.

It will perform the following action:

1. Return error if the user has already input the same vector in Vald
1. Generate the corresponding UUID for internal use.
1. Forward the vector ID and UUID request to the Vald Meta.
1. Forward the vector information (vector ID, vector, and UUID) to Vald Backup Gateway.

#### Vald Meta

Vald Meta is the agent to process the CRUD request of the metadata (vector ID and UUID).
Users can configure which data source to be used in Vald Meta (for example Redis or Cassandra).

### Vald Backup

To support auto-healing functionality and increase performance during disaster recovery, Vald implements the backup mechanism.

#### Vald Compressor

Vald Compressor compresses the vector data and sends to the Vald Backup Manager to process the backup request.

#### Vald Backup Manager

Vald Backup Manager processes the Create/Read/Delete request of the backup request and handles the compressed metadata. Users can configure which data source to be used in Vald Meta (for example Redis or Cassandra).

#### Vald Backup Gateway

Vald Backup Gateway will forward the backup request to the Vald LB Gateway.
It also forwards to Vald Compressor asynchronously with metadata.

### Vald Load Balancing

Load balancing is one of the important concepts in distributed computing, which means it distributes a set of tasks over a set of resources aiming for making the overall processing more efficient.
Vald implements its own load balancing controller.
Vald can load balance the request base on node resources.

#### Vald LB Gateway

Vald LB Gateway loads balance the user request base on the node resources results from the Agent Discoverer.

#### Agent Discoverer

Agent Discoverer discovers active Vald pods and the corresponding node's resources usage via [kube-apiserver](https://github.com/kubernetes/kubernetes/tree/master/cmd/kube-apiserver).

### Vald Core Engine

In this section, we will describe what is Vald Agent and the corresponding components to support Vald Agent.

#### Vald Agent

Vald Agent provides functionalities to perform approximate nearest neighbor search.
Agent-NGT uses [yahoojapan/NGT](https://github.com/yahoojapan/NGT) as a core library.

Each Vald Agent pod has its own vector data space because only several Vald Agents are selected to be inserted/updated in a single insert/update request.

When searching a vector in Vald, each Vald Agent return different results of _k_-nearest neighbors depending on their index, and you'll get the merged result of them.
<img src="../../assets/docs/vector_data_space_explain.svg" />

#### Vald Agent Scheduler

Vald Agent Scheduler is the scheduler of the Vald Agent.
It implements it's own custom scheduling logic to increase the scalability of the Vald Agent.

It schedules Vald Agent base on the Node CPU and memory usage, and the amount of the indexes.

#### Vald Index Manager

Vald Index Manager controls the timing of the indexing inserted vectors on the Vald Agent.
The index is used to increase the performance of the search action.

It retrieves the active Vald Agent pods from the Vald Discoverer and triggers the indexing action on each Vald Agent.

### Vald Replication Manager

Vald Replication Manager manages the healthiness of the Vald Agent.
When the pod is dead, Vald Replication Manager will recover the cache automatically to keeps the reliability of the service.

#### Vald Replication Manager Agent

Vald Replication Manager Agent recovers the specific backup cache to the specific Vald Agent.
It retrieves the target backup from the Vald Compressor and recovers it to the newly created Vald Agent.

#### Vald Replication Manager Controller

Vald Replication Manager Controller keeps track of the active Vald Agent pods.
When the Vald Agent is dead, it will trigger the Vald Replication Manager Agent to recover the backup cache to the auto-healed pods from the backup.

### Kubernetes Components

Vald is base on the Kubernetes platform.
In this section we will explain the Kubernetes component used in Vald and why we need them.

#### Kube-apiserver

Kube-apiserver is a component of Kubernetes.
The main responsibility of Kube-apiserver in Vald is to provide node resource information for Vald agent scalability.

For more information about Kube-apiserver, please refer to [the official document](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/).

#### Custom Resources

Custom Resources in Vald is a [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) implementation.
It provides flexibility for users to control the Vald deployment such as pod startup sequence, etc.
