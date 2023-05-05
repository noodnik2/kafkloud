
APP_NAME=streamer
K8S_NAMESPACE=kafkloud

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | grep -v fgrep | sed 's/\(.*\):.*## \(.*\)/\1 - \2/' | sort

query-k8s: ## queries the status of the components present in Kubernetes
	kubectl get all -n $(K8S_NAMESPACE) -l component=$(APP_NAME) --no-headers

delete-k8s: ## deletes all components present in Kubernetes
	kubectl delete all -n $(K8S_NAMESPACE) -l component=$(APP_NAME)

apply-k8s: ## applies the Kubernetes manifests
	kubectl apply -f k8s

recreate-new-k8s recreate-k8s: delete-k8s apply-k8s ## deletes then recreates all components in Kubernetes

