#!/bin/bash
set -e

VERSION=$1

arch=$(uname -m)
if [[ $arch == x86_64* ]]; then
  ARCH="amd64"
elif  [[ $arch == arm* ]]; then
  ARCH="arm64"
else
  echo "Incompatible architecture: $arch"
  exit 0
fi

platform=$(uname -s)
if [[ $platform == Darwin ]]; then
  PLATFORM="darwin"
else
  PLATFORM="linux"
fi

TARNAME=gql-lint-$PLATFORM-$ARCH.tar.gz

if [[ "$VERSION" == "latest" ]]; then
  URL="https://github.com/sketch-hq/gql-lint/releases/latest/download/$TARNAME"
else
  URL="https://github.com/sketch-hq/gql-lint/releases/download/$VERSION/$TARNAME"
fi

echo Downloading version $VERSION for $PLATFORM-$ARCH
echo $URL
echo ---

curl -SLJO $URL

echo ---
echo Installing (in /usr/local/bin)
echo ---

tar xvzf $TARNAME
rm $TARNAME

echo "sudo required for moving gql-lint to /usr/local/bin"
sudo mv gql-lint /usr/local/bin

echo Done!
