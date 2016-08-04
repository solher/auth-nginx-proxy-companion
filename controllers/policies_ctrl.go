package controllers

import (
	"encoding/json"
	"net/http"

	"git.wid.la/co-net/auth-server/errs"
	"git.wid.la/co-net/auth-server/models"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewPoliciesCtrl)
}

type (
	PoliciesCtrlPoliciesInter interface {
		Find() ([]models.Policy, error)
		FindByName(id string) (*models.Policy, error)
		Create(policy *models.Policy) (*models.Policy, error)
		DeleteByName(id string) (*models.Policy, error)
		UpdateByName(id string, policy *models.Policy) (*models.Policy, error)
	}

	PoliciesCtrlPoliciesValidator interface {
		ValidateCreation(policy *models.Policy) error
		ValidateUpdate(policy *models.Policy) error
	}

	PoliciesCtrl struct {
		i  PoliciesCtrlPoliciesInter
		v  PoliciesCtrlPoliciesValidator
		r  JSONRenderer // Interface used to mock the JSON renderer
		pg ParamsGetter // Interface used to mock request params
	}
)

func NewPoliciesCtrl(
	i PoliciesCtrlPoliciesInter,
	r JSONRenderer, pg ParamsGetter,
	v PoliciesCtrlPoliciesValidator,
) *PoliciesCtrl {
	return &PoliciesCtrl{i: i, r: r, pg: pg, v: v}
}

// Find swagger:route GET /policies Policies PoliciesFind
//
// Find
//
// Finds all the policies from the data source.
//
// Responses:
//  200: PoliciesResponse
//  500: InternalResponse
func (c *PoliciesCtrl) Find(w http.ResponseWriter, r *http.Request) {
	policies, err := c.i.Find()
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusOK, policies)
}

// FindByName swagger:route GET /policies/{name} Policies PoliciesFindByName
//
// Find by name
//
// Finds a policy by name from the data source.
//
// Responses:
//  200: PolicyResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *PoliciesCtrl) FindByName(w http.ResponseWriter, r *http.Request) {
	policy, err := c.i.FindByName(c.pg.GetURLParam(r, "name"))
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, policy)
}

// Create swagger:route POST /policies Policies PoliciesCreate
//
// Create
//
// Creates a policy in the data source.
//
// Responses:
//  201: PolicyResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *PoliciesCtrl) Create(w http.ResponseWriter, r *http.Request) {
	policy := &models.Policy{}

	if err := json.NewDecoder(r.Body).Decode(policy); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.API.BodyDecoding, err)
		return
	}

	if err := c.v.ValidateCreation(policy); err != nil {
		c.r.JSONError(w, 422, errs.API.Validation, err)
		return
	}

	policy, err := c.i.Create(policy)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusCreated, policy)
}

// DeleteByName swagger:route DELETE /policies/{name} Policies PoliciesDeleteByName
//
// Delete by name
//
// Deletes a policy by name from the data source.
//
// Responses:
//  200: PolicyResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *PoliciesCtrl) DeleteByName(w http.ResponseWriter, r *http.Request) {
	policy, err := c.i.DeleteByName(c.pg.GetURLParam(r, "name"))
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		case errs.ErrValidation:
			c.r.JSONError(w, 422, errs.API.Validation, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, policy)
}

// UpdateByName swagger:route PUT /policies/{name} Policies PoliciesUpdateByName
//
// Update by name
//
// Updates a policy by name from the data source.
//
// Responses:
//  200: PolicyResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *PoliciesCtrl) UpdateByName(w http.ResponseWriter, r *http.Request) {
	policy := &models.Policy{}

	if err := json.NewDecoder(r.Body).Decode(policy); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.API.BodyDecoding, err)
		return
	}

	if err := c.v.ValidateUpdate(policy); err != nil {
		c.r.JSONError(w, 422, errs.API.Validation, err)
		return
	}

	policy.Name = nil

	policy, err := c.i.UpdateByName(c.pg.GetURLParam(r, "name"), policy)
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, policy)
}
