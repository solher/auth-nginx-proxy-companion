package repositories

import (
	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewRepository)
}

type (
	DatabaseRunner interface {
		Update(func(tx *bolt.Tx) error) error
		View(func(tx *bolt.Tx) error) error
	}

	Repository struct {
		db DatabaseRunner
	}
)

func NewRepository(db DatabaseRunner) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Update(t func(tx *bolt.Tx) error) error {
	if err := r.db.Update(t); err != nil {
		return errs.Internal.Database
	}

	return nil
}

func (r *Repository) View(t func(tx *bolt.Tx) error) error {
	if err := r.db.View(t); err != nil {
		return errs.Internal.Database
	}

	return nil
}
