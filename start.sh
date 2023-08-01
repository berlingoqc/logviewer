#! /bin/bash

go build -o logviewer -gcflags "all=-N -l"

./logviewer -c ./config.json -i checkout-dev  query --logging-path ./logviewer.log