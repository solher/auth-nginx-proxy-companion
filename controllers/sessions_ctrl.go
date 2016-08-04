package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewSessionsCtrl)
}

type (
	SessionsCtrlSessionsInter interface {
		Find() ([]models.Session, error)
		FindByToken(token string) (*models.Session, error)
		Create(session *models.Session) (*models.Session, error)
		DeleteByToken(token string) (*models.Session, error)
		DeleteByOwnerTokens(ownerToken []string) ([]models.Session, error)
	}

	SessionsCtrlSessionsValidator interface {
		ValidateCreation(session *models.Session) error
	}

	SessionsCtrl struct {
		i  SessionsCtrlSessionsInter
		v  SessionsCtrlSessionsValidator
		r  JSONRenderer // Interface used to mock the JSON renderer
		pg ParamsGetter // Interface used to mock request params
	}
)

func NewSessionsCtrl(
	i SessionsCtrlSessionsInter,
	r JSONRenderer, pg ParamsGetter,
	v SessionsCtrlSessionsValidator,
) *SessionsCtrl {
	return &SessionsCtrl{i: i, r: r, pg: pg, v: v}
}

// Find swagger:route GET /sessions Sessions SessionsFind
//
// Find
//
// Finds all the sessions from the data source.
//
// Responses:
//  200: SessionsResponse
//  500: InternalResponse
func (c *SessionsCtrl) Find(w http.ResponseWriter, r *http.Request) {
	sessions, err := c.i.Find()
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusOK, sessions)
}

// FindByToken swagger:route GET /sessions/{token} Sessions SessionsFindByToken
//
// Find by token
//
// Finds a session by token from the data source.
//
// Responses:
//  200: SessionResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *SessionsCtrl) FindByToken(w http.ResponseWriter, r *http.Request) {
	session, err := c.i.FindByToken(c.pg.GetURLParam(r, "token"))
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, session)
}

// Create swagger:route POST /sessions Sessions SessionsCreate
//
// Create
//
// Creates a session in the data source.
//
// Responses:
//  201: SessionResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *SessionsCtrl) Create(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{}

	if err := json.NewDecoder(r.Body).Decode(session); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.API.BodyDecoding, err)
		return
	}

	if err := c.v.ValidateCreation(session); err != nil {
		c.r.JSONError(w, 422, errs.API.Validation, err)
		return
	}

	session, err := c.i.Create(session)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusCreated, session)
}

// DeleteByToken swagger:route DELETE /sessions/{token} Sessions SessionsDeleteByToken
//
// Delete by token
//
// Deletes a session by token from the data source.
//
// Responses:
//  200: SessionResponse
//  404: NotFoundResponse
//  500: InternalResponse
func (c *SessionsCtrl) DeleteByToken(w http.ResponseWriter, r *http.Request) {
	session, err := c.i.DeleteByToken(c.pg.GetURLParam(r, "token"))
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			c.r.JSONError(w, http.StatusNotFound, errs.API.NotFound, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, session)
}

// DeleteByOwnerToken swagger:route DELETE /sessions Sessions SessionsDeleteByOwnerToken
//
// Delete by owner token
//
// Deletes a session by owner token from the data source.
//
// Responses:
//  200: SessionsResponse
//  500: InternalResponse
func (c *SessionsCtrl) DeleteByOwnerToken(w http.ResponseWriter, r *http.Request) {
	ownerTokens := []string{}

	if err := json.Unmarshal([]byte(r.URL.Query().Get("ownerTokens")), &ownerTokens); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.API.BodyDecoding, err)
		return
	}

	sessions, err := c.i.DeleteByOwnerTokens(ownerTokens)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	c.r.JSON(w, http.StatusOK, sessions)
}
