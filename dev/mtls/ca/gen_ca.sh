#!/bin/sh

# Generates th root CA
openssl req \
    -x509 -sha256 -days 3652 -newkey rsa:2048 -keyout rootCA.key -out rootCA.crt -passout file:passphrase.txt \
    -subj "/C=FR/ST=Reunion Island/L=Saint-Pierre/O=Orion Network/OU=Registry CA/CN=reg.orionet.re"

