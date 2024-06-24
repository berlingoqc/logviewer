#! /bin/bash

WIREMOCK_JAR=${WIREMOCK_JAR:-wiremock.jar}
PORT=9091

java -jar $WIREMOCK_JAR --port $PORT --root-dir ./mock/wiremock --global-response-templating