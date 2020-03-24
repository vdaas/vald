{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "vald.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "vald.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "vald.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

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
{{- end -}}

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
logger: {{ default .default.logger .Values.logger | quote }}
level: {{ default .default.level .Values.level | quote }}
format: {{ default .default.format .Values.format | quote }}
{{- else }}
{{- toYaml .default }}
{{- end }}
{{- end -}}

{{/*
Container ports
*/}}
{{- define "vald.containerPorts" -}}
{{- $livenessEnabled := .default.healths.liveness.enabled }}
{{- if hasKey .Values.healths.liveness "enabled" }}
{{- $livenessEnabled = .Values.healths.liveness.enabled }}
{{- end }}
{{- if $livenessEnabled }}
livenessProbe:
  {{- if .Values.healths.liveness.livenessProbe }}
  httpGet:
    {{- if .Values.healths.liveness.livenessProbe.httpGet }}
    path: {{ default .default.healths.liveness.livenessProbe.httpGet.path .Values.healths.liveness.livenessProbe.httpGet.path }}
    port: {{ default .default.healths.liveness.livenessProbe.httpGet.port .Values.healths.liveness.livenessProbe.httpGet.port }}
    scheme: {{ default .default.healths.liveness.livenessProbe.httpGet.scheme .Values.healths.liveness.livenessProbe.httpGet.scheme }}
    {{- else }}
    {{- toYaml .default.healths.liveness.livenessProbe.httpGet | nindent 4 }}
    {{- end }}
  initialDelaySeconds: {{ default .default.healths.liveness.livenessProbe.initialDelaySeconds .Values.healths.liveness.livenessProbe.initialDelaySeconds }}
  timeoutSeconds: {{ default .default.healths.liveness.livenessProbe.timeoutSeconds .Values.healths.liveness.livenessProbe.timeoutSeconds }}
  successThreshold: {{ default .default.healths.liveness.livenessProbe.successThreshold .Values.healths.liveness.livenessProbe.successThreshold }}
  failureThreshold: {{ default .default.healths.liveness.livenessProbe.failureThreshold .Values.healths.liveness.livenessProbe.failureThreshold }}
  periodSeconds: {{ default .default.healths.liveness.livenessProbe.periodSeconds .Values.healths.liveness.livenessProbe.periodSeconds }}
  {{- else }}
  {{- toYaml .default.healths.liveness.livenessProbe | nindent 2 }}
  {{- end }}
{{- end }}
{{- $readinessEnabled := .default.healths.readiness.enabled }}
{{- if hasKey .Values.healths.readiness "enabled" }}
{{- $readinessEnabled = .Values.healths.readiness.enabled }}
{{- end }}
{{- if $readinessEnabled }}
readinessProbe:
  {{- if .Values.healths.readiness.readinessProbe }}
  httpGet:
    {{- if .Values.healths.readiness.readinessProbe.httpGet }}
    path: {{ default .default.healths.readiness.readinessProbe.httpGet.path .Values.healths.readiness.readinessProbe.httpGet.path }}
    port: {{ default .default.healths.readiness.readinessProbe.httpGet.port .Values.healths.readiness.readinessProbe.httpGet.port }}
    scheme: {{ default .default.healths.readiness.readinessProbe.httpGet.scheme .Values.healths.readiness.readinessProbe.httpGet.scheme }}
    {{- else }}
    {{- toYaml .default.healths.readiness.readinessProbe.httpGet | nindent 4 }}
    {{- end }}
  initialDelaySeconds: {{ default .default.healths.readiness.readinessProbe.initialDelaySeconds .Values.healths.readiness.readinessProbe.initialDelaySeconds }}
  timeoutSeconds: {{ default .default.healths.readiness.readinessProbe.timeoutSeconds .Values.healths.readiness.readinessProbe.timeoutSeconds }}
  successThreshold: {{ default .default.healths.readiness.readinessProbe.successThreshold .Values.healths.readiness.readinessProbe.successThreshold }}
  failureThreshold: {{ default .default.healths.readiness.readinessProbe.failureThreshold .Values.healths.readiness.readinessProbe.failureThreshold }}
  periodSeconds: {{ default .default.healths.readiness.readinessProbe.periodSeconds .Values.healths.readiness.readinessProbe.periodSeconds }}
  {{- else }}
  {{- toYaml .default.healths.readiness.readinessProbe | nindent 2 }}
  {{- end }}
{{- end }}
ports:
  {{- if $livenessEnabled }}
  - name: liveness
    protocol: TCP
    containerPort: {{ default .default.healths.liveness.port .Values.healths.liveness.port }}
  {{- end }}
  {{- if $readinessEnabled }}
  - name: readiness
    protocol: TCP
    containerPort: {{ default .default.healths.readiness.port .Values.healths.readiness.port }}
  {{- end }}
  {{- $restEnabled := .default.servers.rest.enabled }}
  {{- if hasKey .Values.servers.rest "enabled" }}
  {{- $restEnabled = .Values.servers.rest.enabled }}
  {{- end }}
  {{- if $restEnabled }}
  - name: rest
    protocol: TCP
    containerPort: {{ default .default.servers.rest.port .Values.servers.rest.port }}
  {{- end }}
  {{- $grpcEnabled := .default.servers.grpc.enabled }}
  {{- if hasKey .Values.servers.grpc "enabled" }}
  {{- $grpcEnabled = .Values.servers.grpc.enabled }}
  {{- end }}
  {{- if $grpcEnabled }}
  - name: grpc
    protocol: TCP
    containerPort: {{ default .default.servers.grpc.port .Values.servers.grpc.port }}
  {{- end }}
  {{- $pprofEnabled := .default.metrics.pprof.enabled }}
  {{- if hasKey .Values.metrics.pprof "enabled" }}
  {{- $pprofEnabled = .Values.metrics.pprof.enabled }}
  {{- end }}
  {{- if $pprofEnabled }}
  - name: pprof
    protocol: TCP
    containerPort: {{ default .default.metrics.pprof.port .Values.metrics.pprof.port }}
  {{- end }}
  {{- $prometheusEnabled := .default.metrics.prometheus.enabled }}
  {{- if hasKey .Values.metrics.prometheus "enabled" }}
  {{- $prometheusEnabled = .Values.metrics.prometheus.enabled }}
  {{- end }}
  {{- if $prometheusEnabled }}
  - name: prometheus
    protocol: TCP
    containerPort: {{ default .default.metrics.prometheus.port .Values.metrics.prometheus.port }}
  {{- end }}
{{- end -}}

{/*
Service ports
*/}
{{- define "vald.servicePorts" -}}
ports:
  {{- $restEnabled := .default.servers.rest.enabled }}
  {{- if hasKey .Values.servers.rest "enabled" }}
  {{- $restEnabled = .Values.servers.rest.enabled }}
  {{- end }}
  {{- if $restEnabled }}
  - name: rest
    port: {{ default .default.servers.rest.servicePort .Values.servers.rest.servicePort }}
    targetPort: {{ default .default.servers.rest.port .Values.servers.rest.port }}
    protocol: TCP
  {{- end }}
  {{- $grpcEnabled := .default.servers.grpc.enabled }}
  {{- if hasKey .Values.servers.grpc "enabled" }}
  {{- $grpcEnabled = .Values.servers.grpc.enabled }}
  {{- end }}
  {{- if $grpcEnabled }}
  - name: grpc
    port: {{ default .default.servers.grpc.servicePort .Values.servers.grpc.servicePort }}
    targetPort: {{ default .default.servers.grpc.port .Values.servers.grpc.port }}
    protocol: TCP
  {{- end }}
  {{- $readinessEnabled := .default.healths.readiness.enabled }}
  {{- if hasKey .Values.healths.readiness "enabled" }}
  {{- $readinessEnabled = .Values.healths.readiness.enabled }}
  {{- end }}
  {{- if $readinessEnabled }}
  - name: readiness
    port: {{ default .default.healths.readiness.servicePort .Values.healths.readiness.servicePort }}
    targetPort: {{ default .default.healths.readiness.port .Values.healths.readiness.port }}
    protocol: TCP
  {{- end }}
  {{- $pprofEnabled := .default.metrics.pprof.enabled }}
  {{- if hasKey .Values.metrics.pprof "enabled" }}
  {{- $pprofEnabled = .Values.metrics.pprof.enabled }}
  {{- end }}
  {{- if $pprofEnabled }}
  - name: pprof
    port: {{ default .default.metrics.pprof.servicePort .Values.metrics.pprof.servicePort }}
    targetPort: {{ default .default.metrics.pprof.port .Values.metrics.pprof.port }}
    protocol: TCP
  {{- end }}
  {{- $prometheusEnabled := .default.metrics.prometheus.enabled }}
  {{- if hasKey .Values.metrics.prometheus "enabled" }}
  {{- $prometheusEnabled = .Values.metrics.prometheus.enabled }}
  {{- end }}
  {{- if $prometheusEnabled }}
  - name: prometheus
    port: {{ default .default.metrics.prometheus.servicePort .Values.metrics.prometheus.servicePort }}
    targetPort: {{ default .default.metrics.prometheus.port .Values.metrics.prometheus.port }}
    protocol: TCP
  {{- end }}
{{- end -}}

{{/*
Server configures that inserted into server_config
*/}}
{{- define "vald.servers" -}}
servers:
  {{- $restEnabled := .default.servers.rest.enabled }}
  {{- if hasKey .Values.servers.rest "enabled" }}
  {{- $restEnabled = .Values.servers.rest.enabled }}
  {{- end }}
  {{- if $restEnabled }}
  - name: rest
    host: {{ default .default.servers.rest.host .Values.servers.rest.host }}
    port: {{ default .default.servers.rest.port .Values.servers.rest.port }}
    {{- if .Values.servers.rest.server }}
    mode: {{ default .default.servers.rest.server.mode .Values.servers.rest.server.mode }}
    probe_wait_time: {{ default .default.servers.rest.server.probe_wait_time .Values.servers.rest.server.probe_wait_time }}
    http:
      {{- if .Values.servers.rest.server.http }}
      shutdown_duration: {{ default .default.servers.rest.server.http.shutdown_duration .Values.servers.rest.server.http.shutdown_duration }}
      handler_timeout: {{ default .default.servers.rest.server.http.handler_timeout .Values.servers.rest.server.http.handler_timeout }}
      idle_timeout: {{ default .default.servers.rest.server.http.idle_timeout .Values.servers.rest.server.http.idle_timeout }}
      read_header_timeout: {{ default .default.servers.rest.server.http.read_header_timeout .Values.servers.rest.server.http.read_header_timeout }}
      read_timeout: {{ default .default.servers.rest.server.http.read_timeout .Values.servers.rest.server.http.read_timeout }}
      write_timeout: {{ default .default.servers.rest.server.http.write_timeout .Values.servers.rest.server.http.write_timeout }}
      {{- else }}
      {{- toYaml .default.servers.rest.server.http | nindent 6}}
      {{- end }}
    {{- else }}
    {{- toYaml .default.servers.rest.server | nindent 4 }}
    {{- end }}
  {{- end }}
  {{- $grpcEnabled := .default.servers.grpc.enabled }}
  {{- if hasKey .Values.servers.grpc "enabled" }}
  {{- $grpcEnabled = .Values.servers.grpc.enabled }}
  {{- end }}
  {{- if $grpcEnabled }}
  - name: grpc
    host: {{ default .default.servers.grpc.host .Values.servers.grpc.host }}
    port: {{ default .default.servers.grpc.port .Values.servers.grpc.port }}
    {{- if .Values.servers.grpc.server }}
    mode: {{ default .default.servers.grpc.server.mode .Values.servers.grpc.server.mode }}
    probe_wait_time: {{ default .default.servers.grpc.server.probe_wait_time .Values.servers.grpc.server.probe_wait_time | quote }}
    grpc:
      {{- if .Values.servers.grpc.server.grpc }}
      max_receive_message_size: {{ default .default.servers.grpc.server.grpc.max_receive_message_size .Values.servers.grpc.server.grpc.max_receive_message_size }}
      max_send_message_size: {{ default .default.servers.grpc.server.grpc.max_send_message_size .Values.servers.grpc.server.grpc.max_send_message_size }}
      initial_window_size: {{ default .default.servers.grpc.server.grpc.initial_window_size .Values.servers.grpc.server.grpc.initial_window_size }}
      initial_conn_window_size: {{ default .default.servers.grpc.server.grpc.initial_conn_window_size .Values.servers.grpc.server.grpc.initial_conn_window_size }}
      keepalive:
        {{- if .Values.servers.grpc.server.grpc.keepalive }}
        max_conn_idle: {{ default .default.servers.grpc.server.grpc.keepalive.max_conn_idle .Values.servers.grpc.server.grpc.keepalive.max_conn_idle | quote }}
        max_conn_age: {{ default .default.servers.grpc.server.grpc.keepalive.max_conn_age .Values.servers.grpc.server.grpc.keepalive.max_conn_age | quote }}
        max_conn_age_grace: {{ default .default.servers.grpc.server.grpc.keepalive.max_conn_age_grace .Values.servers.grpc.server.grpc.keepalive.max_conn_age_grace | quote }}
        time: {{ default .default.servers.grpc.server.grpc.keepalive.time .Values.servers.grpc.server.grpc.keepalive.time | quote }}
        timeout: {{ default .default.servers.grpc.server.grpc.keepalive.timeout .Values.servers.grpc.server.grpc.keepalive.timeout | quote }}
        {{- else }}
        {{- toYaml .default.servers.grpc.server.grpc.keepalive | nindent 8 }}
        {{- end }}
      write_buffer_size: {{ default .default.servers.grpc.server.grpc.write_buffer_size .Values.servers.grpc.server.grpc.write_buffer_size }}
      read_buffer_size: {{ default .default.servers.grpc.server.grpc.read_buffer_size .Values.servers.grpc.server.grpc.read_buffer_size }}
      connection_timeout: {{ default .default.servers.grpc.server.grpc.connection_timeout .Values.servers.grpc.server.grpc.connection_timeout | quote }}
      max_header_list_size: {{ default .default.servers.grpc.server.grpc.max_header_list_size .Values.servers.grpc.server.grpc.max_header_list_size }}
      header_table_size: {{ default .default.servers.grpc.server.grpc.header_table_size .Values.servers.grpc.server.grpc.header_table_size }}
      interceptors: {{ default .default.servers.grpc.server.grpc.interceptors .Values.servers.grpc.server.grpc.interceptors }}
      {{- else }}
      {{- toYaml .default.servers.grpc.server.grpc | nindent 6 }}
      {{- end }}
    restart: true
    {{- else }}
    {{- toYaml .default.servers.grpc.server | nindent 4 }}
    {{- end }}
  {{- end }}
health_check_servers:
  {{- $livenessEnabled := .default.healths.liveness.enabled }}
  {{- if hasKey .Values.healths.liveness "enabled" }}
  {{- $livenessEnabled = .Values.healths.liveness.enabled }}
  {{- end }}
  {{- if $livenessEnabled }}
  - name: liveness
    host: {{ default .default.healths.liveness.host .Values.healths.liveness.host }}
    port: {{ default .default.healths.liveness.port .Values.healths.liveness.port }}
    {{- if .Values.healths.liveness.server }}
    mode: {{ default .default.healths.liveness.server.mode .Values.healths.liveness.server.mode | quote }}
    probe_wait_time: {{ default .default.healths.liveness.server.probe_wait_time .Values.healths.liveness.server.probe_wait_time | quote }}
    http:
      {{- if .Values.healths.liveness.server.http }}
      shutdown_duration: {{ default .default.healths.liveness.server.http.shutdown_duration .Values.healths.liveness.server.http.shutdown_duration | quote }}
      handler_timeout: {{ default .default.healths.liveness.server.http.handler_timeout .Values.healths.liveness.server.http.handler_timeout | quote }}
      idle_timeout: {{ default .default.healths.liveness.server.http.idle_timeout .Values.healths.liveness.server.http.idle_timeout | quote }}
      read_header_timeout: {{ default .default.healths.liveness.server.http.read_header_timeout .Values.healths.liveness.server.http.read_header_timeout | quote }}
      read_timeout: {{ default .default.healths.liveness.server.http.read_timeout .Values.healths.liveness.server.http.read_timeout | quote }}
      write_timeout: {{ default .default.healths.liveness.server.http.write_timeout .Values.healths.liveness.server.http.write_timeout | quote }}
      {{- else }}
      {{- toYaml .default.healths.liveness.server.http | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.healths.liveness.server | nindent 4 }}
    {{- end }}
  {{- end }}
  {{- $readinessEnabled := .default.healths.readiness.enabled }}
  {{- if hasKey .Values.healths.readiness "enabled" }}
  {{- $readinessEnabled = .Values.healths.readiness.enabled }}
  {{- end }}
  {{- if $readinessEnabled }}
  - name: readiness
    host: {{ default .default.healths.readiness.host .Values.healths.readiness.host }}
    port: {{ default .default.healths.readiness.port .Values.healths.readiness.port }}
    {{- if .Values.healths.readiness.server }}
    mode: {{ default .default.healths.readiness.server.mode .Values.healths.readiness.server.mode | quote }}
    probe_wait_time: {{ default .default.healths.readiness.server.probe_wait_time .Values.healths.readiness.server.probe_wait_time | quote }}
    http:
      {{- if .Values.healths.readiness.server.http }}
      shutdown_duration: {{ default .default.healths.readiness.server.http.shutdown_duration .Values.healths.readiness.server.http.shutdown_duration | quote }}
      handler_timeout: {{ default .default.healths.readiness.server.http.handler_timeout .Values.healths.readiness.server.http.handler_timeout | quote }}
      idle_timeout: {{ default .default.healths.readiness.server.http.idle_timeout .Values.healths.readiness.server.http.idle_timeout | quote }}
      read_header_timeout: {{ default .default.healths.readiness.server.http.read_header_timeout .Values.healths.readiness.server.http.read_header_timeout | quote }}
      read_timeout: {{ default .default.healths.readiness.server.http.read_timeout .Values.healths.readiness.server.http.read_timeout | quote }}
      write_timeout: {{ default .default.healths.readiness.server.http.write_timeout .Values.healths.readiness.server.http.write_timeout | quote }}
      {{- else }}
      {{- toYaml .default.healths.readiness.server.http | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.healths.readiness.server | nindent 4 }}
    {{- end }}
  {{- end }}
metrics_servers:
  {{- $pprofEnabled := .default.metrics.pprof.enabled }}
  {{- if hasKey .Values.metrics.pprof "enabled" }}
  {{- $pprofEnabled = .Values.metrics.pprof.enabled }}
  {{- end }}
  {{- if $pprofEnabled }}
  - name: pprof
    host: {{ default .default.metrics.pprof.host .Values.metrics.pprof.host }}
    port: {{ default .default.metrics.pprof.port .Values.metrics.pprof.port }}
    {{- if .Values.metrics.pprof.server }}
    mode: {{ default .default.metrics.pprof.server.mode .Values.metrics.pprof.server.mode }}
    probe_wait_time: {{ default .default.metrics.pprof.server.probe_wait_time .Values.metrics.pprof.server.probe_wait_time }}
    http:
      {{- if .Values.metrics.pprof.server.http }}
      shutdown_duration: {{ default .default.metrics.pprof.server.http.shutdown_duration .Values.metrics.pprof.server.http.shutdown_duration }}
      handler_timeout: {{ default .default.metrics.pprof.server.http.handler_timeout .Values.metrics.pprof.server.http.handler_timeout }}
      idle_timeout: {{ default .default.metrics.pprof.server.http.idle_timeout .Values.metrics.pprof.server.http.idle_timeout }}
      read_header_timeout: {{ default .default.metrics.pprof.server.http.read_header_timeout .Values.metrics.pprof.server.http.read_header_timeout }}
      read_timeout: {{ default .default.metrics.pprof.server.http.read_timeout .Values.metrics.pprof.server.http.read_timeout }}
      write_timeout: {{ default .default.metrics.pprof.server.http.write_timeout .Values.metrics.pprof.server.http.write_timeout }}
      {{- else }}
      {{- toYaml .default.metrics.pprof.server.http | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.metrics.pprof.server | nindent 4 }}
    {{- end }}
  {{- end }}
  {{- $prometheusEnabled := .default.metrics.prometheus.enabled }}
  {{- if hasKey .Values.metrics.prometheus "enabled" }}
  {{- $prometheusEnabled = .Values.metrics.prometheus.enabled }}
  {{- end }}
  {{- if $prometheusEnabled }}
  - name: prometheus
    host: {{ default .default.metrics.prometheus.host .Values.metrics.prometheus.host }}
    port: {{ default .default.metrics.prometheus.port .Values.metrics.prometheus.port }}
    {{- if .Values.metrics.prometheus.server }}
    mode: {{ default .default.metrics.prometheus.server.mode .Values.metrics.prometheus.server.mode }}
    probe_wait_time: {{ default .default.metrics.prometheus.server.probe_wait_time .Values.metrics.prometheus.server.probe_wait_time }}
    http:
      {{- if .Values.metrics.prometheus.server.http }}
      shutdown_duration: {{ default .default.metrics.prometheus.server.http.shutdown_duration .Values.metrics.prometheus.server.http.shutdown_duration }}
      handler_timeout: {{ default .default.metrics.prometheus.server.http.handler_timeout .Values.metrics.prometheus.server.http.handler_timeout }}
      idle_timeout: {{ default .default.metrics.prometheus.server.http.idle_timeout .Values.metrics.prometheus.server.http.idle_timeout }}
      read_header_timeout: {{ default .default.metrics.prometheus.server.http.read_header_timeout .Values.metrics.prometheus.server.http.read_header_timeout }}
      read_timeout: {{ default .default.metrics.prometheus.server.http.read_timeout .Values.metrics.prometheus.server.http.read_timeout }}
      write_timeout: {{ default .default.metrics.prometheus.server.http.write_timeout .Values.metrics.prometheus.server.http.write_timeout }}
      {{- else }}
      {{- toYaml .default.metrics.prometheus.server.http | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.metrics.prometheus.server | nindent 4 }}
    {{- end }}
  {{- end }}
startup_strategy:
  {{- if $livenessEnabled }}
  - liveness
  {{- end }}
  {{- if $pprofEnabled }}
  - pprof
  {{- end }}
  {{- if $prometheusEnabled }}
  - prometheus
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
full_shutdown_duration: {{ default .default.full_shutdown_duration .Values.full_shutdown_duration }}
tls:
  {{- if .Values.tls }}
  enabled: {{ default .default.tls.enabled .Values.tls.enabled }}
  cert: {{ default .default.tls.cert .Values.tls.cert | quote }}
  key: {{ default .default.tls.key .Values.tls.key | quote }}
  ca: {{ default .default.tls.ca .Values.tls.ca | quote }}
  {{- else }}
  {{- toYaml .default.tls | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
gRPC client configuration
*/}}
{{- define "vald.grpc.client" -}}
addrs: {{ default .default.addrs .Values.addrs }}
connection_pool: {{ default .default.connection_pool .Values.connection_pool }}
health_check_duration: {{ default .default.health_check_duration .Values.health_check_duration | quote }}
backoff:
  {{- if .Values.backoff }}
  initial_duration: {{ default .default.backoff.initial_duration .Values.backoff.initial_duration | quote }}
  backoff_time_limit: {{ default .default.backoff.backoff_time_limit .Values.backoff.backoff_time_limit | quote }}
  maximum_duration: {{ default .default.backoff.maximum_duration .Values.backoff.maximum_duration | quote }}
  jitter_limit: {{ default .default.backoff.jitter_limit .Values.backoff.jitter_limit | quote }}
  backoff_factor: {{ default .default.backoff.backoff_factor .Values.backoff.backoff_factor }}
  retry_count: {{ default .default.backoff.retry_count .Values.backoff.retry_count }}
  enable_error_log: {{ default .default.backoff.enable_error_log .Values.backoff.enable_error_log }}
  {{- else }}
  {{- toYaml .default.backoff | nindent 2 }}
  {{- end }}
call_option:
  {{- if .Values.call_option }}
  wait_for_ready: {{ default .default.call_option.wait_for_ready .Values.call_option.wait_for_ready }}
  max_retry_rpc_buffer_size: {{ default .default.call_option.max_retry_rpc_buffer_size .Values.call_option.max_retry_rpc_buffer_size }}
  max_recv_msg_size: {{ default .default.call_option.max_recv_msg_size .Values.call_option.max_recv_msg_size }}
  max_send_msg_size: {{ default .default.call_option.max_send_msg_size .Values.call_option.max_send_msg_size }}
  {{- else }}
  {{- toYaml .default.call_option | nindent 2 }}
  {{- end }}
dial_option:
  {{- if .Values.dial_option }}
  write_buffer_size: {{ default .default.dial_option.write_buffer_size .Values.dial_option.write_buffer_size }}
  read_buffer_size: {{ default .default.dial_option.read_buffer_size .Values.dial_option.read_buffer_size }}
  initial_window_size: {{ default .default.dial_option.initial_window_size .Values.dial_option.initial_window_size }}
  initial_connection_window_size: {{ default .default.dial_option.initial_connection_window_size .Values.dial_option.initial_connection_window_size }}
  max_msg_size: {{ default .default.dial_option.max_msg_size .Values.dial_option.max_msg_size }}
  max_backoff_delay: {{ default .default.dial_option.max_backoff_delay .Values.dial_option.max_backoff_delay | quote }}
  enable_backoff: {{ default .default.dial_option.enable_backoff .Values.dial_option.enable_backoff }}
  insecure: {{ default .default.dial_option.insecure .Values.dial_option.insecure }}
  timeout: {{ default .default.dial_option.timeout .Values.dial_option.timeout | quote }}
  tcp:
    {{- if .Values.dial_option.tcp }}
    dns:
      {{- if .Values.dial_option.tcp.dns }}
      cache_enabled: {{ default .default.dial_option.tcp.dns.cache_enabled .Values.dial_option.tcp.dns.cache_enabled }}
      refresh_duration: {{ default .default.dial_option.tcp.dns.refresh_duration .Values.dial_option.tcp.dns.refresh_duration | quote }}
      cache_expiration: {{ default .default.dial_option.tcp.dns.cache_expiration .Values.dial_option.tcp.dns.cache_expiration | quote }}
      {{- else }}
      {{- toYaml .default.dial_option.tcp.dns | nindent 6 }}
      {{- end }}
    dialer:
      {{- if .Values.dial_option.tcp.dialer }}
      timeout: {{ default .default.dial_option.tcp.dialer.timeout .Values.dial_option.tcp.dialer.timeout | quote }}
      keep_alive: {{ default .default.dial_option.tcp.dialer.keep_alive .Values.dial_option.tcp.dialer.keep_alive | quote }}
      dual_stack_enabled: {{ default .default.dial_option.tcp.dialer.dual_stack_enabled .Values.dial_option.tcp.dialer.dual_stack_enabled }}
      {{- else }}
      {{- toYaml .default.dial_option.tcp.dialer | nindent 6 }}
      {{- end }}
    tls:
      {{- if .Values.dial_option.tcp.tls }}
      enabled: {{ default .default.dial_option.tcp.tls.enabled .Values.dial_option.tcp.tls.enabled }}
      cert: {{ default .default.dial_option.tcp.tls.cert .Values.dial_option.tcp.tls.cert | quote }}
      key: {{ default .default.dial_option.tcp.tls.key .Values.dial_option.tcp.tls.key | quote }}
      ca: {{ default .default.dial_option.tcp.tls.ca .Values.dial_option.tcp.tls.ca | quote }}
      {{- else }}
      {{- toYaml .default.dial_option.tcp.tls | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.dial_option.tcp | nindent 4 }}
    {{- end }}
  keep_alive:
    {{- if .Values.dial_option.keep_alive }}
    time: {{ default .default.dial_option.keep_alive.time .Values.dial_option.keep_alive.time | quote }}
    timeout: {{ default .default.dial_option.keep_alive.timeout .Values.dial_option.keep_alive.timeout | quote }}
    permit_without_stream: {{ default .default.dial_option.keep_alive.permit_without_stream .Values.dial_option.keep_alive.permit_without_stream }}
    {{- else }}
    {{- toYaml .default.dial_option.keep_alive | nindent 4 }}
    {{- end }}
  {{- else }}
  {{- toYaml .default.dial_option | nindent 2 }}
  {{- end }}
tls:
  {{- if .Values.tls }}
  enabled: {{ default .default.tls.enabled .Values.tls.enabled }}
  cert: {{ default .default.tls.cert .Values.tls.cert | quote }}
  key: {{ default .default.tls.key .Values.tls.key | quote }}
  ca: {{ default .default.tls.ca .Values.tls.ca | quote }}
  {{- else }}
  {{- toYaml .default.tls | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
observability
*/}}
{{- define "vald.observability" -}}
enabled: {{ default .default.enabled .Values.enabled }}
collector:
  {{- if .Values.collector }}
  duration: {{ default .default.collector.duration .Values.collector.duration }}
  metrics:
    {{- if .Values.collector.metrics }}
      enable_version_info: {{ default .default.collector.metrics.enable_version_info .Values.collector.metrics.enable_version_info }}
      enable_memory: {{ default .default.collector.metrics.enable_memory .Values.collector.metrics.enable_memory }}
      enable_goroutine: {{ default .default.collector.metrics.enable_goroutine .Values.collector.metrics.enable_goroutine }}
      enable_cgo: {{ default .default.collector.metrics.enable_cgo .Values.collector.metrics.enable_cgo }}
    {{- else }}
    {{- toYaml .default.collector.metrics | nindent 4 }}
    {{- end }}
  {{- else }}
  {{- toYaml .default.collector | nindent 2 }}
  {{- end }}
trace:
  {{- if .Values.trace }}
  enabled: {{ default .default.trace.enabled .Values.trace.enabled }}
  sampling_rate: {{ default .default.trace.sampling_rate .Values.trace.sampling_rate }}
  {{- else }}
  {{- toYaml .default.trace | nindent 2 }}
  {{- end }}
prometheus:
  {{- if .Values.prometheus }}
    enabled: {{ default .default.prometheus.enabled .Values.prometheus.enabled }}
  {{- else }}
  {{- toYaml .default.prometheus | nindent 2 }}
  {{- end }}
jaeger:
  {{- if .Values.jaeger }}
    enabled: {{ default .default.jaeger.enabled .Values.jaeger.enabled }}
    collector_endpoint: {{ default .default.jaeger.collector_endpoint .Values.jaeger.collector_endpoint | quote }}
    agent_endpoint: {{ default .default.jaeger.agent_endpoint .Values.jaeger.agent_endpoint | quote }}
    username: {{ default .default.jaeger.username .Values.jaeger.username | quote }}
    password: {{ default .default.jaeger.password .Values.jaeger.password | quote }}
    service_name: {{ default .default.jaeger.service_name .Values.jaeger.service_name | quote }}
    buffer_max_count: {{ default .default.jaeger.buffer_max_count .Values.jaeger.buffer_max_count }}
  {{- else }}
  {{- toYaml .default.jaeger | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
initContainers
*/}}
{{- define "vald.initContainers" -}}
{{- range .initContainers -}}
{{- if .type }}
- name: {{ .name }}
  image: {{ .image }}
  {{- if eq .type "wait-for" }}
  command:
    - /bin/sh
    - -e
    - -c
    - |
      {{- if eq .target "compressor" }}
      {{- $compressorReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.compressor.server_config.healths.readiness.port }}
      {{- $compressorReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.compressor.name }}.{{ $.namespace }}.svc.cluster.local:{{ $compressorReadinessPort }}{{ $compressorReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "meta" }}
      {{- $metaReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.meta.server_config.healths.readiness.port }}
      {{- $metaReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.meta.name }}.{{ $.namespace }}.svc.cluster.local:{{ $metaReadinessPort }}{{ $metaReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "discoverer" }}
      {{- $discovererReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.discoverer.server_config.healths.readiness.port }}
      {{- $discovererReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.discoverer.name }}.{{ $.namespace }}.svc.cluster.local:{{ $discovererReadinessPort }}{{ $discovererReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "agent" }}
      {{- $agentReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.agent.server_config.healths.readiness.port }}
      {{- $agentReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.agent.name }}.{{ $.namespace }}.svc.cluster.local:{{ $agentReadinessPort }}{{ $agentReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "manager-backup" }}
      {{- $backupManagerReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.backupManager.server_config.healths.readiness.port }}
      {{- $backupManagerReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.backupManager.name }}.{{ $.namespace }}.svc.cluster.local:{{ $backupManagerReadinessPort }}{{ $backupManagerReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if .untilCondition }}
      until [ {{ .untilCondition }} ]; do
      {{- else if .whileCondition }}
      while [ {{ .whileCondition }} ]; do
      {{- end }}
        echo "waiting for {{ .target }} to be ready..."
        sleep {{ .sleepDuration }};
      done
  {{- else if eq .type "wait-for-mysql" }}
  command:
    - /bin/sh
    - -e
    - -c
    - |
      hosts="{{ include "vald.utils.joinListWithSpace" .mysql.hosts }}"
      options="{{ include "vald.utils.joinListWithSpace" .mysql.options }}"
      for host in $hosts; do
        until [ "$(mysqladmin -h$host $options --show-warnings=false ping | grep alive | awk '{print $3}')" = "alive" ]; do
          echo "waiting for $host to be ready..."
          sleep {{ .sleepDuration }};
        done
      done
  {{- else if eq .type "wait-for-redis" }}
  command:
    - /bin/sh
    - -e
    - -c
    - |
      hosts="{{ include "vald.utils.joinListWithSpace" .redis.hosts }}"
      options="{{ include "vald.utils.joinListWithSpace" .redis.options }}"
      for host in $hosts; do
        until [ "$(redis-cli -h $host $options ping)" = "PONG" ]; do
          echo "waiting for $host to be ready..."
          sleep {{ .sleepDuration }};
        done
      done
  {{- else if eq .type "wait-for-cassandra" }}
  command:
    - /bin/sh
    - -e
    - -c
    - |
      hosts="{{ include "vald.utils.joinListWithSpace" .cassandra.hosts }}"
      options="{{ include "vald.utils.joinListWithSpace" .cassandra.options }}"
      for host in $hosts; do
        until cqlsh $host $options -e "select now() from system.local" > /dev/null; do
          echo "waiting for $host to be ready..."
          sleep {{ .sleepDuration }};
        done
      done
  {{- else if eq .type "limit-vsz" }}
  command:
    - /bin/sh
    - -e
    - -c
    - |
      set -eu
      cgroup_rsslimit="/sys/fs/cgroup/memory/memory.limit_in_bytes"
      if [ -r "$cgroup_rsslimit" ] ; then
        rsslimit=`cat "$cgroup_rsslimit"`
        vszlimit=`expr $rsslimit / 1024`
        ulimit -v $vszlimit
      fi
  {{- end }}
  {{- if .env }}
  env:
    {{- toYaml .env | nindent 4 }}
  {{- end }}
  {{- if .volumeMounts }}
  volumeMounts:
    {{- toYaml .volumeMounts | nindent 4 }}
  {{- end }}
{{- else }}
- {{ . | toYaml | nindent 2 | trim }}
{{- end }}
{{- end }}
{{- end -}}
