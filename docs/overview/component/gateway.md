# Vald Gateway

Vald Gateway handles any requests component in the Vald cluster, as its name suggests.

This component is essential for operating as Vald cluster because it depends on most components.

This page introduces the overview and features of Vald Gateway.

## Responsibility

Vald Gateway is responsible for:

- Pass requests to the other vald components
- Control request timeout
- Sort out search results as user demand

## Features

Vald Gateway has two kinds of components:

1. Vald LB Gateway
   - The main component for Vald cluster and connecting to Kubernetes Ingress.
2. Vald Filter Gateway
   - Bypass between Vald LB Gateway and user-defined Ingress filter or Egress filter components
   - Vald LB Gateway passes to the requests at the right time.

<!-- TODO: insert image of vald lb gateway and vald filter gateway -->

Like the above, we will focus on introducing the features of Vald LB Gateway.

- Pass and Control all requests
  - Vald LB Gateway is the only component that connects to Kubernetes Ingress in the Vald cluster’s components.
  - User requests have passed Kubernetes Ingress reaches to Vald LB Gateway.
  - Vald LB Gateway will suspend processing according to timeout and return a response when the internal processing takes a long time.

- Control inserting vector
  - As its name shows, Vald LB Gateway has the load balancing feature.
  - Vald LB Gateway controls insert vector requests based on `index replica` and each Vald Agent Pod resource usage, which Vald Discoverer provides, to avoid uneven resource usage.

- Broadcast and Gathering search result
  - Vald LB Gateway broadcasts searching requests, e.g., `Search`, `GetObject`, `Exist`, to all Vald Agent Pods and gets their result.
  - In the `Search` phase, Vald LB Gateway gets search results from Vald Agent Pods within the user-defined timeout time and sorts combining results by shortest distance. Then, Vald LB Gateway returns the Top-_N_ of the sorted search results as a result. (_N_ is the user-defined number)

- Work together with Ingress/Egress filter
  - Vald LB Gateway is the only component to connect to Vald Ingress filter or Egress filter via Vald Filter Gateway.
  - When Vald Ingress runs, Vald LB Gateway will pass the object of request to Ingress, then get the vector converted from it.
  - When Vald Egress runs, Vald LB Gateway will pass the search result to Egress and get the filtered result.

<!-- TODO: add the link of configuration page -->
