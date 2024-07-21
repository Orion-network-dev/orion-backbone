#!/bin/sh

systemctl stop oriond
systemctl enable --now oriond

echo "***WARNING***"
echo "This Orion package needs to have authentication certificates installed."
echo "***WARNING***"
