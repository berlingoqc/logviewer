#! /bin/bash
#
#
#

export SHA=$(git rev-parse --short HEAD)
export VERSION=${VERSION:-$SHA}


echo "building version $VERSION"


linuxTarget=(arm64)
macTarget=(arm64)


function build() {
        export GOOS=$1
        export GOARCH=$2

        export WD=./build/

        mkdir -p $WD

        go build -ldflags "-X github.com/berlingoqc/logviewer/cmd.sha1ver=${VERSION}" \
                -o ${WD}/logviewer-${GOOS}-${GOARCH}
}


build "linux" "arm64"
build "linux" "amd64"
build "darwin" "arm64"
build "darwin" "amd64"

