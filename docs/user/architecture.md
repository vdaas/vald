# Vald Architecture <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [Overview](#overview)
- [Data Flow](#data-flow)
  - [Insert](#insert)
  - [Update](#update)
  - [Delete](#delete)
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

This document describes the high-level architecture design of Vald and explains each component in Vald. We need these components to support scalability, high performance and auto-healing in Vald.

<img src="../../design/Vald Future Architecture Overview.svg" />

Vald is based on [Kubernetes](https://kubernetes.io/) architecture. Before you read this document, you must understand the basic concept of Kubernetes.

## Data Flow

### Insert

<img src="./insert_search_flow.svg" />

When the user inserts data into Vald:

```
C1. Vald Ingress receives the request from the user. The request includes the vector and the vector ID.
C2. Vald Ingress will forward the request to Vald Meta Gateway. After the Vald Meta Gateway receive the request, the UUID will be generated for internal use.
C3. Vald Meta Gateway will forward the request with UUID to Vald Backup Gateway.
C4. Vald Backup Gateway will forward the request to Vald LB Gateway.
C5. Vald LB Gateway will decide which Vald Agent instance to process the request base on the node resource usage and forwards the request to the decided Vald Agent.
C6. Vald Agent returns success and the corresponding Vald Agent IP address to Vald LB Gateway.
C7. Vald LB Gateway return success and the IP address to Vald Backup Gateway.
C8. Vald Backup Gateway returns success to Vald Meta Gateway.
C9. Vald Meta Gateway will send the metadata including the UUID and the vector ID to the Vald Meta.
C10. Vald Meta will store the metadata to the persistent layer.
C11. Vald Backup Gateway will send all the data (vector, vector ID and UUID) the Vald Compressor.
C12. Vald Compressor will compress the data received and send to the Vald Backup Manager.
C13. Vald Backup Manager will store the compressed data to the persistent layer.
```

### Update

### Delete

### Search

S1. Vald Ingress receives the search request from the user. The request includes the vector or the vector ID.
S2. Vald Ingress will forward the request to Vald Filter Gateway.
S3. Vald Filter Gateway will forward the request to Vald Ingress Filter.
S4. Vald Ingress Filter will perform the filtering and return the filtering result to the Vald Filter Gateway.
S5. Vald Filter Gateway will forward the request to the Vald Meta Gateway.
S6. Vald Meta Gateway will forward the request to the Vald Backup Gateway.
S7. Vald Backup Gateway will forward the request to the Vald LB Gateway.
S8. Vald LB Gateway will decide which Vald Agent instance to process the request base on the node resource usage and forwards the request to the decided Vald Agent.
S9. Vald Agent returns the searching result to the Vald LB Gateway. The searching result includes the UUID.
S10. Vald LB Gateway returns the searching result to Vald Backup Gateway.
S11. Vald Backup Gateway returns the searching result to the Vald Meta Gateway.
S12. Vald Meta Gateway will forward the searching result to the Vald Meta.
S13. Vald Meta will perform a search for the Vector ID base on the UUID and return the Vector ID to the Vald Meta Gateway.
S14. Vald Meta Gateway returns the searching result with the vector ID to the Vald Filter Gateway.
S15. Vald Filter Gateway will forward the request to Vald Egress Gateway to filter the final result.

## Components

### Vald Filter

Vald Filter is an optional functionality in Vald.
User can implement the custom filtering logic and integrate with Vald.

Vald Filter provides the following functionalities.

- Custom filter base on request query
- Custom filter for the searching result

#### Vald Ingress Filter

Vald Ingress Filter filters the incoming request before processing it.

Users can implement custom filtering logic such as changing the vectors or filtering based on user ID.

#### Vald Egress Filter

Vald Egress Filter filters the response before sending it to the user.

This component can reorder the searching result from multiple Vald Agents based on the user-defined ranking.

#### Vald Filter Gateway

Vald Filter Gateway forwards the request to Vald Ingress Filter before processing it and forwards the request to the Vald Egress Filter before returning the searching result to the user.

### Vald Metadata

In Vald, metadata is the vector data and the corresponding additional data to represent the set of the searching criteria and the result.

Vald Metadata includes the user inputted metadata(vector ID) and the vector, and the internal generated UUID.

#### Vald Meta Gateway

The main responsibility of the Vald Meta Gateway is to process the Vald metadata and forwards the information to Vald Backup Gateway.

It will perform the following action:

1. Return error if the user has already input the same vector in Vald
1. Generate the corresponding UUID for internal use.
1. Forward the metadata (vec_id and UUID) request to the Vald Meta Agent.
1. Forward the vector information (vec_id, vector, and UUID) to Vald Backup Gateway.

#### Vald Meta

Vald Meta is the agent to process the CRUD request of the metadata (vec_id and UUID). Users can configure which data source to be used in Vald Meta (for example Redis or Cassandra).

### Vald Backup

To support auto-healing functionality and increase performance during disaster recovery, Vald implements the backup mechanism.

#### Vald Compressor

Vald Compressor compresses all of the data (metadata and the vector data) and sends to the Vald Backup Manager to process the backup request.

#### Vald Backup Manager

Vald Backup Manager processes the CRD request of the backup request and handles the compressed metadata. Users can configure which data source to be used in Vald Meta (for example Redis or Cassandra).

#### Vald Backup Gateway

Vald Backup Gateway will forward the backup request to the Vald LB Gateway. It also forwards to Vald Compressor asynchronously with metadata.

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

Vald Agent is the core of the Vald. By default Vald use [NGT](https://github.com/yahoojapan/NGT) to provide API for users to insert/update/delete/search vectors.

Each Vald Agent pod holds different vector dimension space, which is constructed by insert/update vectors for searching approximate vectors.

<!-- (todo: diagram to explain different vector dimension space) -->

When you request searching with your vector in Vald, each Vald Agent returns different _k_-nearest neighbors' vectors which are similar to the searching vector.

#### Vald Agent Scheduler

Vald Agent Scheduler is the scheduler of the Vald Agent.
It implements it's own custom scheduling logic to increase the scalability of the Vald Agent.

It schedules Vald Agent base on the Node CPU and memory usage, and the amount of the indexes.

#### Vald Index Manager

Vald Index Manager controls the timing of the indexing inserted vectors on the Vald Agent. The index is used to increase the performance of the search action.

It retrieves the active Vald Agent pods from the Vald Discoverer and triggers the indexing action on each Vald Agent.

### Vald Replication Manager

Vald Replication Manager manages the healthiness of the Vald Agent. When the pod is dead, Vald Replication Manager will auto recover the cache to keeps the reliability of the service.

#### Vald Replication Manager Agent

Vald Replication Manager Agent recovers the specific backup cache to the specific Vald Agent. It retrieves the target backup from the Vald Compressor and recovers it to the newly created Vald Agent.

#### Vald Replication Manager Controller

Vald Replication Manager Controller keeps track of the active Vald Agent pods. When the Vald Agent is dead, it will trigger the Vald Replication Manager Agent to recover the backup cache to the auto-healed pods from the backup.

### Kubernetes Components

Vald is base on the Kubernetes platform. In this section we will explain the Kubernetes component used in Vald and why we need them.

#### Kube-apiserver

Kube-apiserver is a component of Kubernetes. The main responsibility of Kube-apiserver in Vald is to provide node resource information for Vald agent scalability.

For more information about Kube-apiserver, please refer to [the official document](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/).

#### Custom Resources

Custom Resources in Vald is a [Custom Resouce Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) implementation. It provides flexibility for users to control the Vald deployment such as pod startup sequence, etc.
