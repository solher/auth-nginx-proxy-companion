package validators

import (
	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewResourcesValid)
}

type (
	ResourcesValidResourcesRepo interface {
		View(func(tx *bolt.Tx) error) error
	}

	ResourcesValid struct {
		r ResourcesValidResourcesRepo
	}
)

func NewResourcesValid(r ResourcesValidResourcesRepo) *ResourcesValid {
	return &ResourcesValid{r: r}
}

func (v *ResourcesValid) ValidateCreation(resource *models.Resource) error {
	c := make(chan error, 2)

	if resource.Name == nil || len(*resource.Name) == 0 {
		return errs.NewErrValidation("resource name cannot be blank")
	}

	if resource.Hostname == nil || len(*resource.Hostname) == 0 {
		return errs.NewErrValidation("resource hostname cannot be blank")
	}

	go func() {
		if err := v.ValidateHostnameUniqueness(resource); err != nil {
			c <- err
		}
		c <- nil
	}()

	go func() {
		if err := v.ValidateNameUniqueness(resource); err != nil {
			c <- err
		}
		c <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-c; err != nil {
			return err
		}
	}

	return nil
}

func (v *ResourcesValid) ValidateUpdate(resource *models.Resource) error {
	if resource.Name == nil || len(*resource.Name) == 0 {
		return errs.NewErrValidation("resource name cannot be blank")
	}

	if err := v.ValidateNameUniqueness(resource); err != nil {
		return err
	}

	return nil
}

func (v *ResourcesValid) ValidateHostnameUniqueness(resource *models.Resource) error {
	err := v.r.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket([]byte("resources")).Get([]byte(*resource.Hostname))

		if len(raw) != 0 {
			return errs.NewErrValidation("hostname must be unique")
		}

		return nil
	})

	return err
}

func (v *ResourcesValid) ValidateNameUniqueness(resource *models.Resource) error {
	if resource.Name == nil {
		return nil
	}

	err := v.r.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket([]byte("resources")).Get([]byte(*resource.Name))

		if len(raw) != 0 {
			return errs.NewErrValidation("name must be unique")
		}

		return nil
	})

	return err
}
