package structures

// TemporalOperator defines binay operator working on a period compared to a reference period.
// Operators do NOT commute in general.
// It means that first operand HAS TO BE the current period to test whereas second operand is the REFERENCE period.
type TemporalOperator int

// TemporalEquals tests if current equals reference
const TemporalEquals TemporalOperator = 1

// TemporalCommonPoint tests if current and reference have at least a common point
const TemporalCommonPoint TemporalOperator = 2

// TemporalAlwaysAccept always accepts no matter current period
const TemporalAlwaysAccept TemporalOperator = 3

// TemporalAlwaysRefuse always refuses no matter current period
const TemporalAlwaysRefuse TemporalOperator = 4

// TemporalReferenceContains tests if current is included in reference
const TemporalReferenceContains TemporalOperator = 5

// Accepts executes the operator on current and reference (in that order)
func (t TemporalOperator) Accepts(current Period, reference Period) bool {
	switch t {
	case TemporalAlwaysAccept:
		return true
	case TemporalAlwaysRefuse:
		return false
	case TemporalCommonPoint:
		return !current.Intersection(reference).IsEmpty()
	case TemporalEquals:
		return current.Equals(reference)
	case TemporalReferenceContains:
		return current.IsIncludedIn(reference)
	default:
		return false
	}
}
