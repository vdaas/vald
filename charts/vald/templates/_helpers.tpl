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
{{- $startupEnabled := .default.healths.startup.enabled }}
{{- if hasKey .Values.healths.startup "enabled" }}
{{- $startupEnabled = .Values.healths.startup.enabled }}
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
  {{- if not $startupEnabled }}
  initialDelaySeconds: {{ default .default.healths.liveness.livenessProbe.initialDelaySeconds .Values.healths.liveness.livenessProbe.initialDelaySeconds }}
  {{- end }}
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
{{- if $startupEnabled }}
startupProbe:
  {{- if .Values.healths.startup.startupProbe }}
  httpGet:
    {{- if .Values.healths.startup.startupProbe.httpGet }}
    path: {{ default .default.healths.startup.startupProbe.httpGet.path .Values.healths.startup.startupProbe.httpGet.path }}
    port: {{ default .default.healths.startup.startupProbe.httpGet.port .Values.healths.startup.startupProbe.httpGet.port }}
    scheme: {{ default .default.healths.startup.startupProbe.httpGet.scheme .Values.healths.startup.startupProbe.httpGet.scheme }}
    {{- else }}
    {{- toYaml .default.healths.startup.startupProbe.httpGet | nindent 4 }}
    {{- end }}
  initialDelaySeconds: {{ default .default.healths.startup.startupProbe.initialDelaySeconds .Values.healths.startup.startupProbe.initialDelaySeconds }}
  timeoutSeconds: {{ default .default.healths.startup.startupProbe.timeoutSeconds .Values.healths.startup.startupProbe.timeoutSeconds }}
  successThreshold: {{ default .default.healths.startup.startupProbe.successThreshold .Values.healths.startup.startupProbe.successThreshold }}
  failureThreshold: {{ default .default.healths.startup.startupProbe.failureThreshold .Values.healths.startup.startupProbe.failureThreshold }}
  periodSeconds: {{ default .default.healths.startup.startupProbe.periodSeconds .Values.healths.startup.startupProbe.periodSeconds }}
  {{- else }}
  {{- toYaml .default.healths.startup.startupProbe | nindent 2 }}
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
Ingress defaultBackend
*/}
{{- define "vald.ingressDefaultBackend" -}}
{{- $defaultBackend := .default }}
{{- with .Values }}
{{- $defaultBackend = . }}
{{- end }}
{{- if $defaultBackend -}}
defaultBackend:
  {{- if $defaultBackend.resource }}
  resource:
    {{- toYaml $defaultBackend.resource| nindent 4}}
  {{- else }}
  service:
    name: {{ $defaultBackend.service.name }}
    port:
      {{- if $defaultBackend.service.port.number }}
      number: {{ $defaultBackend.service.port.number }}
      {{- else }}
      name: {{ $defaultBackend.service.port.name }}
      {{- end -}}
  {{- end -}}
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
    network: {{ default .default.servers.rest.server.network .Values.servers.rest.server.network | quote }}
    socket_path: {{ default .default.servers.rest.server.socket_path .Values.servers.rest.server.socket_path | quote }}
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
    network: {{ default .default.servers.grpc.server.network .Values.servers.grpc.server.network | quote }}
    socket_path: {{ default .default.servers.grpc.server.socket_path .Values.servers.grpc.server.socket_path | quote }}
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
        min_time: {{ default .default.servers.grpc.server.grpc.keepalive.min_time .Values.servers.grpc.server.grpc.keepalive.min_time | quote }}
        permit_without_stream: {{ default .default.servers.grpc.server.grpc.keepalive.permit_without_stream .Values.servers.grpc.server.grpc.keepalive.permit_without_stream | quote }}
        {{- else }}
        {{- toYaml .default.servers.grpc.server.grpc.keepalive | nindent 8 }}
        {{- end }}
      write_buffer_size: {{ default .default.servers.grpc.server.grpc.write_buffer_size .Values.servers.grpc.server.grpc.write_buffer_size }}
      read_buffer_size: {{ default .default.servers.grpc.server.grpc.read_buffer_size .Values.servers.grpc.server.grpc.read_buffer_size }}
      connection_timeout: {{ default .default.servers.grpc.server.grpc.connection_timeout .Values.servers.grpc.server.grpc.connection_timeout | quote }}
      max_header_list_size: {{ default .default.servers.grpc.server.grpc.max_header_list_size .Values.servers.grpc.server.grpc.max_header_list_size }}
      header_table_size: {{ default .default.servers.grpc.server.grpc.header_table_size .Values.servers.grpc.server.grpc.header_table_size }}
      {{- if .Values.servers.grpc.server.grpc.interceptors }}
      interceptors:
        {{- toYaml .Values.servers.grpc.server.grpc.interceptors | nindent 8 }}
      {{- else if .default.servers.grpc.server.grpc.interceptors }}
      interceptors:
        {{- toYaml .default.servers.grpc.server.grpc.interceptors | nindent 8 }}
      {{- else }}
      interceptors: []
      {{- end }}
      enable_reflection: {{ default .default.servers.grpc.server.grpc.enable_reflection .Values.servers.grpc.server.grpc.enable_reflection }}
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
    network: {{ default .default.healths.liveness.server.network .Values.healths.liveness.server.network | quote }}
    socket_path: {{ default .default.healths.liveness.server.socket_path .Values.healths.liveness.server.socket_path | quote }}
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
    network: {{ default .default.healths.readiness.server.network .Values.healths.readiness.server.network | quote }}
    socket_path: {{ default .default.healths.readiness.server.socket_path .Values.healths.readiness.server.socket_path | quote }}
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
    network: {{ default .default.metrics.pprof.server.network .Values.metrics.pprof.server.network | quote }}
    socket_path: {{ default .default.metrics.pprof.server.socket_path .Values.metrics.pprof.server.socket_path | quote }}
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
  {{- if .Values.tls }}
  enabled: {{ default .default.tls.enabled .Values.tls.enabled }}
  cert: {{ default .default.tls.cert .Values.tls.cert | quote }}
  key: {{ default .default.tls.key .Values.tls.key | quote }}
  ca: {{ default .default.tls.ca .Values.tls.ca | quote }}
  insecure_skip_verify: {{ default .default.tls.insecure_skip_verify .Values.tls.insecure_skip_verify }}
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
circuit_breaker:
  {{- if .Values.circuit_breaker }}
  closed_error_rate: {{ default .default.circuit_breaker.closed_error_rate .Values.circuit_breaker.closed_error_rate }}
  half_open_error_rate: {{ default .default.circuit_breaker.half_open_error_rate .Values.circuit_breaker.half_open_error_rate }}
  min_samples: {{ default .default.circuit_breaker.min_samples .Values.circuit_breaker.min_samples | quote }}
  open_timeout: {{ default .default.circuit_breaker.open_timeout .Values.circuit_breaker.open_timeout | quote }}
  closed_refresh_timeout: {{ default .default.circuit_breaker.closed_refresh_timeout .Values.circuit_breaker.closed_refresh_timeout | quote }}
  {{- else }}
  {{- toYaml .default.circuit_breaker | nindent 2 }}
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
  {{- if .Values.dial_option.interceptors }}
  interceptors:
    {{- toYaml .Values.dial_option.interceptors | nindent 4 }}
  {{- else if .default.dial_option.interceptors }}
  interceptors:
    {{- toYaml .default.dial_option.interceptors | nindent 4 }}
  {{- else }}
  interceptors: []
  {{- end }}
  net:
    {{- if .Values.dial_option.net }}
    dns:
      {{- if .Values.dial_option.net.dns }}
      cache_enabled: {{ default .default.dial_option.net.dns.cache_enabled .Values.dial_option.net.dns.cache_enabled }}
      refresh_duration: {{ default .default.dial_option.net.dns.refresh_duration .Values.dial_option.net.dns.refresh_duration | quote }}
      cache_expiration: {{ default .default.dial_option.net.dns.cache_expiration .Values.dial_option.net.dns.cache_expiration | quote }}
      {{- else }}
      {{- toYaml .default.dial_option.net.dns | nindent 6 }}
      {{- end }}
    dialer:
      {{- if .Values.dial_option.net.dialer }}
      timeout: {{ default .default.dial_option.net.dialer.timeout .Values.dial_option.net.dialer.timeout | quote }}
      keepalive: {{ default .default.dial_option.net.dialer.keepalive .Values.dial_option.net.dialer.keepalive | quote }}
      dual_stack_enabled: {{ default .default.dial_option.net.dialer.dual_stack_enabled .Values.dial_option.net.dialer.dual_stack_enabled }}
      {{- else }}
      {{- toYaml .default.dial_option.net.dialer | nindent 6 }}
      {{- end }}
    tls:
      {{- if .Values.dial_option.net.tls }}
      enabled: {{ default .default.dial_option.net.tls.enabled .Values.dial_option.net.tls.enabled }}
      cert: {{ default .default.dial_option.net.tls.cert .Values.dial_option.net.tls.cert | quote }}
      key: {{ default .default.dial_option.net.tls.key .Values.dial_option.net.tls.key | quote }}
      ca: {{ default .default.dial_option.net.tls.ca .Values.dial_option.net.tls.ca | quote }}
      insecure_skip_verify: {{ default .default.dial_option.net.tls.insecure_skip_verify .Values.dial_option.net.tls.insecure_skip_verify }}
      {{- else }}
      {{- toYaml .default.dial_option.net.tls | nindent 6 }}
      {{- end }}
    socket_option:
      {{- if .Values.dial_option.net.socket_option }}
      reuse_port: {{ default .default.dial_option.net.socket_option.reuse_port .Values.dial_option.net.socket_option.reuse_port }}
      reuse_addr: {{ default .default.dial_option.net.socket_option.reuse_addr .Values.dial_option.net.socket_option.reuse_addr }}
      tcp_fast_open: {{ default .default.dial_option.net.socket_option.tcp_fast_open .Values.dial_option.net.socket_option.tcp_fast_open }}
      tcp_no_delay: {{ default .default.dial_option.net.socket_option.tcp_no_delay .Values.dial_option.net.socket_option.tcp_no_delay }}
      tcp_cork: {{ default .default.dial_option.net.socket_option.tcp_cork .Values.dial_option.net.socket_option.tcp_cork }}
      tcp_quick_ack: {{ default .default.dial_option.net.socket_option.tcp_quick_ack .Values.dial_option.net.socket_option.tcp_quick_ack }}
      tcp_defer_accept: {{ default .default.dial_option.net.socket_option.tcp_defer_accept .Values.dial_option.net.socket_option.tcp_defer_accept }}
      ip_transparent: {{ default .default.dial_option.net.socket_option.ip_transparent .Values.dial_option.net.socket_option.ip_transparent }}
      ip_recover_destination_addr: {{ default .default.dial_option.net.socket_option.ip_recover_destination_addr .Values.dial_option.net.socket_option.ip_recover_destination_addr }}
      {{- else }}
      {{- toYaml .default.dial_option.net.socket_option | nindent 6 }}
      {{- end }}
    {{- else }}
    {{- toYaml .default.dial_option.net | nindent 4 }}
    {{- end }}
  keepalive:
    {{- if .Values.dial_option.keepalive }}
    time: {{ default .default.dial_option.keepalive.time .Values.dial_option.keepalive.time | quote }}
    timeout: {{ default .default.dial_option.keepalive.timeout .Values.dial_option.keepalive.timeout | quote }}
    permit_without_stream: {{ default .default.dial_option.keepalive.permit_without_stream .Values.dial_option.keepalive.permit_without_stream }}
    {{- else }}
    {{- toYaml .default.dial_option.keepalive | nindent 4 }}
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
  insecure_skip_verify: {{ default .default.tls.insecure_skip_verify .Values.tls.insecure_skip_verify }}
  {{- else }}
  {{- toYaml .default.tls | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
observability
*/}}
{{- define "vald.observability" -}}
enabled: {{ default .default.enabled .Values.enabled }}
otlp:
  {{- if .Values.otlp }}
  collector_endpoint: {{ default .default.otlp.collector_endpoint .Values.otlp.collector_endpoint | quote }}
  trace_batch_timeout: {{ default .default.otlp.trace_batch_timeout .Values.otlp.trace_batch_timeout | quote }}
  trace_export_timeout: {{ default .default.otlp.trace_export_timeout .Values.otlp.trace_export_timeout | quote }}
  trace_max_export_batch_size: {{ default .default.otlp.trace_max_export_batch_size .Values.otlp.trace_max_export_batch_size }}
  trace_max_queue_size: {{ default .default.otlp.trace_max_queue_size .Values.otlp.trace_max_queue_size }}
  metrics_export_interval: {{ default .default.otlp.metrics_export_interval .Values.otlp.metrics_export_interval | quote }}
  metrics_export_timeout: {{ default .default.otlp.metrics_export_timeout .Values.otlp.metrics_export_timeout | quote }}
  attribute:
    {{- if .Values.otlp.attribute }}
    namespace: {{ default .default.otlp.attribute.namespace .Values.otlp.attribute.namespace | quote }}
    pod_name: {{ default .default.otlp.attribute.pod_name .Values.otlp.attribute.pod_name | quote }}
    node_name: {{ default .default.otlp.attribute.node_name .Values.otlp.attribute.node_name | quote }}
    service_name: {{ default .default.otlp.attribute.service_name .Values.otlp.attribute.service_name | quote }}
    {{- else }}
    {{- toYaml .default.otlp.attribute | nindent 4 }}
    {{- end }}
  {{- else }}
  {{- toYaml .default.otlp | nindent 2 }}
  {{- end }}
metrics:
  {{- if .Values.metrics }}
  enable_version_info: {{ default .default.metrics.enable_version_info .Values.metrics.enable_version_info }}
  {{- if .Values.metrics.version_info_labels }}
  version_info_labels:
    {{- toYaml .Values.metrics.version_info_labels | nindent 4 }}
  {{- else if .default.metrics.version_info_labels }}
  version_info_labels:
    {{- toYaml .default.metrics.version_info_labels | nindent 4 }}
  {{- else }}
  version_info_labels: []
  {{- end }}
  enable_memory: {{ default .default.metrics.enable_memory .Values.metrics.enable_memory }}
  enable_goroutine: {{ default .default.metrics.enable_goroutine .Values.metrics.enable_goroutine }}
  enable_cgo: {{ default .default.metrics.enable_cgo .Values.metrics.enable_cgo }}
  {{- else }}
  {{- toYaml .default.metrics | nindent 2 }}
  {{- end }}
trace:
  {{- if .Values.trace }}
  enabled: {{ default .default.trace.enabled .Values.trace.enabled }}
  {{- else }}
  {{- toYaml .default.trace | nindent 2 }}
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
      {{- if eq .target "discoverer" }}
      {{- $discovererReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.discoverer.server_config.healths.readiness.port }}
      {{- $discovererReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.discoverer.name }}.{{ $.namespace }}.svc.cluster.local:{{ $discovererReadinessPort }}{{ $discovererReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "agent" }}
      {{- $agentReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.agent.server_config.healths.readiness.port }}
      {{- $agentReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.agent.name }}.{{ $.namespace }}.svc.cluster.local:{{ $agentReadinessPort }}{{ $agentReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if eq .target "gateway-lb" }}
      {{- $lbGatewayReadinessPort := default $.Values.defaults.server_config.healths.readiness.port $.Values.gateway.lb.server_config.healths.readiness.port }}
      {{- $lbGatewayReadinessPath := default $.Values.defaults.server_config.healths.readiness.readinessProbe.httpGet.path .readinessPath }}
      until [ "$(wget --server-response --spider --quiet http://{{ $.Values.gateway.lb.name }}.{{ $.namespace }}.svc.cluster.local:{{ $lbGatewayReadinessPort }}{{ $lbGatewayReadinessPath }} 2>&1 | awk 'NR==1{print $2}')" == "200" ]; do
      {{- else if .untilCondition }}
      until [ {{ .untilCondition }} ]; do
      {{- else if .whileCondition }}
      while [ {{ .whileCondition }} ]; do
      {{- end }}
        echo "waiting for {{ .target }} to be ready..."
        sleep {{ .sleepDuration }};
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

{{/*
Discoverer selector config
*/}}
{{- define "vald.discoverer.selector" -}}
pod:
  labels:
    app.kubernetes.io/component: agent
    app.kubernetes.io/instance: {{ .chart.Release.Name }}
    {{- if .Values.pod.labels }}
      {{- toYaml .Values.pod.labels | nindent 4 }}
    {{- end }}
  fields:
    status.phase: Running
    {{- if .Values.pod.fields }}
      {{- toYaml .Values.pod.fields | nindent 4 }}
    {{- end }}
node:
  {{- if .Values.node.labels }}
  labels:
    {{- toYaml .Values.node.labels | nindent 4 }}
  {{- else }}
  labels: {}
  {{- end }}
  {{- if .Values.node.fields }}
  fields:
    {{- toYaml .Values.node.fields | nindent 4 }}
  {{- else }}
  fields: {}
  {{- end }}
pod_metrics:
  {{- if .Values.pod_metrics.labels }}
  labels:
    {{- toYaml .Values.pod_metrics.labels | nindent 4 }}
  {{- else }}
  labels: {}
  {{- end }}
  fields:
    containers.name: {{ .agent.name }}
    {{- if .Values.pod_metrics.fields }}
      {{- toYaml .Values.pod_metrics.fields | nindent 4 }}
    {{- end }}
node_metrics:
  {{- if .Values.node_metrics.labels }}
  labels:
    {{- toYaml .Values.node_metrics.labels | nindent 4 }}
  {{- else }}
  labels: {}
  {{- end }}
  {{- if .Values.node_metrics.fields }}
  fields:
    {{- toYaml .Values.node_metrics.fields | nindent 4 }}
  {{- else }}
  fields: {}
  {{- end }}
{{- end -}}
