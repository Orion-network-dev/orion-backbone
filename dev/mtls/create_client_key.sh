kid=$(openssl rand -hex 6)

mkdir -p client_key/${kid}/

openssl rand -hex 128 > client_key/${kid}/client.passphrase.txt

cat <<EOF  > client_key/${kid}/openssl.conf
[ req ]
prompt = no
distinguished_name = dn
req_extensions = req_ext

[ dn ]
CN = ${kid}.reg.orionet.re
emailAddress = ${kid}@reg.orionet.re
O = Orionet
OU = DEV PKI mTLS 
L = Saint-Pierre
ST = Réunion
C = FR

[ req_ext ]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = ${kid}.reg.orionet.re

EOF

openssl req -sha512 -nodes -newkey rsa:4096 \
    -keyout client_key/${kid}/client.key -passout file:client_key/${kid}/client.passphrase.txt -out client_key/${kid}/client.csr -extensions req_ext \
    -config <(cat client_key/${kid}/openssl.conf)

openssl x509 -req -CA ca/rootCA.crt -CAkey ca/rootCA.key  \
    -in client_key/${kid}/client.csr -passin file:ca/passphrase.txt -out client_key/${kid}/client.crt -days 365 -CAcreateserial

openssl x509 -in client_key/${kid}/client.crt -out client_key/${kid}/client.der -outform DER

#rm client_key/${kid}/client.csr
