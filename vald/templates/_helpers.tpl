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
Container ports
*/}}
{{- define "vald.containerPorts" -}}
{{- $livenessEnabled := default .default.healths.liveness.enabled .Values.healths.liveness.enabled }}
{{- if $livenessEnabled }}
livenessProbe:
  httpGet:
    path: /liveness
    port: liveness
    scheme: HTTP
  initialDelaySeconds: 5
  timeoutSeconds: 2
  successThreshold: 1
  failureThreshold: 2
  periodSeconds: 3
{{- end }}
{{- $readinessEnabled := default .default.healths.readiness.enabled .Values.healths.readiness.enabled }}
{{- if $readinessEnabled }}
readinessProbe:
  httpGet:
    path: /readiness
    port: readiness
    scheme: HTTP
  initialDelaySeconds: 10
  timeoutSeconds: 2
  successThreshold: 1
  failureThreshold: 2
  periodSeconds: 3
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
  {{- $restEnabled := default .default.servers.rest.enabled .Values.servers.rest.enabled }}
  {{- if $restEnabled }}
  - name: rest
    protocol: TCP
    containerPort: {{ default .default.servers.rest.port .Values.servers.rest.port }}
  {{- end }}
  {{- $grpcEnabled := default .default.servers.grpc.enabled .Values.servers.grpc.enabled }}
  {{- if $grpcEnabled }}
  - name: grpc
    protocol: TCP
    containerPort: {{ default .default.servers.grpc.port .Values.servers.grpc.port }}
  {{- end }}
  {{- $pprofEnabled := default .default.metrics.pprof.enabled .Values.metrics.pprof.enabled }}
  {{- if $pprofEnabled }}
  - name: pprof
    protocol: TCP
    containerPort: {{ default .default.metrics.pprof.port .Values.metrics.pprof.port }}
  {{- end }}
{{- end -}}

{/*
Service ports
*/}
{{- define "vald.servicePorts" -}}
ports:
  {{- $restEnabled := default .default.servers.rest.enabled .Values.servers.rest.enabled }}
  {{- if $restEnabled }}
  - name: rest
    port: {{ default .default.servers.rest.port .Values.servers.rest.port }}
  {{- end }}
  {{- $grpcEnabled := default .default.servers.grpc.enabled .Values.servers.grpc.enabled }}
  {{- if $grpcEnabled }}
  - name: grpc
    port: {{ default .default.servers.grpc.port .Values.servers.grpc.port }}
  {{- end }}
  {{- $pprofEnabled := default .default.metrics.pprof.enabled .Values.metrics.pprof.enabled }}
  {{- if $pprofEnabled }}
  - name: pprof
    port: {{ default .default.metrics.pprof.port .Values.metrics.pprof.port }}
  {{- end }}
{{- end -}}

{{/*
Server configures that inserted into server_config
*/}}
{{- define "vald.servers" -}}
servers:
  {{- $restEnabled := default .default.servers.rest.enabled .Values.servers.rest.enabled }}
  {{- if $restEnabled }}
  - name: rest
    host: {{ default .default.servers.rest.host .Values.servers.rest.host }}
    port: {{ default .default.servers.rest.port .Values.servers.rest.port }}
    mode: REST
    probe_wait_time: 3s
    http:
      shutdown_duration: 5s
      handler_timeout: 5s
      idle_timeout: 2s
      read_header_timeout: 1s
      read_timeout: 1s
      write_timeout: 1s
  {{- end }}
  {{- $grpcEnabled := default .default.servers.grpc.enabled .Values.servers.grpc.enabled }}
  {{- if $grpcEnabled }}
  - name: grpc
    host: {{ default .default.servers.grpc.host .Values.servers.grpc.host }}
    port: {{ default .default.servers.grpc.port .Values.servers.grpc.port }}
    mode: GRPC
    probe_wait_time: "3s"
    grpc:
      max_receive_message_size: 0
      max_send_message_size: 0
      initial_window_size: 0
      initial_conn_window_size: 0
      keepalive:
        max_conn_idle: ""
        max_conn_age: ""
        max_conn_age_grace: ""
        time: ""
        timeout: ""
      write_buffer_size: 0
      read_buffer_size: 0
      connection_timeout: ""
      max_header_list_size: 0
      header_table_size: 0
      interceptors: []
    restart: true
  {{- end }}
health_check_servers:
  {{- $livenessEnabled := default .default.healths.liveness.enabled .Values.healths.liveness.enabled }}
  {{- if $livenessEnabled }}
  - name: liveness
    host: {{ default .default.healths.liveness.host .Values.healths.liveness.host }}
    port: {{ default .default.healths.liveness.port .Values.healths.liveness.port }}
    mode: ""
    probe_wait_time: "3s"
    http:
      shutdown_duration: "5s"
      handler_timeout: ""
      idle_timeout: ""
      read_header_timeout: ""
      read_timeout: ""
      write_timeout: ""
  {{- end }}
  {{- $readinessEnabled := default .default.healths.readiness.enabled .Values.healths.readiness.enabled }}
  {{- if $readinessEnabled }}
  - name: readiness
    host: {{ default .default.healths.readiness.host .Values.healths.readiness.host }}
    port: {{ default .default.healths.readiness.port .Values.healths.readiness.port }}
    mode: ""
    probe_wait_time: "3s"
    http:
      shutdown_duration: "5s"
      handler_timeout: ""
      idle_timeout: ""
      read_header_timeout: ""
      read_timeout: ""
      write_timeout: ""
  {{- end }}
metrics_servers:
  {{- $pprofEnabled := default .default.metrics.pprof.enabled .Values.metrics.pprof.enabled }}
  {{- if $pprofEnabled }}
  - name: pprof
    host: {{ default .default.metrics.pprof.host .Values.metrics.pprof.host }}
    port: {{ default .default.metrics.pprof.port .Values.metrics.pprof.port }}
    mode: REST
    probe_wait_time: 3s
    http:
      shutdown_duration: 5s
      handler_timeout: 5s
      idle_timeout: 2s
      read_header_timeout: 1s
      read_timeout: 1s
      write_timeout: 1s
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
full_shutdown_duration: {{ default .default.full_shutdown_duration .Values.full_shutdown_duration }}
tls:
  enabled: {{ default .default.tls.enabled .Values.tls.enabled }}
  cert: {{ default .default.tls.cert .Values.tls.cert }}
  key: {{ default .default.tls.key .Values.tls.key }}
  ca: {{ default .default.tls.ca .Values.tls.ca }}
{{- end -}}

{{/*
gRPC client configuration
*/}}
{{- define "vald.grpc.client" -}}
addrs: {{ default .default.addrs .Values.addrs }}
health_check_duration: {{ default .default.health_check_duration .Values.health_check_duration }}
backoff:
  {{- if .Values.backoff }}
  initial_duration: {{ default .default.backoff.initial_duration .Values.backoff.initial_duration }}
  backoff_time_limit: {{ default .default.backoff.backoff_time_limit .Values.backoff.backoff_time_limit }}
  maximum_duration: {{ default .default.backoff.maximum_duration .Values.backoff.maximum_duration }}
  jitter_limit: {{ default .default.backoff.jitter_limit .Values.backoff.jitter_limit }}
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
  max_backoff_delay: {{ default .default.dial_option.max_backoff_delay .Values.dial_option.max_backoff_delay }}
  enable_backoff: {{ default .default.dial_option.enable_backoff .Values.dial_option.enable_backoff }}
  insecure: {{ default .default.dial_option.insecure .Values.dial_option.insecure }}
  timeout: {{ default .default.dial_option.timeout .Values.dial_option.timeout }}
  dialer:
    {{- if .Values.dial_option.dialer }}
    dns:
      {{- if .Values.dial_option.dialer.dns }}
      cache_enabled: {{ default .default.dial_option.dialer.dns.cache_enabled .Values.dial_option.dialer.dns.cache_enabled }}
      refresh_duration: {{ default .default.dial_option.dialer.dns.refresh_duration .Values.dial_option.dialer.dns.refresh_duration }}
      cache_expiration: {{ default .default.dial_option.dialer.dns.cache_expiration .Values.dial_option.dialer.dns.cache_expiration }}
      {{- else }}
      {{- toYaml .default.dial_option.dialer.dns | nindent 6 }}
      {{- end }}
    dialer:
      {{- if .Values.dial_option.dialer.dialer }}
      timeout: {{ default .default.dial_option.dialer.dialer.timeout .Values.dial_option.dialer.dialer.timeout }}
      keep_alive: {{ default .default.dial_option.dialer.dialer.keep_alive .Values.dial_option.dialer.dialer.keep_alive }}
      dual_stack_enabled: {{ default .default.dial_option.dialer.dialer.dual_stack_enabled .Values.dial_option.dialer.dialer.dual_stack_enabled }}
      {{- else }}
      {{- toYaml .default.dial_option.dialer.dialer | nindent 6 }}
      {{- end }}
    tls:
      {{- if .Values.dial_option.dialer.tls }}
      enabled: {{ default .default.dial_option.dialer.tls.enabled .Values.dial_option.dialer.tls.enabled }}
      cert: {{ default .default.dial_option.dialer.tls.cert .Values.dial_option.dialer.tls.cert }}
      key: {{ default .default.dial_option.dialer.tls.key .Values.dial_option.dialer.tls.key }}
      ca: {{ default .default.dial_option.dialer.tls.ca .Values.dial_option.dialer.tls.ca }}
      {{- else }}
      {{- toYaml .default.dial_option.dialer.tls | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.dial_option.dialer | nindent 4 }}
    {{- end }}
  keep_alive:
    {{- if .Values.dial_option.keep_alive }}
    time: {{ default .default.dial_option.keep_alive.time .Values.dial_option.keep_alive.time }}
    timeout: {{ default .default.dial_option.keep_alive.timeout .Values.dial_option.keep_alive.timeout }}
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
  cert: {{ default .default.tls.cert .Values.tls.cert }}
  key: {{ default .default.tls.key .Values.tls.key }}
  ca: {{ default .default.tls.ca .Values.tls.ca }}
  {{- else }}
  {{- toYaml .default.tls | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
initContainers
*/}}
{{- define "vald.initContainers" -}}
{{- range .initContainers }}
- name: {{ .name }}
  image: {{ .image }}
  {{- if eq .type "waitFor" }}
  command:
    - /bin/sh
    - -c
    - >
      set -x;
      {{- if eq .target "compressor" }}
      {{- $compressorReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.compressor.server_config.healths.readiness.port }}
      while [ $(curl -sw '%{http_code}' "http://{{ $.Values.compressor.name }}.{{ $.namespace }}.svc.cluster.local:{{ $compressorReadinessPort }}" -o /dev/null) -ne 200]; do
      {{- else if eq .target "meta" }}
      {{- $metaReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.meta.server_config.healths.readiness.port }}
      while [ $(curl -sw '%{http_code}' "http://{{ $.Values.meta.name }}.{{ $.namespace }}.svc.cluster.local:{{ $metaReadinessPort }}" -o /dev/null) -ne 200]; do
      {{- else if eq .target "discoverer" }}
      {{- $discovererReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.discoverer.server_config.healths.readiness.port }}
      while [ $(curl -sw '%{http_code}' "http://{{ $.Values.discoverer.name }}.{{ $.namespace }}.svc.cluster.local:{{ $discovererReadinessPort }}" -o /dev/null) -ne 200]; do
      {{- else if eq .target "agent" }}
      {{- $agentReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.agent.server_config.healths.readiness.port }}
      while [ $(curl -sw '%{http_code}' "http://{{ $.Values.agent.name }}.{{ $.namespace }}.svc.cluster.local:{{ $agentReadinessPort }}" -o /dev/null) -ne 200]; do
      {{- else if eq .target "manager-backup" }}
      {{- $backupManagerReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.backupManager.server_config.healths.readiness.port }}
      while [ $(curl -sw '%{http_code}' "http://{{ $.Values.backupManager.name }}.{{ $.namespace }}.svc.cluster.local:{{ $backupManagerReadinessPort }}" -o /dev/null) -ne 200]; do
      {{- end }}
        sleep {{ .sleepDuration }};
      done
  {{- end }}
{{- end }}
{{- end -}}
