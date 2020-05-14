# About Vald <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [What is Vald?](#what-is-vald)
  - [Use cases](#use-cases)
- [Why Vald?](#why-vald)
- [How does Vald work?](#how-does-vald-work)
- [Try Vald](#try-vald)

### What is Vald?

Please refer to [this page](https://github.com/vdaas/vald#what-is-vald) for more details.

#### Use cases

Vald support similarity searching.

- Related image search
- Speech recognition
- Everything you can vectorize :)

### Why Vald?

Vald is based on Kubernetes and Cloud-Native architecture, which means Vald is highly scalable. You can easily scale Vald by changing Vald's configuration.

Vald uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors by default, but users can switch to another vector searching engine in Vald to support the best performance for your use case.

Also, Vald support auto-healing, to reduce running and maintenance cost. Vald implements the backup mechanism to support diaster recovery. Whenever Vald Agent is down, the new Vald Agent instance will be created automatically and the data will be recovered automatically.

### How does Vald work?

Vald implements its custom resource and custom controller to integrate with Kubernetes; you can take all the benefits from Kubernetes.

Please refer to the [architecture overview](./architecture.md) for more details about the architecture and how each component in Vald works together.

### Try Vald

Please refer to [Get Started](./get-started.md) to try Vald :)
