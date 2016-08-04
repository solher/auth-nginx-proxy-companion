package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionsCtrlSessionsInter struct {
	errDB, errNotFound bool
}

func (i *sessionsCtrlSessionsInter) Find() ([]models.Session, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	sessions := []models.Session{{}, {}, {}}

	return sessions, nil
}

func (i *sessionsCtrlSessionsInter) FindByToken(token string) (*models.Session, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	session := &models.Session{}

	return session, nil
}

func (i *sessionsCtrlSessionsInter) Create(session *models.Session) (*models.Session, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	return session, nil
}

func (i *sessionsCtrlSessionsInter) DeleteByToken(token string) (*models.Session, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	session := &models.Session{}

	return session, nil
}

func (i *sessionsCtrlSessionsInter) DeleteByOwnerTokens(ownerToken []string) ([]models.Session, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	sessions := []models.Session{{}, {}}

	return sessions, nil
}

type sessionsCtrlSessionsValid struct {
	errValid bool
}

func (v *sessionsCtrlSessionsValid) ValidateCreation(session *models.Session) error {
	if v.errValid {
		return errs.NewErrValidation("validation error")
	}

	return nil
}

// TestSessionsCtrlFind runs tests on the SessionsCtrl Find method.
func TestSessionsCtrlFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &sessionsCtrlSessionsInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewSessionsCtrl(inter, render, params, nil)
	sessionsOut := []models.Session{}

	// No error, 3 sessions are returned
	ctrl.Find(recorder, utils.FakeRequest("GET", "http://foo.bar/sessions", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(&sessionsOut)
	r.NoError(err)
	a.Len(sessionsOut, 3)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.Find(recorder, utils.FakeRequest("GET", "http://foo.bar/sessions", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestSessionsCtrlFindByToken runs tests on the SessionsCtrl FindByToken method.
func TestSessionsCtrlFindByToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &sessionsCtrlSessionsInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewSessionsCtrl(inter, render, params, nil)
	sessionOut := &models.Session{}

	// No error, a session is returned
	ctrl.FindByToken(recorder, utils.FakeRequest("GET", "http://foo.bar/sessions/jhHgchgV", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(sessionOut)
	r.NoError(err)
	a.NotNil(sessionOut)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.FindByToken(recorder, utils.FakeRequest("GET", "http://foo.bar/sessions/jhHgchgV", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	// Session not found
	inter.errDB = false
	inter.errNotFound = true
	ctrl.FindByToken(recorder, utils.FakeRequest("GET", "http://foo.bar/sessions/jhHgchgV", nil))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestSessionsCtrlCreate runs tests on the SessionsCtrl Create method.
func TestSessionsCtrlCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &sessionsCtrlSessionsInter{}
	valid := &sessionsCtrlSessionsValid{}
	recorder := httptest.NewRecorder()
	ctrl := NewSessionsCtrl(inter, render, params, valid)
	sessionIn := &models.Session{}
	sessionOut := &models.Session{}

	valid.errValid = true

	// Validation error
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/sessions", sessionIn))
	r.Equal(422, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Validation, render.APIError)
	utils.Clear(params, render, recorder)

	valid.errValid = false

	// No error, one session is created
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/sessions", sessionIn))
	r.Equal(201, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(sessionOut)
	r.NoError(err)
	a.NotNil(sessionOut)
	utils.Clear(params, render, recorder)

	// Null body decoding error
	ctrl.Create(recorder, utils.FakeRequestRaw("POST", "http://foo.bar/sessions", nil))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	// Body decoding error
	ctrl.Create(recorder, utils.FakeRequestRaw("POST", "http://foo.bar/sessions", []byte{'{'}))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/sessions", sessionIn))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestSessionsCtrlDeleteByToken runs tests on the SessionsCtrl DeleteByToken method.
func TestSessionsCtrlDeleteByToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &sessionsCtrlSessionsInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewSessionsCtrl(inter, render, params, nil)
	sessionOut := &models.Session{}

	// No error, a session is returned
	ctrl.DeleteByToken(recorder, utils.FakeRequest("DELETE", "http://foo.bar/sessions/jhHgchgV", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(sessionOut)
	r.NoError(err)
	a.NotNil(sessionOut)
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.DeleteByToken(recorder, utils.FakeRequest("DELETE", "http://foo.bar/sessions/jhHgchgV", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	inter.errDB = false
	inter.errNotFound = true

	// Session not found
	ctrl.DeleteByToken(recorder, utils.FakeRequest("DELETE", "http://foo.bar/sessions/jhHgchgV", nil))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestSessionsCtrlDeleteByOwnerToken runs tests on the SessionsCtrl DeleteByOwnerToken method.
func TestSessionsCtrlDeleteByOwnerToken(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &sessionsCtrlSessionsInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewSessionsCtrl(inter, render, params, nil)
	sessionsOut := []models.Session{}

	// Error: invalid query params
	ctrl.DeleteByOwnerToken(recorder, utils.FakeRequest("DELETE", `http://foo.bar/sessions?ownerTokens=toto`, nil))
	r.Equal(400, render.Status)
	utils.Clear(params, render, recorder)

	// Error: invalid query params
	ctrl.DeleteByOwnerToken(recorder, utils.FakeRequest("DELETE", `http://foo.bar/sessions`, nil))
	r.Equal(400, render.Status)
	utils.Clear(params, render, recorder)

	// No error, deleted sessions are returned
	ctrl.DeleteByOwnerToken(recorder, utils.FakeRequest("DELETE", `http://foo.bar/sessions?ownerTokens=["foobar"]`, nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(&sessionsOut)
	r.NoError(err)
	a.NotEqual(0, len(sessionsOut))
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.DeleteByOwnerToken(recorder, utils.FakeRequest("DELETE", `http://foo.bar/sessions?ownerTokens=["foobar"]`, nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}
