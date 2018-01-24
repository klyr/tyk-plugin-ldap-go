package main

import (
	"log"

	coprocess "github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"golang.org/x/net/context"
)

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(ctx context.Context, object *coprocess.Object) (*coprocess.Object, error) {
	log.Println("Receiving object: ", object)
	log.Println("Metadata: ", object.GetMetadata())

	switch object.HookName {
	case "AddMeta":
		log.Println("Calling 'HookAddMeta' hook...")
		return HookAddMeta(object)
	}

	log.Println("Unknown hook: ", object.HookName)
	return object, nil
}

func (d *Dispatcher) DispatchEvent(ctx context.Context, event *coprocess.Event) (*coprocess.EventReply, error) {
	return &coprocess.EventReply{}, nil
}
