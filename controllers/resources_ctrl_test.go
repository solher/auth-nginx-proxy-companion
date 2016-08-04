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

type resourcesCtrlResourcesInter struct {
	errDB, errNotFound bool
}

func (i *resourcesCtrlResourcesInter) Find() ([]models.Resource, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	resources := []models.Resource{{}, {}, {}}

	return resources, nil
}

func (i *resourcesCtrlResourcesInter) FindByHostname(hostname string) (*models.Resource, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	resource := &models.Resource{}

	return resource, nil
}

func (i *resourcesCtrlResourcesInter) Create(resource *models.Resource) (*models.Resource, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	return resource, nil
}

func (i *resourcesCtrlResourcesInter) DeleteByHostname(hostname string) (*models.Resource, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	resource := &models.Resource{}

	return resource, nil
}

func (i *resourcesCtrlResourcesInter) UpdateByHostname(hostname string, resource *models.Resource) (*models.Resource, error) {
	if i.errDB {
		return nil, errs.Internal.Database
	}

	if i.errNotFound {
		return nil, errs.Internal.NotFound
	}

	return resource, nil
}

type resourcesCtrlResourcesValid struct {
	errValid bool
}

func (v *resourcesCtrlResourcesValid) ValidateCreation(resource *models.Resource) error {
	if v.errValid {
		return errs.NewErrValidation("validation error")
	}

	return nil
}

func (v *resourcesCtrlResourcesValid) ValidateUpdate(resource *models.Resource) error {
	if v.errValid {
		return errs.NewErrValidation("validation error")
	}

	return nil
}

// TestResourcesCtrlFind runs tests on the ResourcesCtrl Find method.
func TestResourcesCtrlFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &resourcesCtrlResourcesInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewResourcesCtrl(inter, render, params, nil)
	resourcesOut := []models.Resource{}

	// No error, 3 resources are returned
	ctrl.Find(recorder, utils.FakeRequest("GET", "http://foo.bar/resources", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(&resourcesOut)
	r.NoError(err)
	a.Len(resourcesOut, 3)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.Find(recorder, utils.FakeRequest("GET", "http://foo.bar/resources", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestResourcesCtrlFindByHostname runs tests on the ResourcesCtrl FindByHostname method.
func TestResourcesCtrlFindByHostname(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &resourcesCtrlResourcesInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewResourcesCtrl(inter, render, params, nil)
	resourceOut := &models.Resource{}

	// No error, a resource is returned
	ctrl.FindByHostname(recorder, utils.FakeRequest("GET", "http://foo.bar/resources/host.com", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.FindByHostname(recorder, utils.FakeRequest("GET", "http://foo.bar/resources/host.com", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	// Resource not found
	inter.errDB = false
	inter.errNotFound = true
	ctrl.FindByHostname(recorder, utils.FakeRequest("GET", "http://foo.bar/resources/host.com", nil))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestResourcesCtrlCreate runs tests on the ResourcesCtrl Create method.
func TestResourcesCtrlCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &resourcesCtrlResourcesInter{}
	valid := &resourcesCtrlResourcesValid{}
	recorder := httptest.NewRecorder()
	ctrl := NewResourcesCtrl(inter, render, params, valid)
	resourceIn := &models.Resource{}
	resourceOut := &models.Resource{}

	valid.errValid = true

	// Validation error
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/resources", resourceIn))
	r.Equal(422, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Validation, render.APIError)
	utils.Clear(params, render, recorder)

	valid.errValid = false

	// No error, one resource is created
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/resources", resourceIn))
	r.Equal(201, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut)
	utils.Clear(params, render, recorder)

	// Null body decoding error
	ctrl.Create(recorder, utils.FakeRequestRaw("POST", "http://foo.bar/resources", nil))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	// Body decoding error
	ctrl.Create(recorder, utils.FakeRequestRaw("POST", "http://foo.bar/resources", []byte{'{'}))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.Create(recorder, utils.FakeRequest("POST", "http://foo.bar/resources", resourceIn))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestResourcesCtrlDeleteByHostname runs tests on the ResourcesCtrl DeleteByHostname method.
func TestResourcesCtrlDeleteByHostname(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &resourcesCtrlResourcesInter{}
	recorder := httptest.NewRecorder()
	ctrl := NewResourcesCtrl(inter, render, params, nil)
	resourceOut := &models.Resource{}

	// No error, a resource is returned
	ctrl.DeleteByHostname(recorder, utils.FakeRequest("DELETE", "http://foo.bar/resources/host.com", nil))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.DeleteByHostname(recorder, utils.FakeRequest("DELETE", "http://foo.bar/resources/host.com", nil))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	// Resource not found
	inter.errDB = false
	inter.errNotFound = true
	ctrl.DeleteByHostname(recorder, utils.FakeRequest("DELETE", "http://foo.bar/resources/host.com", nil))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}

// TestResourcesCtrlUpdateByHostname runs tests on the ResourcesCtrl UpdateByHostname method.
func TestResourcesCtrlUpdateByHostname(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	params := utils.NewFakeParamsGetter()
	render := utils.NewFakeRender()
	inter := &resourcesCtrlResourcesInter{}
	valid := &resourcesCtrlResourcesValid{}
	recorder := httptest.NewRecorder()
	ctrl := NewResourcesCtrl(inter, render, params, valid)
	resourceIn := &models.Resource{
		Hostname: utils.StrCpy("foo.bar.com"),
	}
	resourceOut := &models.Resource{}

	valid.errValid = true

	// Validation error
	ctrl.UpdateByHostname(recorder, utils.FakeRequest("PUT", "http://foo.bar/resources/1", resourceIn))
	r.Equal(422, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Validation, render.APIError)
	utils.Clear(params, render, recorder)

	valid.errValid = false

	// No error, a resource is returned
	ctrl.UpdateByHostname(recorder, utils.FakeRequest("PUT", "http://foo.bar/resources/1", resourceIn))
	r.Equal(200, render.Status)
	err := json.NewDecoder(recorder.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut)
	a.Nil(resourceOut.Hostname)
	utils.Clear(params, render, recorder)

	// Null body decoding error
	ctrl.UpdateByHostname(recorder, utils.FakeRequestRaw("PUT", "http://foo.bar/resources/1", nil))
	r.Equal(400, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.BodyDecoding, render.APIError)
	utils.Clear(params, render, recorder)

	// The interactor returns a database error
	inter.errDB = true
	ctrl.UpdateByHostname(recorder, utils.FakeRequest("PUT", "http://foo.bar/resources/1", resourceIn))
	r.Equal(500, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.Internal, render.APIError)
	utils.Clear(params, render, recorder)

	// Resource not found
	inter.errDB = false
	inter.errNotFound = true
	ctrl.UpdateByHostname(recorder, utils.FakeRequest("PUT", "http://foo.bar/resources/1", resourceIn))
	r.Equal(404, render.Status)
	r.NotEmpty(recorder.Body.Bytes())
	r.NotNil(render.APIError)
	a.IsType(errs.API.NotFound, render.APIError)
	utils.Clear(params, render, recorder)
}
