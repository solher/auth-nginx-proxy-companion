package controllers

import (
	"encoding/json"
	"net/http"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewResourcesCtrl)
}

type (
	ResourcesCtrlResourcesInter interface {
		Find() ([]models.Resource, error)
		FindByHostname(hostname string) (*models.Resource, error)
		Create(resource *models.Resource) (*models.Resource, error)
		DeleteByHostname(hostname string) (*models.Resource, error)
		UpdateByHostname(hostname string, resource *models.Resource) (*models.Resource, error)
	}

	ResourcesCtrlResourcesValidator interface {
		ValidateCreation(resource *models.Resource) error
		ValidateUpdate(resource *models.Resource) error
	}

	ResourcesCtrl struct {
		i  ResourcesCtrlResourcesInter
		v  ResourcesCtrlResourcesValidator
		r  JSONRenderer // Interface used to mock the JSON renderer
		pg ParamsGetter // Interface used to mock request params
	}
)

func NewResourcesCtrl(
	i ResourcesCtrlResourcesInter,
	r JSONRenderer, pg ParamsGetter,
	v ResourcesCtrlResourcesValidator,
) *ResourcesCtrl {
	return &ResourcesCtrl{i: i, r: r, pg: pg, v: v}
}

// Find swagger:route GET /resources Resources ResourcesFind
//
// Find
//
// Finds all the resources from the data source.
//
// Responses:
//  200: ResourcesResponse
//  500: InternalResponse
func (c *ResourcesCtrl) Find(w http.ResponseWriter, r *http.Request) {
	resources, err := c.i.Find()
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusOK, resources)
}

// FindByHostname swagger:route GET /resources/{hostname} Resources ResourcesFindByHostname
//
// Find by hostname
//
// Finds a resource by hostname from the data source.
//
// Responses:
//  200: ResourceResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *ResourcesCtrl) FindByHostname(w http.ResponseWriter, r *http.Request) {
	resource, err := c.i.FindByHostname(c.pg.GetURLParam(r, "hostname"))
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, resource)
}

// Create swagger:route POST /resources Resources ResourcesCreate
//
// Create
//
// Creates a resource in the data source.
//
// Responses:
//  201: ResourceResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *ResourcesCtrl) Create(w http.ResponseWriter, r *http.Request) {
	resource := &models.Resource{}

	if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.API.BodyDecoding, err)
		return
	}

	if err := c.v.ValidateCreation(resource); err != nil {
		c.r.JSONError(w, 422, errs.API.Validation, err)
		return
	}

	resource, err := c.i.Create(resource)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusCreated, resource)
}

// DeleteByHostname swagger:route DELETE /resources/{hostname} Resources ResourcesDeleteByHostname
//
// Delete by hostname
//
// Deletes a resource by hostname from the data source.
//
// Responses:
//  200: ResourceResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *ResourcesCtrl) DeleteByHostname(w http.ResponseWriter, r *http.Request) {
	resource, err := c.i.DeleteByHostname(c.pg.GetURLParam(r, "hostname"))
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, resource)
}

// UpdateByHostname swagger:route PUT /resources/{hostname} Resources ResourcesUpdateByHostname
//
// Update by hostname
//
// Updates a resource by hostname from the data source.
//
// Responses:
//  200: ResourceResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *ResourcesCtrl) UpdateByHostname(w http.ResponseWriter, r *http.Request) {
	resource := &models.Resource{}

	if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.API.BodyDecoding, err)
		return
	}

	if err := c.v.ValidateUpdate(resource); err != nil {
		c.r.JSONError(w, 422, errs.API.Validation, err)
		return
	}

	resource.Hostname = nil

	resource, err := c.i.UpdateByHostname(c.pg.GetURLParam(r, "hostname"), resource)
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, resource)
}
