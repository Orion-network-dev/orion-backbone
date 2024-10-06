package session

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
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
	Event *proto.RouterLogin,
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

	// Parse the user-given certificate
	cert, err := internal.ParsePEMCertificate(Event.Certificate)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to parse the user leaf certificate")
		return err
	}

	// Create a new pool from the user-given PEM trust chain
	intermediates, err := internal.CreateCertPoolFromPEM(Event.Certificate)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to parse the intermediary certificates")
		return err
	}

	identifier := fmt.Sprintf("%d.member.orionet.re", Event.Identity.MemberId)

	// Verifying the certificate validity using the root certificate and user-provided
	// intermediary certificates. This checks that the certificate is signed and allowed to use
	// the name `{member_id}.mem.orionet.re` which specifies a member member for the member_id {member_id}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:         RootCertPool,
		Intermediates: intermediates,
		DNSName:       identifier,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		log.Debug().
			Err(err).
			Msg("user supplied an orion-invalid certificate")
		return err
	}

	if cert.Subject.CommonName != fmt.Sprintf("%s:oriond", identifier) {
		err := fmt.Errorf("this certificate is not valid for oriond")
		log.Error().
			Err(err).
			Msg("user supplied an orion-invalid certificate")
		return err
	}

	// Calculate the hash given in order to check the client signature
	nonce := internal.CalculateNonceBytes(Event.Identity, Event.TimestampSigned)

	// Verify that the user-provided data matches the signature created using the client root key
	successful := ecdsa.VerifyASN1(cert.PublicKey.(*ecdsa.PublicKey), nonce, Event.Signed)
	if !successful {
		err := fmt.Errorf("this signature does not seem to be a valid ECDSA signature")
		log.Debug().
			Err(err).
			Msg("the user authentication failed, invalid signature")
		return err
	}

	log.Info().
		Msgf("User %s auth with certificate with serial: %s", identifier, cert.SerialNumber)

	// the user is authenticated, we start listening for global events

	log.Debug().Msg("client created")
	// registering in the manager
	c.meta = Event.Identity
	c.sessionManager.sessions[internal.IdentityFromRouter(Event.Identity)] = c

	log.Debug().Msg("broadcasting the new client message")
	c.sessionManager.newClients.Notify(
		&proto.RouterConnectedEvent{
			Router: &proto.Router{
				FriendlyName: Event.Identity.FriendlyName,
				MemberId:     Event.Identity.MemberId,
				RouterId:     Event.Identity.RouterId,
			},
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
		c.streamSend.Broadcast(&proto.ServerToPeers{
			Event: &proto.ServerToPeers_SessionId{
				SessionId: &proto.SessionIDAssigned{
					SessionId: sessionId,
				},
			},
		})
	}()
	id := internal.IdentityFromRouter(c.meta)
	c.sessionManager.sessionIdsMap[sessionId] = &id

	log.Debug().Msg("starting listeners")
	go c.eventListeners()

	return nil
}
