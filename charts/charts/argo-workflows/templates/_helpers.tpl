{{- define "argo-workflows.serverFullname" -}}
{{- printf "%s-argo-workflows-server" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "argo-workflows.controllerFullname" -}}
{{- printf "%s-argo-workflows-controller" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "argo-workflows.labels" -}}
app.kubernetes.io/name: argo-workflows
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: agentic-operator
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}
