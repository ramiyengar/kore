## **Quick Start**

The following is a quick start guide for running Kore locally to provision clusters on cloud platforms.

![Demo Video](images/demo.gif)

### Contents
- [Supported Cloud Providers](#supported-cloud-providers)
- [What is required?](#what-is-required)
- [Running Kore](#running-kore)
- [Using the CLI](#using-the-cli)

### Supported Cloud Providers

Kore enables teams to provision clusters. Supported cloud providers include:

+ Google Cloud Provider (GCP)
+ Azure - `Coming Soon`
+ AWS - `Coming Soon`

### What is required?

Before continuing with Kore, you will need:
+ Docker: install instructions can be found [here](https://docs.docker.com/install/)
+ Docker Compose: installation instructions can found [here](https://docs.docker.com/compose/install/)
+ A cloud account, (GCP currently)
+ A service account with GKE permissions
+ An External facing IDP, to authenticate users against

We will be enabling quick configuration of your cloud provider to reduce the barrier for entry in getting going with Kore, but this is to come!

If you don't have any of the above, please read the [pre-req](pre-req.md) document

### Running Kore

This will run several components using docker compose:

+ Kubernetes API in Docker
+ Kubernetes controller manager
+ Mysql Database
+ Kore API
+ Redis
+ Kore Web frontend

To launch the Kore server, from the root directory, run

```shell
make demo
```

### Using the CLI

By defining a korectl configuration file in `$HOME/.korectl/config`

```
server: http://127.0.0.1:10080
```

Is all that is required, you can then do an auth flow to the API server, which will authorise using the defined IDP.

```
./bin/korectl auth
```

Once authenticated, there are several configuration files, which can be found in: `examples` directory:

1. `credentials.yml` : Which is used to define the service account credential for GKE, this will require a Role for GKE in order to work inside of a project.
```
spec:
  region: CHANGE_ME
  project: CHANGE_ME
  account: "THE JSON SERVICE ACCOUNT"
```

The above will need replacing with the region, project name and the service account JSON payload.

2. `allocation.yml` : This is responsible for allocating the Cloud credential to one or more defined teams. To allow it allocatable to every team, this can be left as a wildcard

3. `gke-cluster.yml` : This defines the shape of the kubernetes clusters, these are essential settings for GKE. You can create many of these with different names i.e. `gke-cluster-dev.yml` and `gke-prod.yml` , make sure to change

### Provisioning a team cluster

Once the admin setup is done, teams can provision a cluster
