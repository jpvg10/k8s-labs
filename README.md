# Kubernetes labs

## Clusters
The applications can be deployed to a local cluster or in Azure (AKS):
- Local cluster: I used [K3s](https://docs.k3s.io/).
- AKS: Follow the instructions in [azure-infra-tf](/azure-infra-tf/).

Create the namespaces:
```bash
kubectl create namespace exercises
kubectl create namespace project
```

## Exercises
- [log-output](./log-output/)
- [ping-pong](./ping-pong/)

To deploy, run in each folder:
```bash
kubectl apply -f manifests/

# If deploying to local:
kubectl apply -f manifests/local/

# If deploying to AKS:
kubectl apply -f manifests/azure/
```

## The Project
- [broadcaster](./broadcaster/)
- [todo-app](./broadcaster/)
- [todo-backend](./todo-backend/)
- [todo-reminder-job](./todo-reminder-job/)

### Deploy manually
For **broadcaster** and **todo-backend**, NATS is installed as part of the deployment. The Helm chart is specified in [shared/nats](./shared/nats/):
```bash
# If deploying to local:
kustomize build --enable-helm kustomize/overlays/local/ | kubectl apply -f -

# If deploying to AKS:
kustomize build --enable-helm kustomize/overlays/azure/ | kubectl apply -f -
```

For **todo-app**:
```bash
# If deploying to local:
kubectl apply -k kustomize/overlays/local/

# If deploying to AKS:
kubectl apply -k kustomize/overlays/azure/
```

For **todo-reminder-job**:
```bash
kubectl apply -f manifests/
```

### Deploy with ArgoCD
```bash
cd argocd
kubectl apply -f broadcaster.yaml
kubectl apply -f todo-backend.yaml
kubectl apply -f todo-app.yaml
```
