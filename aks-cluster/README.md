# Azure Kubernetes Cluster

```bash
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
