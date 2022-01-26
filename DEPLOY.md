# KNE and Ixia-c Deployment

## Overview
This document is divided into following sections:

* [Prerequisites](#prerequisites)
* [Install KNE Command Line Tool](#install-kne-command-line-tool)
* [Deploy Kubernetes Cluster for KNE](#deploy-kubernetes-cluster)
* [Validate KNE operations](#validate-kne-operations)
* [Initialize Ixia Traffic Generator (Athena) subsystem](#initialize-ixia-traffic-generator-athena-subsystem)
* [Run sample tests in KNE with Athena](#run-sample-tests-in-kne-with-athena)
* [Destroy KNE cluster](#destroy-kne-cluster)
* [Updating to the latest code](#updating-to-the-latest-code)

## Prerequisites

1. A Google account with Google Cloud access
2. Your Google account should be granted access to Keysight Athena repository at https://source.cloud.google.com/kt-nts-athena-dev/keysight/ 

[//]: # (TODO what is a proper way to request access to the repo?)

3. Install [Google Cloud SDK](https://cloud.google.com/sdk/docs) and authenticate via

```Shell
gcloud init
````

4. Create a [GCP Service Account](https://console.cloud.google.com/iam-admin/serviceaccounts) and request Keysight Athena team to provide proper access permissions for the account to access Athena artifacts. In this setup I'm using `athena-g@kt-nas-demo.iam.gserviceaccount.com` Sevice Account

```Shell
gcloud iam service-accounts create athena-g@kt-nas-demo.iam.gserviceaccount.com
gcloud projects add-iam-policy-binding kt-nas-demo --member="serviceAccount:athena-g@kt-nas-demo.iam.gserviceaccount.com" --role="roles/owner"
gcloud iam service-accounts keys create athena-g.json --iam-account=athena-g@kt-nas-demo.iam.gserviceaccount.com
````

[//]: # (TODO proper location for athena-g.json file)
[//]: # (TODO GAP what is a proper way to request access to the artifacts?)

6. Install [Go](https://golang.org/dl/) for your platform
7. Install `kubectl`

```Shell
gcloud components install kubectl
````

## Install KNE Command Line Tool

1. Clone KNE repository

```Shell
git clone https://github.com/google/kne.git
````

2. Compile KNE

```Shell
cd kne/kne_cli
go build 
cd ../..
````

## Deploy Kubernetes cluster

In this guide we support the following Kubernetes deployment options for running KNE. Please deploy a cluster using one of these methods and then return to the next step of this guide:

* [Single-machine KIND deployment](#kind-cluster-deployment)
* [kOps-managed cluster in Google Cloud](DEPLOY-kOps.md)

### KIND cluster deployment

1. Install [Docker](https://docs.docker.com/engine/install/)

2. Install KIND Go module

```Shell
GO111MODULE="on" go install sigs.k8s.io/kind@latest
cat >> $HOME/.bashrc << EOF

# Local go modules
if [ -d "\$HOME/go/bin" ] ; then
    PATH="\$PATH:\$HOME/go/bin"
fi

EOF
source $HOME/.bashrc
````

2. To deploy a KIND cluster on the same machine where KNE is installed, run

```Shell
./kne/kne_cli/kne_cli deploy kne/deploy/kne/kind.yaml
````

## Validate KNE operations

1. To validate KNE operations, create a simple two-node topology and validate `2node-host` namespace is present in the cluster

```Shell
./kne/kne_cli/kne_cli create ./kne/examples/2node-host.pb.txt
./kne/kne_cli/kne_cli show ./kne/examples/2node-host.pb.txt
kubectl get pods -n 2node-host
````

2. Check both PODs have two eth interfaces each, `eth0` and `eth1`. Note status of interface `eth1`, it should be `<BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN>`

```Shell
kubectl exec -it vm-1 -n 2node-host -- ip a
kubectl exec -it vm-2 -n 2node-host -- ip a
````

3. Bring `eth1` on the first node down, and validate that on the second node `eth1` status changes to `<NO-CARRIER,BROADCAST,MULTICAST,UP,M-DOWN>`

```Shell
kubectl exec -it vm-1 -n 2node-host -- ip link set eth1 down
kubectl exec -it vm-2 -n 2node-host -- ip a
````

4. Bring `eth1` on the first node up, and validate that on the second node `eth1` status changes to `<BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN>`

```Shell
kubectl exec -it vm-1 -n 2node-host -- ip link set eth1 up
kubectl exec -it vm-2 -n 2node-host -- ip a
````

5. Destroy the two-node topology

```Shell
./kne/kne_cli/kne_cli delete ./kne/examples/2node-host.pb.txt
kubectl get pods -n 2node-host
````

This concludes KNE validation steps. As part of the validation, we confirmed Meshnet CNI "wire" up/down operations between two nodes.

## Initialize Ixia-c Traffic Generator subsystem

1. Clone `keysight` repository from Ixia Athena development project

[//]: # (TODO cd to top directory)
[//]: # (TODO this should be moved up into prereq)

```Shell
gcloud source repos clone keysight --project=kt-nts-athena-dev
````

2. Create a namespace and a secret for Ixia K8s Operator

```Shell
EMAIL=<your email>
kubectl create ns ixiatg-op-system
kubectl create secret -n ixiatg-op-system docker-registry ixia-pull-secret \
  --docker-server=us-central1-docker.pkg.dev \
  --docker-username=_json_key \
  --docker-password="$(cat athena-g.json)" \
  --docker-email=$EMAIL
kubectl annotate secret ixia-pull-secret -n ixiatg-op-system secretsync.ixiatg.com/replicate='true'
````

3. Deploy Ixia Operator which will watch for CRD creation by KNE. This step requires a secret (explained in previous point) to successfully deploy the operator.

```Shell
kubectl apply -f keysight/athena/operator/ixiatg-operator.yaml
kubectl get pods -n ixiatg-op-system
````

4. Deploy `gosnappi` POD for running test packages from inside the KNE cluster

````
kubectl apply -f kne-demo/configs/gosnappi.yaml
watch kubectl get pods
````

  Copy Athena sample tests and test utilities to that container

````
kubectl cp keysight/athena/sample-tests gosnappi:/go/
kubectl exec gosnappi -- /bin/bash -c "go get github.com/open-traffic-generator/snappi/gosnappi@v0.7.3"
kubectl exec gosnappi -- /bin/bash -c "apt update && apt-get install libpcap-dev -y"
````

## [Run sample tests in KNE with Ixia-c](SAMPLE_TESTS.md)

## Destroy KNE cluster

* To delete K8s cluster created by kOps, use

```Shell
kops delete cluster $USER.k8s.local --yes
````

## Updating to the latest code

1. Bring KNE local copy to the latest

```Shell
cd kne
git fetch origin
git pull origin
#git checkout hines
cd kne_cli
go build
````

2. Bring Athena local copy to the latest

```Shell
cd keysight
git fetch origin
git pull origin
````

