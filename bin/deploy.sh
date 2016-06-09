#!/bin/bash

# location of foundation
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/foundation

if [ "$BRANCH" == 'develop' ]; then
    GIT_COMMIT=`echo $COMMIT | head -c 5`
    # grab the latest code
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && git stash && git checkout develop && git fetch --prune && git pull"
    # build locally after successful merge to develop
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && go get -t ./... && bash bin/make.sh -tags mongo"
    # backup the current binary and put the newly built binary into service
    ssh ottemo@$REMOTE_HOST "sudo /etc/init.d/ottemo stop && cp $SRCDIR/foundation ~/foundation/foundation && sudo /etc/init.d/ottemo start"
fi
