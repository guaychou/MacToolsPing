#!/bin/bash
set -e
mkdir -p ~/Applications/PingService.app/Contents/MacOS
go build -O PingService
cp PingService ~/Applications/PingService.app/Contents/MacOS
cat << EOF > ~/Applications/PingService.app/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>NSHighResolutionCapable</key>
        <string>True</string>
        <!-- avoid showing the app on the Dock -->
        <key>LSUIElement</key>
        <string>1</string>
</dict>
</plist>
EOF
echo "Build done . . . "