{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "pro-carousel.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "pro-carousel.fullname" -}}
{{- if .Values.dontUseReleaseName -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- else if .Values.fullnameOverride -}}
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
{{- define "pro-carousel.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "pro-carousel.labels" -}}
app.kubernetes.io/name: {{ include "pro-carousel.name" . }}
helm.sh/chart: {{ include "pro-carousel.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}


{{/*
Generate the host prefix using the env global var
*/}}
{{- define "pro-carousel.hostPrefix" -}}
{{- $name := include "pro-carousel.name" . -}}
{{- if eq "reg" .Values.globals.env -}}
{{- printf "%s.%s.%s" $name .Release.Namespace .Values.globals.env -}}
{{- else -}}
{{- printf "%s.%s" $name .Values.globals.env -}}
{{- end -}}
{{- end -}}
