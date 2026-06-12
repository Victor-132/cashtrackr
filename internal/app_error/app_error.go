package apperror

type AppError struct {
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func New(e string) AppError {
	return AppError{e}
}
