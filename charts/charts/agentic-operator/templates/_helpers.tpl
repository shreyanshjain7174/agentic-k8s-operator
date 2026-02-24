{{- define "agentic-operator-sub.fullname" -}}
{{- printf "%s-agentic-operator" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "agentic-operator-sub.labels" -}}
app.kubernetes.io/name: agentic-operator
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: operator
app.kubernetes.io/part-of: agentic-operator
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}

{{- define "agentic-operator-sub.selectorLabels" -}}
app.kubernetes.io/name: agentic-operator
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
