#!/bin/sh

echo "running goapp get to fetch dependencies..."

# golang
cd $TRAVIS_BUILD_DIR/machitools
echo $GO_APPENGINE_PATH
"$GO_APPENGINE_PATH/goapp" get

# static
cd ./static
bower install
npm install

echo "dependencies fetched."
exit 0
