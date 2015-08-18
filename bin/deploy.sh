#!/bin/bash

# location of foundation
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/foundation

if [ "$BRANCH" == 'develop' ]; then
    GIT_COMMIT=`echo $COMMIT | head -c 5`
    # grab the latest code
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && git fetch --prune && git pull"
    # build locally after successful merge to develop
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && ./bin/make.sh -tags mongo"
    # backup the current binary and put the newly built binary into service
    ssh ottemo@$REMOTE_HOST "sudo /etc/init.d/ottemo stop && cp ~/foundation/foundation ~/foundation/backup/foundation-$(date +%Y%m%d)"
    ssh ottemo@$REMOTE_HOST "cp $SRCDIR/foundation ~/foundation/foundation && sudo /etc/init.d/ottemo start"
fi
