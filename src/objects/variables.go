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

type VariableBuilder interface {
	WithName(string) VariableBuilder
	WithAcceptedClass(class Class) VariableBuilder
	WithoutAcceptedClass(class Class) VariableBuilder
	Errors() error
	Build() (Variable, error)
}

type baseVariable struct {
	name            string
	acceptedClasses map[Class]bool
}

func (v *baseVariable) Accepts() []Class {
	result := make([]Class, 0, len(v.acceptedClasses))
	for acceptedClass := range v.acceptedClasses {
		result = append(result, acceptedClass)
	}

	slices.Sort(result)
	return result
}

func (v *baseVariable) AcceptsOneOf(acceptedClasses ...Class) bool {
	for _, acceptedClass := range acceptedClasses {
		if v.acceptedClasses[acceptedClass] {
			return true
		}
	}
	return false
}

func (v *baseVariable) Id() string {
	return v.name
}

func (v *baseVariable) DeclaringClass() Class {
	return CLASS_VARIABLE
}

func (v *baseVariable) Name() string {
	return v.name
}

func (v *baseVariable) isLinkable() bool {
	return true
}

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

type baseVariableBuilder struct {
	name            string
	acceptedClasses map[Class]bool
	globalErrors    error
}

func (b *baseVariableBuilder) WithName(name string) VariableBuilder {
	if name == "" {
		b.globalErrors = errors.New("variable name cannot be empty")
		return b
	}

	b.name = name
	return b
}

func (b *baseVariableBuilder) WithAcceptedClass(class Class) VariableBuilder {
	var empty Class
	if class == empty {
		b.globalErrors = errors.New("accepted class cannot be nil")
		return b
	}

	if b.acceptedClasses == nil {
		b.acceptedClasses = make(map[Class]bool)
	}

	b.acceptedClasses[class] = true
	return b
}
func (b *baseVariableBuilder) WithoutAcceptedClass(class Class) VariableBuilder {
	delete(b.acceptedClasses, class)
	return b
}

func (b *baseVariableBuilder) Errors() error {
	return b.globalErrors

}

func (b *baseVariableBuilder) Build() (Variable, error) {
	if b.globalErrors != nil {
		return nil, b.globalErrors
	} else if b.name == "" {
		return nil, errors.New("variable name cannot be empty")
	} else if len(b.acceptedClasses) == 0 {
		return nil, errors.New("variable must accept at least one class")
	}

	return &baseVariable{
		name:            b.name,
		acceptedClasses: b.acceptedClasses,
	}, nil
}

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
