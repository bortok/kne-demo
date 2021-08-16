#!/bin/bash
helm repo add stable https://charts.helm.sh/stable --force-update
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install -f configs/anthena_ingress-nginx.yaml ingress-nginx ingress-nginx/ingress-nginx