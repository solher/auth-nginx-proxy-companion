package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/solher/zest"
)

type FakeModelsGetter struct {
	RedirectURL        string
	GrantAll           bool
	SessionValidity    time.Duration
	SessionTokenLength int
}

func NewFakeModelsGetter() *FakeModelsGetter {
	return &FakeModelsGetter{}
}

func (g *FakeModelsGetter) GetGrantAll() bool {
	return g.GrantAll
}

func (g *FakeModelsGetter) GetRedirectURL() string {
	return g.RedirectURL
}

func (g *FakeModelsGetter) GetSessionValidity() time.Duration {
	return g.SessionValidity
}

func (g *FakeModelsGetter) GetSessionTokenLength() int {
	return g.SessionTokenLength
}

type FakeRender struct {
	Status   int
	APIError *zest.APIError
	Err      error
	render   *zest.Render
}

func NewFakeRender() *FakeRender {
	return &FakeRender{render: zest.NewRender()}
}

func (r *FakeRender) JSONError(w http.ResponseWriter, status int, apiError *zest.APIError, err error) {
	r.Status = status
	r.APIError = apiError
	r.Err = err

	r.render.JSONError(w, status, apiError, err)
}

func (r *FakeRender) JSON(w http.ResponseWriter, status int, object interface{}) {
	r.Status = status

	r.render.JSON(w, status, object)
}

func (r *FakeRender) Clear() {
	r.Status = 0
	r.APIError = nil
	r.Err = nil
}

type FakeParamsGetter struct {
	params map[string]string
}

func NewFakeParamsGetter() *FakeParamsGetter {
	return &FakeParamsGetter{params: make(map[string]string)}
}

func (g *FakeParamsGetter) GetURLParam(req *http.Request, key string) string {
	return g.params[key]
}

func (g *FakeParamsGetter) SetURLParam(key, value string) {
	g.params[key] = value
}

func (g *FakeParamsGetter) ClearURLParams() {
	g.params = make(map[string]string)
}

func (g *FakeParamsGetter) Clear() {
	g.params = make(map[string]string)
}

func Clear(params *FakeParamsGetter, render *FakeRender, recorder *httptest.ResponseRecorder) {
	if params != nil {
		params.Clear()
	}

	if render != nil {
		render.Clear()
	}

	*recorder = *httptest.NewRecorder()
}

func FakeRequest(method, url string, body interface{}) *http.Request {
	m, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(m))
	return req
}

func FakeRequestRaw(method, url string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	return req
}
