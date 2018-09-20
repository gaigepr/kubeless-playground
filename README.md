# kubeless-playground

## install stuff

### rbac

You may need permission to do all the things.

```bash
kubectl create clusterrolebinding <you>-cluster-admin-binding --clusterrole=cluster-admin --user=<you>@cyrusbio.com
```

### nginx ingress

#### apply some yaml

```bash
kubectl apply -f k8s/nginx-ingress-mandatory.yaml
kubectl apply -f k8s/nginx-ingress-cloud-generic.yaml
```

#### do the dns

*TODO*

### install kubeless

```bash
kubectl apply -f k8s/kubeless-v1.0.0-alpha.8.yaml
```

### install kafka

```bash
kubectl apply -f k8s/kafka-zookeeper-v1.0.0-beta.0.yaml
```
