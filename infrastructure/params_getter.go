package infrastructure

import (
	"net/http"

	"github.com/go-zoo/bone"
)

type ClassicParamsGetter struct{}

func NewParamsGetter() *ClassicParamsGetter {
	return &ClassicParamsGetter{}
}

func (g *ClassicParamsGetter) GetURLParam(req *http.Request, key string) string {
	return bone.GetValue(req, key)
}
