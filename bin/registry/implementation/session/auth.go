package session

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/imusmanmalik/randomizer"
	"github.com/rs/zerolog/log"
)

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

	identifier := fmt.Sprintf("%d.member.orionet.re", Event.MemberId)

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
	nonce := internal.CalculateNonceBytes(Event.MemberId, Event.FriendlyName, Event.TimestampSigned)

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

	// registering in the manager
	c.meta = &SessionMeta{
		memberId:     Event.MemberId,
		friendlyName: Event.FriendlyName,
	}
	c.sessionManager.sessions[Event.MemberId] = c
	c.sessionManager.newClients.Broadcast(
		&proto.NewMemberEvent{
			FriendlyName: Event.FriendlyName,
			PeerId:       Event.MemberId,
		},
	)
	sessionId, _ := randomizer.RandomString(64)
	c.sID = sessionId
	c.streamSend <- &proto.RPCServerEvent{
		Event: &proto.RPCServerEvent_SessionId{
			SessionId: &proto.SessionIDIssued{
				SessionId: sessionId,
			},
		},
	}
	c.sessionManager.sessionIdsMap[sessionId] = &c.meta.memberId

	go c.eventListeners()

	return nil
}
