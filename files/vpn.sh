#!/bin/bash
set -e

START="$1"
OPENVPN_CONF="$2"
echo "starting... $START"

if [ -f "$OPENVPN_CONF" ]; then
  sudo openvpn "$OPENVPN_CONF" | tee "$HOME/openvpn.log" &
  while [ `tail "$HOME/openvpn.log" | grep "Initialization Sequence Completed" | wc -l` == "0" ];
  do
    echo "waiting for openvpn to start..."
    sleep 2
  done
else
  echo "no openvpn config file found"
fi

## return traffic that goes through the VPN
gw=$(ip route | awk '/default/ {print $3}')
sudo ip route add to ${LOCAL_NETWORK} via $gw dev eth0

# sudo privoxy --no-daemon
sudo privoxy

$START