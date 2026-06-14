package objects

import (
	"errors"
	"slices"
)

// Variable defines a variable that can be used to store values of different classes.
// It is used in pattern matching to represent a placeholder for a value of a specific class.
// By design, it may take only values of classes (NOT TRAITS) that are explicitly specified when creating the variable.
// Two variables can be Same if they have the same name and accept the same classes.
type Variable interface {
	Element
	// Accepts returns a list of classes that this variable can accept.
	// For instance, a variable for class Instance can accept only instances values.
	Accepts() []Class
	// AcceptsOneOf returns true if this variable can accept at least one of the given classes.
	AcceptsOneOf(acceptedClasses ...Class) bool
	// Name returns the name of the variable.
	Name() string
}

// VariableBuilder is used to build a Variable from scratch.
type VariableBuilder interface {
	// WithName sets the name of the variable being built and returns the updated VariableBuilder instance.
	WithName(string) VariableBuilder
	// WithAcceptedClass adds a class to the list of accepted classes for the variable being built.
	//It returns the updated VariableBuilder instance.
	WithAcceptedClass(class Class) VariableBuilder
	// WithoutAcceptedClass removes a class to the list of accepted class for the variable being built.
	// It returns the updated VariableBuilder instance.
	WithoutAcceptedClass(class Class) VariableBuilder
	// Errors returns any errors encountered during the building process.
	// Errors are cumulative and can be retrieved after each method call.
	Errors() error
	// Build makes the variable from the builder.
	// It also resets the builder to its initial state.
	Build() (Variable, error)
}

// baseVariable is a base implementation of the Variable interface.
// It basically is a name and a set of accepted classes.
type baseVariable struct {
	// name of the variable
	name string
	// acceptedClasses is a set of accepted classes for the variable.
	// It is classes and not traits on purpose : traits will be managed by another replacement mechanism.
	// By setting classes, we make a first exclusion mechanism.
	acceptedClasses map[Class]bool
}

// Accepts returns the list of accepted classes for the variable.
func (v *baseVariable) Accepts() []Class {
	result := make([]Class, 0, len(v.acceptedClasses))
	for acceptedClass := range v.acceptedClasses {
		result = append(result, acceptedClass)
	}

	slices.Sort(result)
	return result
}

// AcceptsOneOf returns true if the variable accepts at least one of the given classes.
func (v *baseVariable) AcceptsOneOf(acceptedClasses ...Class) bool {
	for _, acceptedClass := range acceptedClasses {
		if v.acceptedClasses[acceptedClass] {
			return true
		}
	}
	return false
}

// Id returns the name of the variable.
func (v *baseVariable) Id() string {
	return v.name
}

// DeclaringClass returns the class of the variable : A CLASS VARIABLE
func (v *baseVariable) DeclaringClass() Class {
	return CLASS_VARIABLE
}

// Name returns the name of the variable.
func (v *baseVariable) Name() string {
	return v.name
}

// isLinkable returns true if the variable is linkable.
// It implements the sealed interface mechanism.
func (v *baseVariable) isLinkable() bool {
	return true
}

// Same returns true if the variable is the same as the other variable :
// same name, same accepted classes.
func (v *baseVariable) Same(other Element) bool {
	if v == nil && other == nil {
		return true
	} else if v == nil || other == nil {
		return false
	} else if other.DeclaringClass() != CLASS_VARIABLE {
		return false
	}

	if otherVariable, ok := other.(Variable); ok {
		if v.name != otherVariable.Name() {
			return false
		}

		otherAccepts := otherVariable.Accepts()
		if len(otherAccepts) != len(v.acceptedClasses) {
			return false
		}

		for _, acceptedClass := range otherAccepts {
			if !v.acceptedClasses[acceptedClass] {
				return false
			}
		}
		return true
	}
	return false
}

// baseVariableBuilder is a builder for a variable.
// A valid variable has a non empty name and at least one accepted class.
type baseVariableBuilder struct {
	// name of the variable to build
	name string
	// acceptedClasses is a map of accepted classes for the variable.
	// Bool value should be true, it is basically a set implementation.
	acceptedClasses map[Class]bool
	// globalErrors is an error that is returned if the variable is not valid.
	// It accumulates errors, no error clean if fixed.
	globalErrors error
}

// WithName sets the name of the variable to build.
func (b *baseVariableBuilder) WithName(name string) VariableBuilder {
	if name == "" {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("variable name cannot be empty"))
		return b
	}

	b.name = name
	return b
}

// WithAcceptedClass sets the accepted class for the variable to build.
func (b *baseVariableBuilder) WithAcceptedClass(class Class) VariableBuilder {
	var empty Class
	if class == empty {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("accepted class cannot be nil"))
		return b
	}

	if b.acceptedClasses == nil {
		b.acceptedClasses = make(map[Class]bool)
	}

	b.acceptedClasses[class] = true
	return b
}

// WithoutAcceptedClass removes the accepted class for the variable to build.
func (b *baseVariableBuilder) WithoutAcceptedClass(class Class) VariableBuilder {
	delete(b.acceptedClasses, class)
	if b.acceptedClasses == nil || len(b.acceptedClasses) == 0 {
		b.globalErrors = errors.Join(b.globalErrors, errors.New("variable must accept at least one class"))
	}

	return b
}

// Errors returns the errors that occurred during the building of the variable.
func (b *baseVariableBuilder) Errors() error {
	return b.globalErrors

}

// Build builds the variable and resets inner structure
func (b *baseVariableBuilder) Build() (Variable, error) {
	if b.globalErrors != nil {
		return nil, b.globalErrors
	} else if b.name == "" {
		return nil, errors.New("variable name cannot be empty")
	} else if len(b.acceptedClasses) == 0 {
		return nil, errors.New("variable must accept at least one class")
	}

	accepted := make(map[Class]bool)
	for class, value := range b.acceptedClasses {
		if value {
			accepted[class] = true
		}
	}

	result := &baseVariable{
		name:            b.name,
		acceptedClasses: accepted,
	}

	b.acceptedClasses = nil
	b.acceptedClasses = make(map[Class]bool)
	b.globalErrors = nil
	if b.name == "" {
		b.globalErrors = errors.New("variable name cannot be empty")
	}

	return result, nil
}

// NewVariableBuilder creates a new variable builder with the given name.
// Name should not be empty.
func NewVariableBuilder(name string) VariableBuilder {
	result := &baseVariableBuilder{
		name:            name,
		acceptedClasses: make(map[Class]bool),
	}

	if name == "" {
		result.globalErrors = errors.New("variable name cannot be empty")
	}

	return result
}
