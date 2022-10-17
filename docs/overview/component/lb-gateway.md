# Vald LB Gateway

Vald LB Gateway is the component to handle requests in the Vald cluster.

This component is essential to operate the Vald cluster because all requests should pass it, and other components depend on it.

This page introduces the overview and features of Vald Gateway.

## Responsibility

Vald LB Gateway is responsible for:

- Pass and load balancing requests to other Vald components
- Control processing with a timeout setting
- Aggregate search results from each Vald Agent
- Sort and combine into one search result

The next chapter shows the main features.

## Features

Vald LB Gateway is a mandatory component for the Vald cluster.

Each following section introduces each feature.

### Pass and Control all requests

Vald LB Gateway is the only component that connects to Kubernetes Ingress in the Vald cluster’s components in basic.

User requests have passed Kubernetes Ingress reaches to Vald LB Gateway.

Then, Vald LB Gateway will suspend processing according to timeout and return a response when the internal processing takes a long time.

<div class="note">

If Vald Filter Gateway runs and the request has the filter gateway option, Kubernetes Ingress will pass the request to Vald Filter Gateway.

Then, Vald Filter Gateway sends the processed request to the Vald LB Gateway.

</div>

### Control insert vectors

As its name shows, Vald LB Gateway has the load balancing feature.

Vald LB Gateway controls insert vector requests based on `index replica` and each Vald Agent Pod resource usage, which [Vald Discoverer](./discoverer.md) provides, to avoid uneven resource usage.

### Broadcast search request and aggregate search result

Vald LB Gateway broadcasts searching requests, e.g., `Search`, `GetObject`, `Exist`, to all Vald Agent Pods and gets their result.

In the `Search` phase, Vald LB Gateway receives search results from all Vald Agent Pods within the user-defined timeout time and sorts combining results by shortest distance.
Then, Vald LB Gateway returns the Top-_N_ of the sorted search results (_N_ is the user-defined number).

<div class="note">

When Vald Agent Pod could not give search result by the timeout limit, Vald LB Gateway will cancel the searching requests to the Vald Agent Pod.

</div>

### Work together with Vald Filter Gateway

Vald LB Gateway is the only component to connect to the ingress (or egress) filter component via Vald Filter Gateway.

When the ingress component runs, Vald LB Gateway will pass the object of request to the Ingress, then get the vector converted from it.
Vald LB Gateway will send the search result to the Egress when the egress component runs and get the filtered result.

<div class="note">

For more information about Vald Filter Gateway, please refer to [Vald Filter Gateway Overview](./filter-gateway.md).

</div>

<!-- TODO: add the link of configuration page -->
