package interactors

import (
	"testing"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type resourcesInterResourcesRepo struct {
	err bool
}

func (r *resourcesInterResourcesRepo) Update(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

func (r *resourcesInterResourcesRepo) View(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

type resourcesInterPoliciesInter struct {
	err bool
}

func (r *resourcesInterPoliciesInter) DeleteCascade(resource *models.Resource) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

// TestResourcesInterFind runs tests on the ResourcesInter Find method.
func TestResourcesInterFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &resourcesInterResourcesRepo{}
	inter := NewResourcesInter(repo, nil)

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

// TestResourcesInterFindByHostname runs tests on the ResourcesInter FindByHostname method.
func TestResourcesInterFindByHostname(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &resourcesInterResourcesRepo{}
	inter := NewResourcesInter(repo, nil)

	// Not found
	result, err := inter.FindByHostname("")
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.FindByHostname("")
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestResourcesInterCreate runs tests on the ResourcesInter Create methods.
func TestResourcesInterCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &resourcesInterResourcesRepo{}
	inter := NewResourcesInter(repo, nil)

	// Success
	result, err := inter.Create(&models.Resource{})
	r.NoError(err)
	a.NotNil(result)

	// Nil error
	result, err = inter.Create(nil)
	r.Error(err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.Create(&models.Resource{})
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestResourcesInterDeleteByHostname runs tests on the ResourcesInter DeleteByHostname method.
func TestResourcesInterDeleteByHostname(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &resourcesInterResourcesRepo{}
	policiesInter := &resourcesInterPoliciesInter{}
	inter := NewResourcesInter(repo, policiesInter)

	// Not found
	result, err := inter.DeleteByHostname("")
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.DeleteByHostname("")
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)

	repo.err = false
	policiesInter.err = true

	// Database error when cascade
	result, err = inter.DeleteByHostname("")
	r.Error(err)
	// a.IsType(errs.Internal.Database, err)
	a.IsType(errs.Internal.NotFound, err) // Can't mock BoltDB...
	a.Nil(result)
}

// TestResourcesInterUpdateByHostname runs tests on the ResourcesInter UpdateByHostname method.
func TestResourcesInterUpdateByHostname(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &resourcesInterResourcesRepo{}
	inter := NewResourcesInter(repo, nil)

	// Not found
	result, err := inter.UpdateByHostname("", &models.Resource{})
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	// Nil error
	result, err = inter.UpdateByHostname("", nil)
	r.Error(err)
	a.Nil(result)

	// Database error
	repo.err = true
	result, err = inter.UpdateByHostname("", &models.Resource{})
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}
