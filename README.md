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
oc create cm hadolint-config-cm --from-file=hadolint.yaml=conf/hadolint.yaml -n fuse-demo

for i in pipeline-resources spring-build-pvc hadolint-task spring-build-task kustomize-deployment-task spring-maven-task spring-nexus-task postgres-already-deployed-condition spring-maven-pipeline; do
  oc create -f tekton/$i.yaml -n fuse-demo
done
```

3. Run the pipeline

```bash
$ oc create -f tekton/spring-maven-pipelinerun.yaml -n fuse-demo
```

![OCP Pipeline Run](/assets/pipeline.png)

And wait for the deployment to show

```
$ oc get deployment -n fuse-demo
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

### GitHub Triggers

This repo can be setup to trigger a build on every github push, using Tekton Triggers.
To try triggers, you need to clone this repo and update all resources to match the new repo name.

1. Deploy TriggerTemplates, TriggerBindings and EventListener Custom Resources

```bash
$ oc create -f tekton/triggers.yaml -n fuse-demo
```

2. Expose the EventListener WebHook controller:

```bash
$ oc get svc -n fuse-demo
NAME                              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
el-fuse-demo-event-listener       ClusterIP   172.30.198.191   <none>        8080/TCP            13m
spring-fuse-demo-service-github   ClusterIP   172.30.201.196   <none>        80/TCP              96m
spring-postgres-service-github    ClusterIP   172.30.103.193   <none>        5432/TCP,5433/TCP   148m
```

expose the service via an HTTP Route:

```bash
$ oc expose svc el-fuse-demo-event-listener -n fuse-demo

$ oc get routes -n fuse-demo
NAME                          HOST/PORT                                                                        PATH   SERVICES                      PORT            TERMINATION   WILDCARD
el-fuse-demo-event-listener   el-fuse-demo-event-listener-fuse-demo.apps.lab01.gpslab.club          el-fuse-demo-event-listener   http-listener                 None
```

3. Set up Webhook invocation in GitHub, under the settings tab of the cloned repo:

![GitHub Webhook Setup](/assets/webhook.png)

### Sample Post Call

linux# ``wget --post-data '{"userId":7, "title":"This is a title", "body":"Blah"}' --header 'content-type:application/json' http://fuse.apps.kubernetes.local/camel/api/post``

PS> ``wget -Method Post -Body '{"userId":7, "title":"This is a title", "body":"Blah"}' -ContentType 'application/json' http://fuse.apps.kubernetes.local/camel/api/post``

### Cleaning up

If you no longer need the example, you can just run

```bash
$ oc delete -k deployments/fuse-demo
$ oc delete -k deployments/fuse-postgres
```

## ADDITIONAL CONFIGURATION

Depending on cluster setup and administrator policies, the 'pipeline' serviceaccount may need further permissions to run this demo:

1. Add the 'privileged' scc to the service account. This is needed because otherwise the buildah pod that runs the build task will fail.

```bash
$ oc adm add-scc-to-user privileged -z system:serviceaccount:fuse-demo:pipeline
```

2. Grant the monitoring-edit role to the service account. This is needed to allow the SA to create/destroy ServiceMonitor objects.

```bash
$ oc adm policy add-role-to-user monitoring-edit system:serviceaccount:fuse-demo:pipeline  -n fuse-demo
```

## TODO

1. Fix bugs here and there

## Related repositories

- [Kubernetes/OCP Demo Application](https://github.com/mcaimi/k8s-demo-app): for some Dockerfiles, Jenkins deployment, Nexus & SonarQube deployments
- [OCP4 Component Deployment via ArgoCD](https://github.com/mcaimi/ocp4-argocd): for deployment manifests of Openshift Logging and User Workload Monitoring
- [Custom Agents](https://hub.docker.com/u/mcaimi): pre-built container images used in some tekton tasks

