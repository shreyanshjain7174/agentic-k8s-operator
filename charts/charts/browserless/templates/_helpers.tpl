{{- define "browserless.fullname" -}}
{{- printf "%s-browserless" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "browserless.labels" -}}
app.kubernetes.io/name: browserless
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: cdp-pool
app.kubernetes.io/part-of: agentic-operator
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}

{{- define "browserless.selectorLabels" -}}
app.kubernetes.io/name: browserless
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
