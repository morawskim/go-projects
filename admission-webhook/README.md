# K8S admission controller

From Kubernetes [docs](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/):

> An admission controller is a piece of code that intercepts requests to the Kubernetes API server prior to persistence of the object, but after the request is authenticated and authorized.

In other words we can use admission controller to define some rules to enforce how a cluster is used.

In this demo, admission controller enforces that pod has required label `ContractPerson` 
and is able to modify request by adding a label `foo` to pod definition. 

## Create certification

If you don't want to use certs provided within this repo.
Don't use those certs on production!

```
openssl genrsa -out certs/tls.key 2048
openssl req -new -key certs/tls.key -out certs/tls.csr -subj "/CN=admission-webhook-demo.admission-demo.svc"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:admission-webhook-demo.admission-demo.svc") -in certs/tls.csr -signkey certs/tls.key -out certs/tls.crt
```

Execute command `cat certs/tls.crt | base64 | tr -d '\n'` and update caBundle in `k8s-resources.yaml` file.

Execute command `kubectl create secret tls admission-webhook-server-tls --dry-run=client -oyaml  --cert "certs/tls.crt" --key "certs/tls.key" -n admission-demo`
and replace definition of existing secret with a new one.

## Usage

1. Start k8s cluster - `minikube start`
2. Build container image `make build`
3. Create required k8s resources `kubectl apply -f ./k8s-resources.yaml`
4. Try to create a pod without required label `kubectl run -it -n test-admission --rm --restart=Never nginx --image=nginx -- sh`.  
You should see message and the pod should not be created:  
>Error from server: admission webhook "demo.example.com" denied the request: Label ContactPerson not found  
5. Create pod with required label `kubectl run -it -n test-admission --rm --restart=Never nginx --image=nginx --labels ContactPerson=User -- sh`  
Pod should be created `kubectl get pod -n test-admission`
6. To get logs from admission controller - `kubectl logs -f -n admission-demo pods/admission-webhook-demo-[RANDOM-STRING]`
7. To create a deployment with required label for pod - `kubectl apply -f ./example-deployment.yaml`
