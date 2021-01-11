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
Ingress port
*/}
{{- define "vald.ingressPort" -}}
port:
  {{- if regexMatch "^()([1-9]|[1-5]?[0-9]{2,4}|6[1-4][0-9]{3}|65[1-4][0-9]{2}|655[1-2][0-9]|6553[1-5])$" .Values.servicePort -}}
  number: {{ .Values.servicePort }}
  {{- else }}
  name: {{ .Values.servicePort }}
  {{- end -}}
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
gRPC client addr configuration
*/}}
{{- define "vald.grpc.client.addrs" -}}
{{- if .Values -}}
addrs:
  {{- toYaml .Values | nindent 2 }}
{{- else if .default -}}
addrs:
  {{- toYaml .default | nindent 2 }}
{{- else -}}
addrs: []
{{- end -}}
{{- end -}}

{{/*
gRPC client configuration
*/}}
{{- define "vald.grpc.client" -}}
health_check_duration: {{ default .default.health_check_duration .Values.health_check_duration | quote }}
connection_pool:
  {{- if .Values.connection_pool }}
  enable_dns_resolver: {{ default .default.connection_pool.enable_dns_resolver .Values.connection_pool.enable_dns_resolver }}
  enable_rebalance: {{ default .default.connection_pool.enable_rebalance .Values.connection_pool.enable_rebalance }}
  rebalance_duration: {{ default .default.connection_pool.rebalance_duration .Values.connection_pool.rebalance_duration | quote }}
  size: {{ default .default.connection_pool.size .Values.connection_pool.size }}
  old_conn_close_duration: {{ default .default.connection_pool.old_conn_close_duration .Values.connection_pool.old_conn_close_duration }}
  {{- else }}
  {{- toYaml .default.connection_pool | nindent 2 }}
  {{- end }}
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
  backoff_max_delay: {{ default .default.dial_option.backoff_max_delay .Values.dial_option.backoff_max_delay | quote }}
  backoff_base_delay: {{ default .default.dial_option.backoff_base_delay .Values.dial_option.backoff_base_delay | quote }}
  backoff_multiplier: {{ default .default.dial_option.backoff_multiplier .Values.dial_option.backoff_multiplier }}
  backoff_jitter: {{ default .default.dial_option.backoff_jitter .Values.dial_option.backoff_jitter }}
  min_connection_timeout: {{ default .default.dial_option.min_connection_timeout .Values.dial_option.min_connection_timeout | quote }}
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
    {{- if .Values.collector.metrics.version_info_labels }}
    version_info_labels:
      {{- toYaml .Values.collector.metrics.version_info_labels | nindent 6 }}
    {{- else if .default.collector.metrics.version_info_labels }}
    version_info_labels:
      {{- toYaml .default.collector.metrics.version_info_labels | nindent 6 }}
    {{- else }}
    version_info_labels: []
    {{- end }}
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
stackdriver:
  {{- if .Values.stackdriver }}
  project_id: {{ default .default.stackdriver.project_id .Values.stackdriver.project_id | quote }}
  client:
    {{- if .Values.stackdriver.client }}
    api_key: {{ default .default.stackdriver.client.api_key .Values.stackdriver.client.api_key | quote }}
    {{- if .Values.stackdriver.client.audiences }}
    audiences:
      {{- toYaml .Values.stackdriver.client.audiences | nindent 6 }}
    {{- else if .default.stackdriver.client.audiences }}
    audiences:
      {{- toYaml .default.stackdriver.client.audiences | nindent 6 }}
    {{- else }}
    audiences: []
    {{- end }}
    credentials_file: {{ default .default.stackdriver.client.credentials_file .Values.stackdriver.client.credentials_file | quote }}
    credentials_json: {{ default .default.stackdriver.client.credentials_json .Values.stackdriver.client.credentials_json | quote }}
    endpoint: {{ default .default.stackdriver.client.endpoint .Values.stackdriver.client.endpoint | quote }}
    quota_project: {{ default .default.stackdriver.client.quota_project .Values.stackdriver.client.quota_project | quote }}
    request_reason: {{ default .default.stackdriver.client.request_reason .Values.stackdriver.client.request_reason | quote }}
    {{- if .Values.stackdriver.client.scopes }}
    scopes:
      {{- toYaml .Values.stackdriver.client.scopes | nindent 6 }}
    {{- else if .default.stackdriver.client.scopes }}
    scopes:
      {{- toYaml .default.stackdriver.client.scopes | nindent 6 }}
    {{- else }}
    scopes: []
    {{- end }}
    user_agent: {{ default .default.stackdriver.client.user_agent .Values.stackdriver.client.user_agent | quote }}
    telemetry_enabled: {{ default .default.stackdriver.client.telemetry_enabled .Values.stackdriver.client.telemetry_enabled }}
    authentication_enabled: {{ default .default.stackdriver.client.authentication_enabled .Values.stackdriver.client.authentication_enabled }}
    {{- else }}
    {{- toYaml .default.stackdriver.client | nindent 4 }}
    {{- end }}
  exporter:
    {{- if .Values.stackdriver.exporter }}
    monitoring_enabled: {{ default .default.stackdriver.exporter.monitoring_enabled .Values.stackdriver.exporter.monitoring_enabled }}
    tracing_enabled: {{ default .default.stackdriver.exporter.tracing_enabled .Values.stackdriver.exporter.tracing_enabled }}
    location: {{ default .default.stackdriver.exporter.location .Values.stackdriver.exporter.location | quote }}
    bundle_delay_threshold: {{ default .default.stackdriver.exporter.bundle_delay_threshold .Values.stackdriver.exporter.bundle_delay_threshold | quote }}
    bundle_count_threshold: {{ default .default.stackdriver.exporter.bundle_count_threshold .Values.stackdriver.exporter.bundle_count_threshold }}
    trace_spans_buffer_max_bytes: {{ default .default.stackdriver.exporter.trace_spans_buffer_max_bytes .Values.stackdriver.exporter.trace_spans_buffer_max_bytes }}
    metric_prefix: {{ default .default.stackdriver.exporter.metric_prefix .Values.stackdriver.exporter.metric_prefix | quote }}
    skip_cmd: {{ default .default.stackdriver.exporter.skip_cmd .Values.stackdriver.exporter.skip_cmd }}
    timeout: {{ default .default.stackdriver.exporter.timeout .Values.stackdriver.exporter.timeout | quote }}
    reporting_interval: {{ default .default.stackdriver.exporter.reporting_interval .Values.stackdriver.exporter.reporting_interval | quote }}
    number_of_workers: {{ default .default.stackdriver.exporter.number_of_workers .Values.stackdriver.exporter.number_of_workers }}
    {{- else }}
    {{- toYaml .default.stackdriver.exporter | nindent 4 }}
    {{- end }}
  profiler:
    {{- if .Values.stackdriver.profiler }}
    enabled: {{ default .default.stackdriver.profiler.enabled .Values.stackdriver.profiler.enabled }}
    service: {{ default .default.stackdriver.profiler.service .Values.stackdriver.profiler.service | quote }}
    service_version: {{ default .default.stackdriver.profiler.service_version .Values.stackdriver.profiler.service_version | quote }}
    debug_logging: {{ default .default.stackdriver.profiler.debug_logging .Values.stackdriver.profiler.debug_logging }}
    mutex_profiling: {{ default .default.stackdriver.profiler.mutex_profiling .Values.stackdriver.profiler.mutex_profiling }}
    cpu_profiling: {{ default .default.stackdriver.profiler.cpu_profiling .Values.stackdriver.profiler.cpu_profiling }}
    alloc_profiling: {{ default .default.stackdriver.profiler.alloc_profiling .Values.stackdriver.profiler.alloc_profiling }}
    heap_profiling: {{ default .default.stackdriver.profiler.heap_profiling .Values.stackdriver.profiler.heap_profiling }}
    goroutine_profiling: {{ default .default.stackdriver.profiler.goroutine_profiling .Values.stackdriver.profiler.goroutine_profiling }}
    alloc_force_gc: {{ default .default.stackdriver.profiler.alloc_force_gc .Values.stackdriver.profiler.alloc_force_gc }}
    api_addr: {{ default .default.stackdriver.profiler.api_addr .Values.stackdriver.profiler.api_addr | quote }}
    instance: {{ default .default.stackdriver.profiler.instance .Values.stackdriver.profiler.instance | quote }}
    zone: {{ default .default.stackdriver.profiler.zone .Values.stackdriver.profiler.zone | quote }}
    {{- else }}
    {{- toYaml .default.stackdriver.profiler | nindent 4 }}
    {{- end }}
  {{- else }}
  {{- toYaml .default.stackdriver | nindent 2 }}
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
      {{- $compressorReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.manager.compressor.server_config.healths.readiness.port }}
      {{- $compressorReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.manager.compressor.name }}.{{ $.namespace }}.svc.cluster.local:{{ $compressorReadinessPort }}{{ $compressorReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
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
      {{- $backupManagerReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.manager.backup.server_config.healths.readiness.port }}
      {{- $backupManagerReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.manager.backup.name }}.{{ $.namespace }}.svc.cluster.local:{{ $backupManagerReadinessPort }}{{ $backupManagerReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "gateway-backup" }}
      {{- $backupGatewayReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.gateway.backup.server_config.healths.readiness.port }}
      {{- $backupGatewayReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.gateway.backup.name }}.{{ $.namespace }}.svc.cluster.local:{{ $backupGatewayReadinessPort }}{{ $backupGatewayReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "gateway-lb" }}
      {{- $lbGatewayReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.gateway.lb.server_config.healths.readiness.port }}
      {{- $lbGatewayReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.gateway.lb.name }}.{{ $.namespace }}.svc.cluster.local:{{ $lbGatewayReadinessPort }}{{ $lbGatewayReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "gateway-meta" }}
      {{- $metaGatewayReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.gateway.meta.server_config.healths.readiness.port }}
      {{- $metaGatewayReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.gateway.meta.name }}.{{ $.namespace }}.svc.cluster.local:{{ $metaGatewayReadinessPort }}{{ $metaGatewayReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
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

{{/*
Affinity rules
*/}}
{{- define "vald.affinity" -}}
nodeAffinity:
  {{- if .nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution }}
  preferredDuringSchedulingIgnoredDuringExecution:
    {{- toYaml .nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | nindent 4 }}
  {{- else }}
  preferredDuringSchedulingIgnoredDuringExecution: []
  {{- end }}
  {{- if .nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms }}
  requiredDuringSchedulingIgnoredDuringExecution:
    nodeSelectorTerms:
      {{- toYaml .nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | nindent 6 }}
  {{- end }}
podAffinity:
  {{- if .podAffinity.preferredDuringSchedulingIgnoredDuringExecution }}
  preferredDuringSchedulingIgnoredDuringExecution:
    {{- toYaml .podAffinity.preferredDuringSchedulingIgnoredDuringExecution | nindent 4 }}
  {{- else }}
  preferredDuringSchedulingIgnoredDuringExecution: []
  {{- end }}
  {{- if .podAffinity.requiredDuringSchedulingIgnoredDuringExecution }}
  requiredDuringSchedulingIgnoredDuringExecution:
    {{- toYaml .podAffinity.requiredDuringSchedulingIgnoredDuringExecution | nindent 4 }}
  {{- else }}
  requiredDuringSchedulingIgnoredDuringExecution: []
  {{- end }}
podAntiAffinity:
  {{- if .podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution }}
  preferredDuringSchedulingIgnoredDuringExecution:
    {{- toYaml .podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | nindent 4 }}
  {{- else }}
  preferredDuringSchedulingIgnoredDuringExecution: []
  {{- end }}
  {{- if .podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution }}
  requiredDuringSchedulingIgnoredDuringExecution:
    {{- toYaml .podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | nindent 4 }}
  {{- else }}
  requiredDuringSchedulingIgnoredDuringExecution: []
  {{- end }}
{{- end -}}
