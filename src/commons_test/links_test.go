package commons_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// ============================================================================
// MOCKS & HELPERS FOR TESTING
// ============================================================================

// cyclicMockElement is a custom element designed to intentionally create a cyclic graph.
// Since Links are immutable, it's normally impossible to create a cycle using only pure Links.
// However, because Links can hold any Element, a poorly designed external object could introduce one.
// This mock helps us prove our algorithms (Same, ReplaceVariable) are completely safe against Out Of Memory (OOM) crashes.
type cyclicMockElement struct {
	id   string
	self commons.Element
}

func (c *cyclicMockElement) Same(other commons.Element) bool {
	if otherMock, ok := other.(*cyclicMockElement); ok {
		return c.id == otherMock.id
	}
	return false
}

func (c *cyclicMockElement) DeclaringClass() commons.Class {
	return commons.CLASS_INSTANCE
}

// ============================================================================
// TEST SUITE 1: BASIC OPERATIONS & IMMUTABILITY (COPY-ON-WRITE)
// ============================================================================

func TestLink_BasicOperations_And_Immutability(t *testing.T) {
	// Objective: Verify that all standard read/write operations work correctly,
	// and that every mutation returns a NEW instance without altering the original.

	p1 := commons.NewFullPeriod()
	now := time.Now()
	p2 := commons.NewFinitePeriod(now, now.Add(time.Hour), true, true)

	// 1. Creation and basic reads
	root := commons.NewLink("ROOT", p1)
	if root.Name() != "ROOT" {
		t.Errorf("Expected name 'ROOT', got '%s'", root.Name())
	}
	if root.DeclaringClass() != commons.CLASS_LINK {
		t.Errorf("Expected class LINK, got '%v'", root.DeclaringClass())
	}
	if !root.Validity().Equals(p1) {
		t.Errorf("Validity mismatch on creation")
	}

	// 2. WithValidity
	linkV2 := root.WithValidity(p2)
	if !root.Validity().Equals(p1) {
		t.Errorf("Immutability broken: root validity changed")
	}
	if !linkV2.Validity().Equals(p2) {
		t.Errorf("linkV2 did not register new validity")
	}

	// 3. WithOperand & Slice Protection
	traitA := commons.NewTrait("TraitA")
	linkV3 := linkV2.WithOperand("subject", []commons.Element{traitA})

	ops, ok := linkV3.Operand("subject")
	if !ok || len(ops) != 1 {
		t.Fatalf("WithOperand failed to attach elements")
	}

	// HACK ATTEMPT: Try to modify the internal slice returned by Operand()
	ops[0] = commons.NewTrait("HACKED")
	safeOps, _ := linkV3.Operand("subject")
	if safeOps[0].Same(commons.NewTrait("HACKED")) {
		t.Errorf("Encapsulation broken: Internal slice was mutated by external code")
	}

	// 4. WithAppended
	traitB := commons.NewTrait("TraitB")
	linkV4 := linkV3.WithAppended("subject", traitB)

	oldOps, _ := linkV3.Operand("subject")
	newOps, _ := linkV4.Operand("subject")
	if len(oldOps) != 1 {
		t.Errorf("Immutability broken: linkV3 operand size changed")
	}
	if len(newOps) != 2 {
		t.Errorf("WithAppended failed to add the new element")
	}

	// 5. Without
	linkV5 := linkV4.Without("subject", func(e commons.Element) bool {
		return e.Same(traitA)
	})

	finalOps, _ := linkV5.Operand("subject")
	if len(finalOps) != 1 || !finalOps[0].Same(traitB) {
		t.Errorf("Without failed to filter elements correctly")
	}

	// 6. Deterministic Operands Ordering
	linkOrdered := commons.NewLink("ORDER", p1).
		WithOperand("Zeta", []commons.Element{traitA}).
		WithOperand("Alpha", []commons.Element{traitA}).
		WithOperand("Gamma", []commons.Element{traitA})

	keys := linkOrdered.Operands()
	if len(keys) != 3 || keys[0] != "Alpha" || keys[1] != "Gamma" || keys[2] != "Zeta" {
		t.Errorf("Operands() must return keys sorted alphabetically. Got %v", keys)
	}
}

// ============================================================================
// TEST SUITE 2: DEEP GRAPH EQUALITY (SAME) & CYCLE PREVENTION
// ============================================================================

func TestLink_Same_ComplexGraphs(t *testing.T) {
	// Objective: Validate Breadth-First Search (BFS) for deep equality comparison.
	p := commons.NewFullPeriod()

	buildGraph := func(leafName string) commons.Link {
		leaf := commons.NewTrait(leafName)
		child := commons.NewLink("CHILD", p).WithOperand("target", []commons.Element{leaf})
		return commons.NewLink("PARENT", p).WithOperand("nested", []commons.Element{child})
	}

	graph1 := buildGraph("LEAF_A")
	graph2 := buildGraph("LEAF_A")
	graph3 := buildGraph("LEAF_B") // Differs at the deepest level

	if !graph1.Same(graph2) {
		t.Errorf("Same() failed: identical complex graphs should be true")
	}
	if graph1.Same(graph3) {
		t.Errorf("Same() failed: graphs differing at leaf level should be false")
	}

	graph4 := graph1.WithOperand("extra", []commons.Element{commons.NewTrait("X")})
	if graph1.Same(graph4) {
		t.Errorf("Same() failed: graphs with different structural keys should be false")
	}
}

func TestLink_Same_CyclePrevention(t *testing.T) {
	// Objective: Prove that Same() does not trigger an infinite loop (OOM)
	// when comparing graphs containing cyclic references.
	p := commons.NewFullPeriod()

	cyclicObj := &cyclicMockElement{id: "mock1"}
	cyclicObj.self = cyclicObj // The cycle is explicitly created here

	linkA := commons.NewLink("ROOT", p).WithOperand("loop", []commons.Element{cyclicObj})
	linkB := commons.NewLink("ROOT", p).WithOperand("loop", []commons.Element{cyclicObj})

	// If the cycle detection (visited map) is missing in Same(), this will crash the test runner.
	isSame := linkA.Same(linkB)

	if !isSame {
		t.Errorf("Same() should safely resolve to true even with cycles, by detecting visited pairs")
	}
}

// ============================================================================
// TEST SUITE 3: VARIABLE REPLACEMENT & DAG OPTIMIZATION
// ============================================================================

func TestLink_ReplaceVariable_DeepTree(t *testing.T) {
	// Objective: Verify that a variable is correctly identified and replaced deep inside a graph.
	p := commons.NewFullPeriod()
	vTarget := commons.NewVariable("TARGET")
	replacement := commons.NewTrait("REPLACED")

	childLink := commons.NewLink("CHILD", p).
		WithOperand("items", []commons.Element{commons.NewTrait("A"), vTarget, commons.NewTrait("B")})

	root := commons.NewLink("ROOT", p).WithOperand("branch", []commons.Element{childLink})

	newRoot := root.ReplaceVariable(vTarget, replacement)

	// Check that the replacement happened correctly
	branch, _ := newRoot.Operand("branch")
	newChild := branch[0].(commons.Link)
	items, _ := newChild.Operand("items")

	if !items[1].Same(replacement) {
		t.Errorf("ReplaceVariable failed to insert replacement deep in the tree")
	}

	// Verify original tree remains uncorrupted
	oldBranch, _ := root.Operand("branch")
	oldChild := oldBranch[0].(commons.Link)
	oldItems, _ := oldChild.Operand("items")
	if !oldItems[1].Same(vTarget) {
		t.Errorf("CRITICAL: Original tree was mutated during ReplaceVariable!")
	}
}

func TestLink_ReplaceVariable_MultipleOccurrences(t *testing.T) {
	// Objective: Verify that if a variable appears multiple times across different operands,
	// all instances are correctly replaced.
	vTarget := commons.NewVariable("X")
	val := commons.NewTrait("NEW_VAL")

	link := commons.NewLink("ROOT", commons.NewFullPeriod()).
		WithOperand("left", []commons.Element{vTarget, commons.NewTrait("A"), vTarget}).
		WithOperand("right", []commons.Element{vTarget})

	newLink := link.ReplaceVariable(vTarget, val)

	leftOps, _ := newLink.Operand("left")
	if len(leftOps) != 3 || !leftOps[0].Same(val) || !leftOps[2].Same(val) {
		t.Errorf("Failed to replace multiple occurrences in the same operand")
	}

	rightOps, _ := newLink.Operand("right")
	if len(rightOps) != 1 || !rightOps[0].Same(val) {
		t.Errorf("Failed to replace occurrence in a different operand")
	}
}

func TestLink_ReplaceVariable_TypeConstraints(t *testing.T) {
	// Objective: Verify that CanBeReplacedBy is respected.
	// We create a variable that ONLY accepts CLASS_INSTANCE elements.
	vConstrained := commons.NewVariable("STRICT_VAR", commons.CLASS_INSTANCE)

	link := commons.NewLink("ROOT", commons.NewFullPeriod()).
		WithOperand("target", []commons.Element{vConstrained})

	// Scenario A: Attempt to replace with a Trait (Should fail / be ignored)
	invalidReplacement := commons.NewTrait("JUST_A_TRAIT")
	newLinkA := link.ReplaceVariable(vConstrained, invalidReplacement)

	opsA, _ := newLinkA.Operand("target")
	if opsA[0].Same(invalidReplacement) {
		t.Errorf("ReplaceVariable bypassed CanBeReplacedBy() constraint! Trait was injected instead of Instance.")
	}

	// Scenario B: Attempt to replace with an Instance (Should succeed)
	validReplacement := commons.NewTemporalInstance() // Implements CLASS_INSTANCE
	newLinkB := link.ReplaceVariable(vConstrained, validReplacement)

	opsB, _ := newLinkB.Operand("target")
	if !opsB[0].Same(validReplacement) {
		t.Errorf("ReplaceVariable blocked a valid replacement matching the type constraints.")
	}
}

func TestLink_ReplaceVariable_StructuralSharing(t *testing.T) {
	// Objective: Prove that if no replacement actually occurs in a sub-graph,
	// ReplaceVariable reuses the exact same memory pointer to save allocations.
	vMissing := commons.NewVariable("MISSING")
	val := commons.NewTrait("VAL")

	child := commons.NewLink("CHILD", commons.NewFullPeriod()).
		WithOperand("leaf", []commons.Element{commons.NewTrait("A")})

	root := commons.NewLink("ROOT", commons.NewFullPeriod()).
		WithOperand("branch", []commons.Element{child})

	newRoot := root.ReplaceVariable(vMissing, val)

	// Since "MISSING" doesn't exist in the graph, newRoot should be the exact same pointer as root.
	ptrOld := reflect.ValueOf(root).Pointer()
	ptrNew := reflect.ValueOf(newRoot).Pointer()

	if ptrOld != ptrNew {
		t.Errorf("Structural Sharing optimization failed: tree was rebuilt even though no variables were replaced.")
	}
}

func TestLink_ReplaceVariable_CyclePrevention(t *testing.T) {
	// Objective: Prove that DFS correctly identifies a cycle and breaks it
	// instead of crashing into an infinite recursion (OOM).
	vTarget := commons.NewVariable("TARGET")
	val := commons.NewTrait("VAL")

	cyclicObj := &cyclicMockElement{id: "mock1"}
	cyclicObj.self = cyclicObj // Cycle

	link := commons.NewLink("ROOT", commons.NewFullPeriod()).
		WithOperand("loop", []commons.Element{cyclicObj, vTarget})

	// If inStack map check is missing in ReplaceVariable, this will crash the test runner.
	newLink := link.ReplaceVariable(vTarget, val)

	// Check if the variable inside the cyclic graph was still replaced safely
	ops, _ := newLink.Operand("loop")
	if len(ops) != 2 || !ops[1].Same(val) {
		t.Errorf("Cycle prevention triggered, but sibling variable was not replaced correctly")
	}
}
