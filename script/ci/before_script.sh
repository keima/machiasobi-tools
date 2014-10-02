#!/bin/sh

echo "running goapp get to fetch dependencies..."

cd $TRAVIS_BUILD_DIR/machitools

echo $GO_APPENGINE_PATH

"$GO_APPENGINE_PATH/goapp" get

echo "dependencies fetched."
exit 0