package interactors

import (
	"testing"
	"time"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testResource = &models.Resource{
	Name:     utils.StrCpy("Foobar"),
	Hostname: utils.StrCpy("foo.bar.com"),
}

var testSession = &models.Session{
	Token:    utils.StrCpy("F00bAr"),
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
			Resource: utils.StrCpy("Foobar2"), // Will be ignored
			Paths:    []string{"/something"},
		},
		{
			Resource: utils.StrCpy("Foobar"),
			Paths:    []string{"/bar", "/bar2"},
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

type authInterPoliciesInter struct {
	errDB, errNotFound                         bool
	disableGuestPolicy, disableGuestPermission bool
	policy                                     *models.Policy
}

func (r *authInterPoliciesInter) FindByName(name string) (*models.Policy, error) {
	if r.errDB {
		return nil, errs.Internal.Database
	}

	if r.errNotFound {
		return nil, errs.Internal.NotFound
	}

	switch name {
	case "guest":
		return guestPolicy, nil
	case "Foo":
		return testPolicy1, nil
	case "Bar":
		return testPolicy2, nil
	}

	return r.policy, nil
}

type authInterResourcesInter struct {
	errDB, errNotFound bool
	resource           *models.Resource
}

func (r *authInterResourcesInter) FindByHostname(hostname string) (*models.Resource, error) {
	if r.errDB {
		return nil, errs.Internal.Database
	}

	if r.errNotFound {
		return nil, errs.Internal.NotFound
	}

	switch hostname {
	case "foo.bar.com":
		return testResource, nil
	}

	return r.resource, nil
}

type authInterSessionsInter struct {
	errDB, errNotFound bool
	session            *models.Session
}

func (r *authInterSessionsInter) FindByToken(token string) (*models.Session, error) {
	if r.errDB {
		return nil, errs.Internal.Database
	}

	if r.errNotFound {
		return nil, errs.Internal.NotFound
	}

	switch token {
	case "F00bAr":
		return testSession, nil
	}

	return r.session, nil
}

// TestAuthInterAuthorizeToken runs tests on the AuthInter AuthorizeToken method.
func TestAuthInterAuthorizeToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	policiesInter := &authInterPoliciesInter{}
	resourcesInter := &authInterResourcesInter{}
	sessionsInter := &authInterSessionsInter{}
	inter := NewAuthInter(
		policiesInter,
		resourcesInter,
		sessionsInter,
	)
	hostname := "foo.bar.com"
	path := ""
	token := "F00bAr"

	// Success: root
	granted, session, err := inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.True(granted)
	a.NotNil(session)

	path = "/foo/bar"

	// Success: weight system
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.True(granted)
	a.NotNil(session)

	path = "/foo/bar/"

	// Success: trailing slash
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.True(granted)
	a.NotNil(session)

	path = "/foo/"

	// Success: trailing slash
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.True(granted)
	a.NotNil(session)

	path = "/bar"

	// Multipath denied
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.False(granted)

	path = "/bar2"

	// Multipath denied
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.False(granted)

	path = "/foo/foo"

	// Denied
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.False(granted)

	testResource.Public = utils.BoolCpy(true)

	// Success: public resource
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.True(granted)
	a.Nil(session)

	testResource.Public = utils.BoolCpy(false)
	sessionsInter.errNotFound = true

	// Success: guest policy
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.True(granted)
	a.Nil(session)

	guestPolicy.Enabled = utils.BoolCpy(false)

	// Denied: guest policy is disabled
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.False(granted)
	a.Nil(session)

	guestPolicy.Enabled = utils.BoolCpy(true)
	guestPolicy.Permissions[0].Enabled = utils.BoolCpy(false)

	// Denied: guest policy permissions are disabled
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.NoError(err)
	a.False(granted)
	a.Nil(session)

	guestPolicy.Permissions[0].Enabled = utils.BoolCpy(true)
	policiesInter.errNotFound = true

	// Error: guest policy
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.Error(err)
	a.False(granted)
	a.Nil(session)

	sessionsInter.errNotFound = false

	// Not found error
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.False(granted)

	policiesInter.errNotFound = false
	resourcesInter.errNotFound = true

	// Not found error
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.Error(err)
	a.IsType(errs.Internal.NotFound, err)
	a.False(granted)

	resourcesInter.errNotFound = false
	policiesInter.errDB = true

	// Database error
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.False(granted)

	policiesInter.errDB = false
	resourcesInter.errDB = true

	// Database error
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.False(granted)

	resourcesInter.errDB = false
	sessionsInter.errDB = true

	// Database error
	granted, session, err = inter.AuthorizeToken(hostname, path, token)
	r.Error(err)
	a.IsType(errs.Internal.Database, err)
	a.False(granted)
}
