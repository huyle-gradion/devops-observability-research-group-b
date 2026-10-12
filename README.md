# env/huy — desired state for devops-app (huyle-gradion)

Single source of truth for the `devops-app` workload, watched by Argo CD (Application
`devops-app-huy`, path `k8s/app`). Contains **only** the app manifests — no app source,
no workflow, no Namespace (the namespace is bootstrapped from the source branch).

- `k8s/app/configmap.yaml` — `APP_ENV=cluster-huy`
- `k8s/app/deployment.yaml` — hardened Deployment, image by digest
- `k8s/app/service.yaml` — ClusterIP
