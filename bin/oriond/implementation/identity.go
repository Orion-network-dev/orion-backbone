package implementation

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog/log"
)

var (
	memberIdRegex    = regexp.MustCompile(`(.*?).member.orionet.re`)
	memberIdOverride = flag.Uint("override-member-id", 0, "An override of the memberID of this instance")
	asn              = flag.Int("override-asn", 0, "An override of the ASN number used by this instance")
)

func (c *OrionClientDaemon) resolveIdentity() error {
	// If we have a memver-id override
	if *memberIdOverride != 0 {
		c.memberId = uint32(*memberIdOverride)
	} else {
		certificateFile, err := os.ReadFile(*internal.CertificatePath)
		if err != nil {
			log.Error().
				Err(err).
				Str("file", *internal.CertificatePath).
				Msg("Cannot open the certificate path")
			return err
		}

		certificate, err := internal.ParsePEMCertificate(
			certificateFile,
		)
		if err != nil {
			log.Error().
				Err(err).
				Str("file", *internal.CertificatePath).
				Msg("Failed to load the certificate data")
			return err
		}

		// The certificate should only have a single DNSName in his list
		if len(certificate.DNSNames) != 1 {
			err = fmt.Errorf("the certificate is authorized for multiple dns names, please use -override-member-id to specify the member-id")
			log.Error().
				Err(err).
				Str("file", *internal.CertificatePath).
				Msg(err.Error())
			return err
		}

		// Set the memberId to the ont in the certificate.
		matches := memberIdRegex.FindStringSubmatch(certificate.DNSNames[0])
		if len(matches) == 2 {
			number, err := strconv.ParseInt(matches[1], 10, 32)
			if err != nil {
				log.Error().
					Err(err).
					Str("file", *internal.CertificatePath).
					Msg("the member_id field in the certificate couldn't be parsed, please use -override-member-id to specify the member-id")
				return err
			}
			c.memberId = uint32(number)
		} else {
			err = fmt.Errorf("the member_id couldn't be extracted from the certicate")
			log.Error().
				Err(err).
				Str("file", *internal.CertificatePath).
				Msg("the member_id field in the certificate couldn't be parsed, please use -override-member-id to specify the member-id")
			return err
		}
	}

	// We check if we have a asn number override
	if *asn != 0 {
		c.asn = uint32(*asn)
	} else {
		// The regular orion as allocation
		c.asn = 64511 + c.memberId
	}

	return nil
}
