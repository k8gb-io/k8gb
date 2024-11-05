resource "aws_route53_zone" "k8gb" {
  name = var.dns_zone_name
}

resource "aws_iam_policy" "externaldns_policy" {
  name        = "externaldns-policy"
  description = "Policy for external DNS management in Route 53"

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "route53:ChangeResourceRecordSets"
        ],
        "Resource" : [
          "arn:aws:route53:::hostedzone/*"
        ]
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "route53:ListHostedZones",
          "route53:ListResourceRecordSets",
          "route53:ListTagsForResource"
        ],
        "Resource" : [
          "*"
        ]
      }
    ]
  })
}

resource "aws_iam_user" "externaldns" {
  name = "externaldns"
}

resource "aws_iam_user_policy_attachment" "externaldns_user_policy_attachment" {
  user       = aws_iam_user.externaldns.name
  policy_arn = aws_iam_policy.externaldns_policy.arn
}
