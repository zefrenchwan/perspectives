package objects

import (
	"slices"
	"strings"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// The idea here is to generate a hash string for each object type based on its unique properties.
// Hash function is assumed to be deterministic, consistent and injective most of the time (with super high probability).

// hashDynamicValues returns a hash string for the given DynamicValues.
func hashDynamicValues(dv DynamicValues) string {
	elements := make([]string, 0)
	for period, value := range dv.Range {
		elements = append(elements, period.AsRawString()+"=>"+primitiveValueToString(value))
	}

	slices.Sort(elements)
	return commons.HashString(strings.Join(elements, ","))
}

// hashInstance returns a hash string for the given Instance.
func hashInstance(instance Instance) string {
	content := make([]string, 0)
	content = append(content, instance.Id())
	content = append(content, instance.Activity().AsRawString())
	for attr, value := range instance.Values() {
		content = append(content, attr+"=>"+hashDynamicValues(value))
	}

	slices.Sort(content)
	return commons.HashString(strings.Join(content, ","))
}

// hashTrait returns a hash string for the given Trait.
func hashTrait(trait Trait) string {
	content := make([]string, 0)
	content = append(content, trait.Id())
	content = append(content, trait.Name())
	for attr, value := range trait.Attributes() {
		content = append(content, attr+"=>"+value)
	}

	slices.Sort(content)
	return commons.HashString(strings.Join(content, ","))
}

// hashLink returns a hash string for the given Link.
func hashLink(link Link) string {
	content := make([]string, 0)
	content = append(content, link.Name()+":"+link.Activity().AsRawString())
	for attr, value := range link.Range {
		content = append(content, attr+"=>"+value.toHashString())
	}

	slices.Sort(content)
	return commons.HashString(strings.Join(content, ","))
}
