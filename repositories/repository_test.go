package repositories

import (
	"errors"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"
)

type database struct {
	err bool
}

func (d *database) Update(func(tx *bolt.Tx) error) error {
	if d.err {
		return errors.New("ERROR")
	}

	return nil
}

func (d *database) View(func(tx *bolt.Tx) error) error {
	if d.err {
		return errors.New("ERROR")
	}

	return nil
}

// TestRepository runs tests on the Repository.
func TestRepository(t *testing.T) {
	r := require.New(t)
	db := &database{}
	repo := NewRepository(db)

	err := repo.Update(func(tx *bolt.Tx) error { return nil })
	r.NoError(err)

	err = repo.View(func(tx *bolt.Tx) error { return nil })
	r.NoError(err)

	db.err = true

	err = repo.Update(func(tx *bolt.Tx) error { return nil })
	r.Error(err)

	err = repo.View(func(tx *bolt.Tx) error { return nil })
	r.Error(err)
}
