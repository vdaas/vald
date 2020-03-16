Vald
===

This is a Helm chart to install Vald components.


Install
---

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Run the following command to install the chart,

    $ helm install --generate-name vald/vald


Configuration
---

| Parameter | Description | Default |
|-----------|-------------|---------|
| `<component>.time_zone` | Time zone | `UTC` |
| `<component>.logging.logger` | logger name | `glg` |
| `<component>.logging.level` | logging level | `debug` |
| `<component>.logging.format` | logging format | `raw` |
| `<component>.image.repository` | image repository | `vdaas/vald-<component>` |
| `<component>.image.tag` | image tag | version or nightly |
| `<component>.image.pullPolicy` | image pull policy | `Always` |

TBW
