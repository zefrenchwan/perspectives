package values_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

func TestLocalState(t *testing.T) {
	lastName := values.NewStringLocalMapping(map[string]periods.Period{"Doe": periods.NewFullPeriod()})
	firstName := values.NewStringLocalMapping(map[string]periods.Period{"John": periods.NewFullPeriod()})
	attributes := map[string]values.ImmutableValuesMapping[values.PrimitiveValue]{"firstName": firstName, "lastName": lastName}
	state := values.NewLocalState("id", periods.NewFullPeriod(), attributes, nil)

	full := periods.NewFullPeriod()
	if state.Id() != "id" {
		t.Error("wrong id")
	} else if state.ToHashString() == "" {
		t.Error("wrong hash : default value")
	} else if !state.Activity().Equals(full) {
		t.Error("wrong activity : should be full period")
	}

	for key, _ := range state.Roles() {
		t.Errorf("wrong role, expected nothing, got %s", key)
	}

	for key, mapper := range state.Attributes() {
		switch key {
		case "firstName":
			for duration, value := range mapper.Range() {
				if value.Content() != "John" {
					t.Errorf("wrong value, expected John, got %s", value)
				} else if !duration.Equals(full) {
					t.Errorf("wrong duration, expected full period, got %v", duration)
				}
			}
		case "lastName":
			for duration, value := range mapper.Range() {
				if value.Content() != "Doe" {
					t.Errorf("wrong value, expected Doe, got %s", value)
				} else if !duration.Equals(full) {
					t.Errorf("wrong duration, expected full period, got %v", duration)
				}
			}
		default:
			t.Errorf("unexpected attribute, expected firstName or lastName, got %s", key)
		}
	}
}

func TestLocalStateHash(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	// same values on purpose, one is primitive values, the other is reference values
	// primitives
	lastName := values.NewStringLocalMapping(map[string]periods.Period{"Doe": periods.NewFullPeriod()})
	firstName := values.NewStringLocalMapping(map[string]periods.Period{"John": periods.NewFullPeriod()})
	// references
	lastNameRole := values.NewReferenceLocalMapping(map[string]periods.Period{"Doe": periods.NewFullPeriod()})
	firstNameRole := values.NewReferenceLocalMapping(map[string]periods.Period{"John": periods.NewFullPeriod()})

	// note that it makes no sense, but we test that we distinguish primitive and reference values
	attributes := map[string]values.ImmutableValuesMapping[values.PrimitiveValue]{"firstName": firstName, "lastName": lastName}
	roles := map[string]values.ImmutableValuesMapping[values.ReferenceValue]{"firstName": firstNameRole, "lastName": lastNameRole}

	emptyState := values.NewLocalState("id", periods.NewFullPeriod(), nil, nil)
	partialState := values.NewLocalState("id", periods.NewPeriodSince(now, true), nil, nil)
	attributesState := values.NewLocalState("id", periods.NewFullPeriod(), attributes, nil)
	roleState := values.NewLocalState("id", periods.NewFullPeriod(), nil, roles)

	// distinguish periods
	if emptyState.ToHashString() == partialState.ToHashString() {
		t.Error("wrong hash : empty state should not be equal to partial state")
	}

	// split values for roles and attributes
	if attributesState.ToHashString() == roleState.ToHashString() {
		t.Error("wrong hash : attributes state should not be equal to role state")
	}
}
