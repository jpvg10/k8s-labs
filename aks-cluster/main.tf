terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.52.0"
    }
  }
}

variable "subscription_id" {
  type = string
}

provider "azurerm" {
  features {}
  subscription_id = var.subscription_id
}

resource "azurerm_resource_group" "k8s_labs" {
  name     = "k8s-labs"
  location = "Sweden Central"
}

resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "k8s-labs-cluster"
  location            = azurerm_resource_group.k8s_labs.location
  resource_group_name = azurerm_resource_group.k8s_labs.name
  dns_prefix          = "k8slabscluster"
  sku_tier = "Free"

  default_node_pool {
    name       = "default"
    node_count = 3
    vm_size    = "Standard_A2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  web_app_routing {
    dns_zone_ids = []
  }
}

output "client_certificate" {
  value     = azurerm_kubernetes_cluster.cluster.kube_config[0].client_certificate
  sensitive = true
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.cluster.kube_config_raw
  sensitive = true
}
