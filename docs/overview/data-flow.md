# Data Flow

On this page, we describe the data flow inside Vald and how vector indexes are stored in the Vald Cluster.
It will help you to understand what Vald does from users' requests.

The below image is the basic Vald architecture.

<img src="../../assets/docs/overview/vald_basic_architecture.svg" />

We will explain using this image in the following sections.

- [Data Flow](#data-flow)
  - [Insert](#insert)
  - [Search](#search)
  - [Update](#update)
  - [Upsert](#upsert)
  - [Delete](#delete)

## Insert

Users can insert a vector to Vald cluster using Insert API. When the insert command is processed, Vald will intelligently insert the index into the most suitable Vald Agent(s) based on the memory usage of the Vald Agent and the node.

To obtain the memory usage of Vald Agent and the node, Vald Discoverer is required to rank which is the most suitable agent to perform the insert request by talking with the kube-apiserver.

To make the insert command effective, the `CreateIndex` instruction is required by the Vald Index Manager or Vald Agent itself to update the index to provide extremely high searching performance to users.

Please note that only one embedding space is supported in a single Vald Cluster, if you want to support multiple embedded spaces, you may need to consider to deploy multiple Vald cluster to support this use case.

<img src="../../assets/docs/overview/insert_flow.svg" />

When the user inserts data into Vald:

1. Vald Ingress receives the request from the user. The request includes the vector and the vector ID. The vector ID is the user-defined unique ID for each vector. Vald Ingress will forward the request to the Vald LB Gateway to process the request.
2. Vald LB Gateway will determine which Vald Agent(s) to process the request based on the resource usage of the nodes and pods, and the number of vector replicas.
3. Vald LB Gateway will generate the UUID, and forward the generated UUID and the vector data to the selected Vald Agents in parallel. Vald Agent will insert the vector and UUID in an on-memory vector queue.
4. A vector queue will be committed to an ANN graph index by a `CreateIndex` instruction executed by the Vald Index Manager.
5. Vald Agent will start save index data into file by `SaveIndex` instruction after Vald Agent successfully `CreateIndex`.

## Search

Users can perform a _k_ nearest neighbor vector searching in Vald. The searching performance is extremely high because of the on-memory graph index structure and optimized Vald structure.

The search request will be broadcast to all Vald Agents to search the top _k_ nearest neighbor vectors from each Vald Agents, and the result will be combined and sorted by the distance of the target vector by the Vald LB Gateway.

<img src="../../assets/docs/overview/search_flow.svg" />

When the user searches a vector from Vald:

1. Vald Ingress receives a search request from the user. Vald provides 2 searching interfaces to the user, the user can search by vector or search by the vector ID. Vald Ingress will forward the request to the Vald LB Gateway to process the request.
2. Vald LB Gateway will forward the request to all Vald Agents in parallel. Each Vald Agent will search the _k_ nearest neighbor vectors in an on memory graph index.
3. Vald Agent returns the searching result to the Vald LB Gateway. The searching result includes the UUID, the vector distance, and the vector. The number of the each result will be the same as requested.
4. Vald LB Gateway will aggregate all searching results from all Vald Agents, and rank the result by the vector distance.
5. Vald LB Gateway will return the aggregated searching result to the Vald Ingress. Vald Ingress will return the aggregated searching result to the user.

## Update

Users can update a vector by sending an update request to Vald. Vald will perform delete and insert requests to perform a single update command.

<img src="../../assets/docs/overview/update_flow.svg" />

When the user updates a vector from Vald:

1. Vald Ingress receives the request from the user. The request includes the existing vector ID(s) and the new vector(s) to be updated. Vald Ingress will forward the request to the Vald LB Gateway to process the request.
2. Vald LB Gateway will generate the UUID(s), and broadcast the delete request with the generated UUID(s) to the Vald Agents. Each Vald Agent will delete the vector data and the metadata if the corresponding UUID(s) is found in the in-memory graph index.
3. Each Vald Agent will return success to the Vald LB Gateway if completed to delete the requested data successfully.
4. The insertion step will start after the deletion steps. It is the same as insert flow, please refer to [Insert](#insert) section.
5. If Vald Agent successfully inserts the request data, it will return success (e.g. IP address of pod) to the Vald LB Gateway.
6. Vald LB Gateway will return success to the Vald Ingress, then Vald Ingress will return success to the user.

## Upsert

Upsert request updates the existing vector if the same vector ID exists, or inserts the vector into Vald.

<img src="../../assets/docs/overview/upsert_flow.svg" />

When the user upserts a vector to Vald:

1. Vald Ingress receives the request from the user. The request includes the vector ID(s) and the vector(s). Vald Ingress will forward the request to the Vald LB Gateway to process the request.
2. Vald LB Gateway will broadcast an existing check request to the Vald Agent(s) to check if the vector exists.
3. Vald Agent returns the existing check result to Vald LB Gateway.
4. If the vector with the same vector ID exists, Vald LB Gateway will send an update request to Vald Agent(s) same as the [update flow](#update) step 3 to step 5. If the vector does not exist, Vald LB Gateway will process the [insert flow](#insert) from step 3 to step 4. Vald Index Manager will send `CreateIndex` request to each Vald Agent at regular intervals.
5. Vald Agent(s) return success to Vald LB Gateway.
6. Vald Filter Gateway will return success to the Vald Ingress. Vald Ingress will return success to the user.

## Delete

Delete request will delete the vector in Vald cluster. Vald will broadcast the delete request to all Vald Agents to delete the vector inside the cluster. To make the delete command effective, the `CreateIndex` command is required by the Vald Index Manager or Vald Agent itself to update the vector index.

<img src="../../assets/docs/overview/delete_flow.svg" />

When the user deletes a vector that is indexed in Vald Agent:

1. Vald Ingress receives the delete request from the user. The request includes the vector ID(s), which is specified by the user. Vald Ingress will forward the request to the Vald LB Gateway.
2. Vald LB Gateway will broadcast the delete request with UUID(s) to the Vald Agents. Each Vald Agent will delete the vector data and the metadata if the corresponding UUID(s) is found in the in-memory graph index.
3. If Vald Agent successfully deletes the request data, it will return success to the Vald LB Gateway.
4. Vald LB Gateway will return success to the Vald Ingress. Vald Ingress will return success to the user.
