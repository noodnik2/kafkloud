# Kafkloud Portal 

## Origins

This React app is based upon [NextJS Boilerplate](https://vercel.com/templates/next.js/nextjs-boilerplate)
project (source repo [examples/nextjs](https://github.com/vercel/vercel/tree/main/examples/nextjs)),
and was converted to use Typescript by following the recipe in
[this](https://upmostly.com/next-js/how-to-convert-your-next-js-app-to-typescript) article.

## Development

- Creation of images using AI was done with [DeepAI](https://deepai.org/machine-learning-model/text2img)
- Creation of the "Favorite Icon" was done with [icoconvert.com](https://icoconvert.com/)

## Running

See useful things to do in the [`Makefile`](./Makefile) targets.

### In Kubernetes

Until we implement a proper ingress controller, port-forwarding is used; e.g.:
```shell
$ kubectl port-forward -n kafkloud $(kubectl get pods -n kafkloud -l component=portal --no-headers -o jsonpath='{.items[0].metadata.name}') 30080:3000
```

Or, if running locally:
- in Docker Desktop, check the status of the `vpnkit-controller` in the `kube-system` namespace;
  so long as it's running, we should be able to access `portal` via its `NodePort` address; e.g.,
  by browsing to `http://localhost:30080`
- in Minikube browse to `http://$(minikube ip):30080`
- etc.