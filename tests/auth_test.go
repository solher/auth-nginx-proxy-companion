// +build integration

package tests

import (
	"net/http"
	"testing"

	"git.wid.la/co-net/auth-server/app"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthRedirect runs integration tests on the Redirect method.
func TestAuthRedirect(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/redirect"

	req := utils.FakeRequestRaw("GET", testURL, nil)
	req.Header.Set("Request-Url", "http://foo.bar.com")

	// Success: redirected
	res, err := http.DefaultTransport.RoundTrip(req)
	r.NoError(err)
	a.Equal(307, res.StatusCode)
	a.Equal("http://www.google.com?redirectUrl=http://foo.bar.com", res.Header.Get("Location"))
	a.Equal("http://foo.bar.com", res.Header.Get("Redirect-Url"))
}

// TestAuthAuthorizeToken runs integration tests on the AuthorizeToken method.
func TestAuthAuthorizeToken(t *testing.T) {
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/auth"

	client := &http.Client{}
	req := utils.FakeRequest("GET", testURL, nil)
	req.Header.Set("Request-URL", "http://foo.bar.com/test")

	// Access granted: guest
	res, err := client.Do(req)
	r.NoError(err)
	r.Equal(204, res.StatusCode)

	req = utils.FakeRequest("GET", testURL, nil)
	req.Header.Set("Request-URL", "http://foo.bar.2.com/test")

	// Access denied: invalid token
	res, err = client.Do(req)
	r.NoError(err)
	r.Equal(403, res.StatusCode)

	req = utils.FakeRequest("GET", testURL, nil)
	req.Header.Set("Request-URL", "http://foo.bar.2.com/test")
	req.Header.Add("Auth-Server-Token", "F00bAr")

	// Access granted: valid session
	res, err = client.Do(req)
	r.NoError(err)
	r.Equal(204, res.StatusCode)

	req = utils.FakeRequest("GET", testURL, nil)
	req.Header.Set("Request-URL", "http://foo.bar.3.com/test")

	// Access denied: not existing hostname
	res, err = client.Do(req)
	r.NoError(err)
	r.Equal(403, res.StatusCode)
}
