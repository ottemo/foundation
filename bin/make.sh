#!/bin/bash

WORKDIR=`pwd`
OTTEMODIR="$(cd "$(dirname "$0")" && pwd)"
OTTEMOPKG="github.com/ottemo/foundation"

if [[ "$OSTYPE" == "darwin"*  ]] || [ -f /etc/alpine-release ] || [ -f /.dockerenv ] ; then
    AWK=gawk
else
    AWK=awk
fi

DATE=`${AWK} 'BEGIN{ print gensub(/(..)$/, ":\\\1", "g", strftime("%Y-%m-%dT%H:%M:%S%z")); exit }'`
TAGS=""
BUILD=`git -C $OTTEMODIR rev-list origin/develop --count`
BRANCH=`git -C $OTTEMODIR rev-parse --abbrev-ref HEAD`
HASH=`git -C $OTTEMODIR rev-parse --short --verify HEAD`

GOVERSION=`go version | ${AWK} '{print $3}' | awk '{sub(/go/,""); print}'`

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

LDFLAGS="-ldflags '"
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildDate=$DATE\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildTags=$TAGS\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildNumber=$BUILD\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildBranch=$BRANCH\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildHash=$HASH\" "
LDFLAGS+="'"

# need to convert GOVERSION string to number
if [ "`${AWK} "BEGIN{ if (($GOVERSION +0) < 1.5) print 1 }"`" == "1" ]; then
  LDFLAGS=${LDFLAGS//=/\" \"}
fi

if [ -z "$GOPATH" ]; then 
REPLACE="/src/$OTTEMOPKG"
export GOPATH="${OTTEMODIR/$REPLACE/}"
fi


if [ -n "$TAGS" ]; then
TAGS=$(echo $TAGS| sed 's/,/ /g')
TAGS="-tags '$TAGS'"
fi

# install glide to $GOPATH and project dependencies
go get github.com/Masterminds/glide
$GOPATH/bin/glide install

CMD="go build -a $TAGS $LDFLAGS $OTTEMOPKG"
eval CGO_ENABLED=0 GOOS=linux $CMD
