#!/bin/bash
#
#go install github.com/fyne-io/fyne-cross@latest
#
~/go/bin/fyne-cross windows -arch=amd64
~/go/bin/fyne-cross windows -arch=arm64
~/go/bin/fyne-cross linux -arch=amd64
~/go/bin/fyne-cross linux -arch=arm64
mv fyne-cross/dist/linux-arm64/fyne-mines.tar.xz fyne-cross/dist/fyne-mines-arm64.tar.xz
mv fyne-cross/dist/linux-amd64/fyne-mines.tar.xz fyne-cross/dist/fyne-mines-amd64.tar.xz
mv fyne-cross/dist/windows-arm64/fyne-mines.exe.zip fyne-cross/dist/fyne-mines-arm64.exe.zip
mv fyne-cross/dist/windows-amd64/fyne-mines.exe.zip fyne-cross/dist/fyne-mines-amd64.exe.zip