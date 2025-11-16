package api

const (
	ErrTeamExistsMsg    = "team already exists"
	ErrTeamNotExistsMsg = "team not exists"
	ErrUserNotFoundMsg  = "user not exists"
	ErrPRExistsMsg      = "PR id already exists"
	ErrPRNotExistsMsg   = "PR not exists"
	ErrNotAssignedMsg   = "reviewer is not assigned to this PR"
	ErrNoCandidateMsg   = "no active replacement candidate in team"
	ErrMergedMsg        = "cannot reassign on merged PR "
)

func NewError(code ErrorResponseErrorCode, msg string) ErrorResponse {
	var e ErrorResponse
	e.Error.Code = code
	e.Error.Message = msg
	return e
}
