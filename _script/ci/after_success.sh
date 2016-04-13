#!/bin/sh

# GAE_OAUTH is provided by travis.yml
export GAE_DIR=$GO_APPENGINE_PATH

cd $TRAVIS_BUILD_DIR/src
python $GAE_DIR/appcfg.py --oauth2_refresh_token=$GAE_OAUTH update ./
if test $? -ne 0 ; then
    echo "Deploy is failed!" 1>&2
    exit 1
fi
