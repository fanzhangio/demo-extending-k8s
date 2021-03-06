#!/bin/bash

NO_COLOR=${NO_COLOR:-""}
if [ -z "$NO_COLOR" ]; then
  header=$'\e[1;33m'
  reset=$'\e[0m'
else
  header=''
  reset=''
fi

function header_text {
  echo ""
  echo ""
  echo "$header$*$reset"
}

header_text "       Starting Demo. Project is from Repo : github.com/fanzhangio/demo-extending-k8s"
header_text "=====> [ Run Make to build, test, generate manifests from project ]"
echo "$ make"
make

header_text "=====> [ Run Make Install to install CRD ]"
echo "$ make install"
make install

header_text "=====> Database CRD has been installed into cluster"
echo "$ kubectl get crd"
kubectl get crd

header_text "=====>[ Create mysql object ] "
cat config/mysql-database.yaml
echo ""
echo "$ kubectl apply -f config/mysql-database.yaml"
kubectl apply -f config/mysql-database.yaml

header_text "=====>[ Deployment does not work without controller running ] "
echo "$ kubectl get deployment"
kubectl get deployment

header_text "=====>[ Run Controller locally ] "
echo "$ make run"
make run


