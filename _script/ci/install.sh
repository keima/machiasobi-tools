#!/bin/sh
FILENAME=go_appengine_sdk_linux_amd64-1.9.27.zip
SDK_URL=https://storage.googleapis.com/appengine-sdks/featured/${FILENAME}

echo "installing libs..."

cd ~
curl -O ${SDK_URL} && unzip -q ${FILENAME}

echo "installing libs fetched."
exit 0
