output "service_account_email" {
    value = google_service_account.dns_admin.email
    description = "email of the service account"
}
