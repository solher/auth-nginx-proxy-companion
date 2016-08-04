package models

type Resource struct {
	// The resource name. Must be unique.
	// required: true
	Name *string `json:"name,omitempty" yaml:"name"`
	// The resource host name. Ex: 'resource.example.com'
	// required: true
	Hostname *string `json:"hostname,omitempty" yaml:"hostname"`
	// Disable the authentication for that resource.
	Public *bool `json:"public,omitempty" yaml:"public"`
	// The redirection URL when access is denied to the resource.
	RedirectURL *string `json:"redirectUrl,omitempty" yaml:"redirectUrl"`
}

// swagger:response ResourcesResponse
type resourcesResponse struct {
	// in: body
	Body []Resource
}

// swagger:response ResourceResponse
type resourceResponse struct {
	// in: body
	Body Resource
}

// swagger:parameters ResourcesFindByHostname ResourcesDeleteByHostname ResourcesUpdateByHostname
type resourcesHostnameParam struct {
	// Resource hostname
	//
	// required: true
	// in: path
	Hostname string
}

// swagger:parameters ResourcesCreate ResourcesUpdateByHostname
type resourcesBodyParam struct {
	// required: true
	// in: body
	Body Resource
}
