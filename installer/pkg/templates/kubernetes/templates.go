// Code generated by go generate; DO NOT EDIT.
// using data from templates/kubernetes
package kubernetes

func TemplatesMap() map[string]string {
	templatesMap := make(map[string]string)

	templatesMap["1_sa.yaml"] = `apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
  namespace: {{ .Namespace }}`

	templatesMap["2_cluster_role.yaml"] = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
rules:
  - apiGroups:
      - argoproj.io
    resources:
      - applications
      - appprojects
    verbs:
      - get
      - list
      - watch
`

	templatesMap["3_cluster_role_binding.yaml"] = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cf-argocd-agent
subjects:
  - kind: ServiceAccount
    name: cf-argocd-agent
    namespace: {{ .Namespace }}`

	templatesMap["4_secret.yaml"] = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: cf-argocd-agent
  namespace: {{ .Namespace }}
data:
  codefresh.token: {{ .Codefresh.Token }}
  argo.token: {{ .Argo.Token }}
  kube.bearertoken: {{ .Kube.BearerToken }}`

	templatesMap["5_deployment.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: cf-argocd-agent
  name: cf-argocd-agent
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: cf-argocd-agent
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
        app: cf-argocd-agent
    spec:
      serviceAccountName: cf-argocd-agent
      containers:
      - env:
        - name: AGENT_VERSION
          value: "{{ .Agent.Version }}"
        - name: ARGO_HOST
          value: {{ .Argo.Host }}
        - name: ARGO_USERNAME
          value: {{ .Argo.Username }}
        - name: ARGO_PASSWORD
          value: {{ .Argo.Password }}
        - name: ARGO_TOKEN
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent
              key: argo.token
        - name: CODEFRESH_HOST
          value: {{ .Codefresh.Host }}
        - name: CODEFRESH_TOKEN
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent
              key: codefresh.token
        - name: IN_CLUSTER
          value: "{{ .Kube.InCluster }}"
        - name: MASTERURL
          value: "{{ .Kube.MasterUrl }}"
        - name: BEARERTOKEN
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent
              key: kube.bearertoken
        - name: SYNC_MODE
          value: "{{ .Codefresh.SyncMode }}"
        - name: APPLICATIONS_FOR_SYNC
          value: "{{ .Codefresh.ApplicationsForSync }}"
        - name: CODEFRESH_INTEGRATION
          value: {{ .Codefresh.Integration }}
        image: codefresh/argocd-agent:stable
        imagePullPolicy: Always
        name: cf-argocd-agent
      restartPolicy: Always`

	return templatesMap
}
