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
package main // import "git.wid.la/co-net/auth-server"

import "git.wid.la/co-net/auth-server/app"

func main() {
	app.Run(nil)
}
