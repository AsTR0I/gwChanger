#!/bin/bash

# Directories
SRC_DIR="$(pwd)"
BIN_DIR="../public/bin"
BUILD_DIR="../public/builds"
CONFIG_FILE="./config.json"
SIPC_BIN="$(pwd)/sipc"

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

            # Check if the archive exists, and remove it if it does
            if [ -f "$BUILD_ARCHIVE" ]; then
                echo "⚠️ Archive already exists. Removing the old one."
                rm "$BUILD_ARCHIVE"
            fi

            echo "📦 Archiving to ${BUILD_ARCHIVE}..."
            # Create a new archive with the binary
           # Check if sipc exists and add it together with gwChanger
            if [ -f "$SIPC_BIN" ]; then
                # Create the archive with both gwChanger and sipc
                tar -czvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$BIN_NAME" "$(basename "$SIPC_BIN")"
                echo "📦 Both gwChanger and sipc added to the archive."
            else
                # If sipc does not exist, just add gwChanger
                tar -czvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$BIN_NAME"
                echo "⚠️ sipc not found, only gwChanger added to the archive."
            fi
        else
            echo "❌ Compilation error for ${PLATFORM}/${ARCH}"
        fi
    done
done

echo "✅ Process completed."
