# Components <!-- omit in toc -->

This document will give you an overview of all the components developed and used in Vald.

## Table of Contents <!-- omit in toc -->

- [Vald Load Balancing](#vald-load-balancing)
  - [Vald LB Gateway](#vald-lb-gateway)
  - [Agent Discoverer](#agent-discoverer)
- [Vald Core Engine](#vald-core-engine)
  - [Vald Agent](#vald-agent)
  - [Vald Agent Scheduler](#vald-agent-scheduler)
  - [Vald Index Manager](#vald-index-manager)
- [Vald Replication Manager](#vald-replication-manager)
- [Kubernetes Components](#kubernetes-components)
  - [Kube-apiserver](#kube-apiserver)
  - [Custom Resources](#custom-resources)

## Vald Load Balancing

Load balancing is one of the important concepts in distributed computing, which means it distributes a set of tasks over a set of resources aiming for making the overall processing more efficient.
Vald implements its own load balancing controller.
Vald can load balance the request base on node resources.

### Vald LB Gateway

Vald LB Gateway loads balance the user request base on the node resources results from the Agent Discoverer.

### Agent Discoverer

Agent Discoverer discovers active Vald pods and the corresponding node's resources usage via [kube-apiserver](https://github.com/kubernetes/kubernetes/tree/master/cmd/kube-apiserver).

## Vald Core Engine

In this section, we will describe what is Vald Agent and the corresponding components to support Vald Agent.

### Vald Agent

Vald Agent provides functionalities to perform approximate nearest neighbor search.
Agent-NGT uses [yahoojapan/NGT](https://github.com/yahoojapan/NGT) as a core library.

Each Vald Agent pod has its own vector data space because only several Vald Agents are selected to be inserted/updated in a single insert/update request.

When searching a vector in Vald, each Vald Agent return different results of _k_-nearest neighbors depending on their index, and you'll get the merged result of them.

<img src="../../../assets/docs/vector_data_space_explain.svg" />

### Vald Agent Scheduler

Vald Agent Scheduler is the scheduler of the Vald Agent.
It implements it's own custom scheduling logic to increase the scalability of the Vald Agent.

It schedules Vald Agent base on the Node CPU and memory usage, and the amount of the indexes.

### Vald Index Manager

Vald Index Manager controls the timing of the indexing inserted vectors on the Vald Agent.
The index is used to increase the performance of the search action.

It retrieves the active Vald Agent pods from the Vald Discoverer and triggers the indexing action on each Vald Agent.

## Vald Replication Manager

Vald Replication Manager manages the healthiness of the Vald Agent.
When the pod is dead, Vald Replication Manager will recover the cache automatically to keeps the reliability of the service.

## Kubernetes Components

Vald is base on the Kubernetes platform.
In this section we will explain the Kubernetes component used in Vald and why we need them.

### Kube-apiserver

Kube-apiserver is a component of Kubernetes.
The main responsibility of Kube-apiserver in Vald is to provide node resource information for Vald agent scalability.

For more information about Kube-apiserver, please refer to [the official document](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/).

### Custom Resources

Custom Resources in Vald is a [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) implementation.
It provides flexibility for users to control the Vald deployment such as pod startup sequence, etc.
