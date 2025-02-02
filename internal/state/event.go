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

func reverseMap[M ~map[K]V, K comparable, V comparable](m M) map[V]K {
	reversedMap := make(map[V]K)
	for key, value := range m {
		reversedMap[value] = key
	}
	return reversedMap
}

var (
	TypesEvents = map[string]reflect.Type{
		MessageKindHello:                              reflect.TypeOf(Hello{}),
		MessageKindRouterConnect:                      reflect.TypeOf(RouterConnectEvent{}),
		MessageKindRouterEdgeConnectInitializeRequest: reflect.TypeOf(RouterInitiateRequest{}),
		MessageKindCreateEdgeRequest:                  reflect.TypeOf(CreateEdgeRequest{}),
		MessageKindCreateEdgeResponse:                 reflect.TypeOf(CreateEdgeResponse{}),
		MessageKindSeedEdgeRequest:                    reflect.TypeOf(SeedEdgeRequest{}),
		MessageKindRouterEdgeTeardown:                 reflect.TypeOf(RouterEdgeRemovedEvent{}),
	}
	TypesEventsReverse = reverseMap(TypesEvents)
)

func UnmarshalEvent(
	event JsonEvent,
) (Event, error) {
	type_ := TypesEvents[event.Kind]
	if type_ == nil {
		return nil, fmt.Errorf("cannot unmarshal type: unknown message kind")
	}
	data := reflect.New(type_).Interface()
	if err := json.Unmarshal(event.Content, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func MarshalEvent(
	event Event,
) (*JsonEvent, error) {
	type_ := TypesEventsReverse[reflect.TypeOf(event)]
	if type_ == "" {
		return nil, fmt.Errorf("cannot marshal type: unknown message kind %T", event)
	}
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return &JsonEvent{
		Kind:    type_,
		Content: bytes,
	}, nil
}
