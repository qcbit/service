SHELL := /bin/bash

# 	For full Kind v0.19 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.19.0

# ==============================================================================
# Define dependencies
KIND            := kindest/node:v1.27.2
KIND_CLUSTER    := starter-cluster

APP             := sales
NAMESPACE       := sales-system
BASE_IMAGE_NAME := qcbit/service
SERVICE_NAME    := sales-api
# VERSION         := 0.0.1
VERSION         := 1.0
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

run:
	go run app/services/sales-api/main.go

run-help:
	go run app/services/sales-api/main.go --help

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

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-load:
	kind load docker-image $(SERVICE_NAME):$(VERSION) --name $(KIND_CLUSTER)
#	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 
#	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe:
	kubectl describe nodes
	kubectl describe svc

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply