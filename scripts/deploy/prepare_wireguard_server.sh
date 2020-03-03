#!/bin/bash

set -xe

sudo bash -c 'echo "deb http://deb.debian.org/debian/ unstable main" > /etc/apt/sources.list.d/unstable.list'
sudo bash -c "printf 'Package: *\nPin: release a=unstable\nPin-Priority: 90\n' > /etc/apt/preferences.d/limit-unstable"
sudo apt-get update && sudo apt-get install -y  --no-install-recommends --no-install-suggests wireguard iproute2 linux-headers-$(uname -r)
sudo sysctl -w net.ipv4.ip_forward=1
