package session

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func (c *Session) Authenticate(
	Event *proto.InitializeRequest,
	RootCertPool *x509.CertPool,
) error {
	// Verify that the date only has a variation inferior to 2s
	time := time.Now().Unix()
	if time-Event.TimestampSigned > 2 {
		err := fmt.
			Errorf("the verification timestamp is too far from the current time")
		log.Debug().
			Err(err).
			Msg("user supplied an invalid date/time")
		return err
	}
	intermediates := x509.NewCertPool()
	var userCertificate *x509.Certificate

	for block, rest := pem.Decode(Event.Certificate); block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "CERTIFICATE" {
			certificate, err := x509.ParseCertificate(block.Bytes)
			if certificate.IsCA && err == nil {
				intermediates.AddCert(certificate)
			} else {
				userCertificate = certificate
			}
		}
	}

	if _, err := userCertificate.Verify(x509.VerifyOptions{
		Roots:         RootCertPool,
		Intermediates: intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		log.Debug().
			Err(err).
			Msg("user supplied an orion-invalid certificate")
		return err
	}

	if userCertificate.Subject.CommonName == fmt.Sprintf("%d:oriond", Event.MemberId) {
		err := fmt.Errorf("this certificate is not valid for oriond")
		log.Error().
			Err(err).
			Msg("user supplied an orion-invalid certificate")
		return err
	}

	// Calculate the hash given in order to check the client signature
	nonce := internal.CalculateNonceBytes(Event.MemberId, Event.FriendlyName, Event.TimestampSigned)

	// Verify that the user-provided data matches the signature created using the client root key
	successful := ecdsa.VerifyASN1(userCertificate.PublicKey.(*ecdsa.PublicKey), nonce, Event.Signed)
	if !successful {
		err := fmt.Errorf("this signature does not seem to be a valid ECDSA signature")
		log.Debug().
			Err(err).
			Msg("the user authentication failed, invalid signature")
		return err
	}

	log.Info().
		Msgf("User %s (%d) auth with certificate with serial: %s", Event.FriendlyName, Event.MemberId, userCertificate.SerialNumber)

	// the user is authenticated, we start listening for global events

	log.Debug().Msg("client created")
	// registering in the manager
	c.meta = &SessionMeta{
		memberId:     Event.MemberId,
		friendlyName: Event.FriendlyName,
	}
	c.sessionManager.sessions[Event.MemberId] = c

	log.Debug().Msg("broadcasting the new client message")
	c.sessionManager.newClients.Notify(
		&proto.NewMemberEvent{
			FriendlyName: Event.FriendlyName,
			PeerId:       Event.MemberId,
		},
	)

	log.Debug().Msg("random session id generation")
	sessionId, err := generateRandomString(64)
	if err != nil {
		return err
	}
	c.sID = sessionId
	// Since the registry is not handling the channel while login, we simply wait by launching a goroutine
	go func() {
		c.streamSend.Broadcast(&proto.RPCServerEvent{
			Event: &proto.RPCServerEvent_SessionId{
				SessionId: sessionId,
			},
		})
	}()
	c.sessionManager.sessionIdsMap[sessionId] = &c.meta.memberId

	log.Debug().Msg("starting listeners")
	go c.eventListeners()

	return nil
}
