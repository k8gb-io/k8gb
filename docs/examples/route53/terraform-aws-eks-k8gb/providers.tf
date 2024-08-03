terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.61.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.31.0"
    }
  }
}
