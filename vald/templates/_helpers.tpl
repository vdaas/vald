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
Server configures that inserted into server_config
*/}}
{{- define "vald.servers" -}}
servers:
  {{- $restEnabled := default .default.servers.rest.enabled .Values.servers.rest.enabled }}
  {{- if $restEnabled }}
  - name: {{ .Values.prefix }}-rest
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
  - name: {{ .Values.prefix }}-grpc
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
        max_idle_conn_idle: O
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
  - name: {{ .Values.prefix }}-liveness
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
  - name: {{ .Values.prefix }}-readiness
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
  - name: {{ .Values.prefix }}-pprof
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
  - {{ .Values.prefix }}-liveness
  {{- end }}
  {{- if $pprofEnabled }}
  - {{ .Values.prefix }}-pprof
  {{- end }}
  {{- if $grpcEnabled }}
  - {{ .Values.prefix }}-grpc
  {{- end }}
  {{- if $restEnabled }}
  - {{ .Values.prefix }}-rest
  {{- end }}
  {{- if $readinessEnabled }}
  - {{ .Values.prefix }}-readiness
  {{- end }}
full_shutdown_duration: {{ default .default.full_shutdown_duration .Values.full_shutdown_duration }}
tls:
  enabled: {{ default .default.tls.enabled .Values.tls.enabled }}
  cert: {{ default .default.tls.cert .Values.tls.cert }}
  key: {{ default .default.tls.key .Values.tls.key }}
  ca: {{ default .default.tls.ca .Values.tls.ca }}
{{- end -}}
