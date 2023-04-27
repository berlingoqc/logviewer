#! /bin/bash
#
#
#

export SHA=$(git rev-parse --short HEAD)
export VERSION=${VERSION:-$SHA}


linuxTarget=(arm64)
macTarget=(arm64)


function build() {
        export GOOS=$1
        export GOARCH=$2

        export WD=./build/${GOOS}/${GOARCH}

        mkdir -p $WD

        go build -ldflags "-X git.tmaws.io/tmconnect/logexplorer/cmd.sha1ver=${VERSION}" \
                -o ${WD}/logexplorer
}


build "linux" "arm64"
build "linux" "amd64"
build "darwin" "arm64"
build "darwin" "amd64"

