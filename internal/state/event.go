package state

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Meta interface holding the different events
type Event interface{}

type JsonEvent struct {
	Kind    string          `json:"k"`
	Content json.RawMessage `json:"e"`
}

func UnmarshalEvent(
	event JsonEvent,
) (Event, error) {
	switch event.Kind {
	case MessageKindHello:
		hello := Hello{}
		if err := json.Unmarshal(event.Content, &hello); err != nil {
			return nil, err
		}
		return hello, nil
	case MessageKindRouterConnect:
		hello := RouterConnectEvent{}
		if err := json.Unmarshal(event.Content, &hello); err != nil {
			return nil, err
		}
		return hello, nil
	case MessageKindRouterEdgeConnectInitializeRequest:
		hello := RouterInitiateRequest{}
		if err := json.Unmarshal(event.Content, &hello); err != nil {
			return nil, err
		}
		return hello, nil
	case MessageKindRouterEdgeTeardown:

		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("unknown event for deserialize")
	}
}

func MarshalEvent(
	event Event,
) (*JsonEvent, error) {
	switch event := event.(type) {
	case Hello:
		bytes, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		return &JsonEvent{
			Kind:    MessageKindHello,
			Content: bytes,
		}, nil
	case RouterConnectEvent:
		bytes, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		return &JsonEvent{
			Kind:    MessageKindRouterConnect,
			Content: bytes,
		}, nil
	case RouterInitiateRequest:
		bytes, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		return &JsonEvent{
			Kind:    MessageKindRouterEdgeConnectInitializeRequest,
			Content: bytes,
		}, nil
	case CreateEdgeRequest:
		bytes, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		return &JsonEvent{
			Kind:    MessageKindCreateEdgeRequest,
			Content: bytes,
		}, nil
	case CreateEdgeResponse:
		bytes, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		return &JsonEvent{
			Kind:    MessageKindCreateEdgeResponse,
			Content: bytes,
		}, nil
	case SeedEdgeRequest:

	}

	return nil, fmt.Errorf("event serialization not supported %s", reflect.TypeOf(event))
}
