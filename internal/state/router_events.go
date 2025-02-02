package state

type RouterConnectEvent struct {
	Router *Router
}

type RouterInitiateRequest struct {
	Identity *RouterIdentity
}

type RouterEdgeRemovedEvent struct {
	Edge *Edge
}

type Hello struct {
	Message  string         `json:"message"`
	Identity RouterIdentity `json:"identity"`
	Version  string         `json:"version"`
	Commit   string         `json:"commit"`
	Session  string         `json:"session"`
}
