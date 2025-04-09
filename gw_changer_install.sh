#!/bin/bash

# Welcome message
echo "🚀 Welcome, I am the gwChanger, starting installation..."

# Determining the OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)
TIME_STAMP=$(date +"old.%y-%m-%d-%H-%M-%S")
echo "🔍 OS detected: $OS"
echo "🔍 Architecture detected: $ARCH"

# Define the base URL for downloads
BASE_URL="https://raw.githubusercontent.com/AsTR0I/gwChanger/refs/heads/main/public/builds"
# SIPC_URL="https://raw.githubusercontent.com/AsTR0I/gwChanger/refs/heads/main/src/sipc"

if [ "$OS" == "Linux" ]; then
        INSTALL_PATH="/opt/gwChanger"
        BACKUP_DIR="$INSTALL_PATH/$TIME_STAMP"
        command -v curl >/dev/null 2>&1 || { echo "❌ curl utility is not installed. Install curl and try again."; exit 1; }
        DOWNLOAD_CMD="curl -L -o gwChanger.tar.gz"
         case "$ARCH" in
            "x86_64" | "amd64")
                BIN_URL="${BASE_URL}/linux/amd64/gwChanger.tar.gz"
                ;;
            "i386" | "i686")
                BIN_URL="${BASE_URL}/linux/386/gwChanger.tar.gz"
                ;;
            *)
                echo "❌ Unknown architecture for Linux: $ARCH. Exiting script."
                exit 1
                ;;
        esac
elif [ "$OS" == "FreeBSD" ]; then
    INSTALL_PATH="/opt/gwChanger"
    BACKUP_DIR="$INSTALL_PATH/$TIME_STAMP"
    command -v fetch >/dev/null 2>&1 || { echo "❌ fetch utility is not installed. Install fetch and try again."; exit 1; }
    DOWNLOAD_CMD="fetch -o gwChanger.tar.gz"
    case "$ARCH" in
        "x86_64" | "amd64")
            BIN_URL="${BASE_URL}/freebsd/amd64/gwChanger.tar.gz"
            ;;
        "i386")
            BIN_URL="${BASE_URL}/freebsd/386/gwChanger.tar.gz"
            ;;
        *)
            echo "❌ Unknown architecture for FreeBSD: $ARCH. Exiting script."
            exit 1
            ;;
    esac
else
    echo "❌ Only Linux and FreeBSD are supported. Exiting script."
    exit 1
fi

command -v tar >/dev/null 2>&1 || { echo "❌ tar utility is not installed. Install tar and try again."; exit 1; }

# Print the download link
echo "🔗 Download link for the archive: $BIN_URL"

# Download the archive
echo "📥 Downloading the archive..."
$DOWNLOAD_CMD "$BIN_URL" 2>&1 | tee -a install_log.txt

if [ $? -ne 0 ]; then
    echo "❌ Error while downloading the archive."
    exit 1
fi

# Create a backup folder
mkdir -p "$BACKUP_DIR"

# Moving old configs and logs
if [ -f "$INSTALL_PATH/config.json" ]; then
    echo "❗ Moving old config.json config to backup folder..."
    cp "$INSTALL_PATH/config.json" "$BACKUP_DIR/config.json"
    echo "✅ Old config moved to $BACKUP_DIR"
fi


# If the program is already installed, remove it:
if [ -f "$INSTALL_PATH/gwChanger" ]; then
    echo "❗ Program already installed. Removing old version..."
    # Stop the old version of the program (if it's running)
    PID=$(pgrep -f "$INSTALL_PATH/gwChanger")
    if [ -n "$PID" ]; then
        kill -9 "$PID"
        echo "✅ Old version stopped."
    fi
    # Remove the old version of the program
    rm "$INSTALL_PATH/gwChanger"
    echo "✅ Old version removed."
else
    echo "✅ No old version found, proceeding with installation."
fi

# Check if download was successful
if [ $? -ne 0 ]; then
    echo "❌ Error while downloading the archive."
    exit 1
fi

echo "✅ Archive successfully downloaded."

# Create the installation directory
if [ ! -d "$INSTALL_PATH" ]; then
    mkdir -p "$INSTALL_PATH"
fi

echo "📦 Extracting the archive..."
tar -xzvf gwChanger.tar.gz -C "$INSTALL_PATH" 2>&1 | tee -a install_log.txt
echo "✅ Archive successfully extracted."

# Check if the config exists
CONFIG_PATH="$INSTALL_PATH/config.json"
if [ ! -f "$CONFIG_PATH" ]; then
echo "❌ config.json config not found. Creating a new config..."
    cat <<EOF > "$CONFIG_PATH"
{
    "hostname_machine": "",
    "hosts": [
        { "hostname": "yandex.ru", "ip": "77.88.55.88" },
        { "hostname": "yandex.ru", "ip": "77.88.55.88" }
    ],
    "target_hostname": "voip.test voip.test2",
    "sipc_path": "",
    "mail": {
        "from": "",
        "to": "voip@cocobri.ru",
        "smtp_server": "",
        "smtp_server_port": ""
    }
}
EOF
    echo "✅ config.json created/updated."
else
    echo "✅ config.json already exists. Skipping creation."
fi

# Make the file executable
chmod +x "$INSTALL_PATH/gwChanger"
chmod +x "$INSTALL_PATH/sipc"
chmod -R 755 "$INSTALL_PATH"
echo "✅ Installation complete."

# Add to cron for auto-start
CRON_CMD="$INSTALL_PATH/gwChanger"

# Check if the command already exists in /etc/crontab
if ! grep -q "$CRON_CMD" /etc/crontab; then
    # Add the command to /etc/crontab
    echo "*/1 * * * * root $CRON_CMD" | tee -a /etc/crontab > /dev/null
        
    # Check that the command has actually been added
    if ! grep -q "$CRON_CMD" /etc/crontab; then
        echo "❌ Cron job added error."
    else
        echo "✅ Cron job added for auto-start."
    fi
else
    echo "🔄 Cron job already exists."
fi

echo "Cleaning up temporary files..."
rm gwChanger.tar.gz 
rm install_log.txt
echo "Cleanup complete, enjoy using the program!"