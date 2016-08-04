package interactors

import (
	"errors"
	"testing"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type policiesInterPoliciesRepo struct {
	err bool
}

func (r *policiesInterPoliciesRepo) Update(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

func (r *policiesInterPoliciesRepo) View(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

type policiesInterSessionsInter struct {
	err bool
}

func (r *policiesInterSessionsInter) DeleteCascade(policy *models.Policy) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

type policiesInterPoliciesValid struct {
	errValid bool
}

func (v *policiesInterPoliciesValid) ValidateDeletion(policy *models.Policy) error {
	if v.errValid {
		return errors.New("validation error")
	}

	return nil
}

// TestPoliciesInterFind runs tests on the PoliciesInter Find method.
func TestPoliciesInterFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesInterPoliciesRepo{}
	inter := NewPoliciesInter(repo, nil, nil)

	// Success
	result, err := inter.Find()
	r.NoError(err)
	a.Equal(0, len(result))

	repo.err = true

	// Database error
	result, err = inter.Find()
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestPoliciesInterFindByName runs tests on the PoliciesInter FindByName method.
func TestPoliciesInterFindByName(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesInterPoliciesRepo{}
	inter := NewPoliciesInter(repo, nil, nil)

	// Not found
	result, err := inter.FindByName("")
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.FindByName("")
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestPoliciesInterCreate runs tests on the PoliciesInter Create methods.
func TestPoliciesInterCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesInterPoliciesRepo{}
	inter := NewPoliciesInter(repo, nil, nil)

	// Success
	repo.err = false
	result, err := inter.Create(&models.Policy{})
	r.NoError(err)
	a.NotNil(result)

	// Nil error
	result, err = inter.Create(nil)
	r.Error(err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.Create(&models.Policy{})
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestPoliciesInterDeleteByName runs tests on the PoliciesInter DeleteByName method.
func TestPoliciesInterDeleteByName(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesInterPoliciesRepo{}
	sessionsInter := &policiesInterSessionsInter{}
	valid := &policiesInterPoliciesValid{}
	inter := NewPoliciesInter(repo, sessionsInter, valid)

	valid.errValid = true

	// Validation error
	result, err := inter.DeleteByName("")
	r.Error(err)
	a.Nil(result)

	valid.errValid = false

	// Not found
	result, err = inter.DeleteByName("")
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.DeleteByName("")
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)

	repo.err = false
	sessionsInter.err = true

	// Database error when cascade
	result, err = inter.DeleteByName("")
	r.Error(err)
	// a.IsType(errs.Internal.Database, err)
	a.IsType(errs.Internal.NotFound, err) // Can't mock BoltDB...
	a.Nil(result)
}

// TestPoliciesInterUpdateByName runs tests on the PoliciesInter UpdateByName method.
func TestPoliciesInterUpdateByName(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &policiesInterPoliciesRepo{}
	inter := NewPoliciesInter(repo, nil, nil)

	// Not found
	result, err := inter.UpdateByName("", &models.Policy{})
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	// Nil error
	result, err = inter.UpdateByName("", nil)
	r.Error(err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.UpdateByName("", &models.Policy{})
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}
