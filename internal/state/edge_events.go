package state

type NewEdgeEvent struct {
	Router *Edge
}

type AskForNewEdge struct {
	OtherNode *Router
}

type CreateEdgeRequest struct{}
type CreateEdgeResponse struct{}

type SeedEdgePeer struct {
	PublicKey    string
	PresharedKey string
	DefaultRoute bool
}

type SeedEdgeRequest struct {
	Peers []SeedEdgePeer
}
