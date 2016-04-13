#!/bin/sh

echo "running goapp get to fetch dependencies..."

# golang
cd $TRAVIS_BUILD_DIR/src/machitools
"$GO_APPENGINE_PATH/goapp" get
if test $? -ne 0 ; then
    echo "goapp get is failed!" 1>&2
    exit 1
fi
cd -

# static
cd $TRAVIS_BUILD_DIR/static
bower install
if test $? -ne 0 ; then
    echo "bower install is failed!" 1>&2
    exit 1
fi

npm install
if test $? -ne 0 ; then
    echo "npm install is failed!" 1>&2
    exit 1
fi
cd -

echo "dependencies fetched."
exit 0