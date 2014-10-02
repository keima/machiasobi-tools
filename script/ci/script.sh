#!/bin/sh

echo "running tests and building..."
cd $TRAVIS_BUILD_DIR/machitools

$GO_APPENGINE_PATH/goapp test
$GO_APPENGINE_PATH/goapp build
