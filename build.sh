#!/bin/bash
# -------------------------------------------
#  /u01/devel/gcsconntest/build.sh
#  v1.0.0xg  2025/01/10  XDG
# -------------------------------------------
# Prereqs: golang 1.23+
# Syntax: build.sh 0.1.0xg-20250110

REPO_DEST="library/cloud"
APP_NAME="gcsconntest"
APP_VERSION="$1"
if [ -z "${APP_VERSION}" ];then
  echo "Missing argument: application version, aborting"
  exit 1
fi
ARC_NAME="${APP_NAME}-${APP_VERSION%xg*}"

mkdir -p bin archive distrib
[ -s "bin/${APP_NAME}" ] && mv -fv bin/${APP_NAME} archive/${APP_NAME}.$(date '+%Y%m%d-%H%M%S')

[ ! -s go.mod ] && go mod init mckesson/${APP_NAME}
go mod tidy
go build -v -ldflags "-s -w -X main.Version=${APP_VERSION}" -o ./bin/
chmod 755 ./bin/${APP_NAME}

./bin/${APP_NAME} -version
echo "Generating binary and development archives in [distrib/]"
tar -Jcf distrib/${ARC_NAME}-amd64.tar.xz bin/${APP_NAME}
tar -Jcf distrib/${ARC_NAME}-devel-amd64.tar.xz go.* *.go *.sh

echo "Copying binary archive to repository"
cp -afv distrib/${ARC_NAME}-amd64.tar.xz ${REPO_BASE}/${REPO_DEST}/
