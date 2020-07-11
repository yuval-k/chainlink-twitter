.PHONY: build-docker
build-docker:
	docker build -t ganache-cli ./docker -f docker/Dockerfile.ganache-cli

.PHONY: load-docker
load-docker: build-docker
	kind load docker-image ganache-cli

.PHONY: kind-start
kind-start:
	# always succeed in case cluster exists
	kind create cluster --config=manifests/kind.yaml|| :

.PHONY: deploy-token
deploy-token:
	sed -i -e 's/7545/32000/' LinkToken/truffle-config.js
	cd LinkToken && yarn install && yarn migrate-ganache

.PHONY: deploy-testnet
deploy-testnet: kind-start load-docker
	kubectl apply -f manifests/testnet.yaml

.PHONY: deploy-node
deploy-node: 
	kubectl apply -f manifests/postgresql.yaml
	kubectl apply -f manifests/chainlink.yaml