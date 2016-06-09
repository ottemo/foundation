#!/bin/bash

# location of foundation
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/foundation

if [ "$BRANCH" == 'develop' ]; then
    GIT_COMMIT=`echo $COMMIT | head -c 5`
    # grab the latest code
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && git checkout develop && git fetch --prune && git pull"
    # update packate dependencies
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && go get -t ./..."
    # restart the documenation service
    ssh ottemo@$REMOTE_HOST "sudo service godoc restart"
fi
