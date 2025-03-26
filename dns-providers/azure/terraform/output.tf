output "service_principal_client_id" {
    value = azuread_service_principal.k8gb.client_id
    description = "client id of the service principal"
}

output "service_principal_client_secret" {
    value = azuread_service_principal_password.k8gb.value
    description = "client secret of the service principal"
    sensitive = true
}
