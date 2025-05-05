#!/bin/bash

set -e  # Optional: exit on error for all commands except those in trap

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

cleanup() {
  exit_code=$?
  cd "${SCRIPT_DIR}/opentofu"
  tofu destroy -auto-approve -var="dns_zone_name=k8gb.io"
  echo "Destroyed infrastructure"
  rm -f ../../dns-provider-test/route53/values.yaml
  cd "${SCRIPT_DIR}/../.."
  rm -rf credentials
  make destroy-full-local-setup

  if [ $exit_code -eq 0 ]; then
    echo "---- Test successfull ----"
  else
    echo "---- Test failed ----"
  fi
}
trap cleanup EXIT

echo "---- Creating DNS zone and IAM user in AWS ----"
cd "${SCRIPT_DIR}/opentofu"
tofu init
tofu apply -auto-approve -var="dns_zone_name=k8gb.io"
echo "Created infrastructure"

echo "---- Creating local clusters ----"
cd "${SCRIPT_DIR}/../../"
ZONE_ID=$(aws route53 list-hosted-zones --query "HostedZones[?Name == 'k8gb.io.'].Id" --output text)
EDGE_DNS_SERVER=$(aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" --query "ResourceRecordSets[?Type == 'NS'].ResourceRecords[0]" --output text | sed 's/\.$//')
cp dns-provider-test/route53/values-template.yaml dns-provider-test/route53/values.yaml
sed -i '' "s/DNS_SERVER_TODO/$EDGE_DNS_SERVER/g" dns-provider-test/route53/values.yaml
make create-local-clusters

echo "---- Deploying AWS credentials ----"
SECRET_ACCESS_KEY=$(aws iam create-access-key --user-name "externaldns")
cat <<-EOF > credentials
[default]
aws_access_key_id = $(echo $SECRET_ACCESS_KEY | jq -r '.AccessKey.AccessKeyId')
aws_secret_access_key = $(echo $SECRET_ACCESS_KEY | jq -r '.AccessKey.SecretAccessKey')
EOF
kubectl create ns k8gb --context k3d-test-gslb1
kubectl create ns k8gb --context k3d-test-gslb2
kubectl create secret generic external-dns-secret-aws -n k8gb --from-file credentials --context k3d-test-gslb1
kubectl create secret generic external-dns-secret-aws -n k8gb --from-file credentials --context k3d-test-gslb2

echo "---- Deploying k8gb ----"
VALUES_YAML=dns-provider-test/route53/values.yaml K8GB_LOCAL_VERSION=test DEPLOY_APPS=false make deploy-test-version

echo "---- Waiting for external-dns to update Route53 records ----"
sleep 60

echo "---- Checking A records ----"
A_RECORDS_JSON=$(aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" --query "ResourceRecordSets[?Type == 'A']")

EU_IPS=$(echo "$A_RECORDS_JSON" | jq -r '.[] | select(.Name == "gslb-ns-eu-cloud.k8gb.io.") | .ResourceRecords | length')
US_IPS=$(echo "$A_RECORDS_JSON" | jq -r '.[] | select(.Name == "gslb-ns-us-cloud.k8gb.io.") | .ResourceRecords | length')

if [[ "$EU_IPS" -eq 2 && "$US_IPS" -eq 2 ]]; then
  echo "Both A records exist and have two IP addresses each."
else
  echo "ERROR: Required A records or IP addresses missing."
  echo "gslb-ns-eu-cloud.k8gb.io. IP count: $EU_IPS"
  echo "gslb-ns-us-cloud.k8gb.io. IP count: $US_IPS"
  exit 1
fi

echo "---- Checking NS records ----"
NS_RECORDS_JSON=$(aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" --query "ResourceRecordSets[?Type == 'NS']")

NS_COUNT=$(echo "$NS_RECORDS_JSON" | jq -r '[.[] | select(.Name == "k8gb.io." or .Name == "cloud.k8gb.io.")] | length')
echo "NS_COUNT: $NS_COUNT"
K8GB_IO_NS=$(echo "$NS_RECORDS_JSON" | jq -r '.[] | select(.Name == "k8gb.io.")')
echo "K8GB_IO_NS: $K8GB_IO_NS"
CLOUD_K8GB_IO_NS=$(echo "$NS_RECORDS_JSON" | jq -r '.[] | select(.Name == "cloud.k8gb.io.")')
echo "CLOUD_K8GB_IO_NS: $CLOUD_K8GB_IO_NS"

# Check values for cloud.k8gb.io.
CLOUD_NS_VALUES=$(echo "$NS_RECORDS_JSON" | jq -r '.[] | select(.Name == "cloud.k8gb.io.") | .ResourceRecords[].Value')
echo "CLOUD_NS_VALUES: $CLOUD_NS_VALUES"
CLOUD_NS_COUNT=$(echo "$CLOUD_NS_VALUES" | wc -l)
echo "CLOUD_NS_COUNT: $CLOUD_NS_COUNT"
CLOUD_NS_EU=$(echo "$CLOUD_NS_VALUES" | grep -c "^gslb-ns-eu-cloud.k8gb.io$")
echo "CLOUD_NS_EU: $CLOUD_NS_EU"
CLOUD_NS_US=$(echo "$CLOUD_NS_VALUES" | grep -c "^gslb-ns-us-cloud.k8gb.io$")
echo "CLOUD_NS_US: $CLOUD_NS_US"

if [[ -n "$K8GB_IO_NS" && -n "$CLOUD_K8GB_IO_NS" && "$NS_COUNT" -eq 2 && "$CLOUD_NS_COUNT" -eq 2 && "$CLOUD_NS_EU" -eq 1 && "$CLOUD_NS_US" -eq 1 ]]; then
  echo "Both NS records exist and cloud.k8gb.io. contains the correct values."
else
  echo "ERROR: NS record check failed."
  echo "NS record count: $NS_COUNT" but 2 expected
  echo "cloud.k8gb.io. NS values:"
  echo "$CLOUD_NS_VALUES" but gslb-ns-eu-cloud.k8gb.io and gslb-ns-us-cloud.k8gb.io expected
  exit 1
fi

echo "---- Test successfull ----"
