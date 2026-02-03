#!/bin/bash

# 项目名称
BINARY="portflow"
VERSION="1.0.0"
BUILD_DIR="bin"

# 创建编译目录
mkdir -p ${BUILD_DIR}

# 编译平台列表
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    OS=$(echo ${PLATFORM} | cut -d'/' -f1)
    ARCH=$(echo ${PLATFORM} | cut -d'/' -f2)
    OUTPUT="${BUILD_DIR}/${BINARY}-${OS}-${ARCH}"
    
    if [ "${OS}" == "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for ${OS}/${ARCH}..."
    GOOS=${OS} GOARCH=${ARCH} CGO_ENABLED=0 go build -ldflags "-s -w" -o ${OUTPUT} ./cmd/portflow
done

echo "Build complete. Binaries are in ${BUILD_DIR}/"
