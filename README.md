# k8s peruse

A simple dashboard intended to obviate the need to maintain wiki docs with the state of Ingress/Service/Deployments

## Setup

Install `k3d` to manage and maintain the `k3s` kubernetes cluster.

### Cluster Setup

```bash
k3d create --publish 10080:80 --publish 10443:443 --enable-registry --workers 2
export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
```

### Install a thing

```bash
kubectl apply -f ./examples/nginx.yaml
curl localhost:10080
```

# Running Peruse

## In-Cluster Auth

Build the container, and create the necessary service accounts and RBAC

```bash
docker build .
docker tag [hash] registry.local:5000/peruse:latest
docker push registry.local:5000/peruse:latest
```

Load the example deployment

```bash
kubectl apply -f ./examples/peruse.yaml
```
