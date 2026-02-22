package maths

// FloatNumber defines a type constraint for generic functions.
// The tilde (~) symbol ensures that the interface matches not just float32 and float64,
// but also any custom types derived from them (e.g., type MyFloat float64).
type FloatNumber interface {
	~float64 | ~float32
}

// These constants define the precision thresholds for comparisons.
// Since floating-point arithmetic can result in tiny rounding errors (e.g., 0.1 + 0.2 != 0.3),
// we check if the difference between two numbers is smaller than an "Epsilon."
const (
	// SHORT_EPSILON is used for float32, which has about 7 decimal digits of precision.
	SHORT_EPSILON = 1e-5
	// LONG_EPSILON is used for float64, which has about 15-17 decimal digits of precision.
	LONG_EPSILON = 1e-9
)

// isFloat64 is a helper function that uses type assertion to determine
// if the generic value passed is specifically a float64.
func isFloat64[F FloatNumber](f F) bool {
	// We convert to 'any' (interface{}) to perform a type switch on the underlying type.
	switch any(f).(type) {
	case float64:
		return true
	default:
		return false
	}
}

// equalsFloats acts as the main dispatcher. It detects the precision of the
// input types and routes the comparison to the appropriate epsilon-based logic.
func equalsFloats[F FloatNumber](a, b F) bool {
	// Check if both inputs are float64.
	// Note: In this specific implementation, if even one is float32, it defaults to float32 precision.
	isA64, isB64 := isFloat64(a), isFloat64(b)

	switch isA64 && isB64 {
	case true:
		return equalsFloat64(a, b)
	default:
		return equalsFloat32(a, b)
	}
}

// equalsFloat32 compares two numbers using the SHORT_EPSILON (1e-5).
// It calculates the absolute difference between 'a' and 'b' manually.
func equalsFloat32[F FloatNumber](a, b F) bool {
	if a < b {
		return b-a < SHORT_EPSILON
	} else {
		return a-b < SHORT_EPSILON
	}
}

// equalsFloat64 compares two numbers using the LONG_EPSILON (1e-9).
// It provides a much stricter check than the float32 version.
func equalsFloat64[F FloatNumber](a, b F) bool {
	if a < b {
		return b-a < LONG_EPSILON
	} else {
		return a-b < LONG_EPSILON
	}
}
