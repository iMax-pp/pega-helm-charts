Deploying Pega Platform on GKE
===============================

About this document
-------------------

Use this set of procedures to deploy Pega Platform™ on an Google Kubernetes
Engine (GKE) cluster using a PostgreSQL database you configure in Google Cloud
Platform (GCP). It is designed for any level of user, from a System
Administrator to a development engineer who's interested in learning how to
install and deploy Pega Platform onto GKE.

Pega helps enterprises and agencies quickly build business apps that deliver the
outcomes and end-to-end customer experiences they need. Using the procedures in
this guide, a user will be able to install and deploy Pega software onto GKE
without much experience in either GKE configurations or Pega Platform
deployments.

How the deployment works
------------------------

Pega provides customized orchestration tools and Docker images required to
orchestrate a deployment in an GKE cluster you create for the deployment using
the following stepped tasks:

1. Prepare your local system using the appropriate instruction set:

- [Prepare a local Linux system – 45 minutes](#https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-linux.md) – install required applications and configuration files.

- [Prepare a local Windows 10 system – 45 minutes](#https://github.com/pegasystems/pega-helm-charts/blob/master/docs/prepping-local-system-runbook-windows.md) – install required applications and configuration files.

2. [Prepare your GKE resources – 45 minutes](#prepare-your-resources--45-minutes) – request a GKE cluster form Pivotal, an SQL database, and a storage resource in an account such as Google Cloud Platform (GPC).

3. [Deploying Pega Platform using Helm charts – 90 minutes](#installing-and-deploying-pega-platform-using-helm-charts--90-minutes) – customize a configuration file with your GKE details and use kubectl and Helm to install and then deploy Pega Platform onto your GKE cluster.

5. [Logging into Pega Platform – 10 minutes](#logging-into-pega-platform--10-minutes) – configure your network connections in your DNS zone so you can log onto Pega Platform.

You can review the Pega architecture overview to understand how Pega maps
Kubernetes objects with Pega applications and services: [How Pega Platform and
applications are deployed on
Kubernetes](https://community.pega.com/knowledgebase/articles/cloud-choice/how-pega-platform-and-applications-are-deployed-kubernetes).

Document goals
--------------

By following the procedures in this document, you will create a deployment of Pega Platform on which you can implement a scalable Pega application in an GKE cluster. This deployment can be used for a Pega Platform development environment. The procedures cover these required sections:

- Instructions for requesting a demo GKE cluster with suitable resources to support a scalable dev/test environment for your Pega application.

- Instructions for creating a GCP SQL Database that hosts Pega Platform data and rules in your GCP account.

- Instructions for creating a Pega Platform installation Docker image you use to install and deploy Pega Platform onto your GKE cluster.

- Guidance for logging onto a Kubernetes dashboard and exploring areas in the dashboard to help you see the status of your deployment and troubleshoot deployment and installation problems.

    After you follow this document to its completion, you'll have a Pega Platform cluster with a single PostgreSQL Database instance and two clustered virtual machines (VMs).

### Assumptions and Prerequisites

This guide assumes:

- You have a basic familiarity with running commands from a Windows 10 PowerShell with Administrator privileges or a Linux command prompt with root privileges.

- You use opensource packaging tools on Windows or Linux to install applications onto your local system.

The following account and application versions are required for use in this
document:

- A GCP account with a payment method set up to pay for the GCP resources you create using this document. You should also have sufficient GCP account permissions and knowledge to:

- Create a Google Cloud Platform project.

- Create a PostgreSQL Database.

- Select an appropriate location in which to deploy your database resource; the document assumes your location is US East.

- You are responsible for any financial costs incurred for your GCP resources.

- Pega Platform 8.3.1 or later

- Pega Docker images – your deployment requires the use of a custom Docker image to install Pega Platform into a database that is used by your AKS cluster. After you build your image, you must make it available in a private Docker registry. In order to construct your own image from the base image that Pega provides, you must have:

  - A DockerHub account to which you will push your final image to a private DockerHub repository. The image you build with Pega-provided components cannot be shared in a public DockerHub repository.

  - The Docker application downloaded to your local system, either Linux- or Windows-based. Log into your DockerHub account from the Docker applcation on your local system.

- Helm – install Helm 3.0 or later. Helm is only required to use the Helm charts and not to use the Kubernetes yaml examples directly. For detailed usage, refer to the [Helm documentation portal](https://helm.sh/docs/).

- Kubectl –the kubernetes command line tool you must use to connect and to manage your Kubernetes resources.

Create a Google Cloud Platofrm Project - 5 minutes
-----------------------------------------
In order to deploy Pega Platform to a GKE cluster, you must create a Google Cloud project inwhich you will create your Kubernetes cluster resources.

1. Using the browser of your choice, log into [Google Cloud](#https://cloud.google.com/) with your Google Cloud account credentials.

2. Click **Console** next to your profile name to open the Google Cloud Platform console page.

3. In the search tool, search for "Manage resources" and select **Manage resources IAM & admin** to display the **Google Cloud Platform console > Manage resources** page.

Alternatively, you can click the link, [Create a Google Cloud Platform project](#https://console.cloud.google.com/cloud-resource-manager?_ga=2.261938061.-1452042395.1574691090) to open the **Manage resources** page.

4. Click **+CREATE PROJECT**.

5. In the New Project window, enter a unique **Project name**, select a **Location**, if appropiate, and click **CREATE**.

With the new project created, you can proceed with completing the preparation of your local system.

Prepare your local system – 45 minutes
--------------------------------------

This document requires that you prepare a local Linux or Windows 10 system on
which you can run commands with root or Administrator privileges. To prepare
your system use the guide appropriate for your type of system:

- [Preparing your local Linux system – 45 minutes](https://github.com/pega-helm-charts/tree/master/docs/prepping-local-system-runbook-linux.md)

- [Preparing your local Windows 10 system – 45 minutes](https://github.com/pega-helm-charts/tree/master/docs/prepping-local-system-runbook-windows.md)

After you have your local system prepared, you are ready to prepare your GKE resources.

Prepare your resources – 45 minutes
-----------------------------------

This section covers the details necessary to obtain your GKE credentials and configure the required PostgreSQL database in a GCP account. Pega supports creating a PostgreSQL database in any environment as long as the IP address of the database is available to your GKE cluster.

### Request a GKE cluster

In many cases a Pega client can request a demo cluster from Pivotal running GKE 1.5 or later. If your organization is already familiar with the Pivotal
Container Service, you can create your own local PKS environment.

At a minimum, your cluster should be provisioned with at least two worker nodes
that have 8GB of RAM in order to support the typical processing loads in a Pega
Platform deployment. Pivotal supports SSL authentication, which you can request
if your organization requires it.

In order to login to your demo cluster, you must have the following information:

- The target the IP of your PKS API

- The login credentials: username and password

- Whether any SSL information is required to authenticate your login.

### Create a database resource

Create a PostgreSQL database in which to install Pega Platform. GKE deployments require you to install in an SQL database. When you are finished, you will need to obtain the database name and the server which you are creating in this procedure so you can customize your “pega” Helm chart with this information.

1. Use a browser to log onto <https://cloud.google.com/> and navigate to your
    **Console** in the upper right corner.

2. In your **Google Cloud Platform** console, use the **Navigation Menu** to go
    to **Storage** grouping and select **SQL**.

3. In the SQL page, click **+Create Instance**.

4. In the **Choose your database engine** screen, click **Choose PostgreSQL**.

5. In the **Create PostgreSQL instance** page, add details to the following
    required fields for this database server:

6. In **Instance ID**, enter a database server ID. (demo-runbook-sql1)

7. In **Default user password**, enter a “postgres” user password.

8. Select an appropriate **Region** and **Zone** for your database server.

9. In **Database version**, select **PostgreSQL 11**.

10. In **Configuration options \> Connectivity**, select **Public IP**, click
    **+ Add Network**, enter a **Name** and **Network** of one or more IP
    address to whitelist for this postgres database, and click **Done**.

    For clusters that are provisioned by Pivotal: you can launch the Kubernetes
dashboard to view the external IP address of the nodes in your cluster; to add
that IP network to the database whitelist, enter the first three sets of number,
and use 0/24 for the final set in this IP range. For example: 123.123.123.0/24.

11. In Configuration options \> Machine type and storage:

- For **Machine type**, select 4 vCPU **Cores** and 15 GB **Memory**.

- For **Network throughput**, select **SSD (Recommended)**.

- For **Storage capacity**, enter **20 GB** and select **Enable automatic
    storage increases**.

12. Configure the remaining setting using the default settings:

- For **Auto backups and high availability**, backups can be automated 1AM –
    5AM in a single zone.

- For **Flags**, no flags are required.

- **For Maintenance**, any preference is supported.

- For **Labels**, no labels are required. Labels can help clarify billing details for your GKE resources.

13. Click **Create**.

    You can see a deployment progress page displayed until your database server
deployment is complete, which can take up to 5 minutes. When complete, the GCP UI
displays all of the SQL resources in your account, which includes your newly created PostgreSQL database:

![cid:image007.png\@01D5A3B1.62812F70](media/9aa072ea703232c2f6651fe95854e8dc.62812f70)

Installing and deploying Pega Platform using Helm charts – 90 minutes
---------------------------------------------------------------------

In order to deploy Pega Platform using Helm, you must customize the “pega” Helm
chart that holds the specific settings for your deployment needs and then run a series of Helm commands to complete the deployment.

An installation followed by a deployment will take about 90 minutes total, since it takes about an hour for Pega Platform to completely install in your postgres database.

### Update the pega.yaml Helm chart values

To deploy Pega Platform, you must finalize a number of parameters in the pega.yaml Helm chart. Pega maintains a repository of Helm charts that are required to deploy Pega Platform using Helm, including a generic version of this chart. You will use the pega.yaml file to set customization parameters that are specific to your deployment.

To customize this file, you must download it from the repository to your local system, edit it with a text editor, and then save it with the same filename. To simplify the instruction, you can download the file to the \gke-demo folder you have already created on your local system. By customizing these parameters, you are configuring these required conditions:

- Specify that this is an GKE deployment.

- Access your DockerHub account in order to access the required Docker images.

- Access your GCP SQL database you created in [Create an SQL database resource](#_Create_an_SQL).

- Install the version of Pega Platform that you built into your Docker install image in [Build your Docker image from the pega-installer-ready](#_Build_your_Docker) into your Azure SQL database.

- Access the Pega Docker images you specify for your deployments.

- Specify host names for your web and stream tiers.

To finalize these details, follow these steps:

1. To download pega/pega Helm chat to the \<local filepath\>/gke-demo, enter:

`$ helm inspect values pega/pega > pega.yaml`

2. Open the pega.yaml file from this folder in a text editor and update the following settings in the chart based on your GKE requirements:

| Chart parameter name    | Purpose                                   | Your setting |
|-------------------------|-------------------------------------------|--------------|
| provider:               | Specify a GKE deployment.                 | provider:"gke"|
| actions.execute:        | Specify a “deploy” deployment type.       | execute: "deploy"   |
| Jdbc.url:               | Specify the database IP address and database name for your Pega Platform installation.        | url: jdbc:postgresql://**localhost**:5432/**dbName** where **localhost** is the public IP address you configured for your database connectivity and **dbName** is the name you entered for your PostgreSQL database in [Create a database resource](#create-a-database-resource).         |
| Jdbc.driverClass:       | Specify the driver class for a PostgreSQL database. | driverClass: "org.postgresql.Driver"                |
| Jdbc.dbType:            | Specify PostgreSQL database type.         | dbType: "postgres”   |
| Jdbc.driverUri:         | Specify the database driver Pega Platform uses during the deployment.| driverUri: "latest jar file available” For PostgreSQL databases, use the URL of the latest postgreSQL driver file that is publicly available at <https://jdbc.postgresql.org/download.html>.|
| Jdbc: username: password: | Set the security credentials for your database server to allow installation of Pega Platform into your database.   | username: "\<name of your database user\>" password: "\<password for your database user\>" For GCP PostgreSQL databases, the default user is “postgres”. |
| jdbc.rulesSchema: jdbc.dataSchema:  | Set the names of both your rules and the data schema to the values that Pega uses for these two schemas.      | rulesSchema: "rules" dataSchema: "data" |
| docker.registry.url: username: password: | This object maps the hostname of a registry to an object containing the “username” and “password” for that registry. For details, search for “index.docker.io/v1” in [Engine API v1.24](https://docs.docker.com/engine/api/v1.24/). | url: “<https://index.docker.io/v1/>” username: "\<DockerHub account username\>" password: "\< DockerHub account password\>"      |
| docker.pega.image:       | Refer to the latest Page Platform deployment image on Dockerhub.  | Image: "pegasystems/pega:latest" Pega provides these images: <https://hub.docker.com/r/pegasystems/pega-ready/tags>  |
| upgrade:    | Do not set for installations or deployments | upgrade: for non-upgrade, keep the default value. |
| tier.name: ”web” tier.service.domain:| Set a hostname for the pega-web service of the DNS zone. | domain: "\<the hostname for your web service tier\>" You assign this hostname with an external IP address and log into Pega Platform using this hostname in the URL. Your web tier hostname must comply with your networking standards and be available as an external IP address. |
| tier.name: ”stream” tier.service.domain: | Set the hostname for the pega-stream service of the DNS zone.   | domain: "\<the hostname for your stream service tier\>" Your stream tier hostname should comply with your networking standards           |
| installer.image:        | Specify the Docker image you built to install Pega Platform. | Image: "\<your installation Docker image :your tag\>" You created this image in [Constructing your Docker image from the pega-installer-ready](#_Constructing_your_Docker)   |
| installer. adminPassword:                | Specify a password for your initial logon to Pega Platform.    | adminPassword: "\<initial password\>"  |

3. Save the file.

### Deploy Pega Platform using the command line

These steps walk you through connecting your local system to your GKE cluster;
enabling the use of a browser-based Kubernetes dashboard you can use to monitor
your deployment; and performing the Helm commands required to complete your
deployment of Pega Platform on to your GKE environment.

It's easy to confuse a Helm install with a Pega Platform install, but they are
separate processes. The Helm install command uses Helm to install your
deployment as directed in the Helm charts, one in the charts\\addons folder and
one in the charts\\pega folder. In this document, you specify that the Helm
chart always “deploys” by using the setting, actions.execute: “deploy”. For the
following steps, you overwrite this function on your initial Helm install by
specifying “--set global.actions.execute:install-deploy”, which invokes an
installation of Pega Platform using your installation Docker image and then
automatically followed by a deploy. In subsequent deployments using Helm, you should
not use the override argument, “--set global.actions.execute=”, since Pega
Platform is already installed in your database.

1. Do one of the following:

- Open Windows PowerShell running as Administrator on your local system and change the location to the top folder of your gke-demo folder that you created in [Create a local folder to access all of the configuration file](#prepare-your-resources-45-minutes).

`$ cd <local filepath>\gke-demo`

- Open a Linux bash shell and change the location to the top folder of your gke-demo directory that you created in [Create a local folder to access all of the configuration file](#prepare-your-resources-45-minutes).

`$ cd /home/<local filepath>/gke-demo`

2. Use the GKE CLI to log into your account using the TBD?? and login credentials and skip SSL validation.

`$ gcloud `

If you need to validate with SSL, replace the -k with --ca-cert \<PATH TO CERT\>.

3. To use the gke CLI to view your all of your gke cluster names and status of each, enter:

`$ pks clusters`

Your cluster name is displayed in the **Name** field.

4. To use the gke CLI to download the cluster Kubeconfig access credential file, which is specific to your cluster, into your \<local filepath\>/.kube directory, enter:

```yaml
$ gke get-credentials <cluster-name>`
Fetching credentials for cluster pega-platform.
Context set for cluster pega-platform.
```

If you need to use a Bearer Token Access Credentials instead of this credential file, see the Pivotal document, [Accessing Dashboard](https://docs.pivotal.io/pks/1-3/access-dashboard.html).

5. To use the kubectl command to view the VM nodes, including cluster names and status, enter:

`$ kubectl get nodes`

6. Establish a required cluster role binding setting so you can launch the kubernetes dashboard.

`$ kubectl create clusterrolebinding dashboard-admin -n kube-system
--clusterrole=cluster-admin --serviceaccount=kube-system:kubernetes-dashboard`

7. To start the proxy server for the kubernetes dashboard, enter:

`$ kubectl proxy`

8. To access the Dashboard UI, open a browser and navigate to the following:

<http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/>

9. In the **Kubernetes Dashboard** sign in window, choose the appropriate sign in method:

- To use a cluster Kubeconfig access credential file: select **Kubeconfig**, navigate to your \<local filepath\>/.kube directory and select the config file. Click **SIGN IN**.

- To use a cluster a Kubeconfig token: select **Token** and paste your
    Kubeconfig token into the **Enter token** area. Click **SIGN IN**.

    You can now view your deployment details visually using the Kubernetes dashboard. You’ll use this dashboard to review the status of your deployment as you continue. At this point, with no deployment, you only see the GKE resources. Note that the Kubernetes dashboard does not display your GKE cluster name or your resource name. This is expected behavior.

    In order to continue using the Kubernetes dashboard to see the progress of your deployment, keep this PowerShell or Linux command prompt open and open a new one for the remaining steps.

10. Do one of the following:

- Open a new Windows PowerShell running as Administrator on your local system and change the location to the top folder of your gke-demo folder.

`$ cd <local filepath>\gke-demo`

- Open a new Linux bash shell and change the location to the top folder of your gke-demo directory.

`$ cd /home/<local filepath>/gke-demo`

11. Create namespaces for both your Pega deployment and the addons:

```yaml
$ kubectl create namespace mypega
namespace/mypega created
$ kubectl create namespace pegaaddons
namespace/pegaaddons created
```

12. To install the addons chart, which you already updated during the prepping a local system for your deployment, enter:

```yaml
$ helm install addons pega/addons --namespace pegaaddons --values addons.yaml
```

The pegaddons namespace contains the deployment’s load balancer and disables the metric server. A successful pegaaddons deployment returns details of deployment progress. For further verification of your deployment progress, you can refresh the Kubernetes dashboard and look in the pegaaddons Namespace view.
13. To deploy Pega Platform for the first time by specifying to install Pega Platform into the database you specified in the Helm chart, install the pega.yaml Helm chart:

```yaml
helm install mypega pega/pega --namespace mypega --values pega.yaml --set global.actions.execute=install-deploy
```

For subsequent Helm installs, use the command `helm install mypega pega/pega --namespace mypega --values pega.yaml` to deploy Pega Platform and avoid another Pega Platform installation.

A successful Pega deployment immediately returns details that show progress for your deployment.

14. Refresh the Kubernetes dashboard you opened in step 7. If you closed the dashboard, open a new command prompt running as Administrator and relaunch the browser as directed in Step 10.

15. In the dashboard, use the **Namespace** pulldown to change the view to **mypega** and click on the **Pods** view.

![](media/055d24b4ac0c0dfcb9c68cec334ce42a.png)

Note: A deployment takes about 15 minutes for all of the resource configurations to complete; however a full Pega Platform installation into the database can take up to an hour.

To follow the progress of an installation, use the dashboard. For subsequent deployments, you will not need to do this. Initially, some of the resources are making requests to complete the configuration; therefore, you will see red warnings while the configuration is finishing. This is expected behavior.

16. To view the status of an installation, on the Kubernetes dashboard, select **Jobs**, locate the **pega-db-install** job, and click the logs icon located on the right side of that row.

    After you open the logs view, you can click the icon for automatic refresh to see current updates to the install log.

17.  To see the final deployment in the Kubernetes dashboard after about 15 minutes, refresh the **mypega** namespace pods.

![](media/f7779bd94bdf3160ca1856cdafb32f2b.png)

A successful deployment will not show errors across the various workloads. The **mypega** Namespace **Overview** view shows charts of the percentage of complete tiers and resources configurations. A successful deployment will have 100% complete **Workloads**.

![](media/0fb2d07a5a8113a9725b704e686fbfe6.png)

Logging into Pega Platform – 10 minutes
---------------------------------------

After you complete your deployment, you must associate the hostname of the
pega-web tier with the IP address that the deployment load balancer gave to the
tier. This final step ensures that you can log onto Pega Platform using your
hostname, on which you can independently manage security protocols that match
your networking infrastructure standards.

### Logging in using the IP address

To view the pega deployment components, enter:

`$ kubectl get services --namespace mypega`

![](media/f329e9f92feed8cb5959d91db246aa84.png)

The pega-web tier external IP address and port number are displayed. Port 80 is
used for http traffic, which means you can’t use https encryption when access
the web-tier in a browser with

Instead, Pega recommends using the domain name of the web tier.

### Logging in using the domain name of the web tier

You must manually set the IP address of the web tier domain name in order to log
into Pega Platform domain name that is set during the deployment.

The example of a domain name of the web tier used in this demo is
**gke.web.dev.pega.io**, which you set in the values.yaml file here:

```yaml
name: "web"

# Enter the domain name to access web nodes via a load balancer.
# e.g. web.mypega.example.com

service:

domain: "**gke.web.dev.pega.io**"
```

When you set this to be "\<the hostname for your web service tier\>" as directed
in [Update the Helm chart values](#update-the-helm-chart-values), you will
manually associate "\<the hostname for your web service tier\>" with the IP
address of the web tier domain name.

In order to sign into Pega using "\<the hostname for your web service tier\>",
you must assign the domain name with the same IP address that the deployment
load balancer has assigned to the web tier.

1. From your command prompt, review the IP addresses that are in the pega service

`$ kubectl get services --namespace mypega`

![](media/f329e9f92feed8cb5959d91db246aa84.png)

2. TBD.


With the domain name set to this IP address, you can log into Pega Platform with
a browser using the URL: http:// "\<the hostname for your web service
tier\>"/prweb

![](media/25b18c61607e4e979a13f3cfc1b64f5c.png)
