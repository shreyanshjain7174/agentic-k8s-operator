{{- define "litellm.fullname" -}}
{{- printf "%s-litellm" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "litellm.labels" -}}
app.kubernetes.io/name: litellm
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: llm-proxy
app.kubernetes.io/part-of: agentic-operator
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}

{{- define "litellm.selectorLabels" -}}
app.kubernetes.io/name: litellm
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
