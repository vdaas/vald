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
  {{- if .rest.enabled }}
  - name: {{ .prefix }}-rest
    host: {{ .rest.host | default "0.0.0.0" }}
    port: {{ .rest.port | default 8080 }}
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
  {{- if .grpc.enabled }}
  - name: {{ .prefix }}-grpc
    host: {{ .grpc.host | default "0.0.0.0" }}
    port: {{ .grpc.port | default 8082 }}
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
  {{- if .liveness.enabled }}
  - name: {{ .prefix }}-liveness
    host: {{ .liveness.host | default "0.0.0.0" }}
    port: {{ .liveness.port | default 3000 }}
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
  {{- if .readiness.enabled }}
  - name: {{ .prefix }}-readiness
    host: {{ .readiness.host | default "0.0.0.0" }}
    port: {{ .readiness.port | default 3001 }}
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
  {{- if .metrics.enabled }}
  - name: {{ .prefix }}-pprof
    host: {{ .metrics.host | default "0.0.0.0" }}
    port: {{ .metrics.port | default 6060 }}
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
  {{- if .liveness.enabled }}
  - {{ .prefix }}-liveness
  {{- end }}
  {{- if .metrics.enabled }}
  - {{ .prefix }}-pprof
  {{- end }}
  {{- if .grpc.enabled }}
  - {{ .prefix }}-grpc
  {{- end }}
  {{- if .rest.enabled }}
  - {{ .prefix }}-rest
  {{- end }}
  {{- if .readiness.enabled }}
  - {{ .prefix }}-readiness
  {{- end }}
full_shutdown_duration: {{ .full_shutdown_duration | default "600s" }}
tls:
  enabled: {{ .tls.enabled | default false }}
  cert: /path/to/cert
  key: /path/to/key
  ca: /path/to/ca
{{- end -}}
