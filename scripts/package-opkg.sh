#!/usr/bin/env bash


dir="build/ipkbuild/lightkeeper"
if [ $(uname) = 'Darwin' ]; then
    tar=gtar
    sed=gsed
else
    tar=tar
    sed=sed
fi

# We need to replace hyphens with underscores because of
# OPKG version comparison: it splits by the last hyphen.
VERSION=$(git describe --tags --dirty | $sed 's/-/_/2g')
export VERSION="${VERSION#v}"

if [ -z "$VERSION" ]; then
    echo "Unable to determine version from git tag. Exiting" >&2
    exit 1
fi

set -ex

pushd "$dir"/control
envsubst < control.template > control
$tar --numeric-owner --group=0 --owner=0 -czf ../control.tar.gz ./*
popd

pushd "$dir"/data
$tar --numeric-owner --group=0 --owner=0 -czf ../data.tar.gz ./*
popd

pushd "$dir"
$tar --numeric-owner --group=0 --owner=0 -z -cf ../lightkeeper-$VERSION.ipk ./debian-binary ./data.tar.gz ./control.tar.gz
popd

