# kubebuilder-demo

[![Build Status](https://travis-ci.org/saravanakumar-periyasamy/kubebuilder-demo.svg?branch=master)](https://travis-ci.org/saravanakumar-periyasamy/kubebuilder-demo) [![codecov](https://codecov.io/gh/saravanakumar-periyasamy/kubebuilder-demo/branch/master/graph/badge.svg)](https://codecov.io/gh/saravanakumar-periyasamy/kubebuilder-demo)


## Install

#### Prerequisites

* `kubectl`
* `kustomize`

#### Install steps 

* clone the repo
* run the below commands. 
```
kubectl apply -f config/crds
kustomize build config/default | kubectl apply -f -
```

## Test

* `kubectl apply -f config/samples/batch_v1alpha1_job.yaml`
* 