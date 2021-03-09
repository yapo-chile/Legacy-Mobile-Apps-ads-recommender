# ads-recommender

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/ads-recommender)](https://travis.mpi-internal.com/Yapo/ads-recommender)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/ads-recommender)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/ads-recommender?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/ads-recommender)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/ads-recommender?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/ads-recommender)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/ads-recommender)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/ads-recommender)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/ads-recommender)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/ads-recommender)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->

ads-recommender is the official golang microservice template for Yapo.

## How to run the service

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/ads-recommender.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd ads-recommender
  $ make start
  ```

* To get a list of available commands:

  ```
  $ make help
    Targets:
	  clone                Setup a new service repository based on ads-recommender
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

### GET  /recommendations/{carousel}/{listID}?params=[adParams]&limit=[adsLimit]&from=[fromIndex]
Returns recommended ads depending on the chosen carousel

#### Request
The `params` query parameter can contain any ad param name, for example `category, publisherType, estateType`, etc. This will return the same ads, with those fields included if available.

The `limit` query param indicates how many ads to return. For example a limit of 5 will only return 5 ads.

The `from` query param indicates from which index to return the ads. For example a from value of 1 means that the first recommended ad will be skipped, and the next ads will be returned.

The path variable `carousel` can be obtained from the file `resources/suggestion_params.json`. There reside the available carousels and their configurations.

#### Response

```javascript
200 OK
{
  "ads": [
    {
      "id": "4961183",
      "title": "Dodge Journey 2018",
      "price": 50000000,
      "currency": "$",
      "images": {},
      "url": "/arica_parinacota/dodge_journey_2018_4961183",
      "date": "2021-02-08 20:55:45"
    },
    {
      "id": "4961184",
      "title": "Dodge Journey 2018",
      "price": 50000000,
      "currency": "$",
      "images": {},
      "url": "/arica_parinacota/dodge_journey_2018_4961184",
      "date": "2021-02-08 20:55:45"
    },
    ...
  ]
}

//When there are no recommendations for the provided listID
204 No Content

```

#### Error response
```javascript
//When the listID is not valid
500 Internal Server Error
{
  "ErrorMessage": "get ad fails to get it, len: 0"
}

//When the carousel path variable is not valid
500 Internal Server Error
{
  "ErrorMessage": "invalid carousel: '{invalidCarousel}'"
}
```

### Contact
dev@schibsted.cl

## Kubernetes

Kubernetes and Helm have to be installed in your machine.
If you haven't done it yet, you need to create a secret to reach Artifactory.
`kubectl create secret docker-registry containers-mpi-internal-com -n <namespace> --docker-server=containers.mpi-internal.com --docker-username=<okta_username> --docker-password=<artifactory_api_key> --docker-email=<your_email>`

### Helm Charts

1. You need to fill out the ENV variables in the k8s/ads-recommender/templates/configmap.yaml file.
2. You should fill out the *tag*, and *host* under hosts to your namespace.
3. Add this host name to your /etc/hosts file with the correct IP address (127.21.5.11)
4. You run `helm install -n <name_of_your_release> k8s/ads-recommender`
5. Check your pod is running with `kubectl get pods`
6. If you want to check your request log `kubectl logs <name_of_your_pod>`
