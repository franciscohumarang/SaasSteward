# Azure Container Apps deployment
# main.tf

# Define variables
variable "resource_group_name" {}
variable "app_name" {}
variable "location" {}
variable "image_name" {}
variable "container_port" {}

# Create a resource group
resource "azurerm_resource_group" "main" {
  name     = var.resource_group_name
  location = var.location
}

# Create an Azure Container App
resource "azurerm_container_group" "main" {
  name                = var.app_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  ip_address_type     = "public"
  dns_name_label      = var.app_name
  os_type             = "linux"

  container {
    name   = "main"
    image  = var.image_name
    ports {
      port     = var.container_port
      protocol = "TCP"
    }
  }
}

# Output the app URL
output "app_url" {
  value = "http://${azurerm_container_group.main.fqdn}:${var.container_port}"
}
