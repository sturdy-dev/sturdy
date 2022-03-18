package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Identifier_exact_match(t *testing.T) {
	identifier := &Identifier{
		Type:    Users,
		Pattern: "user-id-1",
	}

	identity := Identity{
		ID:   "user-id-1",
		Type: Users,
	}

	matches := identifier.Matches(identity)

	assert.True(t, matches)
}

func Test_Identifier_no_match(t *testing.T) {
	identifier := &Identifier{
		Type:    Users,
		Pattern: "user-id-1",
	}

	identity := Identity{
		ID:   "user-id",
		Type: Users,
	}

	matches := identifier.Matches(identity)

	assert.False(t, matches)
}

func Test_Identifier_part_match(t *testing.T) {
	identifier := &Identifier{
		Type:    Users,
		Pattern: "user-id-*",
	}

	identity := Identity{
		ID:   "user-id-2",
		Type: Users,
	}

	matches := identifier.Matches(identity)

	assert.True(t, matches)
}

func Test_Policy_single_match(t *testing.T) {
	p := Policy{
		Rules: []*Rule{
			{
				ID: "user-1 can write codebase-1",
				Principals: []*Identifier{
					{Type: Users, Pattern: "user-1"},
				},
				Action: ActionWrite,
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
	}

	isAllowed := p.Assert(
		Identity{Type: Users, ID: "user-1"},
		ActionWrite,
		Identity{Type: Codebases, ID: "codebase-1"},
	)

	assert.True(t, isAllowed)
}

func Test_Policy_no_match(t *testing.T) {
	p := Policy{
		Rules: []*Rule{
			{
				ID: "user-1 can read codebase-1",
				Principals: []*Identifier{
					{Type: Users, Pattern: "user-1"},
				},
				Action: Action("read"),
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
	}

	isAllowed := p.Assert(
		Identity{Type: Users, ID: "user-1"},
		ActionWrite,
		Identity{Type: Codebases, ID: "codebase-1"},
	)

	assert.False(t, isAllowed)
}

func Test_Policy_one_match_of_many(t *testing.T) {
	p := Policy{
		Rules: []*Rule{
			{
				ID: "user-1 can read codebase-1",
				Principals: []*Identifier{
					{Type: Users, Pattern: "user-1"},
				},
				Action: Action("read"),
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
			{
				ID: "any user can write codebase-1",
				Principals: []*Identifier{
					{Type: Users, Pattern: "*"},
				},
				Action: ActionWrite,
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
	}

	isAllowed := p.Assert(
		Identity{Type: Users, ID: "user-1"},
		ActionWrite,
		Identity{Type: Codebases, ID: "codebase-1"},
	)

	assert.True(t, isAllowed)
}

func Test_Policy_no_group_match(t *testing.T) {
	p := Policy{
		Rules: []*Rule{
			{
				ID:     "admins can write codebase-1",
				Action: ActionWrite,
				Principals: []*Identifier{
					{Type: Groups, Pattern: "admins"},
				},
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
		Groups: []*Group{
			{
				ID:      "admins",
				Members: []*Identifier{},
			},
		},
	}

	isAllowed := p.Assert(
		Identity{Type: Users, ID: "user-1"},
		ActionWrite,
		Identity{Type: Codebases, ID: "codebase-1"},
	)

	assert.False(t, isAllowed)
}

func Test_Policy_group_match(t *testing.T) {
	p := Policy{
		Rules: []*Rule{
			{
				ID:     "admins can write codebase-1",
				Action: ActionWrite,
				Principals: []*Identifier{
					{Type: Groups, Pattern: "admins"},
				},
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
		Groups: []*Group{
			{
				ID: "admins",
				Members: []*Identifier{
					{Type: Users, Pattern: "user-1"},
				},
			},
		},
	}

	isAllowed := p.Assert(
		Identity{Type: Users, ID: "user-1"},
		ActionWrite,
		Identity{Type: Codebases, ID: "codebase-1"},
	)

	assert.True(t, isAllowed)
}

func Test_Policy_Errors_all_good(t *testing.T) {
	actionWrite := ActionWrite
	p := Policy{
		Rules: []*Rule{
			{
				ID:     "admins can write codebase-1",
				Action: ActionWrite,
				Principals: []*Identifier{
					{Type: Groups, Pattern: "admins"},
				},
				Resources: []*Identifier{
					{Type: ACLs, Pattern: "*"},
				},
			},
		},
		Groups: []*Group{
			{
				ID: "admins",
				Members: []*Identifier{
					{Type: Users, Pattern: "user-1"},
				},
			},
		},
		Tests: []*Test{
			{
				ID:        "user-1 can write to acls",
				Principal: Identity{Type: Users, ID: "user-1"},
				Allow:     &actionWrite,
				Resource:  Identity{Type: ACLs, ID: aclID},
			},
		},
	}

	assert.Len(t, p.Errors(aclID), 0)
}

var (
	actionWrite            = ActionWrite
	aclID                  = "123"
	adminsCanWriteACLsRule = &Rule{
		ID:     "admins can write acls",
		Action: ActionWrite,
		Principals: []*Identifier{
			{Type: Groups, Pattern: "admins"},
		},
		Resources: []*Identifier{
			{Type: ACLs, Pattern: "*"},
		},
	}
	adminsGroup = &Group{
		ID: "admins",
		Members: []*Identifier{
			{Type: Users, Pattern: "user-1"},
		},
	}
	adminsCanWriteACLsTest = &Test{
		ID:        "admins can write acls",
		Allow:     &actionWrite,
		Principal: Identity{Type: Groups, ID: "admins"},
		Resource:  Identity{Type: ACLs, ID: aclID},
	}
)

func Test_Policy_Errors_subgroups_forbidden(t *testing.T) {
	p := Policy{
		Rules: []*Rule{adminsCanWriteACLsRule},
		Groups: []*Group{
			{ID: "randos"},
			{
				ID: "admins",
				Members: []*Identifier{
					{Type: Users, Pattern: "user-1"},
					{Type: Groups, Pattern: "randos"},
				},
			},
		},
		Tests: []*Test{adminsCanWriteACLsTest},
	}

	if errs := p.Errors(aclID); assert.Len(t, errs, 1) {
		assert.ErrorIs(t, errs["groups[\"admins\"]"], ErrSubgroupsForbidden)
	}
}

func Test_Policy_Errors_failing_allow_test(t *testing.T) {
	actionWrite := ActionWrite
	p := Policy{
		Rules: []*Rule{
			adminsCanWriteACLsRule,
			{
				ID:     "admins can write codebase-1",
				Action: ActionWrite,
				Principals: []*Identifier{
					{Type: Groups, Pattern: "admins"},
				},
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
		Groups: []*Group{
			adminsGroup,
		},
		Tests: []*Test{
			adminsCanWriteACLsTest,
			{
				ID:        "user-2 can write to codebase-1",
				Principal: Identity{Type: Users, ID: "user-2"},
				Allow:     &actionWrite,
				Resource:  Identity{Type: Codebases, ID: "codebase-1"},
			},
		},
	}

	if errs := p.Errors(aclID); assert.Len(t, errs, 1) {
		assert.ErrorIs(t, errs["tests[\"user-2 can write to codebase-1\"]"], ErrTestFails)
	}
}

func Test_Policy_Errors_failing_deny_test(t *testing.T) {
	actionWrite := ActionWrite
	p := Policy{
		Rules: []*Rule{
			adminsCanWriteACLsRule,
			{
				ID:     "admins can write codebase-1",
				Action: ActionWrite,
				Principals: []*Identifier{
					{Type: Groups, Pattern: "admins"},
				},
				Resources: []*Identifier{
					{Type: Codebases, Pattern: "codebase-1"},
				},
			},
		},
		Groups: []*Group{
			adminsGroup,
		},
		Tests: []*Test{
			adminsCanWriteACLsTest,
			{
				ID:        "user-1 can write to codebase-1",
				Principal: Identity{Type: Users, ID: "user-1"},
				Deny:      &actionWrite,
				Resource:  Identity{Type: Codebases, ID: "codebase-1"},
			},
		},
	}

	if errs := p.Errors(aclID); assert.Len(t, errs, 1) {
		assert.ErrorIs(t, errs["tests[\"user-1 can write to codebase-1\"]"], ErrTestFails)
	}
}

func Test_Policy_Errors_test_condition_missing(t *testing.T) {
	p := Policy{
		Rules:  []*Rule{adminsCanWriteACLsRule},
		Groups: []*Group{adminsGroup},
		Tests: []*Test{
			adminsCanWriteACLsTest,
			{
				ID:        "user-2 can write to codebase-1",
				Principal: Identity{Type: Users, ID: "user-2"},
				Resource:  Identity{Type: Codebases, ID: "codebase-1"},
			},
		},
	}

	if errs := p.Errors(aclID); assert.Len(t, errs, 1) {
		assert.ErrorIs(t, errs["tests[\"user-2 can write to codebase-1\"]"], ErrTestMustHaveCondition)
	}
}

func Test_Policy_Errors_acl_test_missing(t *testing.T) {
	p := Policy{}

	if errs := p.Errors(aclID); assert.Len(t, errs, 1) {
		assert.Equal(t, errs["tests"], ErrACLTestMissing(aclID))
	}
}

func Test_Policy_Errors_test_group_contains_invalid_resource(t *testing.T) {
	p := Policy{
		Rules: []*Rule{adminsCanWriteACLsRule},
		Groups: []*Group{adminsGroup, {
			ID:      "test",
			Members: []*Identifier{{Type: identityType("invalid")}},
		}},
		Tests: []*Test{adminsCanWriteACLsTest},
	}

	if errs := p.Errors(aclID); assert.Len(t, errs, 1) {
		assert.ErrorIs(t, errs["groups[\"test\"].members[\"invalid\"]"], ErrUnsupportedIdentityType)
	}
}
