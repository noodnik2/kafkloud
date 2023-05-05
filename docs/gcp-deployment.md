# Deploying Kafkloud in GCP

## Prerequisites

Steps were taken to facilitate deployment in the Google Cloud Platform (GCP):

1. [Enable "billing"](https://cloud.google.com/billing/docs/how-to/create-billing-account)
   on the Google Cloud Console
2. Local [installation of the `gcloud` command](https://www.educative.io/answers/how-to-install-google-cloud-cli-on-macos)
   used in many examples

## Example Workflow

Here's a general sequence used to deploy `Kafkloud` in GKE once all the updated
components' Docker images are built and pushed to the repository:

```shell
$ REGION=us-west1
$ gcloud container clusters create-auto kafkloud-cluster --region=$REGION
Note: The Pod address range limits the maximum size of the cluster. Please refer to https://cloud.google.com/kubernetes-engine/docs/how-to/flexible-pod-cidr to learn how to optimize IP address allocation.
Creating cluster kafkloud-cluster in us-west1... Cluster is being health-checked (master is healthy)...done.                                                                                                                    
Created [https://container.googleapis.com/v1/projects/kafkloud/zones/us-west1/clusters/kafkloud-cluster].
To inspect the contents of your cluster, go to: https://console.cloud.google.com/kubernetes/workload_/gcloud/us-west1/kafkloud-cluster?project=kafkloud
kubeconfig entry generated for kafkloud-cluster.
NAME              LOCATION  MASTER_VERSION   MASTER_IP   MACHINE_TYPE  NODE_VERSION     NUM_NODES  STATUS
kafkloud-cluster  us-west1  1.25.7-gke.1000  34.83.9.42  e2-medium     1.25.7-gke.1000  3          RUNNING
$ gcloud container clusters get-credentials kafkloud-cluster --region=$REGION
Fetching cluster endpoint and auth data.
kubeconfig entry generated for kafkloud-cluster.
$ kubectl config get-contexts
CURRENT   NAME                                     CLUSTER                                  AUTHINFO                                 NAMESPACE
*         gke_kafkloud_us-west1_kafkloud-cluster   gke_kafkloud_us-west1_kafkloud-cluster   gke_kafkloud_us-west1_kafkloud-cluster
$ make apply-k8s-kafkloud apply-k8s
$ : wait for all pods to fully spin-up
$ kubectl get pods -w -n kafkloud
$ : after submitting the next command, open the portal application at http://localhost:30080 
$ kubectl port-forward -n kafkloud $(kubectl get pods -n kafkloud -l component=portal --no-headers -o jsonpath='{.items[0].metadata.name}') 30080:3000
^C
$ : make sure to shut everything down to avoid unnecessary billing...
$ gcloud container clusters delete kafkloud-cluster --zone $REGION
$ : if the above fails because of "is currently repairing cluster kafkloud-cluster. Please wait and try again once it is done." \
> try again after a while, or possibly deleting it manually from the GCP Console - that worked previously. \
> Also, try using the "Logs Explorer" to see if anything can be gleaned.
> Doing that revealed an possibly related error message: `The resource 'projects/kafkloud/global/instanceTemplates/gk3-kafkloud-cluster-default-pool-fbfb4d0a' was not found`
```