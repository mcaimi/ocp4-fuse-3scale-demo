# Tekton + RedHat Fuse Demo On Openshift 4.X

This simple demo builds, packages and deploys a quick Camel Route example using Tekton Pipelines.

## Prerequisites
- Openshift 4.7.x
- Git access to the source repo
- Openshift Pipelines installed
- Openshift User Workload Monitoring and Logging
- An instance of Sonatype Nexus
- Openshift and Tekton CLI commands

## Getting started

### Deploying on Openshift

There are a few prerequisites that need to be deployed beforehand:

1. Install and deploy ArgoCD as shown [here](https://github.com/mcaimi/ocp4-argocd)
2. Deploy ocp-pipeline with manifests found in the previously mentioned repo:

```bash
$ oc apply -k apps/ocp-pipelines
```

3. Deploy User Workload Monitoring and logging with manifests found in the previously mentioned repo:

```bash
$ oc apply -k apps/monitoring
$ oc apply -k apps/logging-operator
$ oc apply -k apps/logging-instance
```

4. Deploy Sonatype Nexus following the steps described [here](https://github.com/mcaimi/k8s-demo-app)

### Install and run tekton pipeline manifests

1. Create a new project:

```bash
$ oc new-project fuse-demo
```

2. Install tekton pipeline components:

```bash
oc create cm hadolint-config-cm --from-file=hadolint.yaml=conf/hadolint.yaml

for i in pipeline-resources spring-build-pvc hadolint-task spring-build-task kustomize-deployment-task spring-maven-task spring-nexus-tasl spring-maven-pipeline; do
  oc create -f tekton/$i.yaml
done
```

3. Run the pipeline

```bash
$ oc create -f tekton/spring-maven-pipelinerun.yaml
```

![OCP Pipeline Run](/assets/pipeline.png)

And wait for the deployment to show

```
$ oc get deployment
[...]
NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
spring-fuse-demo-github   0/1     1            0           19m
[...]
```

You can then discover the url by typing

```bash
$ oc get ingress -n fuse-demo
[...]
NAME                      CLASS    HOSTS                        ADDRESS   PORTS   AGE
spring-fuse-demo-github   <none>   fuse.apps.kubernetes.local             80      2m29s
[...]
```

### Sample Post Call

linux# ``wget --post-data '{"userId":7, "title":"This is a title", "body":"Blah"}' --header 'content-type:application/json' http://fuse.apps.kubernetes.local/camel/api/process``

PS> ``wget -Method Post -Body '{"userId":7, "title":"This is a title", "body":"Blah"}' -ContentType 'application/json' http://fuse.apps.kubernetes.local/camel/api/process``

### Cleaning up

If you no longer need the example, you can just run

```bash
$ oc delete -k deployments/fuse-demo
```

## TODO

1. Add backing database to the deployment
2. Fix bugs here and there
