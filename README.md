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

## functions!

### go-producer-1

```bash
cd go-producer-1
```

This progem publishes a json message to a kafka topic specified by an environment variable.

#### test

Run the unit tests:

```bash
go test -v
```

Does the code build:

NOTE: Because of a hack in kubeless for the go runtime, to have this command succeed involves commenting out the following line in the `Gopkg.toml` and rerun `dep ensure`.
```toml
ignored = ["github.com/kubeless/kubeless/pkg/functions"]
```

 At the request of a contributor, I have opened a ticket [here](https://github.com/kubeless/kubeless/issues/911).

```bash
go build -o /dev/null v1.go
```

#### deploy

```bash
kubectl apply -f go-producer-1/function.yaml
```
