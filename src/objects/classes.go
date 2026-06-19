package objects

// Class is the general definition of elements within the system.
// It applies to any element that can be declared with a specific class.
// Class is NOT a declaration type, but rather a categorization of elements.
type Class string

const CLASS_LINK Class = "link"
const CLASS_INSTANCE Class = "instance"
const CLASS_PREDICATE Class = "predicate"

// IsInstanceOfClass checks if an element is declared with a specific class.
// For instance, a link should declare CLASS_LINK
func IsInstanceOfClass(element Element, c Class) bool {
	if element == nil {
		return false
	}
	return element.DeclaringClass() == c
}
