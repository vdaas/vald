{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "vald.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "vald.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels: keep the historical vald-operator label set and merge the
user-supplied additional labels on top.
*/}}
{{- define "vald-operator.labels" -}}
app: {{ .Values.name }}
app.kubernetes.io/name: vald
app.kubernetes.io/component: vald-operator
{{- with .Values.labels }}
{{ toYaml . }}
{{- end }}
{{- end -}}

{{/*
Create the name of the service account to use.
*/}}
{{- define "vald.serviceAccountName" -}}
{{- default .Values.name .Values.serviceAccount.name -}}
{{- end -}}

{{/*
Create the name of roles and rolebindings.
*/}}
{{- define "vald-operator.rbacName" -}}
{{- default .Values.name .Values.rbac.name -}}
{{- end -}}

{{/*
Resolve the container image reference.
*/}}
{{- define "vald-operator.image" -}}
{{- printf "%s:%s" .Values.image.repository (default .Chart.AppVersion .Values.image.tag) -}}
{{- end -}}

{{/*
Resolve the name of the ConfigMap that provides the default VRS template.
*/}}
{{- define "vald-operator.defaultVrsConfigMapName" -}}
{{- default (printf "%s-default-vrs" .Values.name) .Values.default_vrs.existingConfigMap -}}
{{- end -}}

{{/*
server_config rendered into config.yaml. Ports and enabled flags are shared
with the deployment template so that changing a server port also updates the
matching containerPort, probe and service port.
*/}}
{{- define "vald-operator.server_config" -}}
{{- $grpc := .Values.server_config.servers.grpc -}}
{{- $liveness := .Values.server_config.healths.liveness -}}
{{- $readiness := .Values.server_config.healths.readiness -}}
{{- $pprof := .Values.server_config.metrics.pprof -}}
servers:
  {{- if $grpc.enabled }}
  - name: grpc
    host: {{ $grpc.host }}
    port: {{ $grpc.port }}
    grpc:
      {{- toYaml $grpc.server.grpc | nindent 6 }}
    mode: {{ $grpc.server.mode }}
    network: {{ $grpc.server.network }}
    probe_wait_time: {{ $grpc.server.probe_wait_time }}
    restart: {{ $grpc.server.restart }}
    socket_option:
      {{- toYaml $grpc.server.socket_option | nindent 6 }}
    socket_path: {{ $grpc.server.socket_path | quote }}
  {{- end }}
health_check_servers:
  {{- if $liveness.enabled }}
  - name: liveness
    host: {{ $liveness.host }}
    port: {{ $liveness.port }}
    http:
      {{- toYaml $liveness.server.http | nindent 6 }}
    mode: {{ $liveness.server.mode | quote }}
    network: {{ $liveness.server.network }}
    probe_wait_time: {{ $liveness.server.probe_wait_time }}
    socket_option:
      {{- toYaml $liveness.server.socket_option | nindent 6 }}
    socket_path: {{ $liveness.server.socket_path | quote }}
  {{- end }}
  {{- if $readiness.enabled }}
  - name: readiness
    host: {{ $readiness.host }}
    port: {{ $readiness.port }}
    http:
      {{- toYaml $readiness.server.http | nindent 6 }}
    mode: {{ $readiness.server.mode | quote }}
    network: {{ $readiness.server.network }}
    probe_wait_time: {{ $readiness.server.probe_wait_time }}
    socket_option:
      {{- toYaml $readiness.server.socket_option | nindent 6 }}
    socket_path: {{ $readiness.server.socket_path | quote }}
  {{- end }}
metrics_servers:
  {{- if $pprof.enabled }}
  - name: pprof
    host: {{ $pprof.host }}
    port: {{ $pprof.port }}
    http:
      {{- toYaml $pprof.server.http | nindent 6 }}
    mode: {{ $pprof.server.mode }}
    network: {{ $pprof.server.network }}
    probe_wait_time: {{ $pprof.server.probe_wait_time }}
    restart: {{ $pprof.server.restart }}
    socket_option:
      {{- toYaml $pprof.server.socket_option | nindent 6 }}
    socket_path: {{ $pprof.server.socket_path | quote }}
  {{- end }}
startup_strategy:
  {{- if $liveness.enabled }}
  - liveness
  {{- end }}
  {{- if $pprof.enabled }}
  - pprof
  {{- end }}
  {{- if $grpc.enabled }}
  - grpc
  {{- end }}
  {{- if $readiness.enabled }}
  - readiness
  {{- end }}
shutdown_strategy:
  {{- if $readiness.enabled }}
  - readiness
  {{- end }}
  {{- if $grpc.enabled }}
  - grpc
  {{- end }}
  {{- if $pprof.enabled }}
  - pprof
  {{- end }}
  {{- if $liveness.enabled }}
  - liveness
  {{- end }}
full_shutdown_duration: {{ .Values.server_config.full_shutdown_duration }}
tls:
  {{- toYaml .Values.server_config.tls | nindent 2 }}
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
