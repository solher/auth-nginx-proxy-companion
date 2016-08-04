// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/solher/auth-nginx-proxy-companion/app"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/solher/auth-nginx-proxy-companion/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionFind runs integration tests on the Session session Find methods.
func TestSessionFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/sessions"

	client := &http.Client{}
	sessionsOut := []models.Session{}
	sessionOut := &models.Session{}

	// Find succeeds
	res, err := client.Do(utils.FakeRequest("GET", testURL, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&sessionsOut)
	r.NoError(err)
	a.Len(sessionsOut, 4)

	// FindbyToken succeeds
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/F00bAr", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(sessionOut)
	r.NoError(err)
	a.Equal("F00bAr", *sessionOut.Token)

	// FindbyToken fails: invalid token
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/doesnt.exist", nil))
	r.NoError(err)
	r.Equal(404, res.StatusCode)
}

// TestSessionCreate runs integration tests on the Session session Create methods.
func TestSessionCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/sessions"

	client := &http.Client{}
	sessionIn := &models.Session{}
	sessionOut := &models.Session{}

	// Validation fails: nil policies
	res, err := client.Do(utils.FakeRequest("POST", testURL, sessionIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	sessionIn.Policies = []string{"1000"}

	// Validation fails: policy does not exists
	res, err = client.Do(utils.FakeRequest("POST", testURL, sessionIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	sessionIn.Policies = []string{"Foo"}
	sessionIn.Token = utils.StrCpy("F00bAr")

	// Validation fails: token must be unique
	res, err = client.Do(utils.FakeRequest("POST", testURL, sessionIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	sessionIn.Token = utils.StrCpy("unique")

	// Creation succeeds
	res, err = client.Do(utils.FakeRequest("POST", testURL, sessionIn))
	r.NoError(err)
	r.Equal(201, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(sessionOut)
	r.NoError(err)
	a.NotNil(sessionOut.Token)

	// Creation can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/"+*sessionOut.Token, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(sessionOut)
	r.NoError(err)
	a.NotNil(sessionOut.Token)
}

// TestSessionDelete runs integration tests on the Session session Delete methods.
func TestSessionDelete(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/sessions"

	client := &http.Client{}
	sessionOut := &models.Session{}
	sessionsOut := []models.Session{}

	// Deletion by ownerToken succeeds
	res, err := client.Do(utils.FakeRequest("DELETE", testURL+`?ownerTokens=["owner1","owner6"]`, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&sessionsOut)
	r.NoError(err)
	a.Len(sessionsOut, 2)

	// Deletion by ownerToken fails: invalid token
	res, err = client.Do(utils.FakeRequest("DELETE", testURL+`?ownerTokens=owner1`, nil))
	r.NoError(err)
	r.Equal(400, res.StatusCode)

	// Deletions can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&sessionsOut)
	r.NoError(err)
	a.Len(sessionsOut, 2)

	// Deletion by token succeeds
	res, err = client.Do(utils.FakeRequest("DELETE", testURL+"/F00bAr4", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(sessionOut)
	r.NoError(err)
	a.NotNil(sessionOut.Token)

	// Deletion can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/F00bAr4", nil))
	r.NoError(err)
	r.Equal(404, res.StatusCode)
}
