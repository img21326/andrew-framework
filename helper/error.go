package helper

type ErrorInterface interface {
	ErrorCode() int
	HttpStatus() int
	Message() string
}

type BaseError struct {
	Code     int
	Message  string
	HttpCode int
}

func (e *BaseError) ErrorMessage(msg string) {
	e.Message = msg
}

func (e *BaseError) Error() string {
	return e.Message
}

func (e *BaseError) ErrorCode() int {
	return e.Code
}

func (e *BaseError) HttpStatus() int {
	return e.HttpCode
}

var ErrorMap = map[int]ErrorInterface{}

func AddError(err ErrorInterface) {
	if _, ok := ErrorMap[err.ErrorCode()]; ok {
		panic("Error code already exists")
	}
	ErrorMap[err.ErrorCode()] = err
}
