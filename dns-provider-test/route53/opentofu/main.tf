resource "aws_route53_zone" "k8gb" {
  name = var.dns_zone_name

  provisioner "local-exec" {
    when    = destroy
    command = <<EOT
      set -euo pipefail
      aws route53 list-resource-record-sets --hosted-zone-id ${self.id} \
        --query "ResourceRecordSets[?!(Type == 'NS' && Name == \`${self.name}.\`) && Type != 'SOA']" --output json | \
      jq -c '.[]' | \
      while IFS= read -r record; do
        # Skip empty lines
        if [ -z "$record" ]; then continue; fi
        # Validate JSON
        echo "$record" | jq empty || { echo "Invalid JSON: $record"; continue; }
        aws route53 change-resource-record-sets --hosted-zone-id ${self.id} --change-batch "$(jq -n --argjson rr "$record" '{Changes: [{Action: "DELETE", ResourceRecordSet: $rr}]}')"
      done
    EOT
  }
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

  provisioner "local-exec" {
    when    = destroy
    command = <<EOT
      set -euo pipefail
      for key in $(aws iam list-access-keys --user-name ${self.name} --query 'AccessKeyMetadata[].AccessKeyId' --output text); do
        aws iam delete-access-key --user-name ${self.name} --access-key-id $key
      done
    EOT
  }
}

resource "aws_iam_user_policy_attachment" "externaldns_user_policy_attachment" {
  user       = aws_iam_user.externaldns.name
  policy_arn = aws_iam_policy.externaldns_policy.arn
}
