#!/usr/bin/env bash

# vars for script
PWD=$PWD
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/foundation

cd $SRCDIR
# Build Foundation
echo "BUILDING FOUNDATION"
./bin/make.sh -tags mongo

# stop the foundation service
echo "STOP FOUNDATION SERVICE"
while sudo service ottemo stop >/dev/null 2>&1; do
    echo "warning: foundation is still running"
done
echo "info: foundation has terminated"

echo "DEPLOYING FOUNDATION"
# Backup binaries and restart foundation
cp ~/foundation/foundation ~/foundation/foundation.bak
cp foundation ~/foundation/
sudo service ottemo start

# restore PWD
cd $PWD

echo "DEPLOY DONE"
