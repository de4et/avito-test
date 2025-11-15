package api

const (
	ErrTeamExistsMsg = "team already exists"
)

func NewError(code ErrorResponseErrorCode, msg string) ErrorResponse {
	var e ErrorResponse
	e.Error.Code = code
	e.Error.Message = msg
	return e
}

