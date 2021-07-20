#!/bin/bash
set -e
set -x

if [[ $# -lt 1 ]]; then
    echo "Must specify a filename!"
    exit 1;
fi

FULLPATH=$(readlink -f $1)
BASENAME=$(basename ${FULLPATH})
DIRNAME=$(dirname ${FULLPATH})
cd ${DIRNAME}

if [ $CIRCLE_TAG ] && [ "$CIRCLE_TAG" != "" ]
then
    VERSION=$(echo $CIRCLE_TAG | awk 'match($0, /^v([0-9]+(\.[0-9]+)*)/, a) {print a[1]}')
    REMAINDER=$(echo $CIRCLE_TAG | awk 'match($0, /^v([0-9]+(\.[0-9]+)*)/, a) {print substr($0, RSTART+RLENGTH)}')
    RELEASE=${REMAINDER#-*}
    RELEASE=$(echo $RELEASE | sed 's/-/./')
    if [ "$VERSION" == "" ]
    then
        VERSION=0.2
        COMMITCOUNT=$(git rev-list --count --no-merges HEAD)
        RELEASE=$COMMITCOUNT
    else
        VERSION="${VERSION:-0}"
        RELEASE="${RELEASE:-1}"
    fi
elif [ "$CIRCLE_BRANCH" == "master" ]
then
    VERSION=0.1
    COMMITCOUNT=$(git rev-list --count --no-merges HEAD)
    RELEASE="$COMMITCOUNT+${CIRCLE_SHA1:0:8}"
elif [ -n $(type -p git) ] && [ -d $DIRNAME/../.git ]; then
    VERSION=$(grep -e "^Version" ${FULLPATH} | awk '{print $2}')
    COMMITCOUNT=$(git rev-list --count --no-merges HEAD ^origin/master)
    RELEASE=$(( $COMMITCOUNT + 1 ))
else
    VERSION=0.0
    RELEASE="0"
fi

echo "VERSION=$VERSION"
echo "RELEASE=$RELEASE"
