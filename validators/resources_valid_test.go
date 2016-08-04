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

type resourcesValidResourcesRepo struct {
	err bool
}

func (r *resourcesValidResourcesRepo) View(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

// TestResourcesValidValidateCreation runs tests on the ResourcesValid ValidateCreation method.
func TestResourcesValidValidateCreation(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &resourcesValidResourcesRepo{}
	valid := NewResourcesValid(repo)
	resource := &models.Resource{Name: utils.StrCpy("Foobar")}

	// Validation error: nil name
	err := valid.ValidateCreation(resource)
	r.NotNil(err)

	resource.Name = utils.StrCpy("Foobar")

	// Validation error: nil hostname
	err = valid.ValidateCreation(resource)
	r.NotNil(err)

	resource.Hostname = utils.StrCpy("foo.bar.com")
	repo.err = true

	// The repo returns a database error
	err = valid.ValidateCreation(resource)
	r.NotNil(err)
	a.IsType(errs.Internal.Database, err)

	repo.err = false

	// Success
	err = valid.ValidateCreation(resource)
	r.Nil(err)
}

// TestResourcesValidValidateUpdate runs tests on the ResourcesValid ValidateUpdate method.
func TestResourcesValidValidateUpdate(t *testing.T) {
	r := require.New(t)
	repo := &resourcesValidResourcesRepo{}
	valid := NewResourcesValid(repo)
	resource := &models.Resource{}

	// Validation error: blank name
	err := valid.ValidateUpdate(resource)
	r.NotNil(err)

	resource.Name = utils.StrCpy("Foobar")

	// Success
	err = valid.ValidateUpdate(resource)
	r.Nil(err)
}
