
APP_NAME=seer
K8S_NAMESPACE=kafkloud
DOCKER_IMAGE=docker.io/mdross510/kafkloud-$(APP_NAME)
KAFKA_IMAGE=confluentinc/cp-kafka:7.3.2

#
#	NOTE: the "-new-" targets always build from scratch
#

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | grep -v fgrep | sed 's/\(.*\):.*## \(.*\)/\1 - \2/' | sort

update-local: ## updates the python packages
	pip install --no-cache-dir --upgrade -r requirements.txt

build-local: ## builds the application using the currently installed python packages
#	npm run build

build-new-local: update-local build-local ## updates npm then builds the application

run-local: ## runs the currently built application locally
	# use 'lsof -ti :8030' to get a list of the PID(s) if needed (e.g., kill -9 `lsof -ti :8030`)
	(. .env; export OPENAI_API_KEY; export CHROMA_DB=localhost:8020; export KAFKA_BROKER=localhost:9092; uvicorn app.fastapi:app --reload --port 8030)

run-new-local: build-new-local run-local ## runs a newly built application locally

build-image: ## (re)builds $(DOCKER_IMAGE) from the currently built application
	docker build -t $(DOCKER_IMAGE) .

build-new-image: build-new-local build-image ## (re)builds $(DOCKER_IMAGE) from a newly built application

push-image: ## pushes the currently built $(DOCKER_IMAGE) to the external repository
	docker push $(DOCKER_IMAGE)

push-new-image: build-new-image push-image ## pushes a newly (re-)built $(DOCKER_IMAGE) to the external repository

initialize-streamer: ## create the needed Kafka topics within the "streamer" component
	for topic in seer-statement seer-question seer-answer; do \
		docker run --network kafkloud-backend -it $(KAFKA_IMAGE) /bin/kafka-topics --bootstrap-server broker:29092 --topic $$topic --create --partitions 1; \
	done

run-image: ## runs the currently built $(DOCKER_IMAGE) locally
	docker run --env CHROMA_DB=chromadb:8020 --env KAFKA_BROKER=localhost:9092 -p 8030:8030 --network kafkloud-backend --env-file .env $(DOCKER_IMAGE)

run-new-image: build-new-image run-image ## runs a newly built $(DOCKER_IMAGE) locally

query-k8s: ## queries the status of the components present in Kubernetes
	kubectl get all -n $(K8S_NAMESPACE) -l component=$(APP_NAME) --no-headers

delete-k8s: ## deletes all components present in Kubernetes
	kubectl delete all -n $(K8S_NAMESPACE) -l component=$(APP_NAME)

apply-k8s: ## applies the Kubernetes manifests
	kubectl apply -f k8s

attach-k8s: ## tries to port-forward to a running pod in Kubernetes
	kubectl get pods -n $(K8S_NAMESPACE) -l component=$(APP_NAME) --no-headers
	@sleep 5
	kubectl port-forward -n $(K8S_NAMESPACE) $$(kubectl get pods -n $(K8S_NAMESPACE) -l component=$(APP_NAME) --no-headers -o jsonpath='{.items[0].metadata.name}') 30030:8030

recreate-k8s: delete-k8s apply-k8s ## deletes then recreates all components in Kubernetes

recreate-new-k8s: push-new-image recreate-k8s ## rebuilds, deletes then recreates all components in Kubernetes
