#!/bin/bash
#
#go install github.com/fyne-io/fyne-cross@latest
#
~/go/bin/fyne-cross windows -arch=amd64
~/go/bin/fyne-cross windows -arch=arm64
~/go/bin/fyne-cross linux -arch=amd64
~/go/bin/fyne-cross linux -arch=arm64