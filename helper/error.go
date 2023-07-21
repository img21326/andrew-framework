package helper

type ErrorInterface interface {
	Code() int
	HttpCode() int
	Message() string
}

type BaseError struct {
	code     int
	message  string
	httpCode int
}

func (e *BaseError) ErrorMessage(msg string) {
	e.message = msg
}

func (e *BaseError) Error() string {
	return e.message
}

func (e *BaseError) Code() int {
	return e.code
}

func (e *BaseError) HttpCode() int {
	return e.httpCode
}

var ErrorMap = map[int]ErrorInterface{}

func AddError(err ErrorInterface) {
	if _, ok := ErrorMap[err.Code()]; ok {
		panic("Error code already exists")
	}
	ErrorMap[err.Code()] = err
}
