// Code generated by go generate; DO NOT EDIT.
// using data from templates/kubernetes
package kubernetes

func TemplatesMap() map[string]string {
	templatesMap := make(map[string]string)

	templatesMap["deployment.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: argocd-agent
  name: argocd-agent
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: argocd-agent
  replicas: 1
  revisionHistoryLimit: 5
  strategy:
    rollingUpdate:
      maxSurge: 50%
      maxUnavailable: 50%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: argocd-agent
    spec:
      containers:
      - env:
        - name: ARGO_HOST
          value: {{ .Argo.Host }}
        - name: ARGO_USERNAME
          value: {{ .Argo.Username }}
        - name: ARGO_PASSWORD
          value: {{ .Argo.Password }}
        - name: CODEFRESH_HOST
          value: {{ .Codefresh.Host }}
        - name: CODEFRESH_TOKEN
          value: {{ .Codefresh.Token }}
        - name: IN_CLUSTER
          value: "true"
        - name: CODEFRESH_INTEGRATION
          value: {{ .Codefresh.Integration }}
        image: codefresh/argocd-agent:stable
        imagePullPolicy: Always
        name: argocd-agent
      restartPolicy: Always
`

	return templatesMap
}
