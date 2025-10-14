package commons

// FormalParameters is the expected parameters definition of an action
type FormalParameters struct {
	minimalSize           int      // minimal TOTAL size (including variables)
	minimalPositionalSize int      // minimal number of expected positional values
	expectedVariables     []string // expected variables by name
}

// NewNamedFormalParameters constructs formal parameters accepting those variables at least
func NewNamedFormalParameters(names []string) FormalParameters {
	variables := SliceDeduplicate(names)
	return FormalParameters{
		expectedVariables: variables,
		minimalSize:       len(variables),
	}
}

// NewPositionalFormalParameters constructs formal parameters expecting size positional elements
func NewPositionalFormalParameters(size int) FormalParameters {
	return FormalParameters{
		minimalSize:           size,
		minimalPositionalSize: size,
	}
}

// NewMostPermissiveFormalParameters accepts everything
func NewMostPermissiveFormalParameters() FormalParameters {
	return FormalParameters{}
}

// Variables returns expected variables
func (fp FormalParameters) Variables() []string {
	var result []string
	result = append(result, fp.expectedVariables...)
	return result
}

// Accepts tests if content matches expected constraints:
// Enough total values, enough positional values, and enough variables.
// Failure to comply to any condition returns false
func (fp FormalParameters) Accepts(c Content) bool {
	size := 0
	if c != nil {
		size = c.Size()
	}

	totalSize := size
	var variables []string
	if c != nil {
		variables = c.Variables()
		totalSize = totalSize + len(variables)
	}

	if fp.minimalSize > totalSize {
		return false
	} else if fp.minimalPositionalSize > size {
		return false
	} else if len(fp.expectedVariables) > len(variables) {
		return false
	} else {
		return SlicesContainsAllFunc(variables, fp.expectedVariables, func(a, b string) bool { return a == b })
	}
}

// Max gets the union of conditions to accept content:
// expects max of sizes, max of positional sizes, all expected variables
func (fp FormalParameters) Max(other FormalParameters) FormalParameters {
	result := FormalParameters{}

	// get max of minimal sizes
	result.minimalSize = max(fp.minimalSize, other.minimalSize)

	// get max of positional sizes
	result.minimalPositionalSize = max(fp.minimalPositionalSize, other.minimalPositionalSize)

	// get union of variables
	variables := make(map[string]bool)
	for _, v := range fp.expectedVariables {
		variables[v] = true
	}

	for _, v := range other.expectedVariables {
		variables[v] = true
	}

	for k := range variables {
		result.expectedVariables = append(result.expectedVariables, k)
	}

	return result
}

// parametersCombine generalizes parameter max to a finite set of parameters
func parametersCombine(parameters []FormalParameters) FormalParameters {
	if len(parameters) == 0 {
		return FormalParameters{}
	}

	value := parameters[0]
	for index, p := range parameters {
		if index >= 1 {
			value = value.Max(p)
		}
	}

	return value
}
