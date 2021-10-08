{{- define "devMode.volumes" -}}
{{ if .Values.devMode.enabled }}
- name: app-dev
  hostPath:
    path: {{ .Values.devMode.appDir }}
- name: k8s-configs
  hostPath:
    path: {{ printf "%s/.k8s" .Values.devMode.homeDir }}
{{ end }}
{{- end -}}

{{- define "devMode.volumeMounts" -}}
{{ if .Values.devMode.enabled }}
- name: app-dev
  mountPath: /app/
- name: k8s-configs
  mountPath: /k8s-configs/
{{ end }}
{{- end -}}
