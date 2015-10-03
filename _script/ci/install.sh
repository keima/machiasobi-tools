#!/bin/sh
FILENAME=go_appengine_sdk_linux_amd64-1.9.27.zip
SDK_URL=https://storage.googleapis.com/appengine-sdks/featured/${FILENAME}

echo "installing libs..."

cd ~
curl -O ${SDK_URL} && unzip -q ${FILENAME}
if test $? -ne 0 ; then
    echo "installing libs error has occured!" 1>&2
    exit 1
fi

echo "installing libs fetched."
exit 0