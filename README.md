# Tekton + RedHat Fuse Demo On Openshift 4.X

This simple demo builds, packages and deploys a quick Camel Route example using Tekton Pipelines.
Deployed API is exposed via the 3Scale API Management gateway with OpenID integration (RedHat SSO)

## Prerequisites
- Openshift 4.13.X
- Git access to the source repo
- Openshift Pipelines installed
- Openshift GitOps
- Openshift Pipelines
- RedHat 3Scale
- RedHat Single Sign On Operator

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

### Deploy 3Scale API Gateway

1- Deploy 3Scale Operator with manifests found [here](https://github.com/mcaimi/k8s-demo-argocd) 

```bash
$ oc apply -k config/ocp-threescale/argocd
```

2- Deploy the API gateway and configure the managed API (tenant, product and service backend)

```bash
$ oc apply -f deployments/3scale/3scale-apimanager.yaml
$ oc apply -f deployments/3scale/3scale-tenant.yaml
$ oc apply -f deployments/3scale/3scale-backend.yaml
$ oc apply -f deployments/3scale/3scale-product.yaml
$ oc apply -f deployments/3scale/3scale-developer.yaml
```

### Deploy RedHat Single Sign On

1- Deploy both the SSO operator and the Keycloak instance with provided manifests

```bash
$ oc apply -k deployments/ocp-sso/operator
$ oc apply -f deployments/ocp-sso/keycloak.yaml
```

2- Retrieve the Admin password from the generated secret

```bash
$ oc get secret credential-demo-keycloak -o jsonpath='{.data.ADMIN_PASSWORD}'| base64 -d
```

### Configure the API Product to use Keycloak as an OpenID token provider

1- Create a new Realm in RedHat SSO and configure a new client for 3Scale

![3Scale-OpenID-Client](/assets/keycloak-3scale-openid-client.png)

2- Configure the new Client with the required permissions

![3Scale-OpenID-Permissions](/assets/keycloak-3scale-openid-permissions.png)

3- Get the 3Scale OpenID client secret 

![3Scale-OpenID-credentials](/assets/keycloak-3scale-openid-credentials.png)

4- In the 3Scale API Admin Portal, configure OpenID flow for the FUSE API product

![3Scale-OpenID-Integration](/assets/3scale-product-openid-integration.png)

### Sample Post Call

linux# ``curl -d '{"userId":7, "title":"This is a title", "body":"Blah"}' -H "Content-Type: application/json" http://fuse.apps.kubernetes.local/camel/api/post``

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

