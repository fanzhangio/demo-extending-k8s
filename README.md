# Extending Kubernetes API by Kubebuilder (Demo)
This repo implements `Database` CRD and Controller which scales mysql deployments by [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Before starting

### Install tools
  - [Go](https://golang.org/)
  - [Kubebuilder](https://book.kubebuilder.io/getting_started/installation_and_setup.html)
  - [Make](https://www.gnu.org/software/make/)
  - [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl) (Optional, needed for operating Kubernetes cluster)
  - [Minikube](https://github.com/kubernetes/minikube) (Optional)
  - [Git](https://git-scm.com/downloads)

### Kubernetes Cluster
CRD needs to be installed into kubernetes cluster for fully functionanlity. You may use remote cluster or local one based on `minikube`

### Setup Development Env
- Setup `GOPATH`
- Start a minikube cluster (optional)

### Get this project
- Clone this project
  - `mkdir -p $GOPATH/src/github.com/fanzhangio/`
  - `cd $GOPATH/src/github.com/fanzhangio/`
  - `git clone https://github.com/fanzhangio/demo-extending-k8s.git`

## Demo
1. Start from project directory
   ```bash
   $ cd $GOPATH/src/github.com/fanzhangio/demo-extending-k8s
   ```
2. Start minikube cluster
   ```bash
   $ minikube start
   $ minikube get node
   ```
3. Run start scripts to starting demo
   ```bash
   $ ./start.sh
   ```
   Open another termnial to run `kubectl` see results
   ```bash
   $ kubectl get deployment
   ```

## Reference Resource
- [Kubebuilder Gitbook](http://book.kubebuilder.io)
