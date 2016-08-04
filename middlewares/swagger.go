package middlewares

import "net/http"

type Swagger struct {
	file string
}

func NewSwagger(file string) *Swagger {
	return &Swagger{file: file}
}

func (s *Swagger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == "/swagger" {
		http.ServeFile(rw, r, s.file)
		return
	}

	next(rw, r)
}
