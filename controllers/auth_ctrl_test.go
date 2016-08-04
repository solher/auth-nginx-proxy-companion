package controllers

import (
	"net/http/httptest"
	"testing"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authCtrlAuthInter struct {
	errDB, errNotFound bool
	sessionNotFound    bool
	denyAccess         bool
	noRedirectURL      bool
}

func (i *authCtrlAuthInter) AuthorizeToken(hostname, path, token string) (bool, *models.Session, error) {
	if i.errDB {
		return false, nil, errs.Internal.Database
	}

	if i.errNotFound {
		return false, nil, errs.Internal.NotFound
	}

	session := &models.Session{
		Payload: utils.StrCpy("{}"),
	}

	if i.sessionNotFound {
		session = nil
	}

	return !i.denyAccess, session, nil
}

func (i *authCtrlAuthInter) GetRedirectURL(hostname string) (string, error) {
	if i.errDB {
		return "", errs.Internal.Database
	}

	if i.errNotFound {
		return "", errs.Internal.NotFound
	}

	if i.noRedirectURL {
		return "", nil
	}

	return "http://foo.bar", nil
}

// TestAuthCtrlAuthorizeToken runs tests on the AuthCtrl AuthorizeToken method.
func TestAuthCtrlAuthorizeToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	render := utils.NewFakeRender()
	getter := utils.NewFakeModelsGetter()
	inter := &authCtrlAuthInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewAuthCtrl(inter, render, getter)

	getter.GrantAll = true

	// Everything is granted
	req := utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	req.Header.Set("Request-URL", "http://foo/bar")
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(204, render.Status)
	a.Len(recorder.Header().Get("Auth-Server-Payload"), 0)
	utils.Clear(nil, render, recorder)

	getter.GrantAll = false

	// No error, access is granted
	req = utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	req.Header.Set("Request-URL", "http://foo/bar")
	req.Header.Set("Auth-Server-Token", "kjgcjgh576cg4")
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(204, render.Status)
	a.NotNil(recorder.Header().Get("Auth-Server-Payload"))
	a.NotNil(recorder.Header().Get("Auth-Server-Session"))
	a.NotNil(recorder.Header().Get("Auth-Server-Token"))
	utils.Clear(nil, render, recorder)

	// No error, data via query paprameters
	req = utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	values := req.URL.Query()
	values.Set("requestUrl", "http://foo/bar")
	values.Set("accessToken", "kjgcjgh576cg4")
	req.URL.RawQuery = values.Encode()
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(204, render.Status)
	a.NotNil(recorder.Header().Get("Auth-Server-Payload"))
	utils.Clear(nil, render, recorder)

	// Error, no request URL
	ctrl.AuthorizeToken(recorder, utils.FakeRequest("GET", "http://foo.bar/auth", nil))
	r.Equal(500, render.Status)
	utils.Clear(nil, render, recorder)

	inter.sessionNotFound = true

	// No error, access is granted thanks to the guest policy
	req = utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	req.Header.Set("Request-URL", "http://foo/bar")
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(204, render.Status)
	a.Equal(0, len(recorder.Header().Get("Auth-Server-Payload")))
	utils.Clear(nil, render, recorder)

	inter.sessionNotFound = false
	inter.denyAccess = true

	// Unauthorized: access is denied by a policy
	req = utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	req.Header.Set("Request-URL", "http://foo/bar")
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(403, render.Status)
	a.Len(recorder.Header().Get("Auth-Server-Payload"), 0)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(nil, render, recorder)

	inter.denyAccess = false
	inter.errNotFound = true

	// Unauthorized: a matching resource was not found
	req = utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	req.Header.Set("Request-URL", "http://foo/bar")
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(403, render.Status)
	a.Len(recorder.Header().Get("Auth-Server-Payload"), 0)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(nil, render, recorder)

	inter.errNotFound = false
	inter.errDB = true

	// The interactor returns a database error
	req = utils.FakeRequest("GET", "http://foo.bar/auth", nil)
	req.Header.Set("Request-URL", "http://foo/bar")
	ctrl.AuthorizeToken(recorder, req)
	r.Equal(500, render.Status)
	a.Len(recorder.Header().Get("Auth-Server-Payload"), 0)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(nil, render, recorder)
}

// TestAuthCtrlRedirect runs tests on the AuthCtrl Redirect method.
func TestAuthCtrlRedirect(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	render := utils.NewFakeRender()
	getter := utils.NewFakeModelsGetter()
	getter.RedirectURL = "http://default.com"
	inter := &authCtrlAuthInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewAuthCtrl(inter, render, getter)

	// Success: a resource is found and a redirect URL is set
	req := utils.FakeRequest("GET", "http://foo.bar/redirect", nil)
	req.Header.Add("Request-Url", "http://request.com")
	ctrl.Redirect(recorder, req)
	r.Equal(307, recorder.Code)
	a.NotEqual(0, len(recorder.Header().Get("Location")))
	a.Contains(recorder.Header().Get("Location"), "?redirectUrl=http://request.com")
	a.Equal("http://request.com", recorder.Header().Get("Redirect-Url"))
	utils.Clear(nil, render, recorder)

	inter.noRedirectURL = true

	// Success: a resource is found and no redirect URL is set
	req = utils.FakeRequest("GET", "http://foo.bar/redirect", nil)
	req.Header.Add("Request-Url", "http://request.com")
	ctrl.Redirect(recorder, req)
	r.Equal(307, recorder.Code)
	a.NotEqual(0, len(recorder.Header().Get("Location")))
	utils.Clear(nil, render, recorder)

	inter.noRedirectURL = false
	inter.errNotFound = true

	// Success: no resource is found
	req = utils.FakeRequest("GET", "http://foo.bar/redirect", nil)
	req.Header.Add("Request-Url", "http://request.com")
	ctrl.Redirect(recorder, req)
	r.Equal(307, recorder.Code)
	a.NotEqual(0, len(recorder.Header().Get("Location")))
	utils.Clear(nil, render, recorder)

	inter.errNotFound = false
	inter.errDB = true

	// Error: internal
	req = utils.FakeRequest("GET", "http://foo.bar/redirect", nil)
	req.Header.Add("Request-Url", "http://request.com")
	ctrl.Redirect(recorder, req)
	r.Equal(500, render.Status)
	a.Equal(0, len(recorder.Header().Get("Location")))
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(nil, render, recorder)
}
