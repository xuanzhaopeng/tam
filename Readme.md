# Test account manager

Test account manager(tam) is a lightweight service to distribute and manage test account for parallel test.

> Scenarios: Run test cases in PROD environment, we cannot generate random PROD accounts in runtime, in this moment you could use TAM to manage all given available test accounts in PROD, and share them among all kinds of tests. 

It provides following features:
* Provide an un-allocated test account
* Release an allocated test account
* Each test account only could be allocated by one request
* Test account will be free after given timeout
* You could define the account with your own account structure

## Quick Start Guide

```bash
docker run -d -t -i -v /tmp/examples:/etc/tam \
 -p 6666:6666 \
 -e TAM_TIMEOUT='20s' \
 -e TAM_KEY='username' \
 --name tam tam:1.0.0
```

## Build locally
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
```