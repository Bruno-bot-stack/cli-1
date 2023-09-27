#!/bin/bash

# Builds cli $1 version $2 for all platforms and packages them for release.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1
VSN=$2

if [ -z "$CLI" ]; then
    echo "require cli as first parameter" 2>&1
    exit 1
fi

if [ -z "$VSN" ]; then
    echo "require version as second parameter" 2>&1
    exit 1
fi

if [ ! -d "$CLI" ]; then
    echo "\"$CLI\" is not a cli" 2>&1
    exit 1
fi

# build
rm -f $ROOT/build/${CLI}_${VSN}*
$ROOT/${CLI}/build-all-platforms.sh "$VSN"

# archive
cd $ROOT/build
for t in ${CLI}_${VSN}_* ; do
    if [[ $t == ${CLI}_*_windows_* ]]; then
        zip -q ${t/.exe/.zip} $t
    else
        tar -czf ${t}.tar.gz $t
    fi
done
cd ..

# dist: debian
rm -rf $ROOT/${CLI}/dist/usr
mkdir -p $ROOT/${CLI}/dist/usr/local/bin
cp $ROOT/build/${CLI}_${VSN}_linux_amd64 $ROOT/${CLI}/dist/usr/local/bin/${CLI}
dpkg-deb -Zgzip --build ${ROOT}/${CLI}/dist build/${CLI}_${VSN}.deb
