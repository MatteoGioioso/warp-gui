#!/usr/bin/env sh

set -e

sudo cp warp-gui.jpeg /usr/share/icons/Humanity/apps/32/warp-gui.jpeg
sudo cp warp-gui.desktop /usr/share/applications/warp-gui.desktop
sudo go build -o /usr/local/bin/warp-gui main.go