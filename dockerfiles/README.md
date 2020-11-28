## Kubernetes Cluster Autoscaler Dockerfile

Build the docker image using one of the OS flavors.
```
cd dockerfiles/<OS>/
docker build . -t autoscalar
```

Download and update the configuration file from `https://raw.githubusercontent.com/Chathuru/kubernetes-cluster-autoscaler/master/bin/conf.yml-sample`

### Run the docker

```
docker run -v /<path_to_congig_file/conf.yml>:/go/bin/conf.yml autoscalar
```
