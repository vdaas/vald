# Vald Architecture

## Overview

This document describe the high level architecture design of Vald and explain each component in Vald. We need to these components to support scalability, high performance and auto-healing in Vald.

<img src="../../design/Vald Future Architecture Overview.svg" />

Vald is based on [Kubernetes](https://kubernetes.io/) architecture. Before you read this document you must understand the basic concept of Kubernetes.

## Data Flow

### Insert

When user insert data into Vald:
1 The request will go through the Vald Ingress
    1. The Vald Ingress will log the request
2.1. Vald Ingress will forward the request to Vald Meta Gateway
    1. Vald Meta Gateway will load balance the request
    1. The request data (metadata) will then forward to Vald Meta
        1. The metadata will store to the presistent layer
3.1 Vald Meta Gateway will forward the request to the Vald Backup Gateway
    1. The Backup Gateway
3.2 Vald Meta Gateway will forward the request to Vald Meta
4.1 Vald Backup Gateway will forward the request to the Vald LB Gateway
4.2 Vald Backup Gateway will forward the request to the Vald Compressor
5.1 Vald LB Gateway will forward the request to Agent Discoverer
5.2 Vald LB Gateway will forward the request to Vald Agent

### Update

### Delete

### Search

## Components

### Vald Agent

Vald Agent is the process engine of Vald. In this section we will describe what is Vald Agent and the corresponding components to support Vald Agent.

#### Vald Agent

Vald Agent is the core of the Vald. By default Vald use [NGT](https://github.com/yahoojapan/NGT) to provide API for users to insert/update/delete/search vectors.

#### Vald Agent Scheduler

Vald Agent Scheduler is the scheduler of the Vald Agent. It schedules Vald Agent base on the Node CPU and memory usage.

#### Vald Index Manager

Vald Index Manager manages the index of vector data in Vald Agent.

### Vald Load Balancing

Load balancing is very important concept in distributed computing, which means the distribute a set of task over set of resources aiming for making the overall processing more efficient.
In Vald, we implement our own load balancing controller. Vald can load balance the request base on node resources.

#### Agent Discoverer

Agent Discoverer discovers Vald pods and the corresponding node resource usage. It talks to the Kube-API and get the corresponding node information.

#### Vald LB Gateway

Vald LB Gateway load balance the user request base on the node resources result from the agent discoverer.

### Vald Metadata

In Vald, metadata is the vector data and the corresponding addition data to represent the set of the searching criteria and the result.

#### Vald Meta

Vald Meta is the agent to process metadata of the vector data. It will insert the metadata to the presistent layer.

#### Vald Meta Gateway

Vald Meta Gateway load balance the metadata request to the Vald Meta Agent.

### Vald Backup

To support auto-healing and incresease the performance during disaster recovery, Vald implement the backup mechanism.

#### Vald Backup Gateway

Vald Backup Gateway load balance the backup request to the Vald Compressor to handle vector backup request.

#### Vald Compressor

Vald Compressor compress the vector data and send to the Vald Backup Manager to process the backup request.

#### Vald Backup Manager

Vald Backup Manager process the backup request and store the vector data to the presistent layer.

### Vald Replication Manager

Vald replication manager manages the Vald Agent replicates. It auto-scale the Vald agent base on the resource usage on the node.

#### Vald Replication Manager Agent

#### Vald Replication Manager Controller


### Vald Filter

Vald Filter have 2 main functionality.
1. Filter request query
1. Filter response data

#### Vald Ingress Filter
#### Vald Egress Filter

Vald Egress Filter filter the response before sending to the user. This component will reorder the response data from set of the Vald Agent base on the ranking and then response the number of data user want.

#### Vald Filter Gateway

Vald Filter Gateway load balance the filter request.

### Kubernetes Components

Vald is base on Kubernetes platform. In this section we will explain the Kubernetes component used in Vald and why we need them. 

#### Kube-API Server

#### Custom Resources
