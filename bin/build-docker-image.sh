#!/bin/bash

# build and push foundation image to registry

for i in "$@"
do
case $i in
    version=*)
    version="${i#*=}"
    shift
    ;;
    indocker=*)
    indocker="${i#*=}"
    shift
    ;;
    *)
            # unknown option
    ;;
esac
done

indocker="${indocker:-false}"

MYDIR=$(cd `dirname ${BASH_SOURCE[0]}` && pwd)
FOUNDATIONREPO="$MYDIR/.."
cd $FOUNDATIONREPO

if ! [ -n "$version" ] ; then
  date=$(date +%Y%m%d-%H%M%S)
  IMAGE="gcr.io/ottemo-kube/foundation:${date}"
else
  IMAGE="gcr.io/ottemo-kube/foundation:$version"
fi
echo "use $IMAGE as image name"

if [ "$indocker" = "true" ] ; then
  echo "build image under docker container"
  echo "generate temporary Dockerfile"
  echo "FROM golang:1.6" >Dockerfile.temporary
  echo 'COPY . /go/src/github.com/ottemo/foundation' >>Dockerfile.temporary
  echo 'RUN cd /go/src/github.com/ottemo/foundation && bin/make.sh -tags mongo,redis' >>Dockerfile.temporary

  echo "build foundation in temporary docker container"
  docker build --file Dockerfile.temporary -t temp-foundation-golang .
  if [ $? -ne 0 ]; then
    echo "error in build foundation in temporary docker container"
    exit 2
  fi
  ID=$(docker run -d temp-foundation-golang)
  docker cp $ID:/go/src/github.com/ottemo/foundation/foundation .
else

  echo "build foundation executable with golang:1.6 docker image"
  echo docker run -v "$FOUNDATIONREPO":/go/src/github.com/ottemo/foundation -w /go/src/github.com/ottemo/foundation -e GOOS=linux -e CGO_ENABLED=0 golang:1.6 bin/make.sh -tags mongo,redis
  if [ $? -ne 0 ]; then
    echo "error in build foundation executable"
    exit 2
  fi
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
