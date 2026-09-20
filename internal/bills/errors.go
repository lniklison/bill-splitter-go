package bills

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var ErrBillNotFound = &Error{
	Code:    "bill_not_found",
	Message: "Bill not found",
}

func validationError(code, message string) error {
	return &Error{Code: code, Message: message}
}
