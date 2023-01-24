# Build application

The application may be built into an OCI-Container using [ko](https://github.com/ko-build/ko):

```bash
ko build ./
```

# Run application on dev cluster
To run the application on a dev cluster start one using [kind](https://kind.sigs.k8s.io/) and the config provided in this repository:

```bash
kind create cluster --config ./devenv/kind.yaml
```

and apply the kubernetes manifests provided in this repository:

```bash
k apply -f ./k8s/
```
