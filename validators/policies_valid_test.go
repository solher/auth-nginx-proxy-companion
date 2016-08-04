package validators

import (
	"testing"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type policiesValidPoliciesRepo struct {
	err bool
}

func (r *policiesValidPoliciesRepo) View(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

// TestPoliciesValidValidateCreation runs tests on the PoliciesValid ValidateCreation method.
func TestPoliciesValidValidateCreation(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesValidPoliciesRepo{}
	valid := NewPoliciesValid(repo)
	policy := &models.Policy{}

	// Validation error: nil name
	err := valid.ValidateCreation(policy)
	r.NotNil(err)

	policy.Name = utils.StrCpy("Foobar")

	// Validation error: nil permissions
	err = valid.ValidateCreation(policy)
	r.NotNil(err)

	// policy.Permissions = []models.Permission{{}}

	// // Validation error: permission with a nil/blank resource
	// err = valid.ValidateCreation(policy)
	// r.NotNil(err)

	policy.Permissions = []models.Permission{{Resource: utils.StrCpy("foo")}}

	// Validation error: resource not found
	err = valid.ValidateCreation(policy)
	r.NotNil(err)

	policy.Permissions = []models.Permission{{Resource: utils.StrCpy("*")}}

	// Validation passes: resource wildcard
	err = valid.ValidateCreation(policy)
	r.Nil(err)

	repo.err = true

	// The repo returns a database error
	err = valid.ValidateCreation(policy)
	r.NotNil(err)
	a.IsType(errs.Internal.Database, err)
}

// TestPoliciesValidValidateUpdate runs tests on the PoliciesValid ValidateUpdate method.
func TestPoliciesValidValidateUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesValidPoliciesRepo{}
	valid := NewPoliciesValid(repo)
	policy := &models.Policy{Name: utils.StrCpy("")}

	// Validation error: blank name
	err := valid.ValidateUpdate(policy)
	r.NotNil(err)

	policy.Name = utils.StrCpy("Foobar")
	// policy.Permissions = []models.Permission{{}}

	// // Validation error: permission with a nil/blank resource
	// err = valid.ValidateUpdate(policy)
	// r.NotNil(err)

	policy.Permissions = []models.Permission{{Resource: utils.StrCpy("foo")}}

	// Validation error: resource not found
	err = valid.ValidateUpdate(policy)
	r.NotNil(err)

	policy.Permissions = []models.Permission{{Resource: utils.StrCpy("*")}}

	// Validation passes: resource wildcard
	err = valid.ValidateUpdate(policy)
	r.Nil(err)

	repo.err = true

	// The repo returns a database error
	err = valid.ValidateUpdate(policy)
	r.NotNil(err)
	a.IsType(errs.Internal.Database, err)
}

// TestPoliciesValidValidateDeletion runs tests on the PoliciesValid ValidateDeletion method.
func TestPoliciesValidValidateDeletion(t *testing.T) {
	r := require.New(t)
	repo := &policiesValidPoliciesRepo{}
	valid := NewPoliciesValid(repo)
	policy := &models.Policy{Name: utils.StrCpy("guest")}

	// Validation error: name is guest
	err := valid.ValidateDeletion(policy)
	r.NotNil(err)

	policy.Name = utils.StrCpy("foo")

	// Success
	err = valid.ValidateDeletion(policy)
	r.Nil(err)
}
