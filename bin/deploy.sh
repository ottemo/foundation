#!/bin/bash

# location of foundation
HOME=/home/ottemo
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/foundation

if [ "$BRANCH" == 'develop' ]; then
    #GIT_COMMIT=$( echo "$COMMIT" | head -c 5 )
    # grab the latest code

    currentBranch=`ssh ottemo@$REMOTE_HOST "cd $SRCDIR && git symbolic-ref --quiet --short HEAD 2> /dev/null || git rev-parse --short HEAD 2> /dev/null || echo '(unknown)'"`
    echo ""
    echo "FOUNDATION BRANCH IS ${currentBranch}"

    if [ "$currentBranch" == 'develop' ]; then
        echo "GRAB THE LATEST CODE"
        ssh ottemo@$REMOTE_HOST "cd $SRCDIR && git stash && git checkout develop && git fetch --prune && git pull"
    fi

    echo ""
    echo "BUILD LOCALLY"
    ssh ottemo@$REMOTE_HOST "cd $SRCDIR && go get -t ./... && bash bin/make.sh -tags mongo"

    echo ""
    echo "BACKUP THE CURRENT BINARY AND PUT THE NEWLY BUILT BINARY INTO SERVICE"
    ssh ottemo@$REMOTE_HOST "sudo service ottemo stop && mv $HOME/foundation/foundation $HOME/foundation/foundation.bak"
    ssh ottemo@$REMOTE_HOST "cp $SRCDIR/foundation ~/foundation/foundation && sudo service ottemo start"
fi
