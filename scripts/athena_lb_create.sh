#!/bin/bash
#  NOTE: --network cannot be specified for regional address
​
# This is an instance group created by kOps
# TODO take this as a parameter
INSTANCE_GROUP=b-master-us-west1-b-albortok-k8s-local
​
gcloud compute instance-groups set-named-ports $INSTANCE_GROUP \
        --zone us-central1-a \
        --named-ports "https:30443"
​
# to keep the IP i have commented it out
# gcloud compute addresses create topo-cluster-lb --region us-central1
​
expose_port() {
  NAME=$1
  INTERNAL_PORT=$2
  EXTERNAL_PORT=$3
​
  HC_NAME=topo-cluster-health-check-${NAME}
  BE_NAME=topo-cluster-backend-service-${NAME}
  FR_NAME=topo-cluster-forwarding-${NAME}
​
  gcloud compute health-checks create tcp ${HC_NAME} --region us-central1 --port $INTERNAL_PORT
  gcloud compute backend-services create ${BE_NAME} --protocol TCP --health-checks ${HC_NAME} --health-checks-region us-central1 --region us-central1 --port-name ${NAME}
  gcloud compute backend-services add-backend ${BE_NAME} --instance-group $INSTANCE_GROUP --instance-group-zone us-central1-a --region us-central1
  gcloud compute forwarding-rules create ${FR_NAME} \
    --load-balancing-scheme external --region us-central1 --ports $EXTERNAL_PORT \
    --address topo-cluster-lb --backend-service ${BE_NAME}
}
​
expose_port "https" 30443 443
