package commons

type PrimitiveType string

const (
	StringType  PrimitiveType = "string"
	IntType     PrimitiveType = "int"
	FloatType   PrimitiveType = "float"
	ListsType   PrimitiveType = "list"
	DefaultType PrimitiveType = "default"
)

func (p PrimitiveType) EqualsPrimitive(a, b any) bool {
	return a == b
}

func (p PrimitiveType) Name() string {
	return string(p)
}

func (p PrimitiveType) EqualsPrimitiveType(other PrimitiveType) bool {
	return p == other
}

func GuessPrimitiveType(value any) PrimitiveType {
	if _, okInt := value.(int); okInt {
		return IntType
	} else if _, okFloat := value.(float64); okFloat {
		return FloatType
	} else if _, okString := value.(string); okString {
		return StringType
	} else if _, okList := value.([]any); okList {
		return ListsType
	}

	return DefaultType
}

func DefaultPrimitiveType() PrimitiveType {
	return DefaultType
}
