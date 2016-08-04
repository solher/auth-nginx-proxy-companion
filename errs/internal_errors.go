package errs

var Internal *internalErrors

type internalError struct {
	Description string
}

func (e internalError) Error() string {
	return e.Description
}

type (
	ErrDatabase   struct{ internalError }
	ErrNotFound   struct{ internalError }
	ErrValidation struct{ internalError }
)

type internalErrors struct {
	Database   ErrDatabase
	NotFound   ErrNotFound
	Validation ErrValidation
}

func init() {
	Internal = &internalErrors{
		Database:   ErrDatabase{internalError{Description: "undefined database error"}},
		NotFound:   ErrNotFound{internalError{Description: "the specified resource was not found"}},
		Validation: ErrValidation{internalError{Description: "validation error"}},
	}
}

func NewErrValidation(description string) ErrValidation {
	err := ErrValidation{}
	err.Description = description
	return err
}
