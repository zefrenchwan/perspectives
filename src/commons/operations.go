package commons

// Operation is the most general definition of a way to change a state
type Operation interface {
}

// Operations groups operations.
// There is no specific order to execute each operation
type Operations struct {
	// Content are the operations to perform
	Content []Operation
}
