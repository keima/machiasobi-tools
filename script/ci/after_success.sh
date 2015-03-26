#!/bin/sh

# GAE_OAUTH is provided by travis.yml
export GAE_DIR=$GO_APPENGINE_PATH
export APP_DIR=$TRAVIS_BUILD_DIR/machitools

echo "PR# $TRAVIS_PULL_REQUEST"
cd $APP_DIR
python $GAE_DIR/appcfg.py --oauth2_refresh_token=$GAE_OAUTH update app-prod.yaml
