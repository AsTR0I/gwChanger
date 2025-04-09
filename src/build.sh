#!/bin/bash

# Директории
SRC_DIR="./"
BIN_DIR="../public/bin"
BUILD_DIR="../public/builds"
CONFIG_FILE="./config.json"  # Путь к файлу конфигурации
SIPC_BIN="./sipc"  # Путь к sipc

ARCHITECTURES=("amd64" "386")
PLATFORMS=("linux" "freebsd")
# Исключаем FreeBSD 386, так как Go не поддерживает его
EXCLUDE_FREEBSD_386=true

for PLATFORM in "${PLATFORMS[@]}"; do
    for ARCH in "${ARCHITECTURES[@]}"; do
        # Пропускаем FreeBSD 386
        if [[ "$PLATFORM" == "freebsd" && "$ARCH" == "386" && "$EXCLUDE_FREEBSD_386" == true ]]; then
            echo "⚠️ Пропускаем компиляцию FreeBSD 386 (не поддерживается)"
            continue
        fi

         # Создаём директории, если их нет
        BIN_PATH="${BIN_DIR}/${PLATFORM}/${ARCH}"
        BUILD_PATH="${BUILD_DIR}/${PLATFORM}/${ARCH}"
        mkdir -p "$BIN_PATH" "$BUILD_PATH"

        # Определяем имя бинарника
        BIN_NAME="gwChanger"
        OUTPUT_PATH="${BIN_PATH}/${BIN_NAME}"
        echo "🚀 Компиляция для ${PLATFORM}/${ARCH}..."

        GOOS=$PLATFORM GOARCH=$ARCH go build -o "$OUTPUT_PATH" -ldflags "-s -w -buildvcs=false" "$SRC_DIR"
         if [ $? -eq 0 ]; then
            echo "✅ Компиляция завершена: ${OUTPUT_PATH}"

            BUILD_ARCHIVE="${BUILD_PATH}/gwChanger.tar.gz"
            echo "📦 Архивирование в ${BUILD_ARCHIVE}..."
            tar -czvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$BIN_NAME"

            # Добавляем в архив программу sipc, если она существует
            if [ -f "$SIPC_BIN" ]; then
                echo "📦 Добавление sipc в архив..."
                tar -rzvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$(basename "$SIPC_BIN")"
            else
                echo "⚠️ Программа sipc не найдена, пропускаем добавление в архив."
            fi
        else
            echo "❌ Ошибка компиляции для ${PLATFORM}/${ARCH}"
        fi
    done
done

echo "✅ Процесс завершен."