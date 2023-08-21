#!/bin/bash

helm repo add aqua https://aquasecurity.github.io/helm-charts/
helm repo update

helm install starboard-operator aqua/starboard-operator \
  -n starboard-operator --create-namespace \
  --values ./k8s/starboard-values.yaml \
  --version 0.10.4
