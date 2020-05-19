# About Vald <!-- omit in toc -->

This document gives an overview of what is Vald and what you can do with Vald. 

## Table of Contents <!-- omit in toc -->

- [What is Vald?](#what-is-vald)
  - [What Vald can do?](#what-vald-can-do)
  - [Use cases](#use-cases)
- [Why Vald?](#why-vald)
- [How does Vald work?](#how-does-vald-work)
- [Try Vald](#try-vald)

### What is Vald?

<!-- copied -->
Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine.
Vald is designed and implemented based on Cloud-Native architecture.
It uses the fastest ANN Algorithm NGT to search neighbors.
Vald has automatic vector indexing and index backup, and horizontal scaling which made for searching from billions of feature vector data.
Vald is easy to use, feature-rich and highly customizable as you needed.

#### What Vald can do?

<!--  copied -->

- Asynchronize Auto Indexing
    - Usually the graph requires locking during indexing, which cause stop-the-world. But Vald uses distributed index graph so it continues to work during indexing.

- Customizable Ingress/Egress Filtering
    - Vald implements it's own highly customizable Ingress/Egress filter.
    - Which can be configured to fit the gRPC interface.
        - Ingress Filter: Ability to Vectorize through filter on request.
        - Egress Filter: rerank or filter the searching result with your own algorithm.

- Cloud-native based vector searching engine
    - Horizontal scalable on memory and cpu for your demand.

- Auto Indexing Backup
    - Vald has auto index backup feature using MySQL + Redis or Cassandra which enables disaster recovery.

- Distributed Indexing
    - Vald distribute vector index to multiple agent, each agent stores different index.

- Index Replication
    - Vald stores each index in multiple agents which enables index replicas.
    - Automatically rebalance the replica when some Vald agent goes down.

- Easy to use
    - Vald can be easily installed in a few steps.

- Highly customizable
    - You can configure the number of vector dimension, the number of replica and etc.

- Multi language supported
    - Golang, Java, Nodejs and python is supported.

#### Use cases

Vald supports similarity searching.

- Related image search
- Speech recognition
- Everything you can vectorize :)

### Why Vald?

Vald is based on Kubernetes and Cloud-Native architecture, which means Vald is highly scalable.
You can easily scale Vald by changing Vald's configuration.

Vald uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors by default, but users can switch to another vector searching engine in Vald to support the best performance for your use case.

Also, Vald supports auto-healing, to reduce running and maintenance costs. Vald implements the backup mechanism to support disaster recovery.
Whenever one of the Vald Agent instances is down, the new Vald Agent instance will be created automatically and the data will be recovered automatically.

### How does Vald work?

Vald implements its custom resource and custom controller to integrate with Kubernetes.
You can take all the benefits from Kubernetes.

Please refer to the [architecture overview](./architecture.md) for more details about the architecture and how each component in Vald works together.

### Try Vald

Please refer to [Get Started](../get-started/get-started.md) to try Vald :)
