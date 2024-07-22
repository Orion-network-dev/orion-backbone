package implementation

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog/log"
)

func Authenticate(
	TimeStamp int64,
	Certificate []byte,
	SignedNonce []byte,
	MemberId uint32,
	FriendlyName string,
	RootCertPool *x509.CertPool,
) error {
	// Verify that the date only has a variation inferior to 2s
	time := time.Now().Unix()
	if time-TimeStamp > 2 {
		err := fmt.Errorf("the verification timestamp is too far from the current time")
		log.Debug().Err(err).Msg("user supplied an invalid date/time")
		return err
	}

	// Parse the user-given certificate
	cert, err := internal.ParsePEMCertificate(Certificate)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse the user leaf certificate")
		return err
	}

	// Create a new pool from the user-given PEM trust chain
	intermediates, err := internal.CreateCertPoolFromPEM(Certificate)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse the intermediary certificates")
		return err
	}

	identifier := fmt.Sprintf("%d.member.orionet.re", MemberId)

	// Verifying the certificate validity using the root certificate and user-provided
	// intermediary certificates. This checks that the certificate is signed and allowed to use
	// the name `{member_id}.mem.orionet.re` which specifies a member member for the member_id {member_id}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:         RootCertPool,
		Intermediates: intermediates,
		DNSName:       identifier,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		log.Debug().Err(err).Msg("user supplied an orion-invalid certificate")
		return err
	}

	if cert.Subject.CommonName != fmt.Sprintf("%s:oriond", identifier) {
		err := fmt.Errorf("this certificate is not valid for oriond")
		log.Error().Err(err).Msg("user supplied an orion-invalid certificate")
		return err
	}

	// Calculate the hash given in order to check the client signature
	nonce := internal.CalculateNonceBytes(MemberId, FriendlyName, TimeStamp)

	// Verify that the user-provided data matches the signature created using the client root key
	successful := ecdsa.VerifyASN1(cert.PublicKey.(*ecdsa.PublicKey), nonce, SignedNonce)
	if !successful {
		err := fmt.Errorf("this signature does not seem to be a valid ECDSA signature")
		log.Debug().Err(err).Msg("the user authentication failed, invalid signature")
		return err
	}

	log.Info().Msgf("User %s auth with certificate with serial: %s", identifier, cert.SerialNumber)

	return nil
}
