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


TARNAME=gql-lint-$VERSION-$PLATFORM-$ARCH.tar.gz

echo Downloading version $VERSION for $PLATFORM-$ARCH
echo https://github.com/sketch-hq/gql-lint/releases/download/$VERSION/$TARNAME
echo ---



curl -SLJO https://github.com/sketch-hq/gql-lint/releases/download/$VERSION/$TARNAME

echo ---
echo Installing
echo ---

tar xvzf $TARNAME
rm $TARNAME

sudo mv gql-lint /usr/local/bin

echo Done!
