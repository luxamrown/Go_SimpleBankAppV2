package response

type ErrorResp struct {
	ErrorCode    string
	ErrorMessage error
}

func NewError(errorCode string, errorMessage error) *ErrorResp {
	return &ErrorResp{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}
