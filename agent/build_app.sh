#!/bin/bash
set -e

APP="SoftwareUpdateHelper.app"
CONTENTS="$APP/Contents"
MACOS="$CONTENTS/MacOS"

make clean && make

mkdir -p "$MACOS"
cp agent "$MACOS/agent"
cp Info.plist "$CONTENTS/Info.plist"

echo "[+] $APP built"
echo "[+] Double-click to run: open $APP"
