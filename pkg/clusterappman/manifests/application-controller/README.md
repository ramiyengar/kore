See issue https://github.com/kubernetes-sigs/application/issues/141

The manifests were taken from here (with the image updated as below):
```
git clone git@github.com:appvia/application.git
kubectl kustomize ./config/ > application-all.yaml
```

The image referred to was generated from the source above with:
```
docker build -t quay.io/appvia/application-controller:$(git rev-parse --short HEAD) .
docker push quay.io/appvia/application-controller:$(git rev-parse --short HEAD) .
```
