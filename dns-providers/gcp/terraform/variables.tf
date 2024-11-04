variable "project_id" {
  type        = string
  default     = "k8gb-440514" # TODO remove this
  description = "Project for all resources"
}

variable "dns_zone_name" {
  type        = string
  description = "Name of the DNS zone (must end with a dot)"
}
