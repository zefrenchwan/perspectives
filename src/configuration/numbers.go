package configuration

// These constants define the precision thresholds for comparisons.
// Since floating-point arithmetic can result in tiny rounding errors (e.g., 0.1 + 0.2 != 0.3),
// we check if the difference between two numbers is smaller than an "Epsilon."
const (
	// SHORT_EPSILON is used for float32, which has about 7 decimal digits of precision.
	SHORT_EPSILON = 1e-5
	// LONG_EPSILON is used for float64, which has about 15-17 decimal digits of precision.
	LONG_EPSILON = 1e-9
)
