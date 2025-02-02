package state

type CreateEdgeRequest struct {
	PeerID RouterIdentity // other user id
	EdgeID EdgeIdentity   // edge object id
}

type Endpoint struct {
	Address    string // v4 address
	PublicPort uint16
}

type CreateEdgeResponse struct {
	PublicEndpoint  Endpoint
	PresharedKeybB4 string
}

type SeedEdgeRequest struct {
	OtherPeer    Endpoint
	PresharedKey string
}
