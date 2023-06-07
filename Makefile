SHELL := /bin/bash

# 	For full Kind v0.19 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.19.0

# ==============================================================================
# Define dependencies
KIND            := kindest/node:v1.27.2
KIND_CLUSTER    := starter-cluster

run:
	go run app/services/sales-api/main.go

tidy:
	go mod tidy
	go mod vendor

# ====================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: sales 

sales:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

dev-up-local:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-up: dev-up-local