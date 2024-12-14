package state

const (
	// Event sent at the beginning for a new connection
	// allows version negotiation
	MessageKindHello = "hello"
	// Sent when a new router joined the network,
	// this is so that a user, can choose to initiate
	// a new connection to this new router.
	MessageKindRouterConnect = "new_router"
	// Sent by a peer that wants to initialize a new peer-to-peer link
	MessageKindRouterEdgeConnectInitializeRequest  = "edge_initialize_request"
	MessageKindRouterEdgeConnectInitializeResponse = "edge_initialize_response"

	// Event emitted once an edge is destroyed
	MessageKindRouterEdgeTeardown = "edge_teardown"
)
