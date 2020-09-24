# __SERVICE__

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/__SERVICE__)](https://travis.mpi-internal.com/Yapo/__SERVICE__)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/__SERVICE__)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/__SERVICE__?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/__SERVICE__)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/__SERVICE__?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/__SERVICE__)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/__SERVICE__)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/__SERVICE__)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/__SERVICE__)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/__SERVICE__)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->

__SERVICE__ needs a description here.

## Checklist: Is my service ready?

* [ ] Configure your github repository
  - Open https://github.mpi-internal.com/Yapo/__SERVICE__/settings
  - Features: Wikis, Restrict editing, Issues, Projects
  - Merge button: Only allow merge commits
  - Open https://github.mpi-internal.com/Yapo/pro-carousel/settings/branches
  - Default branch: master
  - Protected branches: choose master
  - Protect this branch
    + Require pull request reviews
    + Require status checks before merging
      - Require branches to be up to date
      - Quality gate code analysis
      - Quality gate coverage
      - Travis-ci
    + Include administrators
* [ ] Enable TravisCI
  - Go to your service's github settings -> Hooks & Services -> Add Service -> Travis CI
  - Fill in the form with the credentials you obtain from https://travis.mpi-internal.com/profile/
  - Sync your repos and organizations on Travis
  - Create a pull request and make a push on it
  - The push should trigger a build. If it didn't, ensure that it is enabled on the travis service list
  - Enjoy! This should automatically enable quality-gate reports and a few other goodies
* [ ] Get your first PR merged
  - Master should be a protected branch, so the only way to get commits there is via pull request
  - Once the travis build is ok, and you got approval merge it back to master
  - This will allow for the broken badges on top of this readme to display correctly
  - Should them not display after some time, please report it
* [ ] Enable automatic deployment
  - Have your service created and deployed on a stack on Rancher
  - Modify `rancher/deploy/*.json` files to reflect new names
  - Follow the instructions on https://github.mpi-internal.com/Yapo/rancher-deploy
* [ ] Create Helm Charts for Kubernetes deploy
  - Create a new Chart with `helm create k8s/__SERVICE__` cmd
  - Copy configmap.yaml from k8s/pro-carousel/templates/ and change pro-carousel to your __SERVICE__ name.
  - In the k8s/__SERVICE__/deployment.yaml file:
      + Add `/healthcheck` value to livenessProbe and readinessProbe section
      + Copy imagePullSecrets, annotations and envFrom section from pro-carousel example deploment.yaml and change the names to your service name
  - Delete pro-carousel chart
* [ ] Delete this section
  - It's time for me to leave, I've done my part
  - It's time for you to start coding your new service and documenting your endpoints below
  - Seriously, document your endpoints and delete this section

## How to run __SERVICE__

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/__SERVICE__.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd __SERVICE__
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

## Contact
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
