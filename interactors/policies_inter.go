package interactors

import (
	"encoding/json"
	"errors"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewPoliciesInter)
}

type (
	PoliciesInterPoliciesRepo interface {
		Update(func(tx *bolt.Tx) error) error
		View(func(tx *bolt.Tx) error) error
	}

	PoliciesInterSessionsInter interface {
		DeleteCascade(policy *models.Policy) error
	}

	PoliciesInterPoliciesValidator interface {
		ValidateDeletion(policy *models.Policy) error
	}

	PoliciesInter struct {
		r  PoliciesInterPoliciesRepo
		si PoliciesInterSessionsInter
		v  PoliciesInterPoliciesValidator
	}
)

func NewPoliciesInter(
	r PoliciesInterPoliciesRepo,
	si PoliciesInterSessionsInter,
	v PoliciesInterPoliciesValidator,
) *PoliciesInter {
	return &PoliciesInter{r: r, si: si, v: v}
}

func (i *PoliciesInter) Find() ([]models.Policy, error) {
	policies := []models.Policy{}

	err := i.r.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("policies")).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			policy := models.Policy{}
			if err := json.Unmarshal(v, &policy); err != nil {
				return err
			}
			policies = append(policies, policy)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (i *PoliciesInter) FindByName(name string) (*models.Policy, error) {
	var raw []byte

	err := i.r.View(func(tx *bolt.Tx) error {
		raw = tx.Bucket([]byte("policies")).Get([]byte(name))

		return nil
	})

	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, errs.Internal.NotFound
	}

	policy := &models.Policy{}

	if err := json.Unmarshal(raw, policy); err != nil {
		return nil, err
	}

	return policy, nil
}

func (i *PoliciesInter) Create(policy *models.Policy) (*models.Policy, error) {
	if policy == nil {
		return nil, errors.New("nil policy")
	}

	err := i.r.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("policies"))

		raw, _ := json.Marshal(policy)

		return b.Put([]byte(*policy.Name), raw)
	})

	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (i *PoliciesInter) DeleteByName(name string) (*models.Policy, error) {
	policy, err := i.FindByName(name)
	if err != nil {
		return nil, err
	}

	if err := i.v.ValidateDeletion(policy); err != nil {
		return nil, err
	}

	err = i.r.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("policies")).Delete([]byte(name))
	})

	if err != nil {
		return nil, err
	}

	if err := i.si.DeleteCascade(policy); err != nil {
		return nil, err
	}

	return policy, nil
}

func (i *PoliciesInter) DeleteCascade(resource *models.Resource) error {
	if resource == nil {
		return errors.New("nil resource")
	}

	err := i.r.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("policies")).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			policy := models.Policy{}
			if err := json.Unmarshal(v, &policy); err != nil {
				return err
			}

			newPermissions := []models.Permission{}

			for _, permission := range policy.Permissions {
				if *resource.Name == *permission.Resource {
					continue
				}

				newPermissions = append(newPermissions, permission)
			}

			policy.Permissions = newPermissions

			raw, _ := json.Marshal(policy)

			if err := c.Bucket().Put([]byte(*policy.Name), raw); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (i *PoliciesInter) UpdateByName(name string, policy *models.Policy) (*models.Policy, error) {
	if policy == nil {
		return nil, errors.New("nil policy")
	}

	oldPolicy, err := i.FindByName(name)
	if err != nil {
		return nil, err
	}

	policy.Name = oldPolicy.Name

	err = i.r.Update(func(tx *bolt.Tx) error {
		raw, _ := json.Marshal(policy)
		return tx.Bucket([]byte("policies")).Put([]byte(name), raw)
	})

	if err != nil {
		return nil, err
	}

	return policy, nil
}
