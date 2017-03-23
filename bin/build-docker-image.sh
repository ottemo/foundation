#!/bin/bash

# build and push foundation image to registry

MYDIR=$(cd `dirname ${BASH_SOURCE[0]}` && pwd)
FOUNDATIONREPO="$MYDIR/.."
cd $FOUNDATIONREPO

date=$(date +%Y%m%d-%H%M%S)
branch=$(git branch| awk '{print $2}')
branch=$(echo $branch|sed 's/ /_/g')
IMAGE="gcr.io/ottemo-kube/foundation:${branch}-${date}"

echo "build foundation executable with golang:1.6 docker image"
docker run -v "$FOUNDATIONREPO":/go/src/github.com/ottemo/foundation -w /go/src/github.com/ottemo/foundation -e GOOS=linux -e CGO_ENABLED=0 golang:1.6 bin/make.sh -tags mongo,redis
if [ $? -ne 0 ]; then
  echo "error in build foundation executable"
  exit 2
fi

echo "build alpine based foundation container"
docker build -t $IMAGE .
if [ $? -ne 0 ]; then
  echo "error in build foundation alpine based container"
  exit 2
fi

gcloud docker -- push $IMAGE
if [ $? -ne 0 ]; then
  echo "error in push image"
  exit 2
fi
