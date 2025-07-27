package models

import "time"

type Actor struct {
	Entity
}

func NewActor(id string, birthdate time.Time) Actor {
	return Actor{Entity: NewEntity(id, birthdate.Round(1*time.Minute))}
}
