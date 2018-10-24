#!/bin/bash

set +x

kubectl get deployment

kubectl delete deployment my-user

kubectl delete databases.example.demo.k8s.io/mysql

kubectl delete crd databases.example.demo.k8s.io

kubectl get deployment