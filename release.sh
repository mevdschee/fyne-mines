#!/bin/bash
set -e
cd $(dirname $0)
if [ ! -f fyne-cross/dist/fyne-mines-amd64.app.zip ]; then
  echo "no binaries found, run package.sh first" >&2
  exit 1
fi
newTag=v$(grep -oP '(?<=^Version = ")[^"]+' FyneApp.toml)
#gh release delete $newTag
gh release create $newTag --title "${1:-$newTag}" \
  fyne-cross/dist/fyne-mines-amd64.exe.zip \
  fyne-cross/dist/fyne-mines-arm64.exe.zip \
  fyne-cross/dist/fyne-mines-amd64.tar.xz \
  fyne-cross/dist/fyne-mines-arm64.tar.xz \
  fyne-cross/dist/fyne-mines-amd64.app.zip \
  fyne-cross/dist/fyne-mines-arm64.app.zip
