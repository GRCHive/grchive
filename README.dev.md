# GRCHive

This document will assume that the git checkout directory is set in an environment variable `$SRC`.

## Prerequisites

- G++ v?
- JDK8
- zip
- unzip
- go 1.13+ 
- PostgreSQL
- Node v12.10
- Docker
- Libreoffice
- jpegoptim
- imagemagick
- curl
- jq
- git

## Build Variables

The build depends on the presence on a number of properly configured variables to properly generate configuration files.
These variables are stored in `$SRC/build/variables.bzl`.
To generate this file, copy `$SRC/build/variables.bzl.tmpl` to `$SRC/build/variables.bzl` and replace variables with values valid for your configuration.

### RabbitMQ

- `RABBITMQ_USER`: The username to use to connect to the RabbitMQ server in the Docker container.
- `RABBITMQ_PASSWORD`: The password to use to connect to the RabbitMQ server in the Docker container.
- `RABBITMQ_HOST`: The hostname/IP address of the RabbitMQ server to connect to.
- `RABBITMQ_PORT`: The port the RabbitMQ server is listening on. For non-TLS connections, this should be 5672 and for TLS connections this should be 5671.
- `RABBITMQ_TLS`: Whether to use TLS to connect to the RabbitMQ server.

### PostgreSQL

- `POSTGRES_USER` : The username to use to connect to the PostgreSQL server.
- `POSTGRES_PASSWORD` : The password to use to connect to the PostgreSQL server.
- `POSTGRES_HOST` : The hostname/IP address of the PostgreSQL server to connect to.
- `POSTGRES_PORT` : The port the PostgreSQL server is listening on.

### Vault

- `VAULT_HOST` : The hostname/IP address of the Vault server to connect to. This should include the protocol (http/https).
- `VAULT_PORT` : The port the Vault server is listening on.
- `VAULT_USER` : The username to authenticate with the Vault server.
- `VAULT_PASSWORD` : The passowrd to authenticate with the Vault server.

### Okta

- `OKTA_URL`: The Org URL for Okta.
- `OKTA_KEY`: The API Key to use for the given org (API > Tokens).
- `OKTA_CLIENT_ID` : The client application Client ID (Applications > Client Credentials > Client ID).
- `OKTA_CLIENT_SECRET` : The client application Client secret (Applications > Client Credentials > Client secret).

### Webserver

- `SECURITY_KEY_0`: A key used to encrypt webserver cookies.
- `SECURITY_KEY_1`: A key used to encrypt webserver cookies.
- `HMAC_KEY`: A key used to generate HMACs (e.g. for uploading/downloading from Google Cloud Sorage).

### Sendgrid

- `SENDGRID_KEY`: API Key used for Sendgrid.

### GRCHive

- `GRCHIVE_PROJECT`: The Google Cloud project to use.
- `GRCHIVE_URI`: The URI to access the running webserver (with the http/https prefix).
- `GRCHIVE_DOMAIN`: The domain name of the webserver (without the http or https prefix and port number).
- `GRCHIVE_DOC_BUCKET`: The Google Cloud Storage bucket in which to store the documentation files.
- `GRPC_QUERY_RUNNER_HOST`: The hostname of the database_query_runner worker.
- `GRPC_QUERY_RUNNER_PORT`: The port the database_query_runner worker is listening on.
- `GRPC_QUERY_RUNNER_TLS`: Whether or not the query runner should use TLS for communication.
- `ROOT_CA_CRT`: Path to the root certificate used for the self-signed certificates.

### Registry

- `GKE_REGISTRY_USER`: Username for pulling from our Gitlab registry.
- `GKE_REGISTRY_PASSWORD`: Password (token) for pulling from our Gitlab registry.
- `GKE_REGISTRY_VAULT`: Where we store a secret for accessing the Gitlab registry in the form of the Docker config.json.

### Gitea

- `GITEA_SECRET_KEY`: The secret key to use for the Gitea installation.
- `GITEA_TOKEN`: The Vault path at which to store the Gitea API token.
- `GITEA_HOST`: The hostname/IP at which the Gitea instance is accessible.
- `GITEA_PORT`: The port at which the Gitea instance is accessible.
- `GITEA_PROTOCOL`: The protocol (http/https) with which to connect to Gitea.
- `GITEA_GLOBAL_ORG`: A global organization to hold all client repositories.

### Drone

- `DRONE_GITEA_CLIENT_SECRET_VAULT`: The Vault secret path where we store the Gitea OAuth client id and secret.
- `DRONE_DATABASE`: The database to use for Drone CI persistence.
- `DRONE_DATABASE_SECRET`: An encryption key to use to encrypt Drone secrets. This must be 32 bytes long.
- `DRONE_RPC_SECRET`: A shared secret between the Drone server and Drone runners.
- `DRONE_SERVER_PROTO`: Protocol to use to access the Drone server (http/https).
- `DRONE_SERVER_HOST`: Host at which to access the Drone server (do not include the port number).
- `DRONE_SERVER_PORT`: Port at which to listen to Drone HTTP requests.
- `DRONE_TOKEN`: The token with which to access the Drone API.
- `DRONE_RUNNER_TYPE`: The runner type to be used in the Drone pipeline. Should be `docker` locally and `kubernetes` in production.
- `DRONE_RUNNER_IMAGE`: The Docker image to use for running Drone pipelines.
- `DRONE_RUNNER_IMAGE_PULL`: Should be `never` locally and `always` in production.

### Artifactory

- `ARTIFACTORY_DATABASE`: The database name to use for Artifactory.
- `ARTIFACTORY_HOST`: The hostname or IP at which to access Artifactory.
- `ARTIFACTORY_PORT`: The port at which to listen to Artifactory requests.
- `ARTIFACTORY_EXTERNAL_PORT`: The port at which to listen to frontend Artifactory requests.
- `ARTIFACTORY_JOIN_KEY`: The `joinKey` config variable for Artifactory.
- `ARTIFACTORY_MASTER_KEY`: The `masterKey` config variable for Artifactory.
- `ARTIFACTORY_DEPLOY_USER` : Deployment username in Artifactory.
- `ARTIFACTORY_ENCRYPTED_PASSWORD` : Encrypted password of the deployment user in Artifactory.

## Grchive Standard Ports

- Webserver: `8080`
- Nginx: `80` (HTTP), `443` (HTTPS)
- RabbitMQ: `5672` (HTTP), `5671` (HTTPS)
- PostgreSQL: `5432`
- Vault: `8200`
- GRPC Query Runner: `6000`
- Gitea: `3000`
- Drone: `8888` (Server), `8889` (Docker Runner)
- Artifactory: `9998` (Artifacts), `9999` (Frontend)

## Setup Dependencies

- Install Blaze.

    ```
    cd $SRC/dependencies
    ./bootstrap_blaze.sh
    ```
- Add the Blaze output directory to your `$PATH` (`$SRC/dependencies/bazel/output`).
- Download Flyway.

    ```
    cd $SRC/dependencies
    ./download_flyway.sh
    ```
- Add the Flyway directory to your `$PATH` (`$SRC/dependencies/flyway`).
- Download Vault.

    ```
    cd $SRC/dependencies
    ./download_vault.sh
    ```

- Add the Vault directory to your `$PATH` (`$SRC/dependencies/vault`).
- Download libffi.

    ```
    cd $SRC/dependencies
    ./download_libffi.sh
    ```
- Download Python.
    
    ```
    cd $SRC/dependencies
    ./download_python.sh
    ```
- Download Kubectl.

    ```
    cd $SRC/dependencies
    ./download_kubectl.sh
    ```
- Add the Kubectl directory to your `$PATH` (`$SRC/dependencies/kubectl`).
- Download Minikube.

    ```
    cd $SRC/dependencies
    ./download_minikube.sh
    ```
- Add the Minikube directory to your `$PATH` (`$SRC/dependencies/minikube`).

## Setup Google Cloud

Obtain a service account JSON key with the permissions
- `storage.objects.delete`
- `storage.objects.get`
- `storage.objects.create`
and place this file in `$SRC/devops/gcloud` and name it `gcloud-webserver-account.json`.

## Docker Setup

We will use Docker to run all 3rd party vendor applications on our local development machines.
To ensure that these applications can talk to each other, we will place them all on the same Docker bridge network.

- `docker network create c3p0`

## Setup PostgreSQL (Docker)

- Run the PostgreSQL container.
    ```
    docker run \
        --network c3p0 \
        --name psql \
        -v /var/lib/postgres/data:/var/lib/postgresql/data \
        postgres:12.2
    ```
- Modify your `pg_hba.conf` to allow connections from the `c3p0` docker network by adding an entry into your `pg_hba.conf` file (`var/lib/postgres/data/pg_hba.conf`) of the form

    ```
    host all all SUBNET trust
    ```

    where `SUBNET` can be found by using the subnet of the `c3p0` network (`docker network inspect c3p0` under `IPAM.Config`).
    You will need to restart the container for these changes to take effect.
- Set `POSTGRES_HOST` to be the result of `docker inspect -f '{{.NetworkSettings.Networks.c3p0.IPAddress}}' psql`
- Setup PostgreSQL (does this work? I haven't actually tested this)
    ```
    cd $SRC/devops/database
    docker cp init_dev_db.sh psql:/init_dev_db.sh
    docker exec psql /init_dev_db.sh
    ```
- Generate the Flyway configurations:
    ```
    cd $SRC/devops/database/vault
    cp dev-flyway.conf.tmpl dev-flyway.conf
    cd $SRC/devops/database/webserver
    cp dev-flyway.conf.tmpl dev-flyway.conf
    ```
- Replace `localhost` with the value you used for `POSTGRES_HOST`.
- Create the PostgreSQL schema
    ```
    cd $SRC/devops/database/vault
    flyway -configFiles=./flyway/dev-flyway.conf migrate
    cd $SRC/devops/database/webserver
    flyway -configFiles=./flyway/dev-flyway.conf migrate
    ```

## Setup and Unseal Vault (Docker)

Replace `${VAULT_PORT}` with the corresponding value in the `variables.bzl` file.

- `cd $SRC`
- `bazel run //devops/docker/vault:vault`
- `docker run --network=c3p0 --name vault bazel/devops/docker/vault:vault`
- Set `VAULT_HOST` to be the result of `docker inspect -f '{{.NetworkSettings.Networks.c3p0.IPAddress}}' vault` prefixed by `http://`.
- `vault operator init -address="${VAULT_HOST}:${VAULT_PORT}" -n 1 -t 1`
- Store the unseal key and the root token somehwere.
- `vault operator unseal -address="${VAULT_HOST}:${VAULT_PORT}"`
- Paste the previously copied unseal key to unseal the vault.
- `bazel run --action_env VAULT_TOKEN="$YOUR_ROOT_TOKEN" //devops/vault:dev_init_vault`

Note that the Vault server must be started before attempting to run the webserver.
Every time the Vault server is restarted, it must be unsealed.

## Run RabbitMQ (Docker)

Replace `${RABBITMQ_PORT}` with the corresponding value in the `variables.bzl` file.

- `cd $SRC`
- `bazel run //devops/docker/rabbitmq:rabbitmq`
- `docker run --name rabbitmq --hostname rabbitmq --mount source=rabbitmqmnt,target=/var/lib/rabbitmq --network c3p0 bazel/devops/docker/rabbitmq:rabbitmq`
- Set `RABBITMQ_HOST` to be the result of `docker inspect -f '{{.NetworkSettings.Networks.c3p0.IPAddress}}' rabbitmq`.

## Gitea

We use Gitea to store and track changes to client code (data objects, scripts, etc.).
Gitea relies on access to a running PostgreSQL database which you should have setup at this point already.
Additionally, we will be setting up an NFS server to use as the storage volume for the Gitea container.

### NFS Server

- `mkdir /srv/nfs/gitea && chown -R $(whoami) /srv/nfs/gitea`
- `cd $SRC`
- `bazel run //devops/docker/nfs:nfs-server`
- `docker run -v /srv/nfs/gitea:/srv/nfs/gitea --name nfssrv --privileged --network c3p0 bazel/devops/docker/nfs:nfs-server`

Retrieve the IP address and then create an NFS volume in Docker for future use.

- `export NFS_IP=$(docker inspect -f '{{ .NetworkSettings.Networks.c3p0.IPAddress }}' nfssrv)`
- `docker volume create --driver local --opt type=nfs --opt o=vers=4,addr=$NFS_IP,rw --opt device=:/ gitea-nfsvolume`

You may need to run `sudo modprobe nfs` to get this step to work; alternatively, you can try to start the NFS service on your machine.

### Gitea

- `cd $SRC`
- `bazel run //devops/docker/gitea:gitea`
- `docker run --network c3p0 --name gitea --mount source=gitea-nfsvolume,target=/data bazel/devops/docker/gitea:gitea`
- Set `GITEA_HOST` to be the result of `docker inspect -f '{{.NetworkSettings.Networks.c3p0.IPAddress}}' gitea`.

At this point, we will need to create the initial admin user and obtain the access token that we will use throughout the rest of our apps.
Run 
```
bazel run --action_env VAULT_TOKEN="$YOUR_ROOT_TOKEN" //devops/docker/gitea:docker_access_token
```
to obtain the access token and set store it in the Vault server at the path specified by  `GITEA_TOKEN` in the `variables.bzl` file.
Note that this assumes that Gitea Docker container as well as the Vault Docker container are up and running.

## Artifactory

- `cd $SRC`
- `bazel run //devops/docker/artifactory:artifactory`
- `docker run --network c3p0 --name artifactory -v /srv/artifactory:/opt/jfrog/artifactory/var/data/artifactory/filestore bazel/devops/docker/artifactory:artifactory`

You will need to rerun the `docker run` command the first time you run the above command for it to properly pick up all the initial configuration.

At this point, Artifactory should be fully functional.
You will want to point your browser to `http://${ARTIFACTORY_HOST}:${ARTIFACTORY_EXTERNAL_PORT}` and login using the following credentials:

```
Username: admin
Password: password
```

- Go to the user management page (Administration > Identity and Access > Users) and change the admin password (and setup an admin email address).
- Furthermore, create a new user to use for deployment purposes with the default settings and a username and password of your choice.
- Next, create a new Permission (Administration > Identity and Access > Permissions) named `Deploy`.
- Under `Resources`, add the `libs-release-local` and `libs-snapshot-local` repositories.
- Under `Users`, add the newly created user and give the user `Read`, `Annotate`, `Deploy/Cache`, and `Delete/Overwrite` permissions for Repositories.
- Logout as the admin and sign in as the deployment user.
- Click on `Welcome, {USERNAME}` and click `Edit Profile`.
- Enter the password and `Unlock`.
- Show the `Encrypted Password` and save this as `ARTIFACTORY_ENCRYPTED_PASSWORD` in `build/variables.bzl`.
- Save the deployment username as `ARTIFACTORY_DEPLOY_USER` in `build/variables.bzl`.

## Drone CI

- `cd $SRC/devops/docker/drone`
- `./setup_drone.sh`
- `bazel run //devops/docker/drone:drone-build`
- `docker run --network c3p0 --name drone bazel/devops/docker/drone:drone`
- Set `DRONE_SERVER_HOST` to the result of `docker inspect -f '{{.NetworkSettings.Networks.c3p0.IPAddress}}' drone`.

At this point, you need to authorize Drone to access Gitea; this can only be done manually.
Point your browser to `${DRONE_PROTOCOL}://${DRONE_SERVER_HOST}:${DRONE_SERVER_PORT}` and login using the following credentials:

```
Username: grchive-gitea-admin
Password: ${PASSWORD}
```

You can find the right value for `${PASSWORD}` by querying Vault:
```
vault kv get -address="${VAULT_HOST}:${VAULT_PORT}" -field=password secret/gitea/token
```

You will also need to set `DRONE_TOKEN` to be the value found under `Your Personal Token` under `User Settings`.
Finally, 

- `bazel run --action_env VAULT_TOKEN="$YOUR_ROOT_TOKEN" //devops/docker/drone:link_to_gitea`
- `docker stop drone && docker start drone`

### Drone CI Runner

- `cd $SRC`
- `bazel run //devops/docker/drone_runner:drone-runner`
- `bazel run //devops/docker/drone_runner_worker_image:latest`
- `docker run --network c3p0 -v /var/run/docker.sock:/var/run/docker.sock --name drone-runner bazel/devops/docker/drone_runner:drone-runner`

## Build and Run Webserver

Note that this step relies on having an unsealed Vault (Docker), a running RabbitMQ server (Docker), and a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/webserver:webserver`
- `bazel run //src/webserver:webserver`

If you wish to run it in production mode:

- `bazel run -c opt //src/webserver:webserver`

If you wish to run the Docker container:

- `bazel run //devops/docker/webserver:docker_webserver`

## Build and Run Notification Hub.

- `cd $SRC`
- `bazel build //src/notification_hub:hub`
- `bazel run //src/notification_hub:hub`

## Build and Run Preview Generator

The preview generator worker listens on RabbitMQ for an files that clients upload and generates a PDF preview.
Thus this step also relies on having an unsealed Vault (Docker), a running RabbitMQ server (Docker), and a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/preview_generator:frontend`
- `bazel run //src/preview_generator:frontend`

If you wish to run the Docker container:

- `bazel run //devops/docker/preview_generator:docker_preview_generator`

## Build and Run Database Refresh Worker

The database worker is the one that retrieves client database schemas and functions.
This relies on having a running RabbitMQ server (Docker) and a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/database_fetcher:fetcher`
- `bazel run //src/database_fetcher:fetcher`

If you wish to run the Docker container:

- `bazel run //devops/docker/database_fetcher:docker_database_fetcher`

## Build and Run Database Runner Worker

The database worker is the one that runs SQL queries.
This relies on having a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/database_query_runner:runner`
- `bazel run //src/database_query_runner:runner`

If you wish to run the Docker container:

- `bazel run //devops/docker/database_query_runner:docker_database_query_runner`

## Run Tests

- `cd $SRC`
- `bazel test --test_output=errors //test/...`

## Running on Kubernetes (Minikube)

### Setup Minikube

- Ensure that you have a virtualization driver installed. See <https://kubernetes.io/docs/setup/learning-environment/minikube/#specifying-the-vm-driver>. kvm2 is what this guide will use.
- Run `minikube start --vm-driver=kvm2`. 

If you encounter errors when running the command:
- Ensure both libvirt and QEMU are installed.
- Ensure that `lsmod | grep kvm` returns something like `kvm                   790528  1 kvm_amd`.
- Ensure that your current user is added to the `libvirt` group: `sudo usermod --append --groups libvirt $(whoami)`
- Ensure that you have the `ebtables`, `iptables`, and `dnsmasq` packages installed.
- Ensure that `virt-host-validate` passes all tests.

After ensuring all these things, you may need to

- Reboot
- Restart the libvirt service (e.g. `sudo systemctl restart libvirtd.service`)
- Run `minikube delete`

Finally, run `eval $(minikube docker-env)` to ensure docker containers are available to minikube.
At this point you will need to rebuild all Docker containers.

### Additional Build Parameters

You will now need to rebuild most targets with the additional command line of

```
--platforms=//build:k8s
```

e.g.

```
bazel build --platforms=//build:k8s //src/webserver:webserver
```

### Storage

- `cd $SRC/devops/k8s/storage/dev`
- `kubectl apply -f .`

### Self-Signed Certificates

We use self-signed certificates to communicate with backend services.
To generate and deploy them onto your Minikube cluster run

```
$SRC/scripts/deploy/deploy_self_signed_certificate.sh
```

### PostgreSQL

- Modify your `postgresql.conf` file (e.g. `/var/lib/postgres/data/postgresql.conf`) to have

    ```
    listen_addresses = '*'
    ```
- Find the IPv4 address of your current ethernet link using `ip addr` (e.g. `192.168.1.160`).
- Modify your `pg_hba.conf` file (e.g. `/var/lib/postgres/data/pg_hba.conf`) to have

    ```
    host    all             all             192.168.1.160/16        trust
    ``` 

  where `192.168.1.160` is the IP address found by `ip addr`.
- Set the `POSTGRES_HOST` build variable to be `postgresql-dev-service`.
- `cd $SRC/devops/k8s/postgresql`
- `cp endpoint.yaml.tmpl endpoint.yaml`
- Modify `endpoint.yaml` to have the correct IP address (as found by `ip addr` earlier).
- `kubectl apply -f .`

### Vault

- `cd $SRC/devops/k8s/vault`
- `kubectl apply -f ./deployment.dev.yaml`
- `kubectl apply -f ./service-internal.yaml`
- `kubectl apply -f ./service-external.dev.yaml`
- Set the `VAULT_HOST` build variable to be `https://internal-vault-service` (this will require you to rebuild some Docker containers).
- Run `kubectl get services -l app=vault` and get the external port of the service. For example

    ```
    external-vault-service   NodePort    10.96.203.250   <none>        8200:30296/TCP   75s
    ```

  would indicate that the external port is `30296`.
- Check that you can reach the Vault server: `vault status -address="http://$(minikube ip):30296`
- Proceed to init/unseal the Vault server using the steps from the Docker Vault but using the new address.


### RabbitMQ

- `cd $SRC/devops/k8s/rabbitmq`
- `kubectl apply -f ./service.yaml`
- `kubectl apply -f ./statefulset.dev.yaml`
- Set the `RABBITMQ_HOST` build variable to be `rabbitmq-service` (this will require you to rebuild some Docker containers).

### Preview Generator

- `cd $SRC/devops/k8s/preview_generator`
- `kubectl apply -f ./deployment.dev.yaml`

### Database Refresh Worker

- `cd $SRC/devops/k8s/database_fetcher`
- `kubectl apply -f ./deployment.dev.yaml`

### Database Runner Worker

- `cd $SRC/devops/k8s/database_query_runner`
- `kubectl apply -f ./deployment.dev.yaml -f ./service.yaml`
- Set the `GRPC_QUERY_RUNNER_HOST` build variable to be `query-runner-service` (this will require you to rebuild some Docker containers).

### Webserver

- `cd $SRC`
- `bazel run //devops/docker/nginx:nginx`
- `cd $SRC/devops/k8s/webserver`
- `kubectl apply -f ./deployment.dev.yaml`
- `kubectl apply -f ./loadbalancer.yaml`
- Run `minikube service external-webserver-service --url` to obtain the URL to put into your web browser to access the webserver.

## Wireguard Setup

To perform deployments you will need to use Wireguard to access the VPN deployed on the GKE VPC network.
First, generate a private and public key to use with wireguard and store these files somewhere (e.g. `~/.config/wg`).

```
wg genkey | tee privatekey | wg pubkey > publickey
```

Now also create a client configuration file (e.g. `wg0-client.conf`), that has the contents:

```
[Interface]
Address = YOUR_PRIVATE_IP_ADDRESS
PrivateKey = YOUR_PRIVATE_KEY
MTU = 1380

[Peer]
PublicKey = SERVER_PUBLIC_KEY
Endpoint = SERVER_IP_ADDRESS:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
```

- `YOUR_PRIVATE_IP_ADDRESS`: This should be of the form `10.200.200.XXX/32`. Ensure that your IP address has not been chosen by checking the `$SRC/devops/wireguard/wg0.conf.tmpl` peer list.
- `YOUR_PRIVATE_KEY`: This should contain the contents of the `privatekey` file you generated in the previous step.
- `SERVER_PUBLIC_KEY`: This should be the server's public key that can be found in the file found at `$SRC/devops/wireguard/publickey`.
- `SERVER_IP_ADDRESS`: This is the server's public IP address. Find the appropriate IP address for the appropriate project in the Google Cloud console.

Next, modify `$SRC/devops/wireguard/wg0.conf.tmpl` and add your own PublicKey and Address as a peer.
Move the `wg0-client.conf` file to `/etc/wireguard/wg0-client.conf`.
Now, bring up the Wireguard VPN by running `sudo wg-quick up wg0-client`.
Ensure you are connected by running `sudo wg show`.
If you wish to stop using the VPN, run `sudo wg-quick down wg0-client`.
