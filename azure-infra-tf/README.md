# ACR and AKS

The idea is that the Kubernetes cluster will be deleted after use and re-created as needed, to save costs. On the other hand, the container registry will only be created once to avoid having to change the GitHub secrets every time.

```bash
cd acr-and-credential

cp terraform.tfvars.example terraform.tfvars
# Replace values with your own values

az login

terraform init
terraform plan
terraform apply

# Store REGISTRY_LOGIN_SERVER, AZURE_CLIENT_ID, AZURE_TENANT_ID and AZURE_SUBSCRIPTION_ID (without quotes)
# as secrets on the GitHub repository

# To see the output values again:
terraform output
```

For each exercise:
```bash
cd aks-cluster

cp terraform.tfvars.example terraform.tfvars
# Replace values with your own values

az login

terraform init
terraform plan
terraform apply

az aks get-credentials --resource-group k8s-labs --name k8s-labs-cluster
# Answer "y" if it asks to overwrite

terraform destroy
```

When Azure is not needed anymore:
```bash
cd acr-and-credential
terraform destroy
```
