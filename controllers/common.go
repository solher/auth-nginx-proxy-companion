package controllers

import (
	"net/http"

	"github.com/solher/zest"
)

type JSONRenderer interface {
	JSONError(w http.ResponseWriter, status int, apiError *zest.APIError, err error)
	JSON(w http.ResponseWriter, status int, object interface{})
}

type ParamsGetter interface {
	GetURLParam(r *http.Request, key string) string
}
