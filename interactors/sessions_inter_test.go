package interactors

import (
	"testing"
	"time"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionsInterSessionsRepo struct {
	err bool
}

func (r *sessionsInterSessionsRepo) Update(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

func (r *sessionsInterSessionsRepo) View(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

// TestSessionsInterFind runs tests on the SessionsInter Find method.
func TestSessionsInterFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &sessionsInterSessionsRepo{}
	inter := NewSessionsInter(repo, nil)

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

// TestSessionsInterFindByToken runs tests on the SessionsInter FindByToken method.
func TestSessionsInterFindByToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &sessionsInterSessionsRepo{}
	inter := NewSessionsInter(repo, nil)

	// Not found
	result, err := inter.FindByToken("")
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.FindByToken("")
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestSessionsInterCreate runs tests on the SessionsInter Create methods.
func TestSessionsInterCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &sessionsInterSessionsRepo{}
	getter := utils.NewFakeModelsGetter()
	getter.SessionValidity = time.Hour
	getter.SessionTokenLength = 32
	inter := NewSessionsInter(repo, getter)

	// Success
	repo.err = false
	result, err := inter.Create(&models.Session{})
	r.NoError(err)
	a.NotNil(result)
	a.Len(*result.Token, 32)
	a.NotNil(result.ValidTo)

	// Nil error
	result, err = inter.Create(nil)
	r.Error(err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.Create(&models.Session{})
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestSessionsInterDeleteByToken runs tests on the SessionsInter DeleteByToken method.
func TestSessionsInterDeleteByToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &sessionsInterSessionsRepo{}
	inter := NewSessionsInter(repo, nil)

	// Not found
	result, err := inter.DeleteByToken("")
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.Nil(result)

	repo.err = true

	// Database error
	result, err = inter.DeleteByToken("")
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}

// TestSessionsInterDeleteByOwnerTokens runs tests on the SessionsInter DeleteByOwnerTokens method.
func TestSessionsInterDeleteByOwnerTokens(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &sessionsInterSessionsRepo{}
	inter := NewSessionsInter(repo, nil)

	// Do nothing
	result, err := inter.DeleteByOwnerTokens([]string{""})
	r.NoError(err)
	a.Len(result, 0)

	repo.err = true

	// Database error
	result, err = inter.DeleteByOwnerTokens([]string{""})
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.Nil(result)
}
