# dns zone

resource "azurerm_resource_group" "this" {
  name     = var.resource_group_name
  location = var.resource_group_location
}

resource "azurerm_dns_zone" "this" {
  name = var.dns_zone_name
  resource_group_name = azurerm_resource_group.this.name
}

# service principal

resource "azuread_application" "k8gb" {
  display_name               = "k8gb"
}

resource "azuread_service_principal" "k8gb" {
  client_id = azuread_application.k8gb.client_id
}

resource "azuread_service_principal_password" "k8gb" {
  service_principal_id = azuread_service_principal.k8gb.id
  end_date             = "2099-01-01T00:00:00Z"
}

resource "azurerm_role_assignment" "dns_zone_contributor" {
  principal_id         = azuread_service_principal.k8gb.object_id
  role_definition_name = "DNS Zone Contributor"
  scope                = azurerm_dns_zone.this.id
}

resource "azurerm_role_assignment" "resource_group_reader" {
  principal_id         = azuread_service_principal.k8gb.object_id
  role_definition_name = "Reader"
  scope                = azurerm_resource_group.this.id
}
