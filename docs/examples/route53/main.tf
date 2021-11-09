module "k8gb-cluster-eu-west-1" {
  source       = "./terraform-aws-eks-k8gb"
  cluster_name = "test-k8gb-cluster-eu-west-1"
  spot_price   = "0.0456"
  providers = {
    aws = aws.eu-west-1
  }
  ec2_workers        = true
  fargate_workers    = false
  vpc_name           = "existing-vpc-eu-west-1"
  kubernetes_version = "1.21"
  eks_tags = {
    Custom = "Tag"
  }
  private_subnet_tags = {
    ExistingSubnetTag = "private"
  }
  public_subnet_tags = {
    ExistingSubnetType = "public"
  }
}

module "k8gb-cluster-af-south-1" {
  source       = "./terraform-aws-eks-k8gb"
  cluster_name = "test-k8gb-cluster-af-south-1"
  spot_price   = "0.0416"
  providers = {
    aws = aws.af-south-1
  }
  vpc_name           = "existing-vpc-af-south-1"
  ec2_workers        = true
  fargate_workers    = false
  kubernetes_version = "1.21"
  eks_tags = {
    Custom = "Tag"
  }
  private_subnet_tags = {
    ExistingSubnetTag = "private"
  }
  public_subnet_tags = {
    ExistingSubnetType = "public"
  }
}
