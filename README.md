# pro-carousel

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/pro-carousel)](https://travis.mpi-internal.com/Yapo/pro-carousel)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/pro-carousel)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/pro-carousel?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/pro-carousel)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/pro-carousel?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/pro-carousel)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/pro-carousel)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/pro-carousel)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/pro-carousel)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/pro-carousel)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->

pro-carousel is the official golang microservice template for Yapo.

## A few rules

* pro-carousel was built following [Clean Architecture](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) so, please, familiarize yourself with it and let's code great code!

* pro-carousel has great [test coverage](https://quality-gate.mpi-internal.com/#/Yapo/pro-carousel) and [examples](https://github.mpi-internal.com/Yapo/pro-carousel/search?l=Go&q=func+Test&type=&utf8=%E2%9C%93) of how good testing can be done. Please honor the effort and keep your test quality in the top tier.

* pro-carousel is not a silver bullet. If your service clearly doesn't fit in this template, let's have a [conversation](mailto:dev@schibsted.cl)

* [README.md](README.md) is the entrypoint for new users of your service. Keep it up to date and get others to proof-read it.

## How to run the service

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/pro-carousel.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd pro-carousel
  $ make start
  ```

* To get a list of available commands:

  ```
  $ make help
    Targets:
	  clone                Setup a new service repository based on pro-carousel
	  info                 Display basic service info
	  help                 This help message
	  run                  Build and start the service in development mode (detached)
	  start                Build and start the service in development mode (attached)
	  build-dev            Build develoment docker image
	  docker-compose-%     Run docker compose commands with the project configuration
	  test                 Run tests and generate quality reports
	  build-test           Build test docker image
	  cover                Run tests and output coverage reports
	  coverhtml            Run tests and open report on default web browser
	  checkstyle           Run code linter and output report as text
	  docker-publish       Push docker image to containers.mpi-internal.com
	  helm-publish         Upload helm charts for deploying on k8s
	  build                Create production docker image
  ```

* If you change the code:

  ```
  $ make start
  ```

* How to run the tests

  ```
  $ make [cover|coverhtml]
  ```

* How to check format

  ```
  $ make checkstyle
  ```
  

## Creating a new service

* Create a repo for your new service on: https://github.mpi-internal.com/Yapo
* Rename your pro-carousel dir to your service name:
  ```
  $ mv pro-carousel YourService
  ```
* Update origin: 
  ```
  # https://help.github.com/articles/changing-a-remote-s-url/
  $ git remote set-url origin git@github.mpi-internal.com:Yapo/YourService.git
  ```

* Replace every pro-carousel reference to your service's name:
  ```
  $ git grep -l pro-carousel | xargs sed -i.bak 's/pro-carousel/yourservice/g'
  $ find . -name "*.bak" | xargs rm
  ```

* Go through the code examples and implement your service
  ```
  $ git grep -il fibonacci
  README.md
  cmd/pro-carousel/main.go
  pkg/domain/fibonacci.go
  pkg/domain/fibonacci_test.go
  pkg/interfaces/handlers/fibonacci.go
  pkg/interfaces/handlers/fibonacci_test.go
  pkg/interfaces/loggers/fibonacciInteractorLogger.go
  pkg/interfaces/repository/fibonacci.go
  pkg/interfaces/repository/fibonacci_test.go
  pkg/usecases/getNthFibonacci.go
  pkg/usecases/getNthFibonacci_test.go
  ```

* Enable TravisCI
  - Go to your service's github settings -> Hooks & Services -> Add Service -> Travis CI
  - Fill in the form with the credentials you obtain from https://travis.mpi-internal.com/profile/
  - Sync your repos and organizations on Travis
  - Make a push on your service
  - The push should trigger a build. If it didn't ensure that it is enabled on the travis service list
  - Enjoy! This should automatically enable quality-gate reports and a few other goodies

## Endpoints
### GET  /healthcheck
Reports whether the service is up and ready to respond.

> When implementing a new service, you MUST keep this endpoint
and update it so it replies according to your service status!

#### Request
No request parameters

#### Response
* Status: Ok message, representing service health

```javascript
200 OK
{
	"Status": "OK"
}
```

### GET  /fibonacci
Implements the Fibonacci Numbers with Clean Architecture

#### Request
{
	"n": int - Ask for the nth fibonacci number
}

#### Response

```javascript
200 OK
{
	"Result": int - The nth fibonacci number
}
```

#### Error response
```javascript
400 Bad Request
{
	"ErrorMessage": string - Explaining what went wrong
}
```

### GET  /user/basic-data?mail=[user_mail]
Returns the essential user data. It is in communication with the Profile Microservice. The main goal of this endpoint is to be used for a basic Pact Test.

#### Request

No additional parameters

#### Response

```javascript
200 OK
{
    "fullname": Full name of the user,
    "cellphone": The user´s cellphone,
    "gender": The user gender,
    "country": The country where the user lives (Currently only Chile is Available),
    "region": The region where the user lives,
    "commune": The commune where the user lives,
}
```

### Contact
dev@schibsted.cl

## Kubernetes

Kubernetes and Helm have to be installed in your machine.
If you haven't done it yet, you need to create a secret to reach Artifactory.
`kubectl create secret docker-registry containers-mpi-internal-com -n <namespace> --docker-server=containers.mpi-internal.com --docker-username=<okta_username> --docker-password=<artifactory_api_key> --docker-email=<your_email>`

### Helm Charts

1. You need to fill out the ENV variables in the k8s/pro-carousel/templates/configmap.yaml file.
2. You should fill out the *tag*, and *host* under hosts to your namespace.
3. Add this host name to your /etc/hosts file with the correct IP address (127.21.5.11)
4. You run `helm install -n <name_of_your_release> k8s/pro-carousel`
5. Check your pod is running with `kubectl get pods`
6. If you want to check your request log `kubectl logs <name_of_your_pod>`
