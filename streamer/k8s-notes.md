# Kubernetes Notes

## Object Definition Files
- Always have four sections:
  - `apiVersion`
  - `kind`
  - `metadata`
  - `spec`

## Object Types

- `Pod`
- `ReplicaSet`
- `Deployment`
- `ConfigMap`
- `Service`
- `Ingress`

### `Pod`s
- Smallest deployment unit
- One or more containers deployed together, and can use
  `localhost` to communicate with each other
- Unique IP address
- Described using a "Pod Spec"

### `ReplicationController` vs. `ReplicaSet`
- New controller type; this newer controller type has started to replace `ReplicationController`
  - _"The replication controller supports equality based selectors whereas the replica set supports
    equality based as well as set based selectors"_
    - i.e., `ReplicaSet` supports `selector` / `matchLabels`

### Useful Commands
- Example of using the `--dry-run=client` and `-o yaml` to create definition files 
  - Generate an `nginx` deployment definition file with 4 replicas:
    - `kubectl create deployment --image=nginx nginx --dry-run=client --replicas=4 -o yaml > nginx-deployment.yaml`