package state

type NewEdgeEvent struct {
	Router *Edge
}

type AskForNewEdge struct {
	OtherNode *Router
}
