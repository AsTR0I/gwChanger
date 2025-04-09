#!/bin/bash

# Directories
SRC_DIR="$(pwd)"
BIN_DIR="../public/bin"
BUILD_DIR="../public/builds"
CONFIG_FILE="./config.json"
SIPC_BIN="./sipc"

ARCHITECTURES=("amd64" "386")
PLATFORMS=("linux" "freebsd")
# We exclude FreeBSD 386, since Go does not support it
EXCLUDE_FREEBSD_386=true

for PLATFORM in "${PLATFORMS[@]}"; do
    for ARCH in "${ARCHITECTURES[@]}"; do
        # Skipping FreeBSD 386
        if [[ "$PLATFORM" == "freebsd" && "$ARCH" == "386" && "$EXCLUDE_FREEBSD_386" == true ]]; then
            echo "⚠️ Skipping FreeBSD 386 compilation (not supported)"
            continue
        fi

        # Create directories if they don't exist
        BIN_PATH="${BIN_DIR}/${PLATFORM}/${ARCH}"
        BUILD_PATH="${BUILD_DIR}/${PLATFORM}/${ARCH}"
        mkdir -p "$BIN_PATH" "$BUILD_PATH"

        # Define binary name
        BIN_NAME="gwChanger"
        OUTPUT_PATH="${BIN_PATH}/${BIN_NAME}"
        echo "🚀 Compiling for ${PLATFORM}/${ARCH}..."

        GOOS=$PLATFORM GOARCH=$ARCH go build -o "$OUTPUT_PATH" -ldflags "-s -w" "$SRC_DIR/gwChanger.go"
         if [ $? -eq 0 ]; then
            echo "✅ Compilation finished: ${OUTPUT_PATH}"

            BUILD_ARCHIVE="${BUILD_PATH}/gwChanger.tar.gz"
            echo "📦 Archiving to ${BUILD_ARCHIVE}..."
            tar -czvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$BIN_NAME"

            # Add the sipc program to the archive if it exists
            if [ -f "$SIPC_BIN" ]; then
                echo "📦 Adding sipc to the archive..."
                tar -rzvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$(basename "$SIPC_BIN")"
            else
                echo "⚠️ sipc program not found, skipping addition to the archive."
            fi
        else
            echo "❌ Compilation error for ${PLATFORM}/${ARCH}"
        fi
    done
done

echo "✅ Process completed."
