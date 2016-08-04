package middlewares

import "net/http"

type Swagger struct {
}

func NewSwagger() *Swagger {
	return &Swagger{}
}

func (s *Swagger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == "/swagger" {
		http.ServeFile(rw, r, "./swagger.json")
		return
	}

	next(rw, r)
}
