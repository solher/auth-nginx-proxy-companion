package interactors

import (
	"encoding/json"
	"errors"
	"time"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/boltdb/bolt"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewSessionsInter)
}

type (
	SessionsInterSessionsRepo interface {
		Update(func(tx *bolt.Tx) error) error
		View(func(tx *bolt.Tx) error) error
	}

	SessionOptionsGetter interface {
		GetSessionValidity() time.Duration
		GetSessionTokenLength() int
	}

	SessionsInter struct {
		r SessionsInterSessionsRepo
		g SessionOptionsGetter
	}
)

func NewSessionsInter(r SessionsInterSessionsRepo, g SessionOptionsGetter) *SessionsInter {
	return &SessionsInter{r: r, g: g}
}

func (i *SessionsInter) Find() ([]models.Session, error) {
	sessions := []models.Session{}

	err := i.r.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("sessions")).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			session := models.Session{}
			if err := json.Unmarshal(v, &session); err != nil {
				return err
			}

			if session.ValidTo.Before(time.Now()) {
				continue
			}

			sessions = append(sessions, session)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (i *SessionsInter) FindByToken(token string) (*models.Session, error) {
	var raw []byte

	err := i.r.View(func(tx *bolt.Tx) error {
		raw = tx.Bucket([]byte("sessions")).Get([]byte(token))

		return nil
	})

	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, errs.Internal.NotFound
	}

	session := &models.Session{}

	if err := json.Unmarshal(raw, session); err != nil {
		return nil, err
	}

	if session.ValidTo.Before(time.Now()) {
		return nil, errs.Internal.NotFound
	}

	return session, nil
}

func (i *SessionsInter) Create(session *models.Session) (*models.Session, error) {
	if session == nil {
		return nil, errors.New("nil session")
	}

	now := time.Now().UTC()
	session.Created = &now

	if session.Token == nil {
		session.Token = utils.StrCpy(utils.GenToken(i.g.GetSessionTokenLength()))
	}

	if session.ValidTo == nil {
		session.ValidTo = utils.TimeCpy(now.Add(i.g.GetSessionValidity()))
	}

	err := i.r.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sessions"))

		raw, _ := json.Marshal(session)

		return b.Put([]byte(*session.Token), raw)
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *SessionsInter) DeleteByToken(token string) (*models.Session, error) {
	session, err := i.FindByToken(token)
	if err != nil {
		return nil, err
	}

	session.ValidTo = utils.TimeCpy(time.Now().UTC())

	err = i.r.Update(func(tx *bolt.Tx) error {
		raw, _ := json.Marshal(session)
		return tx.Bucket([]byte("sessions")).Put([]byte(token), raw)
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *SessionsInter) DeleteByOwnerTokens(ownerTokens []string) ([]models.Session, error) {
	sessions, err := i.Find()
	if err != nil {
		return nil, err
	}

	deletedSessions := []models.Session{}
	now := time.Now().UTC()

	err = i.r.Update(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("sessions"))

		for _, session := range sessions {
			for _, ownerToken := range ownerTokens {
				if session.OwnerToken == nil || *session.OwnerToken != ownerToken {
					continue
				}

				session.ValidTo = &now

				raw, _ := json.Marshal(session)

				if err := s.Put([]byte(*session.Token), raw); err != nil {
					return err
				}

				deletedSessions = append(deletedSessions, session)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return deletedSessions, nil
}

func (i *SessionsInter) DeleteCascade(policy *models.Policy) error {
	if policy == nil {
		return errors.New("nil policy")
	}

	err := i.r.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("sessions")).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			session := models.Session{}
			if err := json.Unmarshal(v, &session); err != nil {
				return err
			}

			newPolicies := []string{}

			for _, p := range session.Policies {
				if *policy.Name == p {
					continue
				}

				newPolicies = append(newPolicies, p)
			}

			session.Policies = newPolicies

			raw, _ := json.Marshal(session)

			if err := c.Bucket().Put([]byte(*session.Token), raw); err != nil {
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
