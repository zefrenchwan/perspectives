package entities

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// The hashing mechanism aims to generate a deterministic, unique, and consistent string
// representation of each object type based on its properties. This allows O(1) checks
// for functional equivalence using hashes.
//
// DESIGN PRINCIPLES TO AVOID COLLISIONS:
// 1. CLASS PREFIXING: Every hash is prefixed with a specific identifier (e.g., "LINK|")
//    to prevent cross-class collisions (an Instance and a group cannot have the same hash).
// 2. STRUCTURAL SEPARATION: Fixed properties (Id, Name, Activity) are kept in strict order and
//    NEVER sorted together with dynamic properties. Sorting is restricted ONLY to map iterations
//    (like attributes and roles) to guarantee determinism without destroying structural semantics.
// 3. LENGTH PREFIXING: To prevent delimiter injection (e.g., key "A=>B" + value "C" colliding
//    with key "A" + value "B=>C"), we prefix strings with their lengths (e.g., "4:name").

// hashDynamicValues returns a collision-resistant hash string for the given DynamicValues.
func hashDynamicValues(dv DynamicValues) string {
	// IMPLEMENTATION CHOICE : According to DynamicValues semantics in instances.go,
	// an empty collection is equivalent to another empty collection regardless of its
	// underlying DataType. Hence, we must return a constant hash for all empty values
	// to respect the a.Same(b) == (hash(a) == hash(b)) contract.
	if dv == nil || dv.IsEmpty() {
		return commons.HashString("DynamicValues:empty")
	}

	valueType := dv.DataType()

	// We don't know the exact number of periods in advance when using the range iterator,
	// so we start with an empty slice.
	elements := make([]string, 0)

	// Range over the time-dependent values using Go 1.22+ iterator pattern
	for period, value := range dv.Range {
		valueString := primitiveValueToString(value)
		sizeString := strconv.Itoa(len(valueString))

		// Use strict formatting with length prefixing to prevent delimiter injection.
		// Format: [Period]->Type(Length):Value
		mappedString := fmt.Sprintf("[%s]->%s(%s):%s", period.AsRawString(), valueType, sizeString, valueString)
		elements = append(elements, mappedString)
	}

	// Sort ONLY the dynamic elements to ensure a deterministic hash regardless of iteration order.
	slices.Sort(elements)

	var builder strings.Builder
	builder.WriteString("DynamicValues|")
	builder.WriteString(strings.Join(elements, "|"))

	return commons.HashString(builder.String())
}

// hashEntity returns the hash of the given entity.
// It uses a deterministic algorithm to compute the hash of the entity.
func hashEntity(element Entity) string {
	if element == nil {
		return ""
	}

	dynamicContent := make([]string, 0)
	for attr, value := range element.Values {
		hashedValue := hashDynamicValues(value)
		// Embed the length of the attribute name to avoid delimiter injection.
		// hashedValue is already a fixed-length hash, so it doesn't need length prefixing.
		dynamicContent = append(dynamicContent, fmt.Sprintf("%d:%s->%s", len(attr), attr, hashedValue))
	}

	// We must not sort fixed fields like ID or Activity along with attributes.
	slices.Sort(dynamicContent)

	// same algorithm on links
	dynamicLinks := make([]string, 0)
	// Iterate over the roles and their associated elements
	for role, linkedElement := range element.Links {
		linkedElementHash := ""
		if linkedElement != nil {
			// recursive resolution is OK because each element has its own hash, no full walk
			linkedElementHash = linkedElement.ToHashString()
		}
		// Embed the length of the role string to avoid any parsing ambiguities.
		dynamicLinks = append(dynamicLinks, fmt.Sprintf("%d:%s=>%s", len(role), role, linkedElementHash))
	}

	slices.Sort(dynamicLinks)

	var builder strings.Builder
	builder.WriteString("INSTANCE|")

	// Append fixed properties in a strict, unchangeable order.
	idStr := element.Id()
	builder.WriteString(fmt.Sprintf("id:%d:%s|", len(idStr), idStr))

	actStr := element.Activity().AsRawString()
	builder.WriteString(fmt.Sprintf("act:%d:%s|", len(actStr), actStr))

	builder.WriteString("vals:")
	builder.WriteString(strings.Join(dynamicContent, "|"))

	builder.WriteString("links:")
	builder.WriteString(strings.Join(dynamicLinks, "||"))

	return commons.HashString(builder.String())
}
