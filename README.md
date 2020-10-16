<p align="center">
  <img src="./docs/safira.png" width="400" />
</p>

<p align="center">
    <img src="https://img.shields.io/badge/license-Apache%202.0-blue" />
    <img src="https://github.com/vertigobr/safira/workflows/Build%20Release/badge.svg" />
</p>

Safira is a CLI Tool build with [Go](https://golang.org/) that has the objetive to make it easier for the Develops to Build and Deploy it's functions on the runtime cluster. 

It makes use of [Open FaaS](https://www.openfaas.com/), to enable that people out of the DevOps scope can manage it's application without the need to know how to operate the Containers in such a low level.

Safira also helps the local development by using [k3d](https://k3d.io/) to fully provision a Kubernetes cluster instance for testing purposes.

### Minimum required
Safira was made to run on LInux OS. It's pre-requisites are:
- [Docker](https://www.docker.com/)
- [Git](https://git-scm.com/)

Debian distributions also requires the ca-certificate instalation.
```sh
sudo apt-get install -y ca-certificates
```

## Install
Safira can be installed through a shell script or in manual way.

### Shellcript
This installation will bring the latest version of the tool

```sh
curl -fsSL -o get_safira.sh https://raw.githubusercontent.com/vertigobr/safira/master/install.sh
chmod 700 get_safira.sh
./get_safira.sh
```

### Manual
Down load the desired version [here](https://github.com/vertigobr/safira/releases).

Then simplily extract the binary and move it to the bin folder:

```sh
tar -zxvf NOME_DO_ARQUIVO.tar.gz
mv safira /usr/local/bin/safira
```

## Quickstart
Try yourself to use safira following the next steps:

### Init
In order to start using Safira, the first instruction you will have to use is the init:
```sh
sudo -E safira init
```
It will download and install all the required set of tools to your local environment:
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://k3d.io/#installation)
- [helm](https://helm.sh/docs/intro/install/)
- [faas-cli](https://github.com/openfaas/faas-cli)
- [okteto](https://okteto.com/docs/getting-started/installation/index.html)

### Up
As was mentioned earlier, with Safira you can spin-up a a fully functional Kubernetes cluster on your local environment so you can test the Functions in a production-like environment. 

To do that we need the infra up command:
```sh
safira infra up
```
Then we can check the services that were deployed with

```sh
safira infra status
```
With this local environment, we are now able to test the integration between the functions and the services we are about to create.

### Functions
As we are talking here about Serverless architecture with Open FaaS, we can also use safira to help us deliver these functions. 

These functions follows a pattern called Templates. To check the available Templates just use:
```sh
safira template list
```

Currently it supports Java, Node, Python and Nodered templates.

Create a new folder for your project, and from inside it use the function set of commands:
```sh
safira function new [FUNCTION NAME] --lang [TEMPLATE NAME]
```
The function will be created with a Hello World sample on it.

### Deploy
Having the Function and the Local cluster provisioned, we now want to deploy our function and test it. For doing this first we will build our function:

```sh
safira function build-push [FUNCTION NAME]
```

Then we deploy it
```sh
safira function deploy [FUNCTION NAME]
```

Finally we can use infra status again to check the URL which we can access the Function:
```sh
safira infra status
```
It will display something like this:
```sh
SERVICES
NAME                   STATUS       AVAILABILITY       URL
basic-auth-plugin      1/1          Ready              
nats                   1/1          Ready              
queue-worker           1/1          Ready              
kong                   1/1          Ready              ipaas.localdomain:8080
gateway                1/1          Ready              openfaas.ipaas.localdomain:8080
faas-idler             1/1          Ready              
swaggereditor          1/1          Ready              editor.localdomain:8080
konga                  1/1          Ready              konga.localdomain:8080

FUNCTIONS
NAME          STATUS       AVAILABILITY       URL
hello         1/1          Ready              ipaas.localdomain:8080/function/hello
```

## Docs

The documentation can be found in the following links:
- [Portuguese](https://vertigobr.gitlab.io/ipaas/docs/safira/visao_geral)
- English [TODO]

## Contribution

Pull requests/Merge Requests are welcome! Please open an issue first and discuss with us about the proposing changes and be sure to perform tests in a proper way.

## License
Safira is licensed under the [Apache License Version 2.0](https://github.com/vertigobr/safira/blob/master/LICENSE).
