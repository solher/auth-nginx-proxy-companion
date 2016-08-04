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

// TestPolicyFind runs integration tests on the Policy resource Find methods.
func TestPolicyFind(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/policies"

	client := &http.Client{}
	policiesOut := []models.Policy{}
	policyOut := &models.Policy{}

	// Find succeeds
	res, err := client.Do(utils.FakeRequest("GET", testURL, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&policiesOut)
	r.NoError(err)
	a.Len(policiesOut, 3)

	// FindbyName succeeds
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/Foo", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	a.Equal("Foo", *policyOut.Name)

	// FindbyName fails: invalid name
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/qwerty", nil))
	r.NoError(err)
	r.Equal(404, res.StatusCode)
}

// TestPolicyCreate runs integration tests on the Policy resource Create methods.
func TestPolicyCreate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/policies"

	client := &http.Client{}
	policyIn := &models.Policy{}
	policyOut := &models.Policy{}

	// Validation fails: blank name
	res, err := client.Do(utils.FakeRequest("POST", testURL, policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Name = utils.StrCpy("Foobar")

	// Validation fails: blank permissions
	res, err = client.Do(utils.FakeRequest("POST", testURL, policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Permissions = []models.Permission{{}}

	// Validation fails: blank resource
	res, err = client.Do(utils.FakeRequest("POST", testURL, policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Permissions = []models.Permission{{Resource: utils.StrCpy("qwerty")}}

	// Validation fails: resource doesn't exist
	res, err = client.Do(utils.FakeRequest("POST", testURL, policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Permissions = []models.Permission{{Resource: utils.StrCpy("Foobar")}}

	// Creation succeeds
	res, err = client.Do(utils.FakeRequest("POST", testURL, policyIn))
	r.NoError(err)
	r.Equal(201, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut.Name)

	// Validation fails: name must be unique
	res, err = client.Do(utils.FakeRequest("POST", testURL, policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	// Creation can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/"+*policyOut.Name, nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut.Name)
}

// TestPolicyDelete runs integration tests on the Policy resource Delete methods.
func TestPolicyDelete(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/policies"

	client := &http.Client{}
	policyOut := &models.Policy{}

	// Validation error: can't delete the guest policy
	res, err := client.Do(utils.FakeRequest("DELETE", testURL+"/guest", nil))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	// Deletion succeeds
	res, err = client.Do(utils.FakeRequest("DELETE", testURL+"/Foo", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	a.NotNil(policyOut.Name)

	// Deletion can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/Foo", nil))
	r.NoError(err)
	r.Equal(404, res.StatusCode)

	sessionOut := &models.Session{}

	// Cascade deletion can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", url+"/sessions/F00bAr", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(sessionOut)
	r.NoError(err)
	r.NotContains(sessionOut.Policies, "Foo")
}

// TestPolicyUpdate runs integration tests on the Policy resource Update methods.
func TestPolicyUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	appli := app.NewTestApp()
	url, err := appli.Launch()
	r.NoError(err)
	defer appli.Stop()

	testURL := url + "/policies"

	client := &http.Client{}
	policyIn := &models.Policy{Name: utils.StrCpy("qwerty")}
	policyOut := &models.Policy{}

	// Validation fails: everything nil
	res, err := client.Do(utils.FakeRequest("PUT", testURL+"/Foo", policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Permissions = []models.Permission{{}}

	// Validation fails: blank resource
	res, err = client.Do(utils.FakeRequest("PUT", testURL+"/Foo", policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Permissions = []models.Permission{{Resource: utils.StrCpy("qwerty")}}

	// Validation fails: resource doesn't exist
	res, err = client.Do(utils.FakeRequest("PUT", testURL+"/Foo", policyIn))
	r.NoError(err)
	r.Equal(422, res.StatusCode)

	policyIn.Permissions = []models.Permission{{Resource: utils.StrCpy("Foobar2")}}

	// Update succeeds
	res, err = client.Do(utils.FakeRequest("PUT", testURL+"/Foo", policyIn))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	a.NotEqual("qwerty", *policyOut.Name)
	a.Equal("Foobar2", *policyOut.Permissions[0].Resource)

	// Update can be confirmed
	res, err = client.Do(utils.FakeRequest("GET", testURL+"/Foo", nil))
	r.NoError(err)
	r.Equal(200, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(policyOut)
	r.NoError(err)
	a.Equal("Foobar2", *policyOut.Permissions[0].Resource)

	// Invalid name
	res, err = client.Do(utils.FakeRequest("PUT", testURL+"/qwerty", policyIn))
	r.NoError(err)
	r.Equal(404, res.StatusCode)
}
