openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=FR/ST=Réunion/L=Le Tampon/O=Orion/OU=Registry/CN=reg.orionet.re" \
    -keyout cert.key  -out cert.crt
