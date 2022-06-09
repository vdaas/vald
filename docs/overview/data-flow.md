# Data Flow

This page describes the data flow inside Vald and how to store vector indexes in the Vald Cluster.
It will help you to understand what Vald does from users' requests.

The below image is the basic Vald architecture.

<img src="../../assets/docs/overview/vald_basic_architecture.svg" />

We will explain using this image in the following sections.

- [Data Flow](#data-flow)
  - [Insert](#insert)
  - [Search](#search)
  - [Update](#update)
  - [Upsert](#upsert)
  - [Remove](#remove)

## Insert

Users can insert a vector to the Vald cluster using Insert API. When the insert command is processed, Vald will intelligently insert the index into the most suitable Vald Agent(s) based on the memory usage of the Vald Agent and the node.

To obtain the memory usage of Vald Agent and the node, Vald Discoverer must rank which is the most suitable agent to perform the insert request by talking with the kube-apiserver.

It requires the `CreateIndex` instruction to Vald Agent by Vald Index Manager or self-updating the index by Vald Agent to provide high searching performance to users because the insert command is not responsible for indexing.

<div class="warn">
A single Vald cluster supports only one embedding space.<BR>
If you want to support multiple embedded spaces, you may need to consider deploying multiple Vald clusters to support this use case.
</div>

<img src="../../assets/docs/overview/insert_flow.svg" />

When the user inserts data into Vald:

1. Vald LB Gateway receives the request from the user. The request includes the vector and the vector ID. The vector ID is the user-defined unique ID for each vector.
2. Vald LB Gateway will determine which Vald Agent(s) to process the request based on the resource usage of the nodes and pods and the number of vector replicas.
3. Vald LB Gateway will generate the UUID for each inserting vector and forward the generated UUID and the vector data to the selected Vald Agents in parallel. Vald Agent will insert the vector and UUID in an on-memory vector queue.
4. A vector queue will be committed to an ANN graph index by a `CreateIndex` instruction executed by the Vald Index Manager.
5. Vald Agent will start to save index data into the file by `SaveIndex instruction` after Vald Agent successfully `CreateIndex`.

## Search

Users can perform a _k_ nearest neighbor vector searching in Vald. The searching performance is extremely high because of the on-memory graph index structure and optimized Vald structure.

Vald LB Gateway will broadcast the search request to all Vald Agent Pods to search the top _k_ nearest neighbor vectors from each Vald Agent, and the result will be combined and sorted by the distance of the target vector by the Vald LB Gateway.

<img src="../../assets/docs/overview/search_flow.svg" />

When the user searches a vector from Vald:

1. Vald LB Gateway receives a search request from the user. Vald provides two searching interfaces, search by vector and search by the vector ID, to the user.
2. Vald LB Gateway will forward the request to all Vald Agents in parallel. Each Vald Agent will search the _k_ nearest neighbor vectors on a memory graph index.
3. Vald Agent returns the search result to the Vald LB Gateway. The search result includes the UUID, the vector distance, and the vector. The number of each outcome will be the same as requested.
4. Vald LB Gateway will aggregate all searching results from all Vald Agents, then rank the result by the vector distance.
5. Vald LB Gateway will return the aggregated searching result to the user.

## Update

Users can update a vector by sending an update request to Vald.
Vald will execute deleting and inserting instructions to perform a single update command.

<img src="../../assets/docs/overview/update_flow.svg" />

When the user updates a vector from Vald:

1. Vald LB Gateway receives the request from the user. The request includes the existing vector ID(s) and the new vector(s) to be updated.
2. Vald LB Gateway will generate the UUID(s) for each vector and broadcast the remove request with the generated UUID(s) to the Vald Agents. Each Vald Agent will remove the vector data and the metadata if the corresponding UUID(s) exist on the in-memory graph index.
3. Each Vald Agent will return the result to the Vald LB Gateway if completed to remove the requested data successfully.
4. The insertion step will start after the deletion steps. It is the same as insert flow; please refer to the [Insert](#insert) section.
5. If Vald Agent successfully inserts the request data, it will return success (e.g., IP address of pod) to the Vald LB Gateway.
6. Vald LB Gateway will return the result to the user.

## Upsert

Upsert request updates the existing vector if the same vector ID already exists or inserts the vector into Vald.

<img src="../../assets/docs/overview/upsert_flow.svg" />

When the user upserts a vector to Vald:

1. Vald LB Gateway receives the request from the user, including the user's vector ID(s) and the vector(s).
2. Vald LB Gateway will broadcast an existing check request to the Vald Agent(s) to check if the vector exists.
3. Vald Agent returns the existing check result to Vald LB Gateway.
4. If the vector with the same vector ID exists, Vald LB Gateway will send an update request to Vald Agent(s), same as the [update flow](#update) step 2 to step 4. If the vector does not exist, Vald LB Gateway will process the [insert flow](#insert) from step 2 to step 3. Vald Index Manager will send `CreateIndex` requests to each Vald Agent regularly.
5. Vald Agent(s) return the result to Vald LB Gateway.
6. Vald Filter Gateway will return the result to the user.

## Remove

Remove request will remove the vector in the Vald cluster.
Vald will broadcast the remove request to all Vald Agents to remove the vector inside the cluster.
It requires the `CreateIndex` instruction to Vald Agent by the Vald Index Manager or self-updating the index by Vald Agent because the remove command is not responsible for deleting indexes.

<img src="../../assets/docs/overview/remove_flow.svg" />

When the user removes a vector that exists in the in-memory graph index in Vald Agent:

1. Vald LB Gateway receives the remove request from the user. The request includes the vector ID(s), which the user specifies.
2. Vald LB Gateway will broadcast the remove request with UUID(s) to the Vald Agents. Each Vald Agent will remove the vector data and the metadata if the corresponding UUID(s) exists in the in-memory graph index.
3. If Vald Agent successfully removes the request data, it will return success to the Vald LB Gateway.
4. Vald LB Gateway will return the result to the user.
