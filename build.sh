#! /bin/bash
#
#
#

export SHA=$(git rev-parse --short HEAD)

export VERSION=${VERSION:-$SHA}


echo $VERSION
