# KNE and Athena Deployment

## Overview
This deployment of [KNE](https://github.com/google/kne) would be using Google Cloud infrastructure with a custom Kubernetes cluster installed as part of the setup.

## Prerequisites

1. A Google account with Google Cloud access
2. Your Google account should be granted access to Keysight Athena repository at https://source.cloud.google.com/kt-nts-athena-dev/keysight/ 
[//]: # (TODO what is a proper way to request access to the repo?)
3. Create a [GCP Service Account](https://console.cloud.google.com/iam-admin/serviceaccounts) and request Keysight Athena team to provide proper access permissions for the account to access Athena artifacts. In this setup I'm using `athena-g@kt-nas-demo.iam.gserviceaccount.com` Sevice Account
[//]: # (TODO what is a proper way to request access to the artifacts?)
4. Install [Google Cloud SDK](https://cloud.google.com/sdk/docs) and authenticate via `gcloud init`
5. Install [kOps](https://kops.sigs.k8s.io/getting_started/install/)

## Adopting command syntax to your environment

1. Throughout the document, a GCP Project ID parameter `--project=kt-nas-demo` is used for `gcloud` command syntax. Please change `kt-nas-demo` to specify a GCP Project ID you intend to use for the deployment
2. Where applicable, GCP Region `us-west1` (Oregon) and/or Zone `us-west1-b` are used withing the document. Consider changing to a region and zone that fit your deployment via `--region=us-west1` and `--zone=us-west1-b` parameters.

## Initialize the environment

1. Create a VPC network for K8s cluster deployment

| Parameter 						| Value
| --- 									| ---
| Name 									| `kne-demo`
| Description 					| Kubernetes Network Emulation Demo
| Subnets 							| Auto

```Shell
gcloud compute networks create kne-demo --project=kt-nas-demo --description="Kubernetes Network Emulation Demo" --subnet-mode=auto --mtu=1460 --bgp-routing-mode=regional
```

2. Create firewall rules for the VPC - we're going to permit all internal connectivity for now, and SSH access from the outside

```Shell
gcloud compute firewall-rules create kne-demo-allow-internal --project=kt-nas-demo --network=projects/kt-nas-demo/global/networks/kne-demo --description=Allows\ connections\ from\ any\ source\ in\ the\ network\ IP\ range\ to\ any\ instance\ on\ the\ network\ using\ all\ protocols. --direction=INGRESS --priority=65534 --source-ranges=10.128.0.0/9 --action=ALLOW --rules=all
gcloud compute firewall-rules create kne-demo-allow-ssh --project=kt-nas-demo --network=projects/kt-nas-demo/global/networks/kne-demo --description=Allows\ TCP\ connections\ from\ any\ source\ to\ any\ instance\ on\ the\ network\ using\ port\ 22. --direction=INGRESS --priority=65534 --source-ranges=0.0.0.0/0 --action=ALLOW --rules=tcp:22
```

4. Give the service account created in Prerequisites section the following IAM roles

	* Compute Instance Admin
	* Compute Network Admin
	* Storage Admin
	* Service Account User

![Adding IAM roles for the service account](images/IAM_Add_Roles.png)
![After the roles were added](images/IAM_Roles_athena-g.png)


