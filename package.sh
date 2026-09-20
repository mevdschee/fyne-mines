#!/bin/bash
#
#go install github.com/fyne-io/fyne-cross@latest
#
set -e
# darwin needs a copy of the macOS SDK, see the "OSX build" section in README.md
SDK=${MACOSX_SDK:-$HOME/SDKs/MacOSX12.3.sdk}
# zig does not apply the sysroot to the framework search path that fyne-cross
# passes, so point /System/Library/Frameworks at the mounted SDK instead
DARWIN_IMAGE=fyne-cross-darwin-sdk
docker build -q -t $DARWIN_IMAGE - <<'EOF'
FROM fyneio/fyne-cross-images:darwin
RUN mkdir -p /System/Library && ln -s /sdk/System/Library/Frameworks /System/Library/Frameworks
EOF
~/go/bin/fyne-cross windows -arch=amd64
~/go/bin/fyne-cross windows -arch=arm64
~/go/bin/fyne-cross linux -arch=amd64
~/go/bin/fyne-cross linux -arch=arm64
~/go/bin/fyne-cross darwin -arch=amd64 -app-id com.tqdev.fyne-mines -macosx-sdk-path "$SDK" -image $DARWIN_IMAGE
~/go/bin/fyne-cross darwin -arch=arm64 -app-id com.tqdev.fyne-mines -macosx-sdk-path "$SDK" -image $DARWIN_IMAGE
mv fyne-cross/dist/linux-arm64/fyne-mines.tar.xz fyne-cross/dist/fyne-mines-arm64.tar.xz
mv fyne-cross/dist/linux-amd64/fyne-mines.tar.xz fyne-cross/dist/fyne-mines-amd64.tar.xz
mv fyne-cross/dist/windows-arm64/fyne-mines.exe.zip fyne-cross/dist/fyne-mines-arm64.exe.zip
mv fyne-cross/dist/windows-amd64/fyne-mines.exe.zip fyne-cross/dist/fyne-mines-amd64.exe.zip
(cd fyne-cross/dist/darwin-arm64 && zip -qry ../fyne-mines-arm64.app.zip fyne-mines.app)
(cd fyne-cross/dist/darwin-amd64 && zip -qry ../fyne-mines-amd64.app.zip fyne-mines.app)
