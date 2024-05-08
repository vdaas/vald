{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "vald.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "vald.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "vald.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "vald.labels" -}}
app.kubernetes.io/name: {{ include "vald.name" . }}
helm.sh/chart: {{ include "vald.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
joinListWithSpace
*/}}
{{- define "vald.utils.joinListWithSpace" -}}
{{- $local := dict "first" true -}}
{{- range $k, $v := . -}}{{- if not $local.first -}}{{- " " -}}{{- end -}}{{- $v -}}{{- $_ := set $local "first" false -}}{{- end -}}
{{- end -}}

{{/*
joinListWithComma
*/}}
{{- define "vald.utils.joinListWithComma" -}}
{{- $local := dict "first" true -}}
{{- range $k, $v := . -}}{{- if not $local.first -}},{{- end -}}{{- $v -}}{{- $_ := set $local "first" false -}}{{- end -}}
{{- end -}}

{{/*
logging settings
*/}}
{{- define "vald.logging"}}
{{- if .Values -}}
logger: {{ .Values.logger | quote }}
level: {{ .Values.level | quote }}
format: {{ .Values.format | quote }}
{{- end }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "vald.selectorLabels" -}}
app.kubernetes.io/name: {{ include "vald.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "vald.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "vald.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end -}}

{{/*
Server configures that inserted into server_config
*/}}
{{- define "vald.servers" -}}
servers:
  {{- $restEnabled := false }}
  {{- if hasKey .Values.servers.rest "enabled" }}
  {{- $restEnabled = .Values.servers.rest.enabled }}
  {{- end }}
  {{- if $restEnabled }}
  - name: rest
    host: {{ .Values.servers.rest.host }}
    port: {{ .Values.servers.rest.port }}
    {{- if .Values.servers.rest.server }}
    mode: {{ .Values.servers.rest.server.mode }}
    probe_wait_time: {{ .Values.servers.rest.server.probe_wait_time }}
    network: {{ .Values.servers.rest.server.network | quote }}
    socket_path: {{ .Values.servers.rest.server.socket_path | quote }}
    http:
      {{- if .Values.servers.rest.server.http }}
      shutdown_duration: {{ .Values.servers.rest.server.http.shutdown_duration }}
      handler_timeout: {{.Values.servers.rest.server.http.handler_timeout }}
      idle_timeout: {{ .Values.servers.rest.server.http.idle_timeout }}
      read_header_timeout: {{ .Values.servers.rest.server.http.read_header_timeout }}
      read_timeout: {{ .Values.servers.rest.server.http.read_timeout }}
      write_timeout: {{ .Values.servers.rest.server.http.write_timeout }}
      {{- end }}
    {{- end }}
  {{- end }}
  {{- $grpcEnabled := false }}
  {{- if hasKey .Values.servers.grpc "enabled" }}
  {{- $grpcEnabled = .Values.servers.grpc.enabled }}
  {{- end }}
  {{- if $grpcEnabled }}
  - name: grpc
    host: {{ .Values.servers.grpc.host }}
    port: {{ .Values.servers.grpc.port }}
    {{- if .Values.servers.grpc.server }}
    mode: {{ .Values.servers.grpc.server.mode }}
    probe_wait_time: {{ .Values.servers.grpc.server.probe_wait_time | quote }}
    network: {{ .Values.servers.grpc.server.network | quote }}
    socket_path: {{ .Values.servers.grpc.server.socket_path | quote }}
    grpc:
      {{- if .Values.servers.grpc.server.grpc }}
      bidirectional_stream_concurrency: {{ .Values.servers.grpc.server.grpc.bidirectional_stream_concurrency }}
      connection_timeout: {{ .Values.servers.grpc.server.grpc.connection_timeout | quote }}
      max_receive_message_size: {{ .Values.servers.grpc.server.grpc.max_receive_message_size }}
      max_send_message_size: {{ .Values.servers.grpc.server.grpc.max_send_message_size }}
      initial_window_size: {{ .Values.servers.grpc.server.grpc.initial_window_size }}
      initial_conn_window_size: {{ .Values.servers.grpc.server.grpc.initial_conn_window_size }}
      keepalive:
        {{- if .Values.servers.grpc.server.grpc.keepalive }}
        max_conn_idle: {{ .Values.servers.grpc.server.grpc.keepalive.max_conn_idle | quote }}
        max_conn_age: {{ .Values.servers.grpc.server.grpc.keepalive.max_conn_age | quote }}
        max_conn_age_grace: {{ .Values.servers.grpc.server.grpc.keepalive.max_conn_age_grace | quote }}
        time: {{ .Values.servers.grpc.server.grpc.keepalive.time | quote }}
        timeout: {{ .Values.servers.grpc.server.grpc.keepalive.timeout | quote }}
        min_time: {{ .Values.servers.grpc.server.grpc.keepalive.min_time | quote }}
        permit_without_stream: {{ .Values.servers.grpc.server.grpc.keepalive.permit_without_stream }}
        {{- end }}
      write_buffer_size: {{ .Values.servers.grpc.server.grpc.write_buffer_size }}
      read_buffer_size: {{ .Values.servers.grpc.server.grpc.read_buffer_size }}
      max_header_list_size: {{ .Values.servers.grpc.server.grpc.max_header_list_size }}
      header_table_size: {{ .Values.servers.grpc.server.grpc.header_table_size }}
      {{- if .Values.servers.grpc.server.grpc.interceptors }}
      interceptors:
        {{- toYaml .Values.servers.grpc.server.grpc.interceptors | nindent 8 }}
      {{- else }}
      interceptors: []
      {{- end }}
      enable_reflection: {{ .Values.servers.grpc.server.grpc.enable_reflection }}
      {{- end }}
    restart: {{ .Values.servers.grpc.server.restart }}
    {{- end }}
  {{- end }}
health_check_servers:
  {{- $livenessEnabled := false }}
  {{- if hasKey .Values.healths.liveness "enabled" }}
  {{- $livenessEnabled = .Values.healths.liveness.enabled }}
  {{- end }}
  {{- if $livenessEnabled }}
  - name: liveness
    host: {{ .Values.healths.liveness.host }}
    port: {{ .Values.healths.liveness.port }}
    {{- if .Values.healths.liveness.server }}
    mode: {{ .Values.healths.liveness.server.mode | quote }}
    probe_wait_time: {{ .Values.healths.liveness.server.probe_wait_time | quote }}
    network: {{ .Values.healths.liveness.server.network | quote }}
    socket_path: {{ .Values.healths.liveness.server.socket_path | quote }}
    http:
      {{- if .Values.healths.liveness.server.http }}
      shutdown_duration: {{ .Values.healths.liveness.server.http.shutdown_duration | quote }}
      handler_timeout: {{ .Values.healths.liveness.server.http.handler_timeout | quote }}
      idle_timeout: {{ .Values.healths.liveness.server.http.idle_timeout | quote }}
      read_header_timeout: {{ .Values.healths.liveness.server.http.read_header_timeout | quote }}
      read_timeout: {{ .Values.healths.liveness.server.http.read_timeout | quote }}
      write_timeout: {{ .Values.healths.liveness.server.http.write_timeout | quote }}
      {{- end }}
    {{- end }}
  {{- end }}
  {{- $readinessEnabled := false }}
  {{- if hasKey .Values.healths.readiness "enabled" }}
  {{- $readinessEnabled = .Values.healths.readiness.enabled }}
  {{- end }}
  {{- if $readinessEnabled }}
  - name: readiness
    host: {{ .Values.healths.readiness.host }}
    port: {{ .Values.healths.readiness.port }}
    {{- if .Values.healths.readiness.server }}
    mode: {{ .Values.healths.readiness.server.mode | quote }}
    probe_wait_time: {{ .Values.healths.readiness.server.probe_wait_time | quote }}
    network: {{ .Values.healths.readiness.server.network | quote }}
    socket_path: {{ .Values.healths.readiness.server.socket_path | quote }}
    http:
      {{- if .Values.healths.readiness.server.http }}
      shutdown_duration: {{ .Values.healths.readiness.server.http.shutdown_duration | quote }}
      handler_timeout: {{ .Values.healths.readiness.server.http.handler_timeout | quote }}
      idle_timeout: {{ .Values.healths.readiness.server.http.idle_timeout | quote }}
      read_header_timeout: {{ .Values.healths.readiness.server.http.read_header_timeout | quote }}
      read_timeout: {{ .Values.healths.readiness.server.http.read_timeout | quote }}
      write_timeout: {{ .Values.healths.readiness.server.http.write_timeout | quote }}
      {{- end }}
    {{- end }}
  {{- end }}
metrics_servers:
  {{- $pprofEnabled := false }}
  {{- if hasKey .Values.metrics.pprof "enabled" }}
  {{- $pprofEnabled = .Values.metrics.pprof.enabled }}
  {{- end }}
  {{- if $pprofEnabled }}
  - name: pprof
    host: {{ .Values.metrics.pprof.host }}
    port: {{ .Values.metrics.pprof.port }}
    {{- if .Values.metrics.pprof.server }}
    mode: {{ .Values.metrics.pprof.server.mode }}
    probe_wait_time: {{ .Values.metrics.pprof.server.probe_wait_time }}
    network: {{ .Values.metrics.pprof.server.network | quote }}
    socket_path: {{ .Values.metrics.pprof.server.socket_path | quote }}
    http:
      {{- if .Values.metrics.pprof.server.http }}
      shutdown_duration: {{ .Values.metrics.pprof.server.http.shutdown_duration }}
      handler_timeout: {{ .Values.metrics.pprof.server.http.handler_timeout }}
      idle_timeout: {{ .Values.metrics.pprof.server.http.idle_timeout }}
      read_header_timeout: {{ .Values.metrics.pprof.server.http.read_header_timeout }}
      read_timeout: {{ .Values.metrics.pprof.server.http.read_timeout }}
      write_timeout: {{ .Values.metrics.pprof.server.http.write_timeout }}
      {{- end }}
    {{- end }}
  {{- end }}
startup_strategy:
  {{- if $livenessEnabled }}
  - liveness
  {{- end }}
  {{- if $pprofEnabled }}
  - pprof
  {{- end }}
  {{- if $grpcEnabled }}
  - grpc
  {{- end }}
  {{- if $restEnabled }}
  - rest
  {{- end }}
  {{- if $readinessEnabled }}
  - readiness
  {{- end }}
full_shutdown_duration: {{ .Values.full_shutdown_duration }}
tls:
  {{- if .Values.tls }}
  enabled: {{ .Values.tls.enabled }}
  cert: {{ .Values.tls.cert | quote }}
  key: {{ .Values.tls.key | quote }}
  ca: {{ .Values.tls.ca | quote }}
  insecure_skip_verify: {{ .Values.tls.insecure_skip_verify }}
  {{- end }}
{{- end -}}

{{/*
observability
*/}}
{{- define "vald.observability" -}}
enabled: {{ .Values.enabled }}
otlp:
  {{- if .Values.otlp }}
  collector_endpoint: {{ .Values.otlp.collector_endpoint | quote }}
  trace_batch_timeout: {{ .Values.otlp.trace_batch_timeout | quote }}
  trace_export_timeout: {{ .Values.otlp.trace_export_timeout | quote }}
  trace_max_export_batch_size: {{ .Values.otlp.trace_max_export_batch_size }}
  trace_max_queue_size: {{ .Values.otlp.trace_max_queue_size }}
  metrics_export_interval: {{ .Values.otlp.metrics_export_interval | quote }}
  metrics_export_timeout: {{ .Values.otlp.metrics_export_timeout | quote }}
  attribute:
    {{- if .Values.otlp.attribute }}
    namespace: {{ .Values.otlp.attribute.namespace | quote }}
    pod_name: {{ .Values.otlp.attribute.pod_name | quote }}
    node_name: {{ .Values.otlp.attribute.node_name | quote }}
    service_name: {{ .Values.otlp.attribute.service_name | quote }}
    {{- end }}
  {{- end }}
metrics:
  {{- if .Values.metrics }}
  enable_version_info: {{ .Values.metrics.enable_version_info }}
  {{- if .Values.metrics.version_info_labels }}
  version_info_labels:
    {{- toYaml .Values.metrics.version_info_labels | nindent 4 }}
  {{- else }}
  version_info_labels: []
  {{- end }}
  enable_memory: {{ .Values.metrics.enable_memory }}
  enable_goroutine: {{ .Values.metrics.enable_goroutine }}
  enable_cgo: {{ .Values.metrics.enable_cgo }}
  {{- end }}
trace:
  {{- if .Values.trace }}
  enabled: {{ .Values.trace.enabled }}
  sampling_rate: {{ .Values.trace.sampling_rate }}
  {{- end }}
{{- end -}}

