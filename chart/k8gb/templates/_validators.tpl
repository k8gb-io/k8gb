# Validates that "duplicated" values are consistent across the chart

# Validates that the geo tag in extdns.txtPrefix/extdns.txtOwnerId contains the same value as the geo tag in k8gb.clusterGeoTag
{{- define "validateGeoTag" -}}
{{- if .Values.extdns.enabled -}}
{{- if not (contains .Values.k8gb.clusterGeoTag .Values.extdns.txtPrefix) -}}
{{- fail (printf "Validation failed: extdns.txtPrefix (%s) does not contain the expected geo tag (%s)" .Values.extdns.txtPrefix .Values.k8gb.clusterGeoTag) -}}
{{- end -}}
{{- if not (contains .Values.k8gb.clusterGeoTag .Values.extdns.txtOwnerId) -}}
{{- fail (printf "Validation failed: extdns.txtOwnerId (%s) does not contain the expected geo tag (%s)" .Values.extdns.txtOwnerId .Values.k8gb.clusterGeoTag) -}}
{{- end -}}
{{- end -}}
{{- end -}}

# Validates that the zones in k8gb.dnsZones match the zones in extdns.domainFilters
{{- define "validateDnsZones" -}}

{{- if .Values.extdns.enabled }}

{{- $parentZones := list -}}
{{- range .Values.k8gb.dnsZones -}}
  {{- $parentZones = append $parentZones .parentZone -}}
{{- end -}}

{{- $extdnsZones := .Values.extdns.domainFilters -}}

{{- if ne (len $parentZones) (len $extdnsZones) -}}
  {{- fail (printf "Validation failed: Number of zones in k8gb.dnsZones (%d) does not match number of domains in extdns.domainFilters (%d)" (len $parentZones) (len $extdnsZones)) -}}
{{- end -}}

{{- range $parentZone := $parentZones -}}
  {{- if not (has $parentZone $extdnsZones) -}}
    {{- fail (printf "Validation failed: Zone '%s' from k8gb.dnsZones is not present in extdns.domainFilters" $parentZone) -}}
  {{- end -}}
{{- end -}}

{{- range $extdnsZone := $extdnsZones -}}
  {{- if not (has $extdnsZone $parentZones) -}}
    {{- fail (printf "Validation failed: Domain '%s' from extdns.domainFilters is not present in k8gb.dnsZones" $extdnsZone) -}}
  {{- end -}}
{{- end -}}

 {{- end }}
{{- end -}}
