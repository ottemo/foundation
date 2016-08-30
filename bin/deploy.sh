#!/bin/bash

# location of foundation
HOME=/home/ottemo
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/foundation

if [ "$BRANCH" == 'develop' ]; then
    #GIT_COMMIT=$( echo "$COMMIT" | head -c 5 )
    # grab the latest code
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && git stash && git checkout develop && git fetch --prune && git pull"
    # build locally after successful merge to develop
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && go get -t ./... && bash bin/make.sh -tags mongo"
    # backup the current binary and put the newly built binary into service
    ssh ottemo@$REMOTE_HOST "sudo service ottemo stop && mv $HOME/foundation/foundation $HOME/foundation/foundation.bak"
    ssh ottemo@$REMOTE_HOST "cp $SRCDIR/foundation ~/foundation/foundation && sudo service ottemo start"
fi
