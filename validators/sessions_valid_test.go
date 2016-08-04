package validators

import (
	"testing"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionsValidSessionsRepo struct {
	err bool
}

func (r *sessionsValidSessionsRepo) View(t func(tx *bolt.Tx) error) error {
	if r.err {
		return errs.Internal.Database
	}

	return nil
}

// TestSessionsValidValidateCreation runs tests on the SessionsValid ValidateCreation method.
func TestSessionsValidValidateCreation(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	repo := &sessionsValidSessionsRepo{}
	valid := NewSessionsValid(repo)
	session := &models.Session{}

	// Validation error: nil policies
	err := valid.ValidateCreation(session)
	r.NotNil(err)

	session.Policies = []string{"1", "2"}
	repo.err = true

	// The repo returns a database error
	err = valid.ValidateCreation(session)
	r.NotNil(err)
	a.IsType(errs.Internal.Database, err)

	repo.err = false

	// Success
	err = valid.ValidateCreation(session)
	r.Nil(err)
}
