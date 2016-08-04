// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"git.wid.la/co-net/auth-server/app"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResourceFind runs integration tests on the Resource resource Find methods.
func TestResourceFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/resources"

	client := &http.Client{}
	resourcesOut := []models.Resource{}
	resourceOut := &models.Resource{}

	// Find succeeds
	res, err := client.Do(utils.FakeRequest("GET", testURL, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&resourcesOut)
	r.NoError(err)
	a.Len(resourcesOut, 2)

	// FindbyHostname succeeds
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/foo.bar.com", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(resourceOut)
	r.NoError(err)
	a.Equal("Foobar", *resourceOut.Name)

	// FindbyHostname fails: invalid hostname
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/doesnt.exist", nil))
	r.NoError(err)
	r.Equal(404, res.StatusCode)
}

// TestResourceCreate runs integration tests on the Resource resource Create methods.
func TestResourceCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/resources"

	client := &http.Client{}
	resourceIn := &models.Resource{}
	resourceOut := &models.Resource{}

	// Validation fails: blank name
	res, err := client.Do(utils.FakeRequest("POST", testURL, resourceIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	resourceIn.Name = utils.StrCpy("Foobar")

	// Validation fails: blank hostname
	res, err = client.Do(utils.FakeRequest("POST", testURL, resourceIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	resourceIn.Hostname = utils.StrCpy("foo.bar.com")

	// Validation fails: hostname must be unique
	res, err = client.Do(utils.FakeRequest("POST", testURL, resourceIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	resourceIn.Hostname = utils.StrCpy("is.unique.com")

	// Creation succeeds
	res, err = client.Do(utils.FakeRequest("POST", testURL, resourceIn))
	r.NoError(err)
	r.Equal(201, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut.Name)

	// Creation can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/"+*resourceOut.Hostname, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut.Name)
}

// TestResourceDelete runs integration tests on the Resource resource Delete methods.
func TestResourceDelete(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/resources"

	client := &http.Client{}
	resourceOut := &models.Resource{}

	// Deletion succeeds
	res, err := client.Do(utils.FakeRequest("DELETE", testURL+"/foo.bar.com", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(resourceOut)
	r.NoError(err)
	a.NotNil(resourceOut.Name)

	// Deletion can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/foo.bar.com", nil))
	r.NoError(err)
	r.Equal(404, res.StatusCode)

	policyOut := &models.Policy{}

	// Cascade deletion can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", url+"/policies/Foo", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	r.Len(policyOut.Permissions, 1)
}

// TestResourceUpdate runs integration tests on the Resource resource Update methods.
func TestResourceUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/resources"

	client := &http.Client{}
	resourceIn := &models.Resource{}
	resourceOut := &models.Resource{}

	// Validation fails: everything nil
	res, err := client.Do(utils.FakeRequest("PUT", testURL+"/foo.bar.com", resourceIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	resourceIn.Name = utils.StrCpy("")

	// Validation fails: blank name
	res, err = client.Do(utils.FakeRequest("PUT", testURL+"/foo.bar.com", resourceIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	resourceIn.Name = utils.StrCpy("New")

	// Update succeeds
	res, err = client.Do(utils.FakeRequest("PUT", testURL+"/foo.bar.com", resourceIn))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(resourceOut)
	r.NoError(err)
	a.Equal("New", *resourceOut.Name)

	// Update can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/foo.bar.com", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(resourceOut)
	r.NoError(err)
	a.Equal("New", *resourceOut.Name)
}
