#!/bin/sh

echo "running goapp get to fetch dependencies..."
cd ./github.com/keima/machitools/machitools
../../../../../go_appengine/goapp get
echo "dependencies fetched."
exit 0