variable "cluster_name" {
  type = string
}

variable "spot_price" {
  type = string
}

variable "ec2_workers" {
  type = bool
}

variable "fargate_workers" {
  type = bool
}

variable "vpc_name" {
  type = string
}

variable "eks_tags" {
  type = map(any)
}

variable "fargate_tags" {
  type    = map(any)
  default = {}
}

variable "private_subnet_tags" {
  type = map(any)
}

variable "public_subnet_tags" {
  type = map(any)
}

variable "kubernetes_version" {
  type = string
}
