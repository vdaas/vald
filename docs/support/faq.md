# FAQ

## Does Vald handle multi-embedding vectors?

Unfortunately, current Vald cannot handle multi embedding spaces with a single Vald cluster directly.
For handling the multi-embedding vectors in Vald, you have to do from 2 options.

1. Deploy multiple Vald cluster
1. Covert vectors in the same space, e.g., padding category label to the vectors.

## Vald Agent NGT crashed at the insert process.

Let’s check your container limit of memory at first.
Vald Agent requires memory for keeping indexing on memory.

## Vald Agent NGT crashed when initContainer.

Vald Agent NGT requires AVX2 processor for running.
Please check your CPU information.

## Is there any support for bulk insert?

Vald provides `MultiInsert` and `StreamInsert` for bulk insert.
Please refer to [the API documentation](../api/insert.md) for more detail.

## What are the pluggable options?

Vald has 3 pluggable options:

1. Backup with the external storage for Vald Agent
   - You can connect the external storage like S3, or GCS, or etc. to Vald Agent Sidecar component for backup.
1. Algorithm of the core engine for Vald Agent
   - We're going to add another algorithm in near future.
1. Filtering with filter gateway
   - you can filter the search results by own defined filter component by connecting to the filter gateway before returning the search result
   - you can convert object data to vector by own defined ingress filter component by connection to filter gateway before inserting

## Vald returns no search result.

It supposes there are 2 reasons.

1. Indexing has not finished in Vald Agent
   - Vald will search the nearest vectors of query from the indexing in Vald Agent.
     If indexing does not finish yet, Vald Agent cancels searching.
1. Too short timeout for searching
   - When the search timeout configuration is too short, Vald LB Gateway stops the searching process before getting search result from Vald Agent.

## What is the ingress filter?

Ingress filter is the user-defined component that connects to the `vald-filter-gateway` component.
It is used for converting object data (e.g. image binary, text, etc.) to the vector using your ML models as pre-process of each request.

## What is the egress filter?

Egress filter is the user-defined component that connects to the vald-filter-gateway component.
It is used for filtering search results, e.g. when you’d like to get a white T-shirts list, use it for removing other colors T-shirts from search results which are created by the vald-lb-gateway component.

## How to backup index data?

There are 3 ways for backup:

1. Using external storage (S3, GCS)
1. Using Persistent Volume
1. Using the external storage and Persistent Volume

Please refer to the sample configuration.

## How to run Vald cluster?

We recommend using helm for running Vald cluster.
You can run by following steps.

1. Install helm(v3~) and prepare Kubernetes cluster
1. Configure helm charts as your demand
1. Deploy by helm command

## Is Vald Index Manager recommended using?

We recommend to use it when you’d like to operate as a cluster.
It helps you to manage indexing timing for each Vald Agent.
