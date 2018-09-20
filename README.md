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

### install kubeless

```bash
kubectl apply -f k8s/kubeless-v1.0.0-alpha.8.yaml
```

### install kafka

```bash
kubectl apply -f k8s/kafka-zookeeper-v1.0.0-beta.0.yaml
```

### install argo

```bash
kubectl apply -f k8s/argo-2.2.0.yaml
```

### install istio

NOTE: This is non-functional right now, don't do this if you want your function pods to enter the `RUNNING` state.

* https://engineering.bitnami.com/articles/serverless-service-mesh-with-kubeless-and-istio.html
* https://istio.io/docs/setup/kubernetes/spec-requirements/

```bash
kubectl apply -f k8s/istio-namespace.yaml
kubectl apply -f k8s/istio-crds.yaml
kubectl apply -f k8s/istio-demo.yaml
```

Label the default namespace such that istio will inject sidecars. 

```bash
kubectl label namespace default istio-injection=enabled
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

#### deploy the function

This function is special, we have a yaml file that adds some environment variables so let's apply that yaml instead of using the `kubeless` command line.

```bash
$ kubectl apply -f function.yaml
```

#### exposing the function over http

```bash
$ kubeless trigger http create go-producer-1-http-trigger \
    --function-name go-producer-1 \
    --hostname kubeless-demo.cyrusbio.com \
    --path go-producer-1
```

#### calling it

```bash
kubeless function call go-producer-1
```

### argo-list

```bash
cd argo-list
```

This function lists all argo workflows.

#### rbac

To make this function work lets do something hacky and wrong and give the `serviceAccount` `default:default` cluster-admin:

```bash
kubectl create clusterrolebinding default-admin --clusterrole=admin --serviceaccount=default:default

```

#### test

Does the code build:

NOTE: Because of a hack in kubeless for the go runtime, to have this command succeed involves commenting out the following line in the `Gopkg.toml` and rerun `dep ensure`.
```toml
ignored = ["github.com/kubeless/kubeless/pkg/functions"]
```

 At the request of a contributor, I have opened a ticket [here](https://github.com/kubeless/kubeless/issues/911).

```bash
go build -o /dev/null v1.go
```

#### deploy the function

```bash
$ kubeless function deploy argo-list \
    --runtime go1.10 \
    --handler kubeless.Handler \
    --from-file v1.go \
    --dependencies Gopkg.toml
```

#### exposing the function over http

```bash
$ kubeless trigger http create argo-list-http-trigger \
    --function-name argo-list \
    --hostname kubeless-demo.cyrusbio.com \
    --path argo-list
```

#### calling it

```bash
kubeless function call argo-list
```

OR

```bash
$ curl --header "Content-Type:application/json" kubeless-demo.cyrusbio.com/argo-list
```

### node-echo

```bash
cd node-echo
```

This function echos whatever input it is given in event.Data.

#### deploy the function

```bash
$ kubeless function deploy node-echo \
    --runtime nodejs8 \
    --handler main.echo \
    --from-file main.js
```

#### exposing the function over http

```bash
$ kubeless trigger http create node-echo-http-trigger \
    --function-name node-echo \
    --hostname kubeless-demo.cyrusbio.com \
    --path node-echo
```

#### creating a kafka topic & trigger for node-echo

This is the topic to which `go-producer-1` pushes messages. 

```bash
$ kubeless topic create node-echo-topic
```

Now, let's create a trigger so that each time a message is publishes to `node-echo-topic`, node-echo will logs the contents of the message.

```bash
$ kubeless trigger kafka create node-echo-kafka-trigger \
    --function-selector node-echo \
    --trigger-topic node-echo-topic
```

### bringing it all together

There are a few ways to trigger these functions we have created. For example, the node-echo function responds to http requests at `http://kubeless-demo.cyrusbio.com/node-echo` as well as messages published to the `node-echo-topic` kafka topic. Let's look at a few ways we can trigger `node-echo`.

Test with `kubeless topic publish`.
```bash
kubeless topic publish --topic node-echo-topic --data 'Hello, friend.'
```

Test with curl.
```bash
curl \
    --data '{"Another": "Echo"}' \
    --header "Content-Type:application/json" \
    kubeless-demo.cyrusbio.com/node-echo
```

This method invokes a function that publishes a message to the `node-echo-topic`. 
```bash
kubeless function call go-producer-1 --data 'Some data for youz'
```

Test with `kubeless function call`
```bash
kubeless function call node-echo --data 'Some different data for youz'
```

