#!/bin/sh

# GAE_OAUTH is provided by travis.yml
export GAE_DIR=$GO_APPENGINE_PATH
export APP_DIR=$TRAVIS_BUILD_DIR/machitools

cd $APP_DIR
rm app.yaml
cp app-prod.yaml app.yaml
python $GAE_DIR/appcfg.py --oauth2_refresh_token=$GAE_OAUTH update ./
