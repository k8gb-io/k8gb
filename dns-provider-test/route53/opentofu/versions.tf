terraform {
  required_version = ">=1.9"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.23.0"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
}
