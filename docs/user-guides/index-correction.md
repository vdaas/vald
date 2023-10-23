# Index Correction

In the Vald cluster, the same Index is replicated to multiple agents due to the `index_replica` setting. However, inconsistencies between replicas may occur due to pod eviction or the occurrence of OOM killer during vector insertions. For example,

1. The timestamp of the index differs between agents (some agents have an old index saved and it has not been updated).
2. The number of replicas does not meet the value set in `index_replica`.

To resolve these inconsistencies, you can use the `Index Correction` feature.

`Index Correction` is implemented as a [`CronJob`](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/), checking the consistency between replicas regularly and resolving any inconsistencies.

## Settings

- enabled  
Turns the index correction feature on/off.
- schedule  
Sets the interval for the job start in cron notation (the default value is `3 6 * * *`, which means 3:06 AM every day).
- suspend  
[Temporary suspension setting](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#schedule-suspension) for CronJob.

```yaml
manager:
  index:
    corrector:
      enabled: true
      schedule: "3 6 * * *"
      suspend: false
```

## Important Notes

- Processing time  
Under conditions of 10 million vectors and agent replica *10, it takes about 10~20 minutes. The process is O(MN) where M is the number of vector items and N is the number of agent replicas.
- concurrencyPolicy  
`Forbid` is set internally, so a new job will not be created while an existing job is running. In other words, if the process does not finish within the interval specified by the schedule, the next job will not be scheduled.
- Index operations during correction  
Vector operations performed after the start of the index correction job are not considered in that job.
