
K8S_NAMESPACE=kafkloud

#
#	NOTE: the "-new-" targets always build from scratch
#

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | grep -v fgrep | sed 's/\(.*\):.*## \(.*\)/\1 - \2/' | sort

query-k8s: ## queries the status of the components present in Kubernetes
	kubectl get all -n $(K8S_NAMESPACE)

apply-k8s-kafkloud: ## applies the base Kubernetes manifests needed by kafkloud
	kubectl apply -f k8s

apply-k8s-streamer: ## applies the Kubernetes manifests for streamer
	(cd streamer; make apply-k8s)

apply-k8s-producer: ## applies the Kubernetes manifests for producer
	(cd producer; make apply-k8s)

apply-k8s-consumer: ## applies the Kubernetes manifests for consumer
	(cd consumer; make apply-k8s)

apply-k8s-portal: ## applies the Kubernetes manifests for portal
	(cd portal; make apply-k8s)

apply-k8s: apply-k8s-kafkloud apply-k8s-streamer apply-k8s-producer apply-k8s-consumer apply-k8s-portal ## applies the Kubernetes manifests for all the Kafkloud components

delete-k8s: ## deletes all components present in Kubernetes
	kubectl delete namespace $(K8S_NAMESPACE)

deploy-new-k8s-streamer: ## rebuilds and re-deploys streamer into Kubernetes
	(cd streamer; make recreate-new-k8s)

deploy-new-k8s-producer: ## rebuilds and re-deploys producer into Kubernetes
	(cd producer; make recreate-new-k8s)

deploy-new-k8s-consumer: ## rebuilds and re-deploys consumer into Kubernetes
	(cd consumer; make recreate-new-k8s)

deploy-new-k8s-portal: ## rebuilds and re-deploys portal into Kubernetes
	(cd portal; make recreate-new-k8s)

deploy-new-k8s: deploy-new-k8s-streamer deploy-new-k8s-producer deploy-new-k8s-consumer deploy-new-k8s-portal ## rebuilds and redeploys kafkloud into Kubernetes
