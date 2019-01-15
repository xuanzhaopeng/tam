#!/usr/bin/env bash

set -e

echo $TRAVIS_REPO_SLUG
echo $1

docker build -t $TRAVIS_REPO_SLUG .
docker tag $TRAVIS_REPO_SLUG $TRAVIS_REPO_SLUG:$1
docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
docker push $TRAVIS_REPO_SLUG:$1