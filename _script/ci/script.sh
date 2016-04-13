#!/bin/sh

echo "running tests and building..."

cd $TRAVIS_BUILD_DIR/src/machitools

$GO_APPENGINE_PATH/goapp test
if test $? -ne 0 ; then
    echo "go test is failed!" 1>&2
    exit 1
fi

$GO_APPENGINE_PATH/goapp build
if test $? -ne 0 ; then
    echo "go build is failed!" 1>&2
    exit 1
fi
cd -

cd $TRAVIS_BUILD_DIR/static
gulp build
if test $? -ne 0 ; then
    echo "gulp build is failed!" 1>&2
    exit 1
fi
