package acl

import (
	"fmt"
)

type Policy struct {
	Rules  []*Rule  `json:"rules,omitempty"`
	Groups []*Group `json:"groups,omitempty"`
	Tests  []*Test  `json:"tests,omitempty"`
}

// List return a list of _typ_ resources that _principal_ can _action_ on.
//
// For example,
//
//   files := List(Identity{Type: Users, ID: "user1"}, ActionWrite, Files)
//
// will return a list of file patterns the user1 can write to.
func (p Policy) List(principal Identity, action Action, typ identityType) []string {
	allowedPatterns := []string{}
	for _, rule := range p.Rules {
		if rule.Action != action {
			continue
		}

		principalFound := false
		for _, p := range resolveGroups(rule.Principals, p.Groups) {
			if p.Matches(principal) {
				principalFound = true
				break
			}
		}

		if !principalFound {
			continue
		}

		for _, resource := range resolveGroups(rule.Resources, p.Groups) {
			if resource.Type != typ {
				continue
			}
			allowedPatterns = append(allowedPatterns, resource.Pattern)
		}
	}
	return allowedPatterns
}

func (p Policy) Assert(principal Identity, action Action, resource Identity) bool {
	for _, acl := range p.Rules {
		if acl.Assert(principal, action, resource, p.Groups) {
			return true
		}
	}
	return false
}

var (
	ErrTestFails               = fmt.Errorf("test fails")
	ErrSubgroupsForbidden      = fmt.Errorf("groups can't have other groups as memebers")
	ErrUnsupportedIdentityType = fmt.Errorf("unsupported identity type")
	ErrTestMustHaveCondition   = fmt.Errorf("test must have either 'allow' or 'deny' condition")
	ErrUnsupportedActionType   = fmt.Errorf("unsupported action type")
	ErrACLTestMissing          = func(id string) error {
		return fmt.Errorf("at least one 'allow write' test must exist for 'acls::%s' resource", id)
	}
)

// Errors returns a non-empty list of errors if the policy is not valid.
func (p Policy) Errors(aclID string) map[string]error {
	errs := map[string]error{}

	var aclTest *Test
	for _, test := range p.Tests {
		// tests must be valid
		if test.Allow == nil && test.Deny == nil {
			errs[fmt.Sprintf("tests[\"%s\"]", test.ID)] = ErrTestMustHaveCondition
		}

		if test.Resource.Type == ACLs && test.Allow != nil && test.Resource.ID == aclID {
			aclTest = test
		}

		// tests must pass
		if test.Allow != nil && *test.Allow == ActionWrite {
			if !p.Assert(test.Principal, *test.Allow, test.Resource) {
				errs[fmt.Sprintf("tests[\"%s\"]", test.ID)] = ErrTestFails
			}
		}

		if test.Deny != nil {
			if p.Assert(test.Principal, *test.Deny, test.Resource) {
				errs[fmt.Sprintf("tests[\"%s\"]", test.ID)] = ErrTestFails
			}
		}
	}

	if aclTest == nil {
		errs["tests"] = ErrACLTestMissing(aclID)
	}

	// groups can't contain other groups
	for _, group := range p.Groups {
		for _, member := range group.Members {
			if member.Type == Groups {
				errs[fmt.Sprintf("groups[\"%s\"]", group.ID)] = ErrSubgroupsForbidden
			}

			if !member.Type.IsValid() {
				bytes, _ := member.MarshalJSON()
				errs[fmt.Sprintf("groups[\"%s\"].members[%s]", group.ID, string(bytes))] = ErrUnsupportedIdentityType
			}
		}
	}

	for _, rule := range p.Rules {
		for _, p := range rule.Principals {
			if !p.Type.IsValid() {
				bytes, _ := p.MarshalJSON()
				errs[fmt.Sprintf("rules[\"%s\"].principals[%s]", rule.ID, string(bytes))] = ErrUnsupportedIdentityType
			}
		}

		for _, r := range rule.Resources {
			if !r.Type.IsValid() {
				bytes, _ := r.MarshalJSON()
				errs[fmt.Sprintf("rules[\"%s\"].resources[%s]", rule.ID, string(bytes))] = ErrUnsupportedIdentityType
			}
		}

		if !rule.Action.IsValid() {
			errs[fmt.Sprintf("rules[\"%s\"].action", rule.ID)] = ErrUnsupportedActionType
		}
	}

	return errs
}

type Group struct {
	ID      string        `json:"id,omitempty"`
	Members []*Identifier `json:"members,omitempty"`
}

type Test struct {
	ID        string   `json:"id"`
	Principal Identity `json:"principal"`
	Allow     *Action  `json:"allow,omitempty"`
	Deny      *Action  `json:"deny,omitempty"`
	Resource  Identity `json:"resource"`
}

func resolveGroups(identifiers []*Identifier, groups []*Group) []*Identifier {
	resolved := make([]*Identifier, 0, len(identifiers))
	for _, i := range identifiers {
		resolved = append(resolved, i)

		if i.Type == Groups {
			for _, group := range groups {
				if i.Matches(Identity{ID: group.ID, Type: Groups}) {
					resolved = append(resolved, group.Members...)
				}
			}
		}
	}
	return resolved
}

type Rule struct {
	ID         string        `json:"id,omitempty"`
	Action     Action        `json:"action,omitempty"`
	Principals []*Identifier `json:"principals,omitempty"`
	Resources  []*Identifier `json:"resources,omitempty"`
}

func (a *Rule) Assert(principal Identity, action Action, resource Identity, groups []*Group) bool {
	if a.Action != action {
		return false
	}
	return a.assertPrincipal(principal, groups) && a.assertResource(resource, groups)
}

func (a *Rule) assertPrincipal(principal Identity, groups []*Group) bool {
	for _, p := range resolveGroups(a.Principals, groups) {
		if p.Matches(principal) {
			return true
		}
	}
	return false
}

func (a *Rule) assertResource(resource Identity, groups []*Group) bool {
	for _, r := range resolveGroups(a.Resources, groups) {
		if r.Matches(resource) {
			return true
		}
	}
	return false
}

type Action string

var supportedActions = map[Action]bool{ActionWrite: true}

func (a Action) IsValid() bool {
	return supportedActions[a]
}

const (
	ActionWrite Action = "write"
)
