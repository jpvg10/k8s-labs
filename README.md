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
For **todo-app**:
```bash
# If deploying to local:
kubectl apply -k kustomize/overlays/local/

# If deploying to AKS:
kubectl apply -k kustomize/overlays/azure/
```

For **broadcaster** and **todo-backend**, NATS is installed as part of the deployment. The Helm chart is specified in [shared/nats](./shared/nats/):
```bash
# If deploying to local:
kustomize build --enable-helm kustomize/overlays/local/ | kubectl apply -f -

# If deploying to AKS:
kustomize build --enable-helm kustomize/overlays/azure/ | kubectl apply -f -
```

For **todo-reminder-job**:
```bash
kubectl apply -f manifests/
```

### Deploy with pipeline
A GitHub Actions [workflow](https://github.com/jpvg10/k8s-labs/actions/workflows/main.yaml) deploys the project to AKS (except **todo-reminder-job**).

On a push to **main** branch, it deploys to the **project** namespace. On a push to another branch, it deploys to a namespace with the name of that branch. When a branch is deleted, the associated namespace is [deleted too](https://github.com/jpvg10/k8s-labs/actions/workflows/delete-namespace.yaml).

Note that this approach doesn't use the container images from Docker Hub. As part of the workflow, it builds the images and pushes them to ACR.

### Deploy with ArgoCD
A third option is to use ArgoCD to deploy the project to AKS. As of now it requires to create the Applications manually.
