variable "resource_group_location" {
  type        = string
  default     = "germanywestcentral"
  description = "Location for all resources"
}

variable "resource_group_name" {
  type        = string
  default     = "rg-k8gb"
  description = "Resource group name to be created"
}

variable "dns_zone_name" {
  type        = string
  description = "Name of the DNS zone"
}
