# Test account manager

[![codecov](https://codecov.io/gh/xuanzhaopeng/tam/branch/master/graph/badge.svg)](https://codecov.io/gh/xuanzhaopeng/tam)
[![Build Status](https://travis-ci.org/xuanzhaopeng/tam.svg?branch=master)](https://travis-ci.org/xuanzhaopeng/tam)

Test account manager(tam) is a lightweight service to distribute and manage test account for parallel test.

> Scenarios: Run test cases in PROD environment, we cannot generate random PROD accounts in runtime, in this moment you could use TAM to manage all given available test accounts in PROD, and share them among all kinds of tests. 

It provides following features:
* Provide a free test account
* Release an allocated test account
* Each test account only could be allocated by one request
* It will release test account after given timeout
* You could define the account with your own account structure
* You could filter your account by one or multiple combined quires, which also support regex

## Quick Start Guide

Prepare accounts.json

```bash
mkdir -p /tmp/tam

$ cat /tmp/tam/accoounts.json
[
  {
    "username": "tester1",
    "password": "pass1",
    "refreshToken": "rt1",
    "region": "US",
    "data": {
      "a1": 1,
      "a2": 2
    }
  },
  {
    "username": "tester2",
    "password": "pass1",
    "refreshToken": "rt1",
    "region": "NL",
    "data": {
      "a1": 1,
      "a2": 3
    }
  }
]
```

Run TAM in docker

```bash
# TAM_KEY: the unique key of your account structure
# TAM_TIMEOUT: the given timeout for auto release

docker run -d -t -i -v /tmp/tam:/etc/tam \
 -p 6666:6666 \
 -e TAM_TIMEOUT='20s' \
 -e TAM_KEY='username' \
 --name tam tam:latest-release
```