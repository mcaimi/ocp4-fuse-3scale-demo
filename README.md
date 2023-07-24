# Tekton + RedHat Fuse Demo On Openshift 4.X

This simple demo builds, packages and deploys a quick Camel Route example using Tekton Pipelines.

## Prerequisites
- Openshift 4.13.X
- Git access to the source repo
- Openshift Pipelines installed
- Openshift GitOps
- Openshift Pipelines

## Getting started

### Deploying on Openshift

There are a few prerequisites that need to be deployed beforehand:

1. Install and deploy ArgoCD as shown [here](https://github.com/mcaimi/k8s-demo-argocd)
2. Deploy ocp-pipeline with manifests found in the previously mentioned repo:

```bash
$ oc apply -k config/ocp-pipelines/argocd
```

### Install and run tekton pipeline manifests

1. Create a new project:

```bash
$ oc new-project fuse-jdbc-demo
```

2. Install tekton pipeline components:

```bash
oc apply -k tekton/
```

3. Run the pipeline via the Web UI (or create a pipelinerun manifest)

![OCP Pipeline Run](/assets/pipeline.png)

4. Deploy ArgoCD applications:

```bash
oc apply -k argocd/postgres
oc apply -k argocd/fuse-spring
```

And wait for the deployment to show

```
$ oc get deployment -n fuse-jdbc-demo
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

linux# ``wget --post-data '{"userId":7, "title":"This is a title", "body":"Blah"}' --header 'content-type:application/json' http://fuse.apps.kubernetes.local/camel/api/post``

PS> ``wget -Method Post -Body '{"userId":7, "title":"This is a title", "body":"Blah"}' -ContentType 'application/json' http://fuse.apps.kubernetes.local/camel/api/post``

### Cleaning up

If you no longer need the example, you can just run

```bash
$ oc delete -k argocd/fuse-spring
$ oc delete -k argocd/postgres
```

## TODO

1. Fix bugs here and there

## Related repositories

- [Kubernetes/OCP Demo Application](https://github.com/mcaimi/quarkus-notes): example quarkus application for K8s
- [OCP4 Component Deployment via ArgoCD](https://github.com/mcaimi/k8s-demo-argocd): for deployment manifests of Openshift components with Openshift GitOps

