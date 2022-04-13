# Troubleshooting

This page shows the troubleshooting for operating a Vald cluster.

## Insert Operation

### Vald Agent NGT crashed at the insert process.

Let’s check your container limit of memory at first.
Vald Agent requires memory for keeping indexing on memory.

## Search Operation

### Vald returns no search result.

It supposes there are two reasons.

1. Indexing has not finished in Vald Agent
   - Vald will search the nearest vectors of query from the indexing in Vald Agent.
     If indexing does not complete yet, Vald Agent cancels searching.
1. Too short timeout for searching
   - When the search timeout configuration is too short, Vald LB Gateway stops the searching process before getting search result from Vald Agent.

## Others

### Vald Agent NGT crashed when initContainer.

Vald Agent NGT requires an AVX2 processor for running.
Please check your CPU information.

---

## Related Document

- [FAQ](../support/faq.md)
