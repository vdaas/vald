# Operation troubleshooting

This page introduces the popular troubleshooting for operating a Vald cluster.

The [Flow chart](../troubleshooting/chart.md) helps you to find the root reason for your problem.

Additionally, if you encounter some errors when using API, [API status code](../api/status.md) helps you, too.

## Insert Operation

### Vald Agent NGT crashed at the insert process.

Let’s check your container limit of memory at first.
Vald Agent requires memory for keeping indexing on memory.

```bash
kubectl describe statefulset vald-agent-ngt
```

If the limit memory is set, please remove it or update value larger.

## Search Operation

### Vald returns no search result.

There are two possible reasons.

1. Indexing has not finished in Vald Agent

    Vald will search the nearest vectors of query from the indexing in Vald Agent.
    If the indexing process is running, Vald Agent returns no search result.
    
    It will resolve when completed indexing instruction, like as `CreateIndex`.

1. Too short timeout for searching

    When the search timeout configuration is too short, Vald LB Gateway stops the searching process before getting search result from Vald Agent.

    In the sense of search operation, you can modify search timeout by [payload config](../api/search.md).

<div class="notice">
It is easy to find out which problem occurs by inspections the log of each Pod, like <a href="https://github.com/wercker/stern">stern</a>.
</div>

## Others

### Vald Agent NGT crashed when initContainer.

Vald Agent NGT requires an AVX2 processor for running.
Please check your CPU information.

---

## Related Document

- [Flow Chart](../troubleshooting/chart.md)
- [API Status](../api/status.md)
- [FAQ](../support/faq.md)
