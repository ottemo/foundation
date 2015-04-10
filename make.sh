#!/bin/bash

WORKDIR=`pwd`
OTTEMODIR="$(cd "$(dirname "$0")" && pwd)"
OTTEMOPKG="github.com/ottemo/foundation"

BRANCH=`git -C $OTTEMODIR rev-parse --abbrev-ref HEAD`
BUILD=`git -C $OTTEMODIR rev-list origin/develop --count`
DATE=`date`
TAGS=""

while test $# -gt 0; do
        case "$1" in
                -tags)
			shift
                        TAGS+=$1
			shift
                        ;;
                -wd|--work-dir)
			shift
			WORKDIR=$1
                        shift
                        ;;
		-od|--ottemo-dir)
			shift
                        WORKDIR=$1
			shift
                        ;;
		*)
			echo $1
			shift
			;;
        esac
done

cd $WORKDIR

LDFLAGS="-ldflags \""
LDFLAGS+="-X github.com/ottemo/foundation/app.buildDate '$DATE' "
LDFLAGS+="-X github.com/ottemo/foundation/app.buildTags '$TAGS' "
LDFLAGS+="-X github.com/ottemo/foundation/app.buildNumber '$BUILD' "
LDFLAGS+="-X github.com/ottemo/foundation/app.buildBranch '$BRANCH'"
LDFLAGS+="\""

if [ -z "$GOPATH" ]; then 
REPLACE="/src/$OTTEMOPKG"
export GOPATH="${OTTEMODIR/$REPLACE/}"
fi


if [ -n "$TAGS" ]; then
TAGS="-tags $TAGS"
fi

CMD="go get $TAGS $OTTEMOPKG"
eval $CMD

CMD="go build -a $TAGS $LDFLAGS $OTTEMOPKG"
eval $CMD
