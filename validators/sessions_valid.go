package validators

import (
	"fmt"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewSessionsValid)
}

type (
	SessionsValidSessionsRepo interface {
		View(func(tx *bolt.Tx) error) error
	}

	SessionsValid struct {
		r SessionsValidSessionsRepo
	}
)

func NewSessionsValid(r SessionsValidSessionsRepo) *SessionsValid {
	return &SessionsValid{r: r}
}

func (v *SessionsValid) ValidateCreation(session *models.Session) error {
	c := make(chan error, 2)

	if session.Policies == nil {
		return errs.NewErrValidation("session policies cannot be blank")
	}

	go func() {
		if err := v.ValidateTokenUniqueness(session); err != nil {
			c <- err
		}
		c <- nil
	}()

	go func() {
		if err := v.ValidatePolicyExistence(session); err != nil {
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

func (v *SessionsValid) ValidateTokenUniqueness(session *models.Session) error {
	if session.Token == nil {
		return nil
	}

	err := v.r.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket([]byte("sessions")).Get([]byte(*session.Token))

		if len(raw) != 0 {
			return errs.NewErrValidation("token must be unique")
		}

		return nil
	})

	return err
}

func (v *SessionsValid) ValidatePolicyExistence(session *models.Session) error {
	err := v.r.View(func(tx *bolt.Tx) error {
		for _, policyID := range session.Policies {
			raw := tx.Bucket([]byte("policies")).Get([]byte(policyID))

			if len(raw) == 0 {
				return errs.NewErrValidation(fmt.Sprintf("policy doesn't exists or is invalid: '%s'", policyID))
			}
		}

		return nil
	})

	return err
}
