# Development / Environment Notes - Log Style

Record of some notes related to investigations / thoughts / things to remember but not 
quite "solid" enough to record into `README` files, etc.

## TODO Lists
- Consider renaming the components based upon the rationale:
  - `producer` to `courier` - it's not the _source_ of the message; it just takes it to the deliverer
  - `streamer` to `deliverer` - since it _ensures messages get delivered_
  - `consumer` to `monitor` - as of this writing, `consumer` just reports on status 

## Reverse Chronological Log

The following are free-form notes related to investigations, thoughts, ideas, etc...
Keep them in reverse-chronological order in order for the sake of the reader...

### May 3rd

#### K8s Tool Try-out

#### Kubeshark
Learned about [Kubeshark](https://kubeshark.co/) which seemed like it could help diagnose
issues with routing between the pods:

```shell
$ brew tap kubeshark/kubeshark
$ brew install kubeshark
```

It was promising, but apparently there's a [new known bug](https://github.com/kubeshark/kubeshark/issues/1345)
because of which installations using `brew` aren't working, so removed it for now - should
try it again after they fix the bug.

```shell
$ brew uninstall kubeshark
Uninstalling /usr/local/Cellar/kubeshark/40.0... (5 files, 44.5MB)
```

#### Kubernetes Management Dashboard

- [kubernetes/dashboard](https://github.com/kubernetes/dashboard)

```shell
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml
$ kubectl get all -n kubernetes-dashboard --no-headers
pod/dashboard-metrics-scraper-64bcc67c9c-6f5fc   1/1   Running   0     76s
pod/kubernetes-dashboard-5c8bd6b59-pjsp8         1/1   Running   0     76s
service/dashboard-metrics-scraper   ClusterIP   10.104.82.150   <none>   8000/TCP   76s
service/kubernetes-dashboard        ClusterIP   10.102.99.24    <none>   443/TCP    76s
deployment.apps/dashboard-metrics-scraper   1/1   1     1     76s
deployment.apps/kubernetes-dashboard        1/1   1     1     76s
replicaset.apps/dashboard-metrics-scraper-64bcc67c9c   1     1     1     76s
replicaset.apps/kubernetes-dashboard-5c8bd6b59         1     1     1     76s
$ kubectl proxy
```

Now that the dashboard is installed, we need to create a token used to authenticate to it
([this video](https://www.youtube.com/watch?v=CICS57XbS9A) was somewhat helpful to learn
how to do this, but is probably outdated - i.e., creation of the "secret" tied to the
service account didn't appear to work, so I searched for and discovered the `create token`
command which might be all that's needed..._  [This page](https://medium.com/@gowthamshankar09/how-to-create-serviceaccount-secret-in-kubernetes-1-24-36a61bdb73ad#:~:text=How%20to%20create%20ServiceAccount%20Secrets,annotation%20section%20as%20shown%20below.&text=As%20indicted%20above%2C%20the%20secret,a%20service%2Daccount%2Dtoken.)
made more sense afterwards (where I first saw the `create token` command).

Also, editing the manifest of the dashboard app to change it to a `LoadBalancer` service
didn't help since I'm not running in the cloud, where a `LoadBalancer` produces an external
address...

```shell
$ kubectl create serviceaccount dashboard -n default
$ kubectl create clusterrolebinding dashboard-admin -n default \
  --clusterrole=cluster-admin --serviceaccount=default:dashboard
$ kubectl create token dashboard -n default
<secret>
```

Browse now to [the dashboard using the Kubernetes API]( http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/).

#### Routing to `NodePorts` Inside Bare-Metal

One thing that's eluded me so far is what sort of k8s "Ingress" controller is needed
or should be provided to access Kafkloud resources.  It seems an overkill to define it
here, as the configuration of an Ingress controller seems more relevant to the target
deployment environment, for which: (1) is ultimately outside the scope of the work
needed to develop and get Kafkloud up and running; and (2) not well understood by me
at this time, and pretty complex.

Researching this today however has led me to several new observations:

1. Kubernetes has its own DNS implementation, which at first glance seems like if
   only I could "plug in" to the DNS used by the O/S, I'd get the behavior I want.
    - Example: I should be able to navigate to something like
      `portal-svc.kubernetes.local` in order to access that service
2. There exist several other possible "solutions" suggested in at least
   [one article](https://stackoverflow.com/questions/44110876/kubernetes-service-external-ip-pending),
   and several of these are particular to the local Kubernetes environment,
   such as using the "Node IP Address" returned by `$(minikube ip)`, etc.
   However, I didn't find any similar built-in solutions suggested for use in
   "Docker Desktop" versions of Kubernetes...
    - One such "solution" is [MetalLB](https://metallb.universe.tf/), which creates
      an internal "Load Balancer" that Kubernetes can use to assign IP addresses,
      sort of like what the Cloud deployment environment providers will do for you.
    - A [related document](https://kubernetes.github.io/ingress-nginx/deploy/baremetal/)
      describes how `Nginx` can help in this sort of setup.
    - [This `nginx` doc](https://kubernetes.github.io/ingress-nginx/deploy/#docker-desktop)
      sounds relevant as well, indicating it's possible to setup an `Nginx` ingress
      controller that will forward everything - e.g., based upon path prefixes, etc.
    - Later I found [netris.io](https://www.netris.io/) which promises to help here as well.

##### Experiment

Trying out a suggestion from `nginx` documentation referenced in `#2` above:

```shell
$ helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
$ kubectl --namespace ingress-nginx get services -o wide -w ingress-nginx-controller
```

That didn't go so well.  It just created a service which had the same problem as before:
only available from within the cluster.

#### `vpnkit-controller`

[This article](https://www.docker.com/blog/how-kubernetes-works-under-the-hood-with-docker-desktop/)
introduced me to a controller which is suggested as assigning `localhost` to the "external IP"
value of `LoadBalancer` services:
- Vpnkit-controller - _"a port forwarding service which opens ports on the host and
  forwards connections transparently to the pods inside the VM"_

I looked at the status of this component which is described as "built-in" to the Docker Desktop
version of Kubernetes:

```shell
$ kubectl get pod vpnkit-controller -n kube-system --no-headers
vpnkit-controller   0/1   Completed   2223 (9d ago)   58d
```

It looks like it's run once then completed - why?  If I can get this thing to run, as I've
seen implied in several other "screenshots", maybe I'll get that behavior _without_ resorting
to having to run the `MetalLB` service?

Trying to restart these pods after consulting [this page](https://www.airplane.dev/blog/using-kubectl-to-restart-a-kubernetes-pod):

```shell
$ kubectl get pod vpnkit-controller -n kube-system -o yaml | k replace --force -f -
$ k get pod vpnkit-controller -n kube-system --no-headers
vpnkit-controller   1/1   Running   0     5m10s
$ curl localhost:30072/status
{"n_consumed":32}
```

Yay - it works!!

While I'm at it, I also restarted the `storage-provisioner` pod which appeared to not be
running very well; e.g.:

```shell
$ k get pod storage-provisioner -n kube-system --no-headers
storage-provisioner 0/1 ContainerStatusUnknown 38 (9d ago) 58d
$ k get pod storage-provisioner -n kube-system -o yaml | k replace --force -f -
$ k get pod storage-provisioner -n kube-system --no-headers
storage-provisioner   1/1   Running   0     92s
```

Now, I'm also able to get to the Kubernetes Dashboard (of course, having to create & use
a new token):
```shell
$ open https://localhost:443
$ kubectl create token dashboard -n default
```

##### `NodePort` vs. `LoadBalancer`

One of the sources I consulted with yesterday had indicated that on "bare metal", a
K8s `LoadBalancer` was essentially the same as a `NodePort`.  I thought to take advantage
of that and convert my `NodePort` for `portal` into a `LoadBalancer` so that it would
continue to leverage `vpnkit-controller` locally (so that I could access `portal` at
`localhost:30080`) and get a nice external address assigned in GKE deployments.  But,
when I did that, I noticed `localhost:30080` no longer responded - bummer!!!  Switching
back to `NodePort` restored operation locally at `localhost:30080` but of course meant
that I needed to use port forwarding in GKE.  

