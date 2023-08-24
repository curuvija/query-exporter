# query-exporter

A Helm chart for query-exporter (export Prometheus metrics from SQL queries)

## Additional Information

[query-exporter](https://github.com/albertodonato/query-exporter) exposes Prometheus metrics based on SQL queries. It supports
different databases. You can find more details about it here https://github.com/albertodonato/query-exporter.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Milos Curuvija | <curuvija@live.com> |  |

## Installing the Chart

You have to create a secret first. Secret contains configuration for query-exporter. Here is an example configuration for `sqlite` database. You can find more examples [here](https://github.com/albertodonato/query-exporter/tree/main/examples).

```yaml
databases:
  db1:
    dsn: sqlite://
    connect-sql:
      - PRAGMA application_id = 123
      - PRAGMA auto_vacuum = 1
    labels:
      region: us1
      app: app1

metrics:
  metric1:
    type: gauge
    description: A sample gauge

queries:
  query1:
    interval: 5
    databases: [db1]
    metrics: [metric1]
    sql: SELECT random() / 1000000000000000 AS metric1
```
Put this configuration into a file named `config.yaml` (this name is important since the key under `data` in created secret will contain exactly `config.yaml` which is later needed as a mount point in deployment `volumeMounts`).

Now create the secret in the namespace you want to deploy `query-exporter` Helm chart. For default namespace just run:

```bash
kubectl create secret generic --from-file=./config.yaml query-exporter-config-secret
```

If you want to change secret name from `query-exporter-config-secret` to something else change `configSecretName` in the values file.

To install the chart run:

```console
$ helm repo add curuvija https://curuvija.github.io/helm-charts/
$ helm repo update
$ helm install query-exporter curuvija/query-exporter
```

## Configure Prometheus scraping

If you use Prometheus operator ServiceMonitor will be created by default to configure your instance to scrape it.

If you don't use Prometheus operator then you can use this configuration to configure scraping (and disable ServiceMonitor creation):

```yaml
    additionalScrapeConfigs:
    - job_name: query-exporter-scrape
      honor_timestamps: true
      scrape_interval: 15s
      scrape_timeout: 10s
      metrics_path: /metrics
      scheme: http
      follow_redirects: true
      relabel_configs:
      - source_labels: [__meta_kubernetes_service_label_app_kubernetes_io_instance, __meta_kubernetes_service_labelpresent_app_kubernetes_io_instance]
        separator: ;
        regex: (query-exporter);true
        replacement: $1
        action: keep
      kubernetes_sd_configs:
      - role: endpoints
```

## Metrics

Pod listens by defayult on port `9560`. You can use port-forward to inspect its output at `http://localhost:9560/metrics`. Here is an example output you can expect with `sqldb` configuration example:

```text
# HELP database_errors_total Number of database errors
# TYPE database_errors_total counter
# HELP queries_total Number of database queries
# TYPE queries_total counter
queries_total{app="app1",database="db1",query="query1",region="us1",status="success"} 10.0
# HELP queries_created Number of database queries
# TYPE queries_created gauge
queries_created{app="app1",database="db1",query="query1",region="us1",status="success"} 1.6928516313686748e+09
# HELP query_latency Query execution latency
# TYPE query_latency histogram
query_latency_bucket{app="app1",database="db1",le="0.005",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.01",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.025",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.05",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.075",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.1",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.25",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.5",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="0.75",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="1.0",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="2.5",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="5.0",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="7.5",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="10.0",query="query1",region="us1"} 10.0
query_latency_bucket{app="app1",database="db1",le="+Inf",query="query1",region="us1"} 10.0
query_latency_count{app="app1",database="db1",query="query1",region="us1"} 10.0
query_latency_sum{app="app1",database="db1",query="query1",region="us1"} 0.00021720037329941988
# HELP query_latency_created Query execution latency
# TYPE query_latency_created gauge
query_latency_created{app="app1",database="db1",query="query1",region="us1"} 1.6928516313685274e+09
# HELP metric1 A sample gauge
# TYPE metric1 gauge
metric1{app="app1",database="db1",region="us1"} 8800.0
```

## Upgrade from version 1.x.x

* move config into secrets as explained in [Installing the Chart](#installing-the-chart) section
* run `helm upgrade` command



## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | configure affinity |
| autoscaling.enabled | bool | `false` | enable or disable autoscaling |
| autoscaling.maxReplicas | int | `100` | maximum number of replicas |
| autoscaling.minReplicas | int | `1` | minimum number of replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | configure at what percentage to trigger autoscalling |
| configSecretName | string | `"query-exporter-config-secret"` | Loads configuration from existing secret |
| fullnameOverride | string | `""` | overrides name without having chartName in front of it |
| image | object | `{"pullPolicy":"IfNotPresent","repository":"adonato/query-exporter","tag":"2.9.0"}` | Image to use for deployment |
| image.pullPolicy | string | `"IfNotPresent"` | define pull policy |
| image.repository | string | `"adonato/query-exporter"` | repository to pull image |
| image.tag | string | `"2.9.0"` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | Image pull secrets if you want to host the image |
| ingress | object | `{"annotations":{},"className":"","enabled":false,"hosts":[{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}],"tls":[]}` | ingress configuration |
| ingress.annotations | object | `{}` | ingress annotations |
| ingress.className | string | `""` | ingress class name |
| ingress.enabled | bool | `false` | enable or disable ingress configuration creation |
| ingress.hosts | list | `[{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}]` | hosts |
| ingress.hosts[0] | object | `{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}` | hostname |
| ingress.hosts[0].paths | list | `[{"path":"/","pathType":"ImplementationSpecific"}]` | paths |
| ingress.hosts[0].paths[0] | object | `{"path":"/","pathType":"ImplementationSpecific"}` | path |
| ingress.hosts[0].paths[0].pathType | string | `"ImplementationSpecific"` | path type |
| ingress.tls | list | `[]` | tls configuration |
| livenessProbe | object | `{"httpGet":{"path":"/","port":9560}}` | configure liveness probe |
| nameOverride | string | `""` | overrides name (partial name override - chartName + nameOverride) |
| nodeSelector | object | `{}` | define node selector to schedule your pod(s) |
| podAnnotations | object | `{}` | pod annotations |
| podSecurityContext | object | `{}` | define pod security context https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| prometheus | object | `{"monitor":{"additionalLabels":{},"enabled":true,"interval":"15s","namespace":[],"path":"/metrics"}}` | configure Prometheus Service monitor to expose metrics |
| prometheus.monitor.additionalLabels | object | `{}` | add additonal labels to service monitoring |
| prometheus.monitor.enabled | bool | `true` | enable or disable creation of service monitor |
| prometheus.monitor.interval | string | `"15s"` | Prometheus scraping interval |
| prometheus.monitor.namespace | list | `[]` | provide namespace where to create this service monitor |
| prometheus.monitor.path | string | `"/metrics"` | path where you want to expose metrics |
| readinessProbe | object | `{"httpGet":{"path":"/","port":9560}}` | configure readiness probe |
| replicaCount | int | `1` | replicaCount - number of pods to run |
| resources | object | `{"limits":{"cpu":"100m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"128Mi"}}` | specify resources |
| resources.limits | object | `{"cpu":"100m","memory":"128Mi"}` | specify resource limits |
| resources.limits.cpu | string | `"100m"` | specify resource limits for cpu |
| resources.limits.memory | string | `"128Mi"` | specify resource limits for memory |
| resources.requests.cpu | string | `"100m"` | specify resource requests for cpu |
| resources.requests.memory | string | `"128Mi"` | specify resource requests for memory |
| securityContext | object | `{"readOnlyRootFilesystem":true,"runAsNonRoot":true,"runAsUser":1000}` | define security context https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-capabilities-for-a-container |
| securityContext.readOnlyRootFilesystem | bool | `true` | Mounts the container's root filesystem as read-only. |
| securityContext.runAsNonRoot | bool | `true` | run docker container as non root user. |
| securityContext.runAsUser | int | `1000` | specify under which user all processes inside container will run. |
| service | object | `{"port":9560,"type":"ClusterIP"}` | service configuration |
| service.port | int | `9560` | service port |
| service.type | string | `"ClusterIP"` | service type |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | provide tolerations |


----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
