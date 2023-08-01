#! /bin/bash
#
dlv attach $(pgrep logviewer) --listen=:2345 --headless --api-version=2 --log