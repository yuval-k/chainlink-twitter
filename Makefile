.PHONY: build-ganache-docker
build-ganache-docker:
	docker build -t yuval.dev/ganache-cli ./docker -f docker/Dockerfile.ganache-cli

.PHONY: load-ganache-docker
load-ganache-docker: build-ganache-docker
	kind load docker-image yuval.dev/ganache-cli

.PHONY: kind-start
kind-start:
	# always succeed in case cluster exists
	kind create cluster --config=manifests/kind.yaml|| :

.PHONY: deploy-token
deploy-token:
	sed -i -e 's/7545/32000/' LinkToken/truffle-config.js
	cd LinkToken && yarn install && yarn migrate-ganache

.PHONY: deploy-testnet
deploy-testnet: kind-start load-ganache-docker
	kubectl apply -f manifests/testnet.yaml

.PHONY: deploy-node
deploy-node: 
	kubectl apply -f manifests/postgresql.yaml
	kubectl apply -f manifests/chainlink.yaml

.PHONY: build-adapter-docker
build-adapter-docker:
	docker build -t yuval.dev/twitter-adapter ./adapter

.PHONY: load-adapter-docker
load-adapter-docker: build-adapter-docker
	kind load docker-image yuval.dev/twitter-adapter

.PHONY: deploy-adapter
deploy-adapter: load-adapter-docker
	kubectl apply -f manifests/adapter.yaml
