#!/bin/bash
# ATTENTION! Service account needs to have Storage Admin role associated!
# Failure to do so will make the node unable to access desired cluster configuration
# which will result in cluster setup failure.
# A narrower role MAY exist but I was too lazy to find it.
# Create google storage bucket, and
export KOPS_STATE_STORE=gs://kne-demo-bucket
export KOPS_FEATURE_FLAGS=AlphaAllowGCE
export CLUSTER=$USER.k8s.local
export SITE="`curl -s ifconfig.me`/32" # ip range you will be accessing cluster from
export ZONES=us-west1-b
export VPC=kne-demo
export SVCACCNT=athena-g@kt-nas-demo.iam.gserviceaccount.com
kops create cluster $CLUSTER \
  --project=kt-nas-demo \
  --node-count 2 \
  --node-size e2-standard-8 \
  --image ubuntu-os-cloud/ubuntu-2004-focal-v20210315 \
  --zones $ZONES \
  --master-zones $ZONES \
  --master-size e2-standard-2 \
  --subnets $ZONES \
  --vpc $VPC \
  --topology private \
  --ssh-access=$SITE \
  --admin-access=$SITE \
  --networking=calico \
  --gce-service-account=$SVCACCNT \
  --associate-public-ip="false" \
  --yes
sleep 300
kops validate cluster $CLUSTER --wait 10m


#  --image-family=ubuntu-2004-lts \
#  --image-project=ubuntu-os-cloud \
