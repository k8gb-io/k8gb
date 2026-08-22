resource "aws_route53_zone" "k8gb" {
  name          = var.dns_zone_name
  force_destroy = true
}

output "zone_id" {
  value = aws_route53_zone.k8gb.zone_id
}
