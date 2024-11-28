package internal

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"reflect"
)

var (
	AuthorityPath = flag.String("tls-authority-path", "/etc/oriond/ca.crt", "Path to the certificate authority file")
	TLSPath       = flag.String("tls-path", "/etc/oriond/identity.key", "Path to the certificate authority file")
)

func LoadPemFile() (*ecdsa.PrivateKey, []*x509.Certificate) {
	bytes, err := os.ReadFile(*TLSPath)
	if err != nil {
		panic(err)
	}

	var privateKey *ecdsa.PrivateKey
	var chain []*x509.Certificate

	for block, rest := pem.Decode(bytes); block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "PRIVATE KEY" {
			pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				panic("cannot read a private key in the pem file")
			}

			if ec, ok := pk.(*ecdsa.PrivateKey); !ok {
				panic(fmt.Errorf("invalid key type: %s", reflect.TypeOf(pk)))
			} else {
				privateKey = ec
			}
		} else if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				panic("cannot read a certificate in the pem file")
			}
			chain = append(chain, cert)
		}
	}
	return privateKey, chain
}

func LoadAuthorityPool() (*x509.CertPool, error) {
	// Load the CA certificate
	trustedCert, err := os.ReadFile(*AuthorityPath)
	if err != nil {
		return nil, err
	}

	// Put the CA certificate into the certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		return nil, fmt.Errorf("the CA certificate couldn't be constructed")
	}

	return certPool, nil
}

func LoadX509KeyPair(privateKey *ecdsa.PrivateKey, chain []*x509.Certificate) tls.Certificate {
	certificate := tls.Certificate{
		PrivateKey:  privateKey,
		Certificate: [][]byte{},
	}
	for _, certificateInChain := range chain {
		certificate.Certificate = append(certificate.Certificate, certificateInChain.Raw)
	}
	return certificate
}
