//go:generate ./generate-doc.sh

// Auth Server
//
// A cool authentication server.
//
// Schemes: http, https
// BasePath: /
// Version: 0.0.3
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main // import "github.com/solher/auth-nginx-proxy-companion"

import "github.com/solher/auth-nginx-proxy-companion/app"

func main() {
	app.Run(nil)
}
