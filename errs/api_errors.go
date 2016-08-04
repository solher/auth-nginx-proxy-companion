package errs

import "github.com/solher/zest"

var API *apiErrors

type apiErrors struct {
	Internal     *zest.APIError
	NotFound     *zest.APIError
	InvalidID    *zest.APIError
	Unauthorized *zest.APIError
	BodyDecoding *zest.APIError
	Validation   *zest.APIError
}

func init() {
	API = &apiErrors{
		Internal:     &zest.APIError{Description: "An internal error occured. Please retry later.", ErrorCode: "INTERNAL_ERROR"},
		NotFound:     &zest.APIError{Description: "The specified resource was not found.", ErrorCode: "NOT_FOUND"},
		InvalidID:    &zest.APIError{Description: "The specified ID is invalid.", ErrorCode: "INVALID_ID"},
		Unauthorized: &zest.APIError{Description: "Authorization Required.", ErrorCode: "AUTHORIZATION_REQUIRED"},
		BodyDecoding: &zest.APIError{Description: "Could not decode the JSON request.", ErrorCode: "BODY_DECODING_ERROR"},
		Validation:   &zest.APIError{Description: "The model validation failed.", ErrorCode: "VALIDATION_ERROR"},
	}
}

// An internal error occured. Please retry later.
// swagger:response InternalResponse
type internalResponse struct {
	// in: body
	Body zest.APIError
}

// The specified resource was not found.
// swagger:response NotFoundResponse
type notFoundResponse struct {
	// in: body
	Body zest.APIError
}

// The specified resource was not found or you do not have sufficient permissions.
// swagger:response UnauthorizedResponse
type unauthorizedResponse struct {
	// in: body
	Body zest.APIError
}

// The specified ID is invalid.
// swagger:response InvalidIDResponse
type invalidIDResponse struct {
	// in: body
	Body zest.APIError
}

// Could not decode the JSON request.
// swagger:response BodyDecodingResponse
type bodyDecodingResponse struct {
	// in: body
	Body zest.APIError
}

// The model validation failed.
// swagger:response ValidationResponse
type validationResponse struct {
	// in: body
	Body zest.APIError
}
