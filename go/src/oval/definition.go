// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

func (od *GOvalDefinitions) getDefinition(s string) *GDefinition {
	for i := range od.Definitions.Definitions {
		if od.Definitions.Definitions[i].ID == s {
			return &od.Definitions.Definitions[i]
		}
	}

	return nil
}

func (od *GOvalDefinitions) getState(s string) interface{} {
	for _, x := range od.States.RPMInfoStates {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.States.TFC54States {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.States.DPKGInfoStates {
		if x.ID == s {
			return &x
		}
	}

	return nil
}

func (od *GOvalDefinitions) getObject(s string) genericObj {
	for i := range od.Objects.RPMInfoObjects {
		x := &od.Objects.RPMInfoObjects[i]
		if x.ID == s {
			return x
		}
	}
	for i := range od.Objects.DPKGInfoObjects {
		x := &od.Objects.DPKGInfoObjects[i]
		if x.ID == s {
			return x
		}
	}
	for i := range od.Objects.TFC54Objects {
		x := &od.Objects.TFC54Objects[i]
		if x.ID == s {
			return x
		}
	}

	return nil
}

func (od *GOvalDefinitions) getTest(s string) genericTest {
	for i := range od.Tests.RPMInfoTests {
		x := &od.Tests.RPMInfoTests[i]
		if x.ID == s {
			return x
		}
	}
	for i := range od.Tests.DPKGInfoTests {
		x := &od.Tests.DPKGInfoTests[i]
		if x.ID == s {
			return x
		}
	}
	for i := range od.Tests.TFC54Tests {
		x := &od.Tests.TFC54Tests[i]
		if x.ID == s {
			return x
		}
	}

	return nil
}

func (od *GDefinition) evaluate(ch chan GOvalResult, p *GOvalDefinitions) {
	var ret GOvalResult

	// We need a lock here as this definition could be selected for
	// evaluation by another definition as part of an extended
	// definition call.
	od.Lock()

	debugPrint("[evaluate] %v\n", od.ID)

	// Evaluate the root criteria item; this will likely result in
	// recursion through various subelements in the definition.
	od.status = od.Criteria.evaluate(p)
	ret.Status = od.status
	ret.Title = od.Metadata.Title
	ret.ID = od.ID

	// If the channel was nil we don't send the result back. This can
	// occur if the definition was called as the result of an
	// extend_definition rule in the OVAL definition being evaluated.
	if ch != nil {
		ch <- ret
	}

	od.Unlock()
}
