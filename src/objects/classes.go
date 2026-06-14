package objects

// Class is the general definition of elements within the system.
// It applies to any element that can be declared with a specific class.
// Class is NOT a declaration type, but rather a categorization of elements.
// To use declaration types, see the traits mechanism.
type Class string

const CLASS_TRAIT Class = "trait"
const CLASS_LINK Class = "link"
const CLASS_INSTANCE Class = "instance"
const CLASS_VARIABLE Class = "variable"
const CLASS_INSTANCES_COLLECTION Class = "collection of instances"

// IsInstanceOfClass checks if an element is declared with a specific class.
// For instance, a link should declare CLASS_LINK
func IsInstanceOfClass(element Element, c Class) bool {
	if element == nil {
		return false
	}
	return element.DeclaringClass() == c
}
