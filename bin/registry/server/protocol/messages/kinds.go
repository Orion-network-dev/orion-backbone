package messages

const (
	// Event sent at the beginning for a new connection
	// allows version negotiation
	MessageKindHello = "hello"
	// Sent when a new router joined the network,
	// this is so that a user, can choose to initiate
	// a new connection to this new router.
	MessageKindRouterConnect = "new_router"
	// Sent when a router session ended meaning
	// a disconnection from the orion registry
	// this will triger a teardown of the client
	// and thus a end of all wireguard tunnels
	MessageKindRouterDisconnect = "router_disconnect"
	// A message sent by a router wanting to connect
	// to another one.
	// step1. (peer1) ---> (registry) ---> (peer2)
	MessageKindRouterEdgeInitConnectRequest = "edge_create_request"
	// A message considered as a response of stage1.
	// step2. (peer2) ---> (registry) ---> (peer1)
	MessageKindRouterEdgeInitConnectRequestResponse = "edge_create_request_ack"
)
