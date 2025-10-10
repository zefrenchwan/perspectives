package models

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
)

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
