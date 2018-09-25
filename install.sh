#!/bin/bash

NAME="yasd"
VERSION="0.1.0"
PREFIX="/usr/local"
GITHUB_USER="tzmfreedom"
TMP_DIR="/tmp"

set -ue

UNAME=$(uname -s)
if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
    echo "Sorry, OS not supported: ${UNAME}. Download binary from https://github.com/${USERNAME}/${NAME}/releases"
    exit 1
fi


if [ "${UNAME}" = "Darwin" ] ; then
  OS="darwin"

  OSX_ARCH=$(uname -m)
  if [ "${OSX_ARCH}" = "x86_64" ] ; then
    ARCH="amd64"
  else
    echo "Sorry, architecture not supported: ${OSX_ARCH}. Download binary from https://github.com/${USERNAME}/${NAME}/releases"
    exit 1
  fi
elif [ "${UNAME}" = "Linux" ] ; then
  OS="linux"

  LINUX_ARCH=$(uname -m)
  if [ "${LINUX_ARCH}" = "i686" ] ; then
    ARCH="386"
  elif [ "${LINUX_ARCH}" = "x86_64" ] ; then
    ARCH="amd64"
  else
    echo "Sorry, architecture not supported: ${LINUX_ARCH}. Download binary from https://github.com/${USERNAME}/${NAME}/releases"
    exit 1
  fi
fi

ARCHIVE_FILE=${NAME}-${VERSION}-${OS}-${ARCH}.tar.gz
BINARY="https://github.com/${GITHUB_USER}/${NAME}/releases/download/v${VERSION}/${ARCHIVE_FILE}"

cd $TMP_DIR
curl -sL -O ${BINARY}

tar xzf ${ARCHIVE_FILE}
mv ${OS}-${ARCH}/${NAME} ${PREFIX}/bin/${NAME}
chmod +x ${PREFIX}/bin/${NAME}
rm -rf ${OS}-${ARCH}
rm -rf ${ARCHIVE_FILE}
