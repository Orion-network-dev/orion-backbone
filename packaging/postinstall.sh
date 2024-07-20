#!/bin/sh

systemctl stop oriond
systemctl enable --now oriond

echo -e "\t\t ***WARNING***"
echo -e "\t\t This Orion package needs to have authentication certificates installed."
echo -e "\t\t ***WARNING***"
