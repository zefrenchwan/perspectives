package models

import (
	"testing"

	"github.com/zefrenchwan/perspectives.git/models"
	"github.com/zefrenchwan/perspectives.git/structures"
)

func TestLinksCreation(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	cheese := models.NewTrait("cheese")
	mary := models.NewObject([]string{"Human"})

	if _, err := models.NewSimpleLink("likes", john, "cheese"); err == nil {
		t.Log("failed to detect wrong object operand")
		t.Fail()
	}

	if l, err := models.NewSimpleLink("likes", john, cheese); err != nil {
		t.Log("failed to use object as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	}

	if l, err := models.NewLink("likes", map[string]any{models.RoleSubject: []*models.Object{john, mary}, models.RoleObject: cheese}, structures.NewFullPeriod()); err != nil {
		t.Log("failed to use group as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	} else if k, err := models.NewSimpleLink("knows", mary, l); err != nil {
		t.Log("failed to use link as operand")
		t.Fail()
	} else if k.Name() != "knows" {
		t.Log("wrong composite name")
		t.Fail()
	}

	group, _ := models.NewObjectsGroup([]*models.Object{john, mary})
	if l, err := models.NewSimpleLink("likes", group, cheese); err != nil {
		t.Log("failed to use group as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	} else if k, err := models.NewSimpleLink("knows", mary, l); err != nil {
		t.Log("failed to use link as operand")
		t.Fail()
	} else if k.Name() != "knows" {
		t.Log("wrong composite name")
		t.Fail()
	}

	group, _ = models.NewGroupOfObjects(john, mary)
	if l, err := models.NewSimpleLink("likes", group, cheese); err != nil {
		t.Log("failed to use group as operand")
		t.Log(err)
		t.Fail()
	} else if l.Name() != "likes" {
		t.Log("wrong name")
		t.Fail()
	} else if k, err := models.NewSimpleLink("knows", mary, l); err != nil {
		t.Log("failed to use link as operand")
		t.Fail()
	} else if k.Name() != "knows" {
		t.Log("wrong composite name")
		t.Fail()
	}
}

func TestLinksSame(t *testing.T) {
	amanda := models.NewObject([]string{"Human"})
	pizza := models.NewObject([]string{"Human"})
	likes, _ := models.NewSimpleLink("likes", amanda, pizza)
	espresso := models.NewObject([]string{"coffee"})
	loves, _ := models.NewSimpleLink("likes", amanda, espresso)
	knows, _ := models.NewSimpleLink("knows", amanda, loves)

	// test same structure, accept all comparison
	if !likes.SameFunc(loves, func(me1, me2 models.ModelEntity) bool { return true }) {
		t.Log("same structures should match")
		t.Fail()
	}

	// sufficient in here, risky elsewhere
	idBasedEntityComparison := func(me1, me2 models.ModelEntity) bool {
		if me1.GetType() != me2.GetType() {
			return false
		} else if me1.GetType() == models.EntityTypeObject {
			o1, _ := me1.AsObject()
			o2, _ := me2.AsObject()
			return o1.Equals(o2)
		}

		return true
	}

	// test same structure, better comparison
	if likes.SameFunc(loves, idBasedEntityComparison) {
		t.Log("same structures should match")
		t.Fail()
	}

	// test different structure, no matter the local comparator
	if likes.SameFunc(knows, func(me1, me2 models.ModelEntity) bool { return true }) {
		t.Log("different structures should never match")
		t.Fail()
	}

}

func TestCloneSimpleLink(t *testing.T) {
	sonia := models.NewObject([]string{"Human"})
	jack := models.NewObject([]string{"Human"})
	married, _ := models.NewSimpleLink("married", sonia, jack)

	// test the run ended
	clone := married.CopyStructure()
	if clone == nil {
		t.Fail()
	}

	if clone.Name() != married.Name() {
		t.Log("bad name for link")
		t.Fail()
	} else if clone.Id() != married.Id() {
		t.Log("copy should keep id")
		t.Fail()
	} else if !clone.Duration().Equals(married.Duration()) {
		t.Log("copy should keep duration")
		t.Fail()
	}

	operands := clone.Operands()
	if len(operands) != 2 {
		t.Log("missing operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != sonia.Id {
		t.Log("subject not mapped")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("object not found")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != jack.Id {
		t.Log("wrong object")
		t.Fail()
	}
}

func TestCloneLongLink(t *testing.T) {
	sonia := models.NewObject([]string{"Human"})
	jack := models.NewObject([]string{"Human"})
	married, _ := models.NewSimpleLink("married", sonia, jack)
	marcel := models.NewObject([]string{"Human"})
	ignores, _ := models.NewSimpleLink("ignores", marcel, married)
	knows, _ := models.NewSimpleLink("knows", jack, ignores)
	// Jack knows that Marcel ignores that Sonia married Jack

	clone := knows.CopyStructure()
	if clone == nil {
		t.Log("clone failed")
		t.Fail()
	} else if clone.Name() != knows.Name() {
		t.Log("knows name failed")
		t.Fail()
	} else if clone.Id() != knows.Id() {
		t.Log("knows id failed")
		t.Fail()
	}

	knowsOperands := clone.Operands()
	if len(knowsOperands) != 2 {
		t.Log("root operands failed")
	} else if s, found := knowsOperands[models.RoleSubject]; !found {
		t.Log("root subject failed")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != jack.Id {
		t.Log("root subject id failed")
		t.Fail()
	} else if o, found := knowsOperands[models.RoleObject]; !found {
		t.Log("failed to find object in root")
		t.Fail()
	} else if o.GetType() != models.EntityTypeLink {
		t.Log("object should be a link")
		t.Fail()
	}

	ignoresLink, _ := knowsOperands[models.RoleObject].AsLink()
	if ignoresLink.Id() != ignores.Id() {
		t.Log("failed first level id")
		t.Fail()
	} else if ignoresLink.Name() != ignores.Name() {
		t.Log("failed first level name")
		t.Fail()
	} else if operands := ignoresLink.Operands(); len(operands) != 2 {
		t.Log("wrong first level operand")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject at first level")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != marcel.Id {
		t.Log("wrong subject at first level")
		t.Fail()
	}

	marriedLink, _ := ignoresLink.Operands()[models.RoleObject].AsLink()
	if marriedLink.Id() != married.Id() {
		t.Log("failed second level id")
		t.Fail()
	} else if marriedLink.Name() != married.Name() {
		t.Log("failed second level name")
		t.Fail()
	} else if operands := marriedLink.Operands(); len(operands) != 2 {
		t.Log("wrong second level operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject at level 2")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != sonia.Id {
		t.Log("wrong subject at level 2")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("missing object at level 2")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != jack.Id {
		t.Log("wrong object at level 2")
		t.Fail()
	}

}

func TestMappingNoChange(t *testing.T) {
	william := models.NewObject([]string{"Human"})
	pizza := models.NewObject([]string{"Food"})
	eats, errEeats := models.NewSimpleLink("eats", william, pizza)
	if errEeats != nil {
		t.Log(errEeats)
		t.Fail()
	}

	if same, err := eats.Morphism(func(me models.ModelEntity) (models.ModelEntity, bool, error) { return nil, false, nil }); err != nil {
		t.Log("failed to map")
		t.Fail()
	} else if l, err := same.AsLink(); err != nil {
		t.Log("failed to create a link")
		t.Fail()
	} else if l.Id() == eats.Id() {
		t.Log("failed to change id")
		t.Fail()
	} else if l.Name() != eats.Name() {
		t.Log("failed to map name")
		t.Fail()
	} else if !l.Duration().Equals(eats.Duration()) {
		t.Log("failed to map duration")
		t.Fail()
	} else if ops := l.Operands(); len(ops) != 2 {
		t.Log("failed to copy operands")
		t.Fail()
	} else if s, found := ops[models.RoleSubject]; !found {
		t.Log("failed to find subject")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != william.Id {
		t.Log("failed to copy subject")
		t.Fail()
	} else if o, found := ops[models.RoleSubject]; !found {
		t.Log("failed to find object")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != pizza.Id {
		t.Log("failed to copy object")
	}
}

func TestMappingRoot(t *testing.T) {
	jenna := models.NewObject([]string{"Human"})
	lorie := models.NewObject([]string{"Human"})
	friends, errSource := models.NewSimpleLink("friends", jenna, lorie)
	if errSource != nil {
		t.Log(errSource)
		t.Fail()
	}

	mappping := func(m models.ModelEntity) (models.ModelEntity, bool, error) {
		if m.GetType() == models.EntityTypeLink {
			if result, err := models.NewSimpleLink("loves", jenna, lorie); err != nil {
				return nil, false, err
			} else {
				return result, true, nil
			}
		} else {
			return nil, false, nil
		}
	}

	if result, err := friends.Morphism(mappping); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.GetType() != models.EntityTypeLink {
		t.Log("wrong mapping")
		t.Fail()
	} else if link, err := result.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else if link.Name() != "loves" {
		t.Log("bad mapping")
		t.Fail()
	} else if operands := link.Operands(); len(operands) != 2 {
		t.Log("missing operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log("wrong subject type")
		t.Fail()
	} else if subject.Id != jenna.Id {
		t.Log("wrong subject")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("missing object")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log("wrong object type")
		t.Fail()
	} else if object.Id != lorie.Id {
		t.Log("wrong object")
		t.Fail()
	}
}

func TestMappingLeaf(t *testing.T) {
	jenna := models.NewObject([]string{"Human"})
	lorie := models.NewObject([]string{"Human"})
	marie := models.NewObject([]string{"Human"})
	friends, errSource := models.NewSimpleLink("friends", jenna, lorie)
	if errSource != nil {
		t.Log(errSource)
		t.Fail()
	}

	mappping := func(m models.ModelEntity) (models.ModelEntity, bool, error) {
		if m.GetType() == models.EntityTypeObject {
			if o, err := m.AsObject(); err != nil {
				return nil, false, err
			} else if o.Id == lorie.Id {
				return marie, true, nil
			}
		}

		return nil, false, nil
	}

	if result, err := friends.Morphism(mappping); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.GetType() != models.EntityTypeLink {
		t.Log("wrong mapping")
		t.Fail()
	} else if link, err := result.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else if link.Id() == friends.Id() {
		t.Log("had to change ids")
		t.Fail()
	} else if link.Name() != "friends" {
		t.Log("bad mapping")
		t.Fail()
	} else if operands := link.Operands(); len(operands) != 2 {
		t.Log("missing operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log("wrong subject type")
		t.Fail()
	} else if subject.Id != jenna.Id {
		t.Log("wrong subject")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("missing object")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log("wrong object type")
		t.Fail()
	} else if object.Id == lorie.Id {
		t.Log("no object change")
		t.Fail()
	} else if object.Id != marie.Id {
		t.Log(object)
		t.Log("wrong object")
		t.Fail()
	}
}

func TestMappingLongLink(t *testing.T) {
	jenna := models.NewObject([]string{"Human"})
	lorie := models.NewObject([]string{"Human"})
	marie := models.NewObject([]string{"Human"})
	friends, errSource := models.NewSimpleLink("friends", jenna, lorie)
	if errSource != nil {
		t.Log(errSource)
		t.Fail()
	}

	knows, errLong := models.NewSimpleLink("knows", marie, friends)
	if errLong != nil {
		t.Log(errLong)
		t.Fail()
	}

	mappping := func(m models.ModelEntity) (models.ModelEntity, bool, error) {
		if m.GetType() == models.EntityTypeObject {
			if o, err := m.AsObject(); err != nil {
				return nil, false, err
			} else if o.Id == lorie.Id {
				return marie, true, nil
			}
		}

		return nil, false, nil
	}

	if result, err := knows.Morphism(mappping); err != nil {
		t.Log(err)
		t.Fail()
	} else if result.GetType() != models.EntityTypeLink {
		t.Log("wrong mapping")
		t.Fail()
	} else if root, err := result.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else if root.Id() == knows.Id() {
		t.Log("failed to change root link")
		t.Fail()
	} else if rootOps := root.Operands(); len(rootOps) != 2 {
		t.Log("wrong root operands")
		t.Fail()
	} else if sroot, err := rootOps[models.RoleSubject].AsObject(); err != nil {
		t.Log("expected object as root subject")
		t.Fail()
	} else if sroot.Id != marie.Id {
		t.Log("failed to map subject for root")
		t.Fail()
	} else if link, err := rootOps[models.RoleObject].AsLink(); err != nil {
		t.Log("faield to map object for root")
		t.Fail()
	} else if link.Id() == friends.Id() {
		t.Log("had to change ids")
		t.Fail()
	} else if link.Name() != "friends" {
		t.Log("bad mapping")
		t.Fail()
	} else if operands := link.Operands(); len(operands) != 2 {
		t.Log("missing operands")
		t.Fail()
	} else if s, found := operands[models.RoleSubject]; !found {
		t.Log("missing subject")
		t.Fail()
	} else if subject, err := s.AsObject(); err != nil {
		t.Log("wrong subject type")
		t.Fail()
	} else if subject.Id != jenna.Id {
		t.Log("wrong subject")
		t.Fail()
	} else if o, found := operands[models.RoleObject]; !found {
		t.Log("missing object")
		t.Fail()
	} else if object, err := o.AsObject(); err != nil {
		t.Log("wrong object type")
		t.Fail()
	} else if object.Id == lorie.Id {
		t.Log("no object change")
		t.Fail()
	} else if object.Id != marie.Id {
		t.Log(object)
		t.Log("wrong object")
		t.Fail()
	}
}

func TestMappingToVariables(t *testing.T) {
	william := models.NewObject([]string{"Human"})
	william.SetValue("first name", "William")
	x := models.NewVariableForObject("x", []string{"Human"})
	// for all x, x and x share the same identity (basic rule for test purpose)
	basicIdentity, errSource := models.NewSimpleLink("identity", x, x)
	if errSource != nil {
		t.Log(errSource)
		t.Fail()
	}

	mapping := func(m models.ModelEntity) (models.ModelEntity, bool, error) {
		if m.GetType() == models.EntityTypeVariable {
			if variable, err := m.AsVariable(); err != nil {
				return nil, false, err
			} else if variable.Name() != "x" {
				return nil, false, nil
			} else {
				result, errMap := variable.MapAs(william)
				if errMap != nil {
					return nil, false, errMap
				} else {
					return result, true, nil
				}
			}
		}

		return nil, false, nil
	}

	instantiation, errInstantiation := basicIdentity.Morphism(mapping)
	if errInstantiation != nil {
		t.Log(errInstantiation)
		t.Fail()
	} else if instantiation.GetType() != models.EntityTypeLink {
		t.Log("bad mapping")
		t.Fail()
	} else if link, err := instantiation.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else if link.Name() != "identity" {
		t.Log("wrong link copy")
		t.Fail()
	} else if ops := link.Operands(); len(ops) != 2 {
		t.Log("bad operands")
		t.Fail()
	} else if s, err := ops[models.RoleSubject].AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if s.Id != william.Id {
		t.Log("bad mapping to subject")
		t.Fail()
	} else if o, err := ops[models.RoleObject].AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if o.Id != william.Id {
		t.Log("bad mapping to object")
		t.Fail()
	}
}

func TestMappingLinkToValue(t *testing.T) {
	pizza := models.NewObject([]string{"food"})
	tiramisu := models.NewObject([]string{"Food"})
	middle, _ := models.NewSimpleLink("is before", pizza, tiramisu)
	burrata := models.NewObject([]string{"Food"})
	starter, _ := models.NewSimpleLink("starter", burrata, middle)
	coffee := models.NewObject([]string{"drink"})

	// replace middle with coffee
	result, errResult := starter.Morphism(func(me models.ModelEntity) (models.ModelEntity, bool, error) {
		if me.GetType() == models.EntityTypeLink {
			link, _ := me.AsLink()
			if link.Id() == middle.Id() {
				return coffee, true, nil
			}
		}

		return nil, false, nil
	})

	if errResult != nil {
		t.Log(errResult)
		t.Fail()
	} else if result.GetType() != models.EntityTypeLink {
		t.Fail()
	} else if link, err := result.AsLink(); err != nil {
		t.Log(err)
		t.Fail()
	} else if link.Name() != starter.Name() {
		t.Log("bad name")
		t.Fail()
	} else if ops := link.Operands(); len(ops) != 2 {
		t.Log("wrong operands")
		t.Fail()
	} else if subject, err := ops[models.RoleSubject].AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if subject.Id != burrata.Id {
		t.Log("bad subject")
		t.Fail()
	} else if object, err := ops[models.RoleObject].AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != coffee.Id {
		t.Log(err)
		t.Fail()
	}
}

func TestMappingValueToLink(t *testing.T) {
	gustave := models.NewObject([]string{"Person"})
	lisa := models.NewObject([]string{"Person"})
	loves, _ := models.NewSimpleLink("loves", lisa, gustave)
	paula := models.NewObject([]string{"Person"})
	variable := models.NewVariableForLink("x")
	knows, _ := models.NewSimpleLink("knows", paula, variable)

	// replace variable with loves.
	// Paula Knows X => Paula Knows Lisa loves Gustave
	replace, errReplace := knows.Morphism(func(me models.ModelEntity) (models.ModelEntity, bool, error) {
		if me.GetType() == models.EntityTypeVariable {
			return loves, true, nil
		}

		return nil, false, nil
	})

	if errReplace != nil {
		t.Log(errReplace)
		t.Fail()
	} else if replace.GetType() != models.EntityTypeLink {
		t.Fail()
	} else if root, errRoot := replace.AsLink(); errRoot != nil {
		t.Log(errRoot)
		t.Fail()
	} else if root.Name() != "knows" {
		t.Log("bad root")
		t.Fail()
	} else if rootOps := root.Operands(); len(rootOps) != 2 {
		t.Fail()
	} else if rootSubject, err := rootOps[models.RoleSubject].AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if rootSubject.Id != paula.Id {
		t.Fail()
	} else if rootObject, err := rootOps[models.RoleObject].AsLink(); err != nil {
		t.Log("bad link")
		t.Fail()
	} else if rootObject.Name() != loves.Name() {
		t.Log("bad mapped link")
		t.Fail()
	} else if ops := rootObject.Operands(); len(ops) != 2 {
		t.Fail()
	} else if subject, err := ops[models.RoleSubject].AsObject(); err != nil {
		t.Fail()
	} else if subject.Id != lisa.Id {
		t.Log("not full link")
		t.Fail()
	} else if object, err := ops[models.RoleObject].AsObject(); err != nil {
		t.Log(err)
		t.Fail()
	} else if object.Id != gustave.Id {
		t.Log("bad object")
		t.Fail()
	}
}

func TestVariableMatchingGlobalVariable(t *testing.T) {
	x := models.NewVariableForLink("x")
	y := models.NewVariableForTrait("y")
	dog := models.NewObject([]string{"dog"})
	cat := models.NewObject([]string{"cat"})
	not, _ := models.NewSimpleLink("is not", dog, cat)

	if mapping, accept := not.IsSpecializationOf(x); !accept {
		t.Log("x accepts a link, not is a link so it should accept")
		t.Fail()
	} else if !models.SameModelEntity(mapping[x.Name()], not) {
		t.Log("wrong mapping for x")
		t.Fail()
	} else if _, accept := not.IsSpecializationOf(y); accept {
		t.Fail()
	}
}

func TestVariableMatchingLeaf(t *testing.T) {
	pizza := models.NewObject([]string{"food"})
	italian := models.NewTrait("italian")
	x := models.NewVariableForObject("x", []string{"food"})
	y := models.NewVariableForTrait("y")
	genericRelation, _ := models.NewSimpleLink("is", x, y)
	baseRelation, _ := models.NewSimpleLink("is", pizza, italian)

	if mapping, accept := baseRelation.IsSpecializationOf(genericRelation); !accept {
		t.Log("generic relation is more general than base relation")
		t.Fail()
	} else if len(mapping) != 2 {
		t.Log("bad variables in mapping")
		t.Fail()
	} else if value := mapping[x.Name()]; !models.SameModelEntity(value, pizza) {
		t.Log("bad allocation for x")
		t.Fail()
	} else if value := mapping[y.Name()]; !models.SameModelEntity(value, &italian) {
		t.Log("bad allocation for y")
		t.Fail()
	}
}

func TestVariableNOTMatchingLeaf(t *testing.T) {
	pizza := models.NewObject([]string{"food"})
	italian := models.NewTrait("italian")
	x := models.NewVariableForObject("x", []string{"food"})
	y := models.NewVariableForTrait("y")
	// cannot match because names won't
	genericRelation, _ := models.NewSimpleLink("is but different than the other", x, y)
	baseRelation, _ := models.NewSimpleLink("is", pizza, italian)

	if _, accept := baseRelation.IsSpecializationOf(genericRelation); accept {
		t.Log("different names, should refuse")
		t.Fail()
	}

	// cannot match because double allocation for x
	genericRelation, _ = models.NewSimpleLink("is", x, x)
	baseRelation, _ = models.NewSimpleLink("is", pizza, italian)

	if _, accept := baseRelation.IsSpecializationOf(genericRelation); accept {
		t.Log("x allocated twice")
		t.Fail()
	}

	// cannot match because variables mismatch
	genericRelation, _ = models.NewSimpleLink("is", y, x)
	baseRelation, _ = models.NewSimpleLink("is", pizza, italian)

	if _, accept := baseRelation.IsSpecializationOf(genericRelation); accept {
		t.Log("bad recognition for types")
		t.Fail()
	}
}

func TestVariableMatchingLinkStructure(t *testing.T) {
	jane := models.NewObject([]string{"Human"})
	bobby := models.NewObject([]string{"Human"})
	martha := models.NewObject([]string{"Human"})
	friends, _ := models.NewSimpleLink("is friend with", jane, bobby)
	knows, _ := models.NewSimpleLink("knows", martha, friends)
	says, _ := models.NewSimpleLink("says", jane, knows)
	// Jane says that Martha knows that Jane is friend with Bobby
	// So, it should match
	// X says that Y knows that Z is friend with Bobby
	x := models.NewVariableForObject("x", []string{"Human"})
	y := models.NewVariableForObject("y", []string{"Human"})
	z := models.NewVariableForObject("z", []string{"Human"})
	genericFriends, _ := models.NewSimpleLink(friends.Name(), x, bobby)
	genericKnows, _ := models.NewSimpleLink(knows.Name(), y, genericFriends)
	genericSays, _ := models.NewSimpleLink(says.Name(), z, genericKnows)

	// test mappings one by one
	if mapping, accept := friends.IsSpecializationOf(genericFriends); !accept {
		t.Log("bad match for links with depth = 1")
		t.Fail()
	} else if len(mapping) != 1 {
		t.Log("bad variable mapping for links with depth = 1")
		t.Fail()
	} else if !models.SameModelEntity(mapping[x.Name()], jane) {
		t.Log("bad mapping for x")
		t.Fail()
	}

	if mapping, accept := knows.IsSpecializationOf(genericKnows); !accept {
		t.Log("bad match for links with depth = 2")
		t.Fail()
	} else if len(mapping) != 2 {
		t.Log("bad variable mapping for links with depth = 2")
		t.Fail()
	} else if !models.SameModelEntity(mapping[x.Name()], jane) {
		t.Log("bad mapping for x")
		t.Fail()
	} else if !models.SameModelEntity(mapping[y.Name()], martha) {
		t.Log("bad mapping for y")
		t.Fail()
	}

	if mapping, accept := says.IsSpecializationOf(genericSays); !accept {
		t.Log("bad match for links with depth = 3")
		t.Fail()
	} else if len(mapping) != 3 {
		t.Log("bad variable mapping for links with depth = 3")
		t.Fail()
	} else if !models.SameModelEntity(mapping[x.Name()], jane) {
		t.Log("bad mapping for x")
		t.Fail()
	} else if !models.SameModelEntity(mapping[y.Name()], martha) {
		t.Log("bad mapping for y")
		t.Fail()
	} else if !models.SameModelEntity(mapping[z.Name()], jane) {
		t.Log("bad mapping for z")
		t.Fail()
	}
}

func TestVariableMatchingLinks(t *testing.T) {
	john := models.NewObject([]string{"Human"})
	lisa := models.NewObject([]string{"Human"})
	lara := models.NewObject([]string{"Human"})
	hates, _ := models.NewSimpleLink("hates", lisa, lara)
	knows, _ := models.NewSimpleLink("knows", john, hates)
	// John knows that Lisa hates Lara
	x := models.NewVariableForLink("x")
	genericKnows, _ := models.NewSimpleLink("knows", john, x)

	if mapping, accept := knows.IsSpecializationOf(genericKnows); !accept {
		t.Log("failed to match")
		t.Fail()
	} else if !models.SameModelEntity(mapping[x.Name()], hates) {
		t.Log("failed to map x")
		t.Fail()
	}

}
