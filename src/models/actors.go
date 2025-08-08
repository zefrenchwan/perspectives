package models

import "github.com/google/uuid"

// Actor is any person represented in the model
type Actor struct {
	Id string
}

// NewActor builds a new actor with a default id
func NewActor() Actor {
	return Actor{Id: uuid.NewString()}
}
