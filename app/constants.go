package app

import "time"

type Constants struct {
	Swagger struct {
		Location string
	}

	App struct {
		Port        int
		ExitTimeout time.Duration
		Config      string
	}

	Auth struct {
		RedirectURL string
		GrantAll    bool
	}

	GC struct {
		Location string
		Freq     time.Duration
	}

	DB struct {
		Location string
		Timeout  time.Duration
	}

	Session struct {
		Validity    time.Duration
		TokenLength int
	}
}

func NewConstants() *Constants {
	return &Constants{}
}

func (c *Constants) GetRedirectURL() string {
	return c.Auth.RedirectURL
}

func (c *Constants) GetGrantAll() bool {
	return c.Auth.GrantAll
}

func (c *Constants) GetSessionValidity() time.Duration {
	return c.Session.Validity
}

func (c *Constants) GetSessionTokenLength() int {
	return c.Session.TokenLength
}
