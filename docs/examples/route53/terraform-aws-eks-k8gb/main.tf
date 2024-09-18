data "aws_eks_cluster" "cluster" {
  name = module.eks-cluster.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks-cluster.cluster_id
}

data "aws_availability_zones" "available" {
}

data "aws_subnet_ids" "private" {
  vpc_id = data.aws_vpc.vpc.id

  tags = var.private_subnet_tags
}

data "aws_subnet_ids" "public" {
  vpc_id = data.aws_vpc.vpc.id

  tags = var.public_subnet_tags
}

data "aws_vpc" "vpc" {
  tags = {
    Name = var.vpc_name
  }
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

resource "aws_security_group" "worker_group_gslb_dns" {
  name_prefix = "worker_group_gslb_dns"
  vpc_id      = data.aws_vpc.vpc.id

  ingress {
    from_port = 53
    to_port   = 53
    protocol  = "udp"

    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }
}

resource "aws_ec2_tag" "eks_shared" {
  for_each    = data.aws_subnet_ids.public.ids
  resource_id = each.value
  key         = "kubernetes.io/cluster/${var.cluster_name}"
  value       = "shared"
}

resource "aws_ec2_tag" "eks_elb" {
  for_each    = data.aws_subnet_ids.public.ids
  resource_id = each.value
  key         = "kubernetes.io/role/elb"
  value       = "1"
}

module "eks-cluster" {
  tags            = var.eks_tags
  source          = "terraform-aws-modules/eks/aws"
  version         = "20.24.1"
  cluster_name    = var.cluster_name
  cluster_version = var.kubernetes_version
  subnets         = data.aws_subnet_ids.private.ids
  vpc_id          = data.aws_vpc.vpc.id
  enable_irsa     = true

  worker_groups = [
    {
      instance_type                 = "t3.medium"
      spot_price                    = var.spot_price
      kubelet_extra_args            = "--node-labels=node.kubernetes.io/lifecycle=spot"
      asg_min_size                  = (var.ec2_workers == true ? length(data.aws_availability_zones.available.names) : 0)
      asg_max_size                  = (var.ec2_workers == true ? length(data.aws_availability_zones.available.names) + 1 : 0)
      asg_desired_capacity          = (var.ec2_workers == true ? length(data.aws_availability_zones.available.names) : 0)
      additional_security_group_ids = [aws_security_group.worker_group_gslb_dns.id]
    }
  ]

  fargate_profiles = (var.fargate_workers == true ? {
    kube-system = {
      namespace = "kube-system"

      tags = var.fargate_tags
    }
    k8gb = {
      namespace = "k8gb"

      tags = var.fargate_tags
    }
    test-gslb = {
      namespace = "test-gslb"

      tags = var.fargate_tags
    }
  } : {})
}
