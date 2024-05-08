# Vald Index Manager

Vald Index Manager is the component that controls the indexing process of Vald Agent Pods in the same Vald cluster.

## Responsibility

Vald Index Manager has a unique and simple role in controlling the indexing timing for all Vald Agent pods in the Vald cluster.

It requires `Vald Discoverer` to fulfill its responsibility.

## Features

This chapter shows the main features to fulfill Vald Index Manager’s role.

### Syncing data from Vald Discoverer

Vald Index Manager gets the IP addresses of each Vald Agent pod from Vald Discoverer when its container starts.

In addition, when IP address changes, Vald Index Manager gets the new IP address from Vald Discoverer.

Vald Index Manager uses these for the controlling indexing process.

### Controlling Indexing process

When Vald Agent Pod creates or saves indexes on its container memory, it blocks all searching requests from the user and returns an error instead of a search result.

Stop-the-world happens when all Vald Agent pods run the function involved in the indexing operation, e.g., `createIndex`, simultaneously.

Vald Index Manager manages the indexing process of all Vald Agent pods to prevent this event.

Vald Index Manager uses a Vald Agent pods' IP address list from Vald Discoverer and index information, including the stored index count, uncommitted index count, creating index count, and saving index count, from each Vald Agent pod.

The control process is Vald Index Manager sends `createIndex` requests for concurrency simultaneously, sends a new request when the job is finished and continues until it sends to all agents.

At the end of each process, Vald Index Manager updates the index information from each Vald Agent pod.

Vald Index Manager runs this process periodically by set time intervals.

<div class="notice">
Concurrency means the number of Vald Agent pods for simultaneously sending requests for the indexing operation.<BR>

When the Vald Agent pod has no uncommitted index or is running the indexing function already, it does not send the request.

</div>
