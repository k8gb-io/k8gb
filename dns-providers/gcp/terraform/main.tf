resource "google_dns_managed_zone" "this" {
  name     = "k8gb"
  dns_name = var.dns_zone_name
  description = "DNS zone to test the K8GB integration"
}

resource "google_service_account" "dns_admin" {
  account_id   = "dns-admin-sa"
  display_name = "DNS Admin Service Account"
}

# Assign the "roles/dns.admin" role to the service account
resource "google_project_iam_member" "dns_admin_role" {
  project = var.project_id
  role    = "roles/dns.admin"
  member  = "serviceAccount:${google_service_account.dns_admin.email}"
}
