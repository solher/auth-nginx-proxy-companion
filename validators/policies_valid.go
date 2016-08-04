package validators

import (
	"encoding/json"
	"fmt"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewPoliciesValid)
}

type (
	PoliciesValidPoliciesRepo interface {
		View(func(tx *bolt.Tx) error) error
	}

	PoliciesValid struct {
		r PoliciesValidPoliciesRepo
	}
)

func NewPoliciesValid(r PoliciesValidPoliciesRepo) *PoliciesValid {
	return &PoliciesValid{r: r}
}

func (v *PoliciesValid) ValidateCreation(policy *models.Policy) error {
	c := make(chan error, 2)

	if policy.Name == nil || len(*policy.Name) == 0 {
		return errs.NewErrValidation("policy name cannot be blank")
	}

	if policy.Permissions == nil {
		return errs.NewErrValidation("policy permissions cannot be blank")
	}

	go func() {
		if err := v.ValidateResourcesExistence(policy); err != nil {
			c <- err
		}
		c <- nil
	}()

	go func() {
		if err := v.ValidateNameUniqueness(policy); err != nil {
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

func (v *PoliciesValid) ValidateDeletion(policy *models.Policy) error {
	if *policy.Name == "guest" {
		return errs.NewErrValidation("guest policy cannot be deleted")
	}

	return nil
}

func (v *PoliciesValid) ValidateUpdate(policy *models.Policy) error {
	if policy.Permissions == nil {
		return errs.NewErrValidation("policy permissions cannot be blank")
	}

	if err := v.ValidateResourcesExistence(policy); err != nil {
		return err
	}

	return nil
}

func (v *PoliciesValid) ValidateNameUniqueness(policy *models.Policy) error {
	if policy.Name == nil {
		return nil
	}

	err := v.r.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket([]byte("policies")).Get([]byte(*policy.Name))

		if len(raw) != 0 {
			return errs.NewErrValidation("name must be unique")
		}

		return nil
	})

	return err
}

func (v *PoliciesValid) ValidateResourcesExistence(policy *models.Policy) error {
	resources := []models.Resource{}

	err := v.r.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("resources")).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			resource := models.Resource{}
			if err := json.Unmarshal(v, &resource); err != nil {
				return err
			}
			resources = append(resources, resource)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, permission := range policy.Permissions {
		if permission.Resource == nil || len(*permission.Resource) == 0 {
			return errs.NewErrValidation("permission resources cannot be blank")
		}

		found := false

		if *permission.Resource == "*" {
			continue
		}

		for _, resource := range resources {
			if *permission.Resource == *resource.Name {
				found = true
				break
			}
		}

		if !found {
			return errs.NewErrValidation(fmt.Sprintf("resource doesn't exists or is invalid: '%s'", *permission.Resource))
		}
	}

	return nil
}
