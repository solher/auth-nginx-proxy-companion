package models

type (
	Policy struct {
		// The policy name.
		// required: true
		Name *string `json:"name,omitempty" yaml:"name"`
		// Can be used to disable a policy.
		Enabled *bool `json:"enabled,omitempty" yaml:"enabled"`
		// An array of resource IDs and their associated right.
		// required: true
		Permissions []Permission `json:"permissions,omitempty" yaml:"permissions"`
	}

	Permission struct {
		// The resource ID concerned by the permission.
		// required: true
		Resource *string `json:"resource,omitempty" yaml:"resource"`
		// The optional paths on which the permission apply.
		Paths []string `json:"paths,omitempty" yaml:"paths"`
		// Can be used to disable a permission.
		Enabled *bool `json:"enabled,omitempty" yaml:"enabled"`
		// Indicates if the permission grants or denies the access on the resource.
		Deny *bool `json:"deny,omitempty" yaml:"deny"`
	}
)

// swagger:response PoliciesResponse
type policiesResponse struct {
	// in: body
	Body []Policy
}

// swagger:response PolicyResponse
type policyResponse struct {
	// in: body
	Body Policy
}

// swagger:parameters PoliciesFindByName PoliciesDeleteByName PoliciesUpdateByName
type policiesIDParam struct {
	// Policy name
	//
	// required: true
	// in: path
	Name string
}

// swagger:parameters PoliciesCreate PoliciesUpdateByName
type policiesBodyParam struct {
	// required: true
	// in: body
	Body Policy
}
