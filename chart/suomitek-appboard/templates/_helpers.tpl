{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "suomitek-appboard.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "suomitek-appboard.fullname" -}}
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
{{- define "suomitek-appboard.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels for additional suomitek-appboard applications. Used on resources whose app name is different
from kubeapps
*/}}
{{- define "suomitek-appboard.extraAppLabels" -}}
chart: {{ include "suomitek-appboard.chart" . }}
release: {{ .Release.Name }}
heritage: {{ .Release.Service }}
helm.sh/chart: {{ template "suomitek-appboard.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/name: {{ include "suomitek-appboard.name" . }}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "suomitek-appboard.labels" -}}
app: {{ include "suomitek-appboard.name" . }}
{{ template "suomitek-appboard.extraAppLabels" . }}
{{- end -}}

{{/*
Render image reference
*/}}
{{- define "suomitek-appboard.image" -}}
{{- $image := index . 0 -}}
{{- $global := index . 1 -}}
{{/*
Helm 2.11 supports the assignment of a value to a variable defined in a different scope,
but Helm 2.9 and 2.10 doesn't support it, so we need to implement this if-else logic.
Also, we can't use a single if because lazy evaluation is not an option
*/}}
{{- if $global -}}
    {{- if $global.imageRegistry -}}
        {{ $global.imageRegistry }}/{{ $image.repository }}:{{ $image.tag }}
    {{- else -}}
        {{ $image.registry }}/{{ $image.repository }}:{{ $image.tag }}
    {{- end -}}
{{- else -}}
    {{ $image.registry }}/{{ $image.repository }}:{{ $image.tag }}
{{- end -}}
{{- end -}}

{{/*
Create a default fully qualified app name for MongoDB dependency.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "suomitek-appboard.mongodb.fullname" -}}
{{- $name := default "mongodb" .Values.mongodb.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/*
Create a default fully qualified app name for PostgreSQL dependency.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "suomitek-appboard.postgresql.fullname" -}}
{{- $name := default "postgresql" .Values.postgresql.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create name for the apprepository-controller based on the fullname
*/}}
{{- define "suomitek-appboard.apprepository.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-apprepository-controller
{{- end -}}

{{/*
Create name for the apprepository pre-upgrade job
*/}}
{{- define "suomitek-appboard.apprepository-job-postupgrade.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-apprepository-job-postupgrade
{{- end -}}

{{/*
Create name for the apprepository cleanup job
*/}}
{{- define "suomitek-appboard.apprepository-jobs-cleanup.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-apprepository-jobs-cleanup
{{- end -}}

{{/*
Create name for the db-secret secret bootstrap job
*/}}
{{- define "suomitek-appboard.db-secret-jobs-cleanup.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-db-secret-jobs-cleanup
{{- end -}}

{{/*
Create name for the kubeapps upgrade job
*/}}
{{- define "suomitek-appboard.suomitek-appboard-jobs-upgrade.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-suomitek-appboard-jobs-upgrade
{{- end -}}

{{/*
Create name for the assetsvc based on the fullname
*/}}
{{- define "suomitek-appboard.assetsvc.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-assetsvc
{{- end -}}

{{/*
Create name for the dashboard based on the fullname
*/}}
{{- define "suomitek-appboard.dashboard.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-dashboard
{{- end -}}

{{/*
Create name for the dashboard config based on the fullname
*/}}
{{- define "suomitek-appboard.dashboard-config.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-dashboard-config
{{- end -}}

{{/*
Create name for the frontend config based on the fullname
*/}}
{{- define "suomitek-appboard.frontend-config.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-frontend-config
{{- end -}}

{{/*
Create proxy_pass for the frontend config based on the useHelm3 flag
*/}}
{{- define "suomitek-appboard.frontend-config.proxy_pass" -}}
{{- if .Values.useHelm3 -}}
http://{{ template "suomitek-appboard.kubeops.fullname" . }}:{{ .Values.kubeops.service.port }}
{{- else -}}
http://{{ template "suomitek-appboard.tiller-proxy.fullname" . }}:{{ .Values.tillerProxy.service.port }}
{{- end -}}
{{- end -}}

{{/*
Create name for the tiller-proxy based on the fullname
*/}}
{{- define "suomitek-appboard.tiller-proxy.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-tiller-proxy
{{- end -}}

{{/*
Create name for kubeops based on the fullname
*/}}
{{- define "suomitek-appboard.kubeops.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-internal-kubeops
{{- end -}}

{{/*
Create name for the kubeops config based on the fullname
*/}}
{{- define "suomitek-appboard.kubeops-config.fullname" -}}
{{ template "suomitek-appboard.fullname" . }}-kubeops-config
{{- end -}}

{{/*
Create name for the secrets related to an app repository
*/}}
{{- define "suomitek-appboard.apprepository-secret.name" -}}
apprepo-{{ .name }}-secrets
{{- end -}}

{{/*
Repositories that include a caCert or an authorizationHeader
*/}}
{{- define "suomitek-appboard.repos-with-orphan-secrets" -}}
{{- range .Values.apprepository.initialRepos }}
{{- if or .caCert .authorizationHeader }}
.name
{{- end }}
{{- end }}
{{- end -}}

{{/*
Frontend service port number
*/}}
{{- define "suomitek-appboard.frontend-port-number" -}}
{{- if .Values.authProxy.enabled -}}
3000
{{- else -}}
8080
{{- end -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "suomitek-appboard.imagePullSecrets" -}}
{{/*
We can not use a single if because lazy evaluation is not an option
*/}}
{{- if .Values.global }}
{{- if .Values.global.imagePullSecrets }}
imagePullSecrets:
{{- range .Values.global.imagePullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Renders a value that contains template.
Usage:
{{ include "suomitek-appboard.tplValue" ( dict "value" .Values.path.to.the.Value "context" $) }}
*/}}
{{- define "suomitek-appboard.tplValue" -}}
    {{- if typeIs "string" .value }}
        {{- tpl .value .context }}
    {{- else }}
        {{- tpl (.value | toYaml) .context }}
    {{- end }}
{{- end -}}
