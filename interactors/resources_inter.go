package interactors

import (
	"encoding/json"
	"errors"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewResourcesInter)
}

type (
	ResourcesInterResourcesRepo interface {
		Update(func(tx *bolt.Tx) error) error
		View(func(tx *bolt.Tx) error) error
	}

	ResourcesInterPoliciesInter interface {
		DeleteCascade(resource *models.Resource) error
	}

	ResourcesInter struct {
		r  ResourcesInterResourcesRepo
		pi ResourcesInterPoliciesInter
	}
)

func NewResourcesInter(r ResourcesInterResourcesRepo, pi ResourcesInterPoliciesInter) *ResourcesInter {
	return &ResourcesInter{r: r, pi: pi}
}

func (i *ResourcesInter) Find() ([]models.Resource, error) {
	resources := []models.Resource{}

	err := i.r.View(func(tx *bolt.Tx) error {
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
		return nil, err
	}

	return resources, nil
}

func (i *ResourcesInter) FindByHostname(hostname string) (*models.Resource, error) {
	var raw []byte

	err := i.r.View(func(tx *bolt.Tx) error {
		raw = tx.Bucket([]byte("resources")).Get([]byte(hostname))

		return nil
	})

	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, errs.Internal.NotFound
	}

	resource := &models.Resource{}

	if err := json.Unmarshal(raw, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (i *ResourcesInter) Create(resource *models.Resource) (*models.Resource, error) {
	if resource == nil {
		return nil, errors.New("nil resource")
	}

	err := i.r.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))

		raw, _ := json.Marshal(resource)

		return b.Put([]byte(*resource.Hostname), raw)
	})

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (i *ResourcesInter) DeleteByHostname(hostname string) (*models.Resource, error) {
	resource, err := i.FindByHostname(hostname)
	if err != nil {
		return nil, err
	}

	err = i.r.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resources")).Delete([]byte(hostname))
	})

	if err != nil {
		return nil, err
	}

	if err := i.pi.DeleteCascade(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (i *ResourcesInter) UpdateByHostname(hostname string, resource *models.Resource) (*models.Resource, error) {
	if resource == nil {
		return nil, errors.New("nil resource")
	}

	oldResource, err := i.FindByHostname(hostname)
	if err != nil {
		return nil, err
	}

	resource.Hostname = oldResource.Hostname

	err = i.r.Update(func(tx *bolt.Tx) error {
		raw, _ := json.Marshal(resource)
		return tx.Bucket([]byte("resources")).Put([]byte(hostname), raw)
	})

	if err != nil {
		return nil, err
	}

	return resource, nil
}
