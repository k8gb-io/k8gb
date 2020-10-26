module "k8gb-cluster-eu-west-1" {
  source       = "./terraform-aws-eks-k8gb"
  cluster_name = "k8gb-cluster-eu-west-1"
  providers = {
    aws = aws.eu-west-1
  }
}

module "k8gb-cluster-us-east-1" {
  source       = "./terraform-aws-eks-k8gb"
  cluster_name = "k8gb-cluster-us-east-1"
  providers = {
    aws = aws.us-east-1
  }
}
