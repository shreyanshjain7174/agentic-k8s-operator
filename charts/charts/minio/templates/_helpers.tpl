{{- define "minio.fullname" -}}
{{- printf "%s-minio" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "minio.labels" -}}
app.kubernetes.io/name: minio
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: object-storage
app.kubernetes.io/part-of: agentic-operator
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}

{{- define "minio.selectorLabels" -}}
app.kubernetes.io/name: minio
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
