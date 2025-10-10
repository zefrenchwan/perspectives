package models

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

func TestEventsCreation(t *testing.T) {
	// test empty
	var e models.Event
	if !e.IsEmpty() {
		t.Log("empty definition failure")
		t.Fail()
	}

	e = models.NewSimpleEvent(nil)
	if !e.IsEmpty() {
		t.Log("empty definition failure")
		t.Fail()
	}

	action := models.EventDeletionById{Id: "test"}
	e = models.NewSimpleEvent(action)
	if e.IsEmpty() {
		t.Log("event cannot be empty because it has an action")
		t.Fail()
	}

}

func TestDeleteByObject(t *testing.T) {
	alexandre := models.NewObject([]string{"Human"})
	arthur := models.NewObject([]string{"Fictional character"})

	event := models.EventDeletionById{Id: arthur.Id()}

	if event.Matches(alexandre) {
		t.Log("id mismatch")
		t.Fail()
	} else if !event.Matches(arthur) {
		t.Log("id should match")
		t.Fail()
	}
}
