package commons

// Link is a relation between elements.
// Think of it as a node in a graph that connects to other elements (Traits, Variables, or nested Links).
type Link interface {
	Element // Inherits Same(Element) and DeclaringClass()

	// --- READ OPERATIONS : State queries ---

	// Name returns the name of the link. For instance, in Loves(subject=John, object=Pizza) => "Loves"
	Name() string
	// Validity returns the time period during which this link is considered true/active.
	Validity() Period
	// Operands returns all operand keys (the names of the relationships, e.g., "subject", "object") deterministically sorted.
	Operands() []string
	// Operand returns the slice of Elements associated with the given operand name.
	// It returns a boolean indicating if the operand key exists.
	Operand(name string) ([]Element, bool)

	// --- FUNCTIONAL MUTATIONS : Copy-on-write operations ---
	// Since Links are immutable, these methods never modify the current instance.
	// They return a newly allocated Link with the requested changes.

	// WithValidity returns a copy of the link with the new validity period.
	WithValidity(p Period) Link

	// WithOperand returns a copy of the link with the given operand forced (overwrites previous values).
	WithOperand(name string, operands []Element) Link

	// WithAppended returns a copy of the link with a new element added to the specified operand.
	WithAppended(name string, operand Element) Link

	// Without returns a copy of the link, filtering out operand values that match the condition.
	// For instance, Loves(John, [Pizza, Salad]) => Without("object", isPizza) => Loves(John, [Salad])
	Without(name string, op func(linkable Element) bool) Link
}

// VisitLink traverses a link graph using Breadth-First Search (BFS),
// calling the provided function on each visited link.
// It explicitly handles cycle detection to prevent infinite loops (OOM panics).
func VisitLink(link Link, onLink func(Link)) {
	if onLink == nil || link == nil {
		return
	}

	linksToExplore := []Link{link}
	visited := make(map[Link]bool)
	visited[link] = true

	for len(linksToExplore) != 0 {
		currentLink := linksToExplore[0]
		linksToExplore = linksToExplore[1:]

		onLink(currentLink)

		for _, name := range currentLink.Operands() {
			operands, _ := currentLink.Operand(name)

			for _, operand := range operands {
				if operand == nil {
					continue
				}

				if IsElementDeclaredInstance(operand, CLASS_LINK) {
					childLink, ok := operand.(Link)
					if !ok {
						continue
					}

					// cycle management
					if !visited[childLink] {
						visited[childLink] = true
						linksToExplore = append(linksToExplore, childLink)
					}
				}
			}
		}
	}
}

// ============================================================================
// PATTERN MATCHING (VIA BFS)
// ============================================================================

// Match checks if a 'target' graph deeply matches a 'pattern' graph.
// The pattern may contain Variables (e.g., ?X, ?Y) acting as wildcards.
// It returns a Substitution map containing the captured variables and a boolean indicating success.
func Match(pattern, target Element) (Substitution, bool) {
	if pattern == nil && target == nil {
		return nil, true
	}
	if pattern == nil || target == nil {
		return nil, false
	}

	bindings := make(Substitution)
	type pair struct {
		p Element
		t Element
	}

	queue := []pair{{pattern, target}}
	visited := make(map[[2]uintptr]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		p, t := curr.p, curr.t

		if p == nil && t == nil {
			continue
		}
		if p == nil || t == nil {
			return nil, false
		}

		// RULE 1: VARIABLE (WILDCARD)
		if IsElementDeclaredInstance(p, CLASS_VARIABLE) {
			v, ok := p.(Variable)
			if !ok {
				return nil, false // Garde fou de type
			}

			if existingValue, exists := bindings[v.name]; exists {
				queue = append(queue, pair{existingValue, t})
			} else {
				bindings[v.name] = t
			}
			continue
		}

		if p.DeclaringClass() != t.DeclaringClass() {
			return nil, false
		}

		// RULE 2: LINKS (GRAPH TRAVERSAL)
		if IsElementDeclaredInstance(p, CLASS_LINK) {
			lp, okP := p.(Link)
			lt, okT := t.(Link)
			if !okP || !okT {
				return nil, false
			}

			// --- SECURE CYCLE DETECTION ---
			ptrP, isPtrP := safePointer(lp)
			ptrT, isPtrT := safePointer(lt)

			if isPtrP && isPtrT {
				memPair := [2]uintptr{ptrP, ptrT}
				if visited[memPair] {
					continue
				}
				visited[memPair] = true
			}

			// Fast Scalar Checks
			if lp.Name() != lt.Name() || !lp.Validity().Equals(lt.Validity()) {
				return nil, false
			}

			opsP := lp.Operands()
			opsT := lt.Operands()
			if len(opsP) != len(opsT) {
				return nil, false
			}

			for i, opName := range opsP {
				if opsT[i] != opName {
					return nil, false
				}
				elsP, _ := lp.Operand(opName)
				elsT, _ := lt.Operand(opName)

				if len(elsP) != len(elsT) {
					return nil, false
				}
				for j := range elsP {
					queue = append(queue, pair{elsP[j], elsT[j]})
				}
			}
		} else {
			// RULE 3: TERMINAL ELEMENTS
			if !p.Same(t) {
				return nil, false
			}
		}
	}
	return bindings, true
}

// ============================================================================
// INSTANTIATION (ITERATIVE DFS)
// ============================================================================

// Instantiate takes a Pattern (Link) and applies the Substitution bindings to it,
// replacing all Variables with their bound targets.
func Instantiate(pattern Link, bindings Substitution) Link {
	if pattern == nil {
		return nil
	}
	if len(bindings) == 0 {
		return pattern
	}

	type frame struct {
		link       Link
		keys       []string
		keyIdx     int
		valIdx     int
		newOps     map[string][]Element
		currVals   []Element
		hasChanges bool
	}

	stack := []*frame{{
		link:   pattern,
		keys:   pattern.Operands(),
		newOps: make(map[string][]Element),
	}}

	memo := make(map[uintptr]Link)
	inStack := make(map[uintptr]bool)

	// Init root pointer tracking securely
	rootPtr, isRootPtr := safePointer(pattern)
	if isRootPtr {
		inStack[rootPtr] = true
	}

	var result Element

	for len(stack) > 0 {
		curr := stack[len(stack)-1]

		// --- ASCENT PHASE ---
		if curr.keyIdx >= len(curr.keys) {
			var resultLink Link

			if !curr.hasChanges {
				resultLink = curr.link
			} else {
				resultLink = &baseLink{
					name:     curr.link.Name(),
					validity: curr.link.Validity(),
					operands: curr.newOps,
				}
			}

			stack = stack[:len(stack)-1]

			// Update DAG caches if it was a pointer
			if ptr, isPtr := safePointer(curr.link); isPtr {
				delete(inStack, ptr)
				memo[ptr] = resultLink
			}

			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.currVals = append(parent.currVals, resultLink)
				if resultLink != curr.link {
					parent.hasChanges = true
				}
			} else {
				result = resultLink
			}
			continue
		}

		// --- DESCENT PHASE ---
		key := curr.keys[curr.keyIdx]
		vals, _ := curr.link.Operand(key)

		if curr.valIdx >= len(vals) {
			if len(curr.currVals) > 0 {
				curr.newOps[key] = curr.currVals
			}
			curr.currVals = nil
			curr.keyIdx++
			curr.valIdx = 0
			continue
		}

		elem := vals[curr.valIdx]
		curr.valIdx++

		// SCENARIO 1: Variable Substitution
		if elem != nil && IsElementDeclaredInstance(elem, CLASS_VARIABLE) {
			if v, ok := elem.(Variable); ok {
				if value, exists := bindings[v.name]; exists && v.CanBeReplacedBy(value) {
					curr.currVals = append(curr.currVals, value)
					curr.hasChanges = true
					continue
				}
			}
			curr.currVals = append(curr.currVals, elem)
			continue
		}

		// SCENARIO 2: Nested Link
		if elem != nil && IsElementDeclaredInstance(elem, CLASS_LINK) {
			if childLink, ok := elem.(Link); ok {

				ptr, isPtr := safePointer(childLink)

				// Apply cycle and DAG logic strictly for pointers
				if isPtr {
					if inStack[ptr] {
						curr.currVals = append(curr.currVals, childLink)
						continue
					}
					if cachedLink, exists := memo[ptr]; exists {
						curr.currVals = append(curr.currVals, cachedLink)
						if cachedLink != childLink {
							curr.hasChanges = true
						}
						continue
					}
					inStack[ptr] = true
				}

				stack = append(stack, &frame{
					link:   childLink,
					keys:   childLink.Operands(),
					newOps: make(map[string][]Element),
				})
				continue
			}
		}

		// SCENARIO 3: Terminal Elements
		curr.currVals = append(curr.currVals, elem)
	}

	if result == nil {
		return nil
	}
	return result.(Link)
}
