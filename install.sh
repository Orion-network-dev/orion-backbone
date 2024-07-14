#!/bin/bash

# Enable error handling
set -e

echo -e "Running apt-get update for accurate packages."
apt-get update -q &> /dev/null

# Check for the `jq` command to be available.
if ! command -v jq &> /dev/null
then
    echo "Ths jq command couldn't be found in your current $PATH... Trying to install it."
    apt-get install -yq jq &> /dev/null
    exit 1
fi

# Check for the `curl` command to be available.
if ! command -v curl &> /dev/null
then
    echo "Ths curl command couldn't be found in your current $PATH... Trying to install it."
    apt-get install -yq curl &> /dev/null
    exit 1
fi

JSON=$(curl -s https://api.github.com/repos/MatthieuCoder/OrionV3/releases/latest)
VERSION=$(echo $JSON | jq -r '.name')
NAME_PREDICATE="contains(\"$(dpkg --print-architecture).deb\")"
URL=$(echo $JSON | jq -r ".assets[] | select(.name | $NAME_PREDICATE) | .browser_download_url")

echo "Downloading version $VERSION for $(dpkg --print-architecture)..."
curl "$URL" -s -L -o "/tmp/orion.deb"
echo "Downloaded version $VERSION... Installing using APT"
apt install -q --allow-downgrades -y /tmp/orion.deb &> /dev/null
echo "Done. Cleaning up."
rm /tmp/orion.deb
