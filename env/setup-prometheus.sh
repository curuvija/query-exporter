#!/bin/bash

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
kubectl create ns monitoring

helm install kube--prometheus-stack prometheus-community/kube-prometheus-stack --version 34.1.1 -n monitoring --values ./k8s/kube-prometheus-stack-values.yaml