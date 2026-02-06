package commons

import "time"

// TIME_FORMAT defines how to serialize and deserialize time data
const TIME_FORMAT = time.RFC3339

// TIME_PRECISION is the accepted thresold to define when two times are the same
const TIME_PRECISION = time.Second

// EPSILON is the accepted margin of error for floating point comparisons.
// It accounts for small precision errors inherent in floating point arithmetic.
const EPSILON = 1e-9
