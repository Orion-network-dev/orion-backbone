package internal

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

// Calculate the nonce bytes
func CalculateNonceBytes(MemberId uint32, FriendlyName string, time int64) []byte {
	return sha512.New().Sum([]byte(fmt.Sprintf("%d:%s:%d", MemberId, FriendlyName, time)))
}

// Used to calculate the nonce and sign it
func CalculateNonce(
	MemberId uint32,
	FriendlyName string,
	Certificate []byte,
	PrivateKey *ecdsa.PrivateKey,
) (*proto.InitializeRequest, error) {
	time := time.Now().Unix()
	authHash := CalculateNonceBytes(MemberId, FriendlyName, time)

	signed, err := ecdsa.SignASN1(rand.Reader, PrivateKey, authHash)

	if err != nil {
		log.Error().Err(err).Msgf("couldn't sign the nonce data")
		return nil, err
	}

	return &proto.InitializeRequest{
		FriendlyName:    FriendlyName,
		TimestampSigned: time,
		MemberId:        MemberId,
		Certificate:     Certificate,
		Signed:          signed,
	}, nil
}
