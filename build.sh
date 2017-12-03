#!/usr/bin/env bash

GOOS=linux
GOARCH=amd64


build() {
    echo "building GOOS=${GOOS};GOARCH=${GOARCH} ..."
    local suffix=""

    if [ ${GOOS} == windows ]
    then
        suffix=".exe"
    fi

    local output="cquant_${GOOS}_${GOARCH}${suffix}"
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${output}
}

usage() {
    local prog="./build.sh"
    echo "Usage: $prog -o [linux|windows|darwin|freebsd] -a [386|amd64|arm]"
    echo "       $prog -h for help."
    exit 1
}

while getopts "a:o:h" arg #选项后面的冒号表示该选项需要参数
do
    case $arg in
         a) GOARCH=$OPTARG;;
         o) GOOS=$OPTARG;;
         h) usage;;
    esac
done

build

exit 0