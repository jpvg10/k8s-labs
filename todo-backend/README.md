NATS installed as part of the deployment.

Local deployment:

```bash
kustomize build --enable-helm kustomize/overlays/local/ | kubectl apply -f -
```

Azure deployment:
```bash
kustomize build --enable-helm kustomize/overlays/azure/ | kubectl apply -f -
```
