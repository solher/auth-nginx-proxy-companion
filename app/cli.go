package app

import (
	"time"

	"github.com/codegangsta/cli"
	"github.com/solher/zest"
)

func SetCli(appli *zest.Zest) {
	c := appli.Cli()

	c.Usage = "Auth Server"
	c.Version = "0.0.3"

	c.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port,p",
			Value:  3000,
			Usage:  "listening port",
			EnvVar: "PORT",
		},
		cli.DurationFlag{
			Name:   "exitTimeout,t",
			Value:  10 * time.Second,
			Usage:  "graceful shutdown timeout (0 for infinite)",
			EnvVar: "EXIT_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "config,c",
			Usage:  "json or yaml config file location (overrides the database)",
			EnvVar: "CONFIG",
		},
		cli.StringFlag{
			Name:   "swaggerLocation",
			Value:  "./swagger.json",
			Usage:  "swager file location (Default: swagger.json)",
			EnvVar: "SWAGGER_LOCATION",
		},
		cli.StringFlag{
			Name:   "dbLocation",
			Value:  "data.db",
			Usage:  "database location (Default: data.db)",
			EnvVar: "DB_LOCATION",
		},
		cli.DurationFlag{
			Name:   "dbTimeout",
			Value:  time.Second,
			Usage:  "Bolt connection timeout",
			EnvVar: "DB_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "gcLocation",
			Value:  "archived.db",
			Usage:  "garbage collector database location (Default: archived.db)",
			EnvVar: "GC_LOCATION",
		},
		cli.DurationFlag{
			Name:   "gcFreq",
			Value:  time.Hour,
			Usage:  "garbage collection frequency",
			EnvVar: "GC_FREQ",
		},
		cli.DurationFlag{
			Name:   "sessionValidity",
			Value:  24 * time.Hour,
			Usage:  "the default duration of a created session",
			EnvVar: "SESSION_VALIDITY",
		},
		cli.IntFlag{
			Name:   "sessionTokenLength",
			Value:  64,
			Usage:  "the default length of generated auth tokens",
			EnvVar: "SESSION_TOKEN_LENGTH",
		},
		cli.StringFlag{
			Name:   "redirectUrl",
			Value:  "http://www.google.com",
			Usage:  "the default redirection URL when access is denied",
			EnvVar: "REDIRECT_URL",
		},
		cli.BoolFlag{
			Name:   "grantAll",
			Usage:  "disables the auth server when set to true",
			EnvVar: "GRANT_ALL",
		},
	}

	appli.SetCli(c)
}
