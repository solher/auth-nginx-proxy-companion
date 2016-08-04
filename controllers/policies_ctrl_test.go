package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/solher/auth-nginx-proxy-companion/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type policiesCtrlPoliciesInter struct {
	errDB, errNotFound, errValidation bool
}

func (i *policiesCtrlPoliciesInter) Find() ([]models.Policy, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	policies := []models.Policy{{}, {}, {}}

	return policies, nil
}

func (i *policiesCtrlPoliciesInter) FindByName(id string) (*models.Policy, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	policy := &models.Policy{}

	return policy, nil
}

func (i *policiesCtrlPoliciesInter) Create(policy *models.Policy) (*models.Policy, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	return policy, nil
}

func (i *policiesCtrlPoliciesInter) DeleteByName(id string) (*models.Policy, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	if i.errValidation {
		return nil, errs.Internal.Validation
	}

	policy := &models.Policy{}

	return policy, nil
}

func (i *policiesCtrlPoliciesInter) UpdateByName(id string, policy *models.Policy) (*models.Policy, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	return policy, nil
}

type policiesCtrlPoliciesValid struct {
	errValid bool
}

func (v *policiesCtrlPoliciesValid) ValidateCreation(policy *models.Policy) error {
	if v.errValid {
		return errs.NewErrValidation("validation error")
	}

	return nil
}

func (v *policiesCtrlPoliciesValid) ValidateUpdate(policy *models.Policy) error {
	if v.errValid {
		return errs.NewErrValidation("validation error")
	}

	return nil
}

// TestPoliciesCtrlFind runs tests on the PoliciesCtrl Find method.
func TestPoliciesCtrlFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &policiesCtrlPoliciesInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewPoliciesCtrl(inter, render, params, nil)
	policiesOut := []models.Policy{}

	// No error, 3 policies are returned
	ctrl.Find(recorder, utils.FakeRequest("GET", "http://foo.bar/policies", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(&policiesOut)
	r.NoError(err)
	a.Len(policiesOut, 3)
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.Find(recorder, utils.FakeRequest("GET", "http://foo.bar/policies", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestPoliciesCtrlFindByName runs tests on the PoliciesCtrl FindByName method.
func TestPoliciesCtrlFindByName(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &policiesCtrlPoliciesInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewPoliciesCtrl(inter, render, params, nil)
	policyOut := &models.Policy{}

	// No error, a policy is returned
	ctrl.FindByName(recorder, utils.FakeRequest("GET", "http://foo.bar/policies/foobar", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut)
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.FindByName(recorder, utils.FakeRequest("GET", "http://foo.bar/policies/foobar", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	inter.errDB = false
	inter.errNotFound = true

	// Policy not found
	ctrl.FindByName(recorder, utils.FakeRequest("GET", "http://foo.bar/policies/foobar", nil))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestPoliciesCtrlCreate runs tests on the PoliciesCtrl Create method.
func TestPoliciesCtrlCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &policiesCtrlPoliciesInter{}
	valid := &policiesCtrlPoliciesValid{}
	recorder := httptest.NewRecorder()
	ctrl := NewPoliciesCtrl(inter, render, params, valid)
	policyIn := &models.Policy{Name: utils.StrCpy("foobar")}
	policyOut := &models.Policy{}

	valid.errValid = true

	// Validation error
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/policies", policyIn))
	r.Equal(422, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Validation, render.APIError)
	utils.Clear(params, render, recorder)

	valid.errValid = false

	// No error, one policy is created
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/policies", policyIn))
	r.Equal(201, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut)
	utils.Clear(params, render, recorder)

	// Null body decoding error
	ctrl.Create(recorder, utils.FakeRequestRaw("POST", "http://foo.bar/policies", nil))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	// Body decoding error
	ctrl.Create(recorder, utils.FakeRequestRaw("POST", "http://foo.bar/policies", []byte{'{'}))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/policies", policyIn))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestPoliciesCtrlDeleteByName runs tests on the PoliciesCtrl DeleteByName method.
func TestPoliciesCtrlDeleteByName(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &policiesCtrlPoliciesInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewPoliciesCtrl(inter, render, params, nil)
	policyOut := &models.Policy{}

	// No error, a policy is returned
	ctrl.DeleteByName(recorder, utils.FakeRequest("DELETE", "http://foo.bar/policies/foobar", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut)
	utils.Clear(params, render, recorder)

	inter.errDB = true

	// The interactor returns a database error
	ctrl.DeleteByName(recorder, utils.FakeRequest("DELETE", "http://foo.bar/policies/foobar", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	inter.errDB = false
	inter.errValidation = true

	// The interactor returns a validation error
	ctrl.DeleteByName(recorder, utils.FakeRequest("DELETE", "http://foo.bar/policies/foobar", nil))
	r.Equal(422, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Validation, render.APIError)
	utils.Clear(params, render, recorder)

	inter.errValidation = false
	inter.errNotFound = true

	// Policy not found
	ctrl.DeleteByName(recorder, utils.FakeRequest("DELETE", "http://foo.bar/policies/foobar", nil))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestPoliciesCtrlUpdateByName runs tests on the PoliciesCtrl UpdateByName method.
func TestPoliciesCtrlUpdateByName(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &policiesCtrlPoliciesInter{}
	valid := &policiesCtrlPoliciesValid{}
	recorder := httptest.NewRecorder()
	ctrl := NewPoliciesCtrl(inter, render, params, valid)
	policyIn := &models.Policy{Name: utils.StrCpy("foobar")}
	policyOut := &models.Policy{}

	valid.errValid = true

	// Validation error
	ctrl.UpdateByName(recorder, utils.FakeRequest("PUT", "http://foo.bar/policies/foobar", policyIn))
	r.Equal(422, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Validation, render.APIError)
	utils.Clear(params, render, recorder)

	valid.errValid = false

	// No error, a policy is returned
	ctrl.UpdateByName(recorder, utils.FakeRequest("PUT", "http://foo.bar/policies/foobar", policyIn))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut)
	a.Nil(policyOut.Name)
	utils.Clear(params, render, recorder)

	// Null body decoding error
	ctrl.UpdateByName(recorder, utils.FakeRequestRaw("PUT", "http://foo.bar/policies/foobar", nil))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.UpdateByName(recorder, utils.FakeRequest("PUT", "http://foo.bar/policies/foobar", policyIn))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	// Policy not found
	inter.errDB = false
	inter.errNotFound = true
	ctrl.UpdateByName(recorder, utils.FakeRequest("PUT", "http://foo.bar/policies/foobar", policyIn))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}
