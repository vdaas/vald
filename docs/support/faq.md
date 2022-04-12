# FAQ

This pages answers the common questions.

Please refer to [trouble shooting page](../user-guides/trouble-shooting.md) when you encounter operation problem.

## Component

### What is the ingress filter?

Ingress filter is the user-defined component that connects to the `vald-filter-gateway` component.
It is used for converting object data (e.g. image binary, text, etc.) to the vector using your ML models as pre-process of each request.

### What is the egress filter?

Egress filter is the user-defined component that connects to the vald-filter-gateway component.
It is used for filtering search results, e.g. when you’d like to get a white T-shirts list, use it for removing other colors T-shirts from search results which are created by the vald-lb-gateway component.

### Is Vald Index Manager recommended using?

We recommend to use it when you’d like to operate as a cluster.
It helps you to manage indexing timing for each Vald Agent.

## Custom options

### What are the pluggable options?

Vald has 3 pluggable options:

1. Backup with the external storage for Vald Agent
   - You can connect the external storage like S3, or GCS, or etc. to Vald Agent Sidecar component for backup.
1. Algorithm of the core engine for Vald Agent
   - We're going to add another algorithm in near future.
1. Filtering with filter gateway
   - you can filter the search results by own defined filter component by connecting to the filter gateway before returning the search result
   - you can convert object data to vector by own defined ingress filter component by connection to filter gateway before inserting

## Deployment

### How to run Vald cluster?

We recommend using helm for running Vald cluster.
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

Unfortunately, current Vald cannot handle multi embedding spaces with a single Vald cluster directly.
For handling the multi-embedding vectors in Vald, you have to do from 2 options.

1. Deploy multiple Vald cluster
1. Covert vectors in the same space, e.g., padding category label to the vectors.

### How to backup index data?

There are 3 ways for backup:

1. Using external storage (S3, GCS)
1. Using Persistent Volume
1. Using the external storage and Persistent Volume

Please refer to the sample configuration.

---

## Related Document

- [Trouble Shooting](../user-guides/trouble-shooting.md)
