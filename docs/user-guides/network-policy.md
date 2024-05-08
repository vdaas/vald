# Network Policy

[Network Policy](https://kubernetes.io/docs/concepts/services-networking/network-policies/) is a Kubernetes feature that controls ingress and egress network traffic for pods. In Vald, you can set network policies as follows.

> Please note that [prerequisites](https://kubernetes.io/docs/concepts/services-networking/network-policies/#prerequisites) are required for using network policies. Even if you configure the following settings in a cluster that does not meet the prerequisites, network policies will not be effective.

# Network Policy in Vald

To enable network policies in a Vald cluster, set `defaults.networkPolicy.enabled` to `true` as follows:

```yaml
defaults:
  networkPolicy:
    enabled: true
```

This sets the following ingress/egress rules between Vald components (these are the minimum required rules for a Vald cluster to work).

| from / to      | agent | discoverer | filter gateway | lb gateway | index manager | kube-system |
| :------------- | :---: | :--------: | :------------: | :--------: | :-----------: | :---------: |
| agent          |  N/A  |     ⛔     |       ⛔       |     ⛔     |      ⛔       |     ✅      |
| discoverer     |  ⛔   |    N/A     |       ⛔       |     ⛔     |      ⛔       |     ✅      |
| filter gateway |  ⛔   |     ⛔     |      N/A       |     ✅     |      ⛔       |     ✅      |
| lb gateway     |  ✅   |     ✅     |       ⛔       |    N/A     |      ⛔       |     ✅      |
| index manager  |  ✅   |     ✅     |       ⛔       |     ⛔     |      N/A      |     ✅      |

# Add a user custom Network Policy

There may be cases where you want to connect a Vald cluster to external components. Specifically, for the following cases:

- Enable egress to `OpenTelemetryCollector` to use [observability features](https://vald.vdaas.org/docs/user-guides/observability-configuration/)
- Enable egress to an external filter component to use [filtering features](https://vald.vdaas.org/docs/user-guides/filtering-configuration/).

To handle such cases, Vald allows you to set user custom network policies using the `defaults.networkPolicy.custom` field as follows:

```yaml
defaults:
  networkPolicy:
    enabled: true
    custom:
      ingress:
        - from:
            - podSelector:
                matchLabels:
                  app.kubernetes.io/name: pyroscope
      egress:
        - to:
            - podSelector:
                matchLabels:
                  app.kubernetes.io/name: opentelemetry-collector-collector
```

Please write down the same notation as the `ingress/egress` field of [NetworkPolicy resource](https://kubernetes.io/docs/concepts/services-networking/network-policies/#networkpolicy-resource) in our `custom` field.

> Currently, these custom network policies are applied to all Vald components.
