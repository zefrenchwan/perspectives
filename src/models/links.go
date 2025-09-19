package models

import (
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// linkValue is just a container for values in link.
// Here is the problem:
// We want to have an unique representation of a node in a link even if values matches.
// For instance, consider the link:
// Is(Cheese, Cheese)
// Values are the same, nodes are not
type linkValue struct {
	// uniqueId, assumed to be unique in the link
	uniqueId string
	// content is the actual value of the node in the link
	content ModelEntity
}

// contentType returns the type of the underlying content
func (v linkValue) contentType() EntityType {
	return v.content.GetType()
}

// newLinkValue builds a new node in a link based on a content
func newLinkValue(content ModelEntity) linkValue {
	return linkValue{
		uniqueId: uuid.NewString(),
		content:  content,
	}
}

// newLinkValueForObjects builds a link value as a group of objects
func newLinkValueForObjects(values []Object) linkValue {
	return linkValue{uuid.NewString(), objectsGroup(values)}
}

// Link will link objects together (0 level links) or links and object (higher level links).
// For instance Likes(Steve, Tiramisu) is a basic link and Knows(Paul, Likes(Steve, Tiramisu)) is an higher level link.
type Link struct {
	// unique id of the link
	id string
	// name defines the link semantic
	name string
	// operands are role based operands.
	// Usually, roles are "subject" or "object" or ...
	operands map[string]linkValue
	// Lifetime is the duration of the link
	lifetime structures.Period
}

// localCopy builds a new link containing copies of its direct values (no recursive walkthrough)
func (l *Link) localCopy() *Link {
	if l == nil {
		return nil
	}

	// build a local copy of the operands to ensure new link value ids
	operandCopy := make(map[string]linkValue)
	for role, value := range l.operands {
		operandCopy[role] = newLinkValue(value.content)
	}

	// and done
	result := new(Link)
	result.id = l.id
	result.name = l.name
	result.lifetime = l.lifetime
	result.operands = operandCopy
	return result
}

// GetType returns EntityTypeLink
func (l *Link) GetType() EntityType {
	return EntityTypeLink
}

// AsLink returns the link
func (l *Link) AsLink() (*Link, error) {
	return l, nil
}

// AsGroup raises an error
func (l *Link) AsGroup() ([]Object, error) {
	return nil, errors.ErrUnsupported
}

// AsObject raises an error
func (l *Link) AsObject() (*Object, error) {
	return nil, errors.ErrUnsupported
}

// AsTrait raises an error
func (l *Link) AsTrait() (Trait, error) {
	return Trait{}, errors.ErrUnsupported
}

// AsVariable raises an error
func (l *Link) AsVariable() (Variable, error) {
	return Variable{}, errors.ErrUnsupported
}

// RoleSubject is the constant value for the subject role
const RoleSubject = "subject"

// RoleObject is the constant value for the object role
const RoleObject = "object"

// NewLink builds a link, valid for a given period
// name is the semantic of that link (for instance "loves" or "knows")
// values are the values (role linked to operand)
// duration is the period the link is valid for
//
// Although formal parameter is any, expected types are:
// 1. Slices of objects
// 2. Objects
// 3. Links
// 4. Traits
// 5. Variables representing previous mentioned types
//
// An error will raise if values do not match that constraint
func NewLink(name string, values map[string]any, duration structures.Period) (Link, error) {
	var link, empty Link
	link.id = uuid.NewString()
	link.name = name
	link.operands = make(map[string]linkValue)
	link.lifetime = duration

	for role, operand := range values {
		if operand == nil {
			continue
		} else if l, ok := operand.(Link); ok {
			link.operands[role] = newLinkValue(&l)
		} else if lp, ok := operand.(*Link); ok {
			link.operands[role] = newLinkValue(lp)
		} else if g, ok := operand.([]Object); ok {
			link.operands[role] = newLinkValueForObjects(g)
		} else if o, ok := operand.(Object); ok {
			link.operands[role] = newLinkValue(&o)
		} else if op, ok := operand.(*Object); ok {
			link.operands[role] = newLinkValue(op)
		} else if t, ok := operand.(Trait); ok {
			link.operands[role] = newLinkValue(t)
		} else if v, ok := operand.(Variable); ok {
			link.operands[role] = newLinkValue(v)
		} else {
			return empty, fmt.Errorf("unsupported type for role %s. Expecting either trait or object or link or group of objects", role)
		}
	}

	return link, nil
}

// NewSimpleLink is a shortcut to declare a link(subject, object) valid for the full time
func NewSimpleLink(link string, subject, object any) (Link, error) {
	return NewLink(link, map[string]any{RoleSubject: subject, RoleObject: object}, structures.NewFullPeriod())
}

// CopyStructure clones a link.
// It copies the structure of the node, but uses the same content for anything but links.
// To rephrase: copies all the links, keep the rest as is.
func (l *Link) CopyStructure() *Link {
	// content links a value to its original id.
	// We only store links in content
	content := make(map[string]linkValue)
	// given a link (by value id), get the parent's link's id (to go up).
	// We use this information to build the link from the leafs to the root
	parents := make(map[string]string)
	// given a link (by value id), get all the links childs ids
	childs := make(map[string][]string)
	// leafs contain the leafs of the link, to build back from those leafs
	leafs := make(map[string]bool)

	// step one: find graph structure
	// make a fake value for l as a linkValue to ease iteration
	root := newLinkValue(l)
	queue := []linkValue{root}
	// and go for a BFS to go through the nodes
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]
		// we only put links, so it is safe to assume that current is indeed a link
		currentLink, _ := current.content.AsLink()
		if currentLink == nil {
			continue
		} else if _, found := content[current.uniqueId]; found {
			continue
		}

		// register the content
		content[current.uniqueId] = current

		var linksChilds []string
		// then, walk in the operands and build the structure
		isLeaf := true
		for _, operand := range currentLink.operands {
			switch operand.contentType() {
			case EntityTypeLink:
				isLeaf = false
				parents[operand.uniqueId] = current.uniqueId
				queue = append(queue, operand)
				linksChilds = append(linksChilds, operand.uniqueId)
			}
		}

		if isLeaf {
			leafs[current.uniqueId] = true
		} else {
			childs[current.uniqueId] = linksChilds
		}
	}

	// step two: reverse build (from the leafs to root)
	// newLinks contain, for each SOURCE id, the MAPPED value
	newLinks := make(map[string]linkValue)
	// currentElements to process
	var currentElements []linkValue
	// start with the leafs
	for leaf := range leafs {
		currentValue := content[leaf]
		currentElements = append(currentElements, currentValue)
	}

	// once done, we reach the root: only case no parent to process
	for len(currentElements) != 0 {
		// parents of the processed links
		var parentsOfProcessedLinks []string
		// for each node, locally map it
		for _, node := range currentElements {
			// build the local copy
			currentLink, _ := node.content.AsLink()
			newLink := currentLink.localCopy()
			// change the links childs to read mapped values
			for role, value := range currentLink.operands {
				// if value is a link, then replace it with the new value processed earlier
				if value.contentType() == EntityTypeLink {
					newLinkValue := newLinks[value.uniqueId]
					newLink.operands[role] = newLinkValue
				}
			}
			// link is updated, build it
			newLinkValue := newLinkValue(newLink)
			// and ensure that we register the new built link
			newLinks[node.uniqueId] = newLinkValue
			// add the parent as a link to potentially process.
			// If all of its childs are done, then this parent will pass to the next step
			if parentId, found := parents[node.uniqueId]; found {
				parentsOfProcessedLinks = append(parentsOfProcessedLinks, parentId)
			}
		}

		// then, reach one level up by processing links as soon as all its childs are done.
		// nextElements are the parents of the nodes we processed.
		var nextElements []linkValue
		for _, parent := range parentsOfProcessedLinks {
			allChildsProcessed := true
			for _, childToProcess := range childs[parent] {
				if _, found := newLinks[childToProcess]; !found {
					allChildsProcessed = false
					break
				}
			}

			if allChildsProcessed {
				nextElements = append(nextElements, content[parent])
			}
		}

		// next step: process parents once we mapped all of its childs
		currentElements = nextElements
	}

	// at this point, we made the root
	processedRoot := newLinks[root.uniqueId]
	rootAsLink, _ := processedRoot.content.AsLink()
	return rootAsLink
}

// Id returns the id for that link
func (l *Link) Id() string {
	return l.id
}

// Name returns the name of the link
func (l *Link) Name() string {
	return l.name
}

// Duration returns the link's active period
func (l *Link) Duration() structures.Period {
	return l.lifetime
}

// findAllMatchingCondition goes through the full link and find elements matching condition
func (l *Link) findAllMatchingCondition(acceptance func(ModelEntity) bool) []ModelEntity {
	matches := make([]ModelEntity, 0)
	linksAlreadyVisited := make(map[string]bool)

	elements := []*Link{l}
	for len(elements) != 0 {
		current := elements[0]
		elements = elements[1:]

		// STEP ONE: DEAL WITH THE WALKTHROUGH
		if current.GetType() == EntityTypeLink {
			link, _ := current.AsLink()
			if linksAlreadyVisited[link.id] {
				continue
			} else {
				linksAlreadyVisited[link.id] = true
			}

			for _, value := range link.operands {
				if value.contentType() == EntityTypeLink {
					childLlink, _ := value.content.AsLink()
					elements = append(elements, childLlink)
				}
			}
		}
		// END OF WALKTHROUGH

		// STEP TWO: TEST MATCH AND ADD IN MATCHES ACCORDINGLY
		if acceptance(current) {
			matches = append(matches, current)
		}

	}

	return matches
}

// Operands returns the operands of the link as a map of roles and linked entities
func (l *Link) Operands() map[string]ModelEntity {
	if l == nil {
		return nil
	}

	result := make(map[string]ModelEntity)
	for role, value := range l.operands {
		result[role] = value.content
	}

	return result
}

// AllObjectsOperands returns the objects appearing recursively in the link.
// It means that if l is a link of links of objects, descendants objects will appear.
// Each object appears once per id
func (l *Link) AllObjectsOperands() []Object {
	acceptValueAsObject := func(v ModelEntity) bool {
		matchingTypes := []EntityType{EntityTypeGroup, EntityTypeObject}
		return slices.Contains(matchingTypes, v.GetType())
	}

	var matches []Object
	values := l.findAllMatchingCondition(acceptValueAsObject)
	for _, value := range values {
		switch value.GetType() {
		case EntityTypeGroup:
			g, _ := value.AsGroup()
			matches = append(matches, g...)
		case EntityTypeObject:
			o, _ := value.AsObject()
			matches = append(matches, *o)
		}
	}

	return structures.SliceDeduplicateFunc(matches, func(a, b Object) bool { return a.Id == b.Id })
}

// LocalLinkValueMapper defines a mapping from a value to another.
// Accepted transformations are:
// IF value is anything but a link, THEN its image is also anything but a link
type LocalLinkValueMapper func(ModelEntity) (ModelEntity, bool, error)

// localLinkCaller calls a mapper but ensures invariants are respected
func localLinkValueCaller(value ModelEntity, mapper LocalLinkValueMapper) (ModelEntity, bool, error) {
	return mapper(value)
}

// Morphism maps a link to another, node per node
func (l *Link) Morphism(mapper LocalLinkValueMapper) (ModelEntity, error) {
	// mappedLinkValues contain all the mapped values, links and not links
	mappedLinkValues := make(map[string]linkValue)
	// WALKTHROUGH THE TREE TO CALCULATE MAPPING
	// for each element
	root := newLinkValue(l) // linkValue fake root matching l
	// try to map the root because if it is, we just return
	if rootMapping, toChangeRoot, errRoot := localLinkValueCaller(l, mapper); errRoot != nil {
		return nil, errRoot
	} else if toChangeRoot {
		return rootMapping, nil
	}

	// elements to run through
	elements := []linkValue{root}
	// for each element (bfs)
	for len(elements) != 0 {
		// pop next element
		current := elements[0]
		elements = elements[1:]
		// currentId of the current node
		currentId := current.uniqueId
		// map it and change if toChange (second result) is true
		if mappedEntity, toChange, errMapping := localLinkValueCaller(current.content, mapper); errMapping != nil {
			return nil, errMapping
		} else if toChange {
			// build a new value
			newLeaf := newLinkValue(mappedEntity)
			// and register the mapping
			mappedLinkValues[currentId] = newLeaf
		} else {
			// for a link, add all its childs as elements to process
			if current.contentType() == EntityTypeLink {
				if currentLink, err := current.content.AsLink(); err != nil {
					return nil, err
				} else {
					// keep going in the tree
					for _, child := range currentLink.operands {
						elements = append(elements, child)
					}

					// make a new link and register it.
					// We make a new link copy because result is a fully independant link (no common descendant)
					currentLinkCopy := currentLink.localCopy()
					// new id because it is not the same link as before (we clone or change)
					currentLinkCopy.id = uuid.NewString()
					newCopy := newLinkValue(currentLinkCopy)
					// register old id -> new clone
					mappedLinkValues[currentId] = newCopy
				}
			} else {
				// make a new value for that content.
				// It should not appear with that same id in two links.
				// So, newLinkValue will allocate a new id
				newCopy := newLinkValue(current.content)
				mappedLinkValues[currentId] = newCopy
			}
		}
	}

	// At this point, we have the values for each node and the tree structure.
	// What we do now is to link mapped links with each other.
	// Starting from the root, we link a mapped node to its mapped childs
	elements = append(elements, root)
	// we walk into the SOURCE link and map equivalent nodes.
	// for each SOURCE node:
	//    find its equivalent and its childs
	//    link equivalent childs to equivalent node
	for len(elements) != 0 {
		// current SOURCE node
		current := elements[0]
		elements = elements[1:]
		// by invariant:
		//  current is a link for sure
		//  its mapped value is a link too with same roles (child values may differ)
		currentId := current.uniqueId
		equivalent := mappedLinkValues[currentId]

		// so we will fill those two
		var currentLink, equivalentLink *Link
		// map current node and equivalent node to links
		if clink, errSource := current.content.AsLink(); errSource != nil {
			return nil, errSource
		} else if elink, errMapping := equivalent.content.AsLink(); errMapping != nil {
			return nil, errMapping
		} else {
			currentLink, equivalentLink = clink, elink
		}
		// OK, so currentLink, equivalentLink are set
		// And now, we map one by one the childs

		// for each role and value, go on and map childs
		for role, sourceChild := range currentLink.operands {
			// find the equivalent
			equivalentChild, found := mappedLinkValues[sourceChild.uniqueId]
			if !found {
				continue
			}

			// both are defined, mapping make sense, so update values
			equivalentLink.operands[role] = equivalentChild
			// Now what we want is to go through the relevant childs
			// test if childs are links
			isSourceChildLink := sourceChild.contentType() == EntityTypeLink
			isDestChildLink := equivalentChild.contentType() == EntityTypeLink
			// maintain the invariant
			if isSourceChildLink && isDestChildLink {
				elements = append(elements, sourceChild)
			}
		}
	}

	// and now we are done, so just return mapped root
	return mappedLinkValues[root.uniqueId].content, nil

}
