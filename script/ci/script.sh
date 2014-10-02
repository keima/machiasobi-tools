#!/bin/sh

echo "running tests and building..."
cd ./github.com/keima/machitools/machitools
../../../../../go_appengine/goapp test
../../../../../go_appengine/goapp build
