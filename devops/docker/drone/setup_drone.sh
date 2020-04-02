#!/bin/bash

export GOPATH="$PWD/drone_workdir"
rm -rf $GOPATH
mkdir -p $GOPATH/src

TAG="v1.7.0"

cd $GOPATH/src
git clone https://github.com/drone/drone.git
cd drone
git checkout ${TAG}

export CGO_ENABLED=0
export GOOS=linux
go install -tags "nolimit" github.com/drone/drone/cmd/drone-agent
go install -tags "nolimit" github.com/drone/drone/cmd/drone-controller
go install -tags "nolimit" github.com/drone/drone/cmd/drone-server

rm -rf $GOPATH/src
sudo rm -rf $GOPATH/pkg
