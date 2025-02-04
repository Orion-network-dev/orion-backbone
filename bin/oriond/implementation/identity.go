package implementation

import (
	"flag"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	memberIdOverride = flag.Uint("override-member-id", 0, "An override of the memberID of this instance")
	asn              = flag.Uint("override-asn", 0, "An override of the ASN number used by this instance")
)

func (c *OrionClientDaemon) resolveIdentity() error {
	// If we have a memver-id override
	if *memberIdOverride != 0 {
		c.memberId = uint32(*memberIdOverride)
	} else {
		name := c.chain[0].Subject.CommonName
		nameParts := strings.Split(name, ":")

		number, err := strconv.ParseInt(nameParts[0], 10, 32)
		if err != nil {
			log.Error().
				Err(err).
				Msg("the member_id field in the certificate couldn't be parsed, please use -override-member-id to specify the member-id")
			return err
		}
		c.memberId = uint32(number)
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
