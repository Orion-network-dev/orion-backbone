package server

import (
	"fmt"
	"net/http"
)

func (c *Server) whoami(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil && len(r.TLS.PeerCertificates) > 0 {
		return
	}
	state := r.TLS

	fmt.Fprint(w, ">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<\n")

	fmt.Fprintf(w, "Version: %x\n", state.Version)
	fmt.Fprintf(w, "HandshakeComplete: %t\n", state.HandshakeComplete)
	fmt.Fprintf(w, "DidResume: %t\n", state.DidResume)
	fmt.Fprintf(w, "CipherSuite: %x\n", state.CipherSuite)
	fmt.Fprintf(w, "NegotiatedProtocol: %s\n", state.NegotiatedProtocol)

	fmt.Fprintf(w, "Certificate chain:\n")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		fmt.Fprintf(w, " %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s\n", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		fmt.Fprintf(w, "   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s\n", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
}
