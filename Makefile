SHELL := /bin/bash

# 	For full Kind v0.19 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.19.0
#
# Other commands to install.
# go install github.com/divan/expvarmon@latest
#
# http://sales-service.sales-system.svc.cluster.local:4000/debug/pprof
# curl -il sales-service.sales-system.svc.cluster.local:4000/debug/vars/

status:
	curl -il sales-service.sales-system.svc.cluster.local:3000/status

# RSA Keys
# To generate private/public key PEM file.
# $ openssl genpkey -algorithm rsa -out private.pem -pkeyopt rsa_keygen_bits:2048
# $ openssl rsa -pubout -in private.pem -out public.pem

jwt:
	go run app/scratch/jwt/main.go

# ==============================================================================
# Define dependencies
KIND            := kindest/node:v1.27.2
KIND_CLUSTER    := starter-cluster
GOLANG          := golang:1.20
ALPINE          := alpine:3.18
POSTGRES        := postgres:15.3
VAULT           := hashicorp/vault:1.13
GRAFANA         := grafana/grafana:9.5.2
PROMETHEUS      := prom/prometheus:v2.44.0
TEMPO           := grafana/tempo:2.1.1
TELEPRESENCE    := datawire/tel2:2.13.3

APP             := sales
NAMESPACE       := sales-system
BASE_IMAGE_NAME := qcbit/service
SERVICE_NAME    := sales-api
# VERSION         := 0.0.1
VERSION         := 1.0
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# ======================
# Install dependencies
dev-docker:
	docker pull $(GOLANG)
	docker pull $(ALPINE)
	docker pull $(KIND)
	docker pull $(POSTGRES)
	docker pull $(VAULT)
	docker pull $(GRAFANA)
	docker pull $(PROMETHEUS)
	docker pull $(TEMPO)
	docker pull $(TELEPRESENCE)

# ===============
# Go Tooling
dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest

test-load-local:
	hey -m GET -c 100 -n 10000 http://localhost:3000/status

test-load:
	hey -m GET -c 100 -n 10000 http://sales-service.sales-system.svc.cluster.local:3000/status

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

run-help:
	go run app/services/sales-api/main.go --help

tidy:
	go mod tidy
	go mod vendor

metrics-local:
	expvarmon -ports=":4000" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

metrics-view:
	expvarmon -ports="sales-service.$(NAMESPACE).svc.cluster.local:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"


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

connect-tel:
	telepresence --context=kind-$(KIND_CLUSTER) connect

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(TELEPRESENCE) --name $(KIND_CLUSTER)
	telepresence --context=kind-$(KIND_CLUSTER) helm install
	connect-tel

dev-down:
	telepresence quit -s
	kind delete cluster --name $(KIND_CLUSTER)

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
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe:
	kubectl describe nodes
	kubectl describe svc

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply