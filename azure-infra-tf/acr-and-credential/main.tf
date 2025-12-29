terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.52.0"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "=3.7.0"
    }
  }
}

variable "subscription_id" {
  type = string
}

variable "github_username" {
  type = string
}

variable "github_repo" {
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

# ACR

resource "azurerm_container_registry" "acr" {
  name                = "k8sLabsContainerRegistry"
  resource_group_name = azurerm_resource_group.k8s_labs.name
  location            = azurerm_resource_group.k8s_labs.location
  sku                 = "Basic"
}

output "REGISTRY_LOGIN_SERVER" {
  value = azurerm_container_registry.acr.login_server
}

# OpenID Connect credential for GitHub Actions deployment

resource "azuread_application" "entra_app" {
  display_name = "github-workflow"
}

output "AZURE_CLIENT_ID" {
  value = azuread_application.entra_app.client_id
}

resource "azuread_service_principal" "service_principal" {
  client_id = azuread_application.entra_app.client_id
}

output "AZURE_TENANT_ID" {
  value = azuread_service_principal.service_principal.application_tenant_id
}

output "SERVICE_PRINCIPAL_ID" {
  value = azuread_service_principal.service_principal.object_id
}

resource "azurerm_role_assignment" "role_assignment" {
  scope                = azurerm_container_registry.acr.id
  role_definition_name = "AcrPush"
  principal_id         = azuread_service_principal.service_principal.object_id
  principal_type       = "ServicePrincipal"
}

resource "azuread_application_flexible_federated_identity_credential" "federated_credential" {
  application_id             = azuread_application.entra_app.id
  display_name               = "github-deploy"
  audience                   = "api://AzureADTokenExchange"
  issuer                     = "https://token.actions.githubusercontent.com"
  claims_matching_expression = "claims['sub'] matches 'repo:${var.github_username}/${var.github_repo}:ref:refs/heads/*'"
}
