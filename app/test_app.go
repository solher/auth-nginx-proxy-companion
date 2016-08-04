package app

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/solher/zest"

	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/solher/auth-nginx-proxy-companion/utils"
)

var testResource = &models.Resource{
	RedirectURL: utils.StrCpy("http://www.google.com"),
	Name:        utils.StrCpy("Foobar"),
	Hostname:    utils.StrCpy("foo.bar.com"),
}

var testResource2 = &models.Resource{
	RedirectURL: utils.StrCpy("http://www.google.com"),
	Name:        utils.StrCpy("Foobar2"),
	Hostname:    utils.StrCpy("foo.bar.2.com"),
}

var testSession1 = &models.Session{
	Token:      utils.StrCpy("F00bAr"),
	ValidTo:    utils.TimeCpy(time.Now().UTC().Add(time.Hour)),
	OwnerToken: utils.StrCpy("owner1"),
	Policies:   []string{"Foo", "Bar"},
}

var testSession2 = &models.Session{
	Token:      utils.StrCpy("F00bAr2"),
	ValidTo:    utils.TimeCpy(time.Now()),
	OwnerToken: utils.StrCpy("owner1"),
	Policies:   []string{"Foo", "Bar"},
}

var testSession3 = &models.Session{
	Token:      utils.StrCpy("F00bAr3"),
	ValidTo:    utils.TimeCpy(time.Now().UTC().Add(time.Hour)),
	OwnerToken: utils.StrCpy("owner1"),
	Policies:   []string{"Foo", "Bar"},
}

var testSession4 = &models.Session{
	Token:      utils.StrCpy("F00bAr4"),
	ValidTo:    utils.TimeCpy(time.Now().UTC().Add(time.Hour)),
	OwnerToken: utils.StrCpy("owner2"),
	Policies:   []string{"Foo", "Bar"},
}

var testSession5 = &models.Session{
	Token:    utils.StrCpy("F00bAr5"),
	ValidTo:  utils.TimeCpy(time.Now().UTC().Add(time.Hour)),
	Policies: []string{"Foo", "Bar"},
}

var testPolicy1 = &models.Policy{
	Name: utils.StrCpy("Foo"),
	Permissions: []models.Permission{
		{
			Resource: utils.StrCpy("Foobar"),
			Paths:    []string{"/foo/*"},
			Deny:     utils.BoolCpy(true),
			Enabled:  utils.BoolCpy(true),
		},
		{
			Resource: utils.StrCpy("Foobar"),
			Paths:    []string{"/foo/bar"},
			Deny:     utils.BoolCpy(false),
		},
		{
			Resource: utils.StrCpy("Foobar2"),
			Paths:    []string{"/test"},
		},
		{
			Resource: utils.StrCpy("Foobar"),
			Paths:    []string{"/bar"},
			Deny:     utils.BoolCpy(true),
		},
		{
			Resource: utils.StrCpy("Foobar"),
		},
	},
}

var testPolicy2 = &models.Policy{
	Name:        utils.StrCpy("Bar"),
	Permissions: []models.Permission{},
}

var guestPolicy = &models.Policy{
	Name:    utils.StrCpy("guest"),
	Enabled: utils.BoolCpy(true),
	Permissions: []models.Permission{
		{
			Resource: utils.StrCpy("Foobar"),
			Paths:    []string{"/*"},
		},
	},
}

type TestApp struct {
	dbLocation, gcLocation string
}

func NewTestApp() *TestApp {
	rnd := strconv.Itoa(utils.RandInt(0, 1000000))

	testApp := &TestApp{
		dbLocation: "auth-server-test-" + rnd + ".db",
		gcLocation: "archived-test-" + rnd + ".db",
	}

	return testApp
}

func (a *TestApp) Launch() (string, error) {
	appPort := utils.RandInt(49152, 65535)
	testURL := "http://localhost:" + strconv.Itoa(appPort)

	if err := a.Seed(); err != nil {
		return "", err
	}

	overrideConst := func(z *zest.Zest) error {
		d := &struct{ Const *Constants }{}

		if err := z.Injector.Get(d); err != nil {
			return err
		}

		d.Const.App.Port = appPort
		d.Const.DB.Location = a.dbLocation
		d.Const.GC.Location = a.gcLocation

		return nil
	}

	go Run(overrideConst)

	for {
		if _, err := http.Get(testURL); err == nil {
			break
		}
	}

	return testURL, nil
}
func (a *TestApp) Seed() error {
	db, err := bolt.Open(a.dbLocation, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return err
	}

	// Test policies creation
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("policies"))
		if err != nil {
			return err
		}

		raw, _ := json.Marshal(testPolicy1)

		if err := b.Put([]byte(*testPolicy1.Name), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(testPolicy2)

		if err := b.Put([]byte(*testPolicy2.Name), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(guestPolicy)

		if err := b.Put([]byte(*guestPolicy.Name), raw); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Test resources creation
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("resources"))
		if err != nil {
			return err
		}

		raw, _ := json.Marshal(testResource)

		if err := b.Put([]byte(*testResource.Hostname), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(testResource2)

		if err := b.Put([]byte(*testResource2.Hostname), raw); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Test sessions creation
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("sessions"))
		if err != nil {
			return err
		}

		raw, _ := json.Marshal(testSession1)

		if err := b.Put([]byte(*testSession1.Token), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(testSession2)

		if err := b.Put([]byte(*testSession2.Token), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(testSession3)

		if err := b.Put([]byte(*testSession3.Token), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(testSession4)

		if err := b.Put([]byte(*testSession4.Token), raw); err != nil {
			return err
		}

		raw, _ = json.Marshal(testSession5)

		if err := b.Put([]byte(*testSession5.Token), raw); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

func (a *TestApp) Stop() error {
	if err := os.Remove(a.dbLocation); err != nil {
		return err
	}

	if err := os.Remove(a.gcLocation); err != nil {
		return err
	}

	return nil
}
