package error_handler

type NewError struct {
	Error      string
	StatusCode int
}

// The function returns a pointer to a new error object with the same error message and status code as
// the input error object.
func HandleError(e NewError) *NewError {
	return &NewError{
		Error:      e.Error,
		StatusCode: e.StatusCode,
	}
}
