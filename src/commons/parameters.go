package commons

// FormalParameters is the expected parameters definition of an action
type FormalParameters interface {
	// Accepts returns true if content matches parameters definition
	Accepts(Content) bool
}

// permissiveParameters accepts any content
type permissiveParameters struct{}

// Accepts returns true no matter what
func (p permissiveParameters) Accepts(c Content) bool {
	return true
}

// NewMostPermissiveFormalParameters returns the most permissive parameters: they accept anything
func NewMostPermissiveFormalParameters() FormalParameters {
	return permissiveParameters{}
}

// expectedVariables defines the minimum set of variables to have.
// For instance, if we expect x and y, and content has x,y,z then we accept
type expectedVariables struct {
	// mandatories contains the set of variables to have at least
	mandatories []string
}

// Accepts returns true if c has the expected variables
func (e expectedVariables) Accepts(c Content) bool {
	if c == nil {
		return len(e.mandatories) == 0
	} else if len(e.mandatories) == 0 {
		return true
	}

	contentVariables := c.Variables()
	if len(contentVariables) == 0 {
		return false
	}

	// we want all the variables from e.mandatories to be included in contentVariables.
	// If contentVariables has more variables, it is OK.
	// But it has to have at least each variable from e.mandatories.
	// So, it means that e.mandatories should be in contentVariables

	return SlicesContainsAllFunc(contentVariables, e.mandatories, func(a, b string) bool { return a == b })
}

// NewVariablesFormalParameters returns a new FormalParameters accepting at least a set of variables
func NewVariablesFormalParameters(names []string) FormalParameters {
	return expectedVariables{mandatories: SliceDeduplicate(names)}
}

// expectedPositionalMinimalSize is the minimal size of positional elements to get
type expectedPositionalMinimalSize struct {
	threshold int
}

// Accepts returns true if c has the minimal size for its positional parameters
func (e expectedPositionalMinimalSize) Accepts(c Content) bool {
	if c == nil {
		return e.threshold == 0
	} else if e.threshold == 0 {
		return true
	}

	return c.Size() >= e.threshold
}

// NewPositionalFormalParameters returns formal parameters accepting at least minimalSize elements
func NewPositionalFormalParameters(minimalSize int) FormalParameters {
	return expectedPositionalMinimalSize{threshold: minimalSize}
}

// acceptsUniqueValueInContent returns true if there is an unique value in the content
type acceptsUniqueValueInContent struct{}

// Accepts returns true if, no matter whether it is positional or named, there is one value in the content
func (a acceptsUniqueValueInContent) Accepts(c Content) bool {
	if c == nil {
		return false
	}

	_, matches := c.Unique()
	return matches
}

// NewUniqueFormalParameters returns a formal parameters accepting only one unique value
func NewUniqueFormalParameters() FormalParameters {
	return acceptsUniqueValueInContent{}
}
