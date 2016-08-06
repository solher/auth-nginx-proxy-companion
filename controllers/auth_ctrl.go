package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"

	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewAuthCtrl)
}

type (
	AuthCtrlAuthInter interface {
		AuthorizeToken(hostname, path, token string) (bool, *models.Session, error)
		GetRedirectURL(hostname string) (string, error)
	}

	AuthOptionsGetter interface {
		GetRedirectURL() string
		GetGrantAll() bool
	}

	AuthCtrl struct {
		i AuthCtrlAuthInter
		r JSONRenderer // Interface used to mock the JSON renderer
		g AuthOptionsGetter
	}
)

func NewAuthCtrl(i AuthCtrlAuthInter, r JSONRenderer, g AuthOptionsGetter) *AuthCtrl {
	return &AuthCtrl{i: i, r: r, g: g}
}

// AuthorizeToken swagger:route GET /auth Auth AuthAuthorizeToken
//
// Authorize token
//
// Authenticates and authorizes a given token.
// In the case of a granted access, the session payload is set in the response header 'Auth-Server-Payload'.
//
// Responses:
//  204: nil
//	401: UnauthorizedResponse
//  500: InternalResponse
func (c *AuthCtrl) AuthorizeToken(w http.ResponseWriter, r *http.Request) {
	token := c.accessToken(r)
	requestURL := c.requestURL(r)

	u, err := url.ParseRequestURI(requestURL)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	authorized, session, err := c.i.AuthorizeToken(u.Host, u.Path, token)
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			// continue
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
			return
		}
	}

	if !authorized && !c.g.GetGrantAll() {
		c.r.JSONError(w, http.StatusForbidden, errs.API.Unauthorized, errors.New("session not found, expired or unauthorized access"))
		return
	}

	if session != nil && session.Payload != nil {
		payload := base64.StdEncoding.EncodeToString([]byte(*session.Payload))
		w.Header().Add("Auth-Server-Payload", payload)
	}

	if session != nil {
		session.Policies = nil
		session.Payload = nil

		s, _ := json.Marshal(session)
		payload := base64.StdEncoding.EncodeToString(s)
		w.Header().Add("Auth-Server-Session", payload)
	}

	if token != "" {
		w.Header().Add("Auth-Server-Token", token)
	}

	c.r.JSON(w, http.StatusNoContent, nil)
}

// Redirect swagger:route GET /redirect Auth AuthRedirect
//
// Redirect
//
// Redirects a requests to the URL set in the default configuration or in the corresponding resource.
//
// Responses:
//  307: nil
//  500: InternalResponse
func (c *AuthCtrl) Redirect(w http.ResponseWriter, r *http.Request) {
	requestURL := c.requestURL(r)

	u, err := url.ParseRequestURI(requestURL)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
		return
	}

	redirectURL, err := c.i.GetRedirectURL(u.Host)
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			// continue
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.API.Internal, err)
			return
		}
	}

	if redirectURL == "" {
		redirectURL = c.g.GetRedirectURL()
	}

	w.Header().Add("Location", redirectURL+"?redirectUrl="+requestURL)
	w.Header().Add("Redirect-Url", requestURL)

	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c *AuthCtrl) accessToken(r *http.Request) string {
	token := ""

	if cookie, err := r.Cookie("access_token"); err == nil {
		token = cookie.Value
	}

	if header := r.Header.Get("Auth-Server-Token"); header != "" {
		token = header
	}

	if t := r.URL.Query().Get("accessToken"); t != "" {
		token = t
	}

	return token
}

func (c *AuthCtrl) requestURL(r *http.Request) string {
	requestURL := r.Header.Get("Request-Url")

	if u := r.URL.Query().Get("requestUrl"); u != "" {
		requestURL = u
	}

	return requestURL
}

// swagger:parameters Auth AuthAuthorizeToken
type tokenParam struct {
	// Access token (can also be set via the 'Auth-Server-Token' header. Ex: 'Auth-Server-Token: jhPd6Gf3jIP2h')
	//
	// in: query
	AccessToken string `json:"accessToken"`
}

// swagger:parameters Auth AuthAuthorizeToken AuthRedirect
type requestURLParam struct {
	// The URL requested for access (can also be set via the 'Request-Url' header. Ex: 'Request-Url: http://foo.com/bar')
	//
	// in: query
	RequestURL string `json:"requestUrl"`
}
