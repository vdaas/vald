# FAQ

This page answers the common questions.

Please refer to [the troubleshooting page](../user-guides/troubleshooting.md) when you encounter operation problems.

## Component

### What is the ingress filter?

The ingress filter is the user-defined component that connects to the `vald-filter-gateway` component.
For example, you can convert object data (e.g., image binary, text, etc.) to the vector using your ML models as pre-process of each request.

### What is the egress filter?

The egress filter is the user-defined component that connects to the vald-filter-gateway component.
You can use it for filtering search results, e.g., when you'd like to get a white T-shirt list, use it to remove other colors of T-shirts from search results that the vald-lb-gateway component returns.

### Is Vald Index Manager recommended using?

We recommend using it when you’d like to operate as a cluster.
It helps you to manage indexing timing for each Vald Agent.

## Custom options

### What are the pluggable options?

Vald has three pluggable options:

1. Backup with the external storage for Vald Agent
   - You can connect the external storage like S3, or GCS, or etc. to Vald Agent Sidecar component for backup.
1. Algorithm of the core engine for Vald Agent
   - We're going to add another algorithm in near future.
1. Filtering with filter gateway
   - you can filter the search results by own defined filter component by connecting to the filter gateway before returning the search result
   - you can convert object data to vector by own defined ingress filter component by connection to filter gateway before inserting

Vald has three pluggable options:
1. Backup with the external storage for Vald Agent
    - You can connect the external storage like S3, GCS, etc., to the Vald Agent Sidecar component for backup.
1. Algorithm of the core engine for Vald Agent
    - We're going to add another algorithm soon.
1. Filtering with the filter gateway
    - You can filter the search results by your own defined filter component by connecting to the filter gateway before returning the search result
    - You can convert object data to vector by own defined ingress filter component by connection to filter gateway before inserting

## Deployment

### How to run the Vald cluster?

We recommend using helm for running the Vald cluster.
You can run by following steps.

1. Install helm(v3~) and prepare Kubernetes cluster
1. Configure helm charts as your demand
1. Deploy by helm command

## API

### Is there any support for bulk insert?

Vald provides `MultiInsert` and `StreamInsert` for bulk insert.
Please refer to [the API documentation](../api/insert.md) for more detail.

## Data

### Can Vald handle multi-embedding vectors?

Unfortunately, the current Vald cannot directly handle multi embedding spaces with a single Vald cluster.
For handling the multi-embedding vectors in Vald, you have to do from 2 options.

1. Deploy multiple Vald cluster
1. Covert vector to new vector in the specific space by some methods

### How to backup index data?

There are three for backup:

1. Using external storage (S3, GCS)
1. Using Persistent Volume
1. Using the external storage and Persistent Volume

Please refer to the sample configuration.

---

## Related Document

- [Troubleshooting](../user-guides/troubleshooting.md)
