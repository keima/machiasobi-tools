#!/bin/sh
SDK_URL=https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_linux_amd64-1.9.12.zip
FILENAME=go_appengine_sdk_linux_amd64-1.9.12.zip

echo "installing libs..."
cd ..
curl -O ${SDK_URL} && unzip -q ${FILENAME}
echo "installing libs fetched."
exit 0
