terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.71.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "3.4.0"
    }
  }
}
