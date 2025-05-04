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

# Validates that the zones in k8gb.edgeDNSZone/k8gb.dnsZones match the zones in extdns.domainFilters
{{- define "validateDnsZones" -}}
{{- $k8gbZones := list -}}
{{- range .Values.k8gb.dnsZones -}}
  {{- $k8gbZones = append $k8gbZones .zone -}}
{{- end -}}
{{- if and (or (not .Values.k8gb.dnsZones) (eq (len .Values.k8gb.dnsZones) 0)) .Values.k8gb.dnsZone .Values.k8gb.edgeDNSZone }}
  {{- $k8gbZones = append $k8gbZones .Values.k8gb.edgeDNSZone -}}
{{- end }}

{{- $extdnsZones := .Values.extdns.domainFilters -}}

{{- if ne (len $k8gbZones) (len $extdnsZones) -}}
  {{- fail (printf "Validation failed: Number of zones in k8gb.edgeDNSZone/k8gb.dnsZones (%d) does not match number of domains in extdns.domainFilters (%d)" (len $k8gbZones) (len $extdnsZones)) -}}
{{- end -}}

{{- range $k8gbZone := $k8gbZones -}}
  {{- if not (has $k8gbZone $extdnsZones) -}}
    {{- fail (printf "Validation failed: Zone '%s' from k8gb.edgeDNSZone/k8gb.dnsZones is not present in extdns.domainFilters" $k8gbZone) -}}
  {{- end -}}
{{- end -}}

{{- range $extdnsZone := $extdnsZones -}}
  {{- if not (has $extdnsZone $k8gbZones) -}}
    {{- fail (printf "Validation failed: Domain '%s' from extdns.domainFilters is not present in k8gb.edgeDNSZone/k8gb.dnsZones" $extdnsZone) -}}
  {{- end -}}
{{- end -}}

{{- end -}}
