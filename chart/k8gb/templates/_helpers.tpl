{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "chart.fullname" -}}
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
{{- define "chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "chart.labels" -}}
helm.sh/chart: {{ include "chart.chart" . }}
{{ include "chart.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "chart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "chart.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "chart.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{- define "k8gb.extdnsOwnerID" -}}
k8gb-{{ index (split ":" (index (split ";" (include "k8gb.dnsZonesString" .)) "_0")) "_1" }}-{{ .Values.k8gb.clusterGeoTag }}
{{- end -}}

{{- define "k8gb.edgeDNSServers" -}}
{{- if .Values.k8gb.edgeDNSServer -}}
{{ .Values.k8gb.edgeDNSServer }}
{{- else -}}
{{ join "," .Values.k8gb.edgeDNSServers }}
{{- end -}}
{{- end -}}

{{- define "k8gb.dnsZonesString" -}}
{{- $entries := list -}}
{{- range .Values.k8gb.dnsZones }}
  {{- $dnsZoneNegTTL := toString (.dnsZoneNegTTL | default "300") }}
  {{- $entry := printf "%s:%s:%s" .zone .domain $dnsZoneNegTTL }}
  {{- $entries = append $entries $entry }}
{{- end }}
{{- if and (or (not .Values.k8gb.dnsZones) (eq (len .Values.k8gb.dnsZones) 0)) .Values.k8gb.dnsZone .Values.k8gb.edgeDNSZone }}
  {{- $extraEntry := printf "%s:%s:%s" .Values.k8gb.edgeDNSZone .Values.k8gb.dnsZone "300" }}
  {{- $entries = append $entries $extraEntry }}
{{- end }}
{{- join ";" $entries }}
{{- end }}


{{- define "k8gb.coredns.extraPlugins" -}}
{{- if .Values.k8gb.coredns.extra_plugins }}
{{- range .Values.k8gb.coredns.extra_plugins }}
        {{ . }}
{{- end }}
{{- end }}
{{- end }}

{{- define "k8gb.metrics_port" -}}
{{ print (split ":" .Values.k8gb.metricsAddress)._1 }}
{{- end -}}
