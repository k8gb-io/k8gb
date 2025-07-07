terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.1.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.37.1"
    }
  }
}
