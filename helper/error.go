package helper

type ErrorInterface interface {
	ErrorCode() int
	HttpStatus() int
	Error() string
	ErrorData() interface{}
}

type BaseError struct {
	Code     int
	Message  string
	HttpCode int
	Data     interface{}
}

func New() ErrorInterface {
	return &BaseError{}
}

func (e *BaseError) SetErrMsg(msg string) {
	e.Message = msg
}

func (e *BaseError) SetErrData(data interface{}) {
	e.Data = data
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

func (e *BaseError) ErrorData() interface{} {
	return e.Data
}

var ErrorMap = map[int]ErrorInterface{}

func AddError(err ErrorInterface) {
	if _, ok := ErrorMap[err.ErrorCode()]; ok {
		panic("Error code already exists")
	}
	ErrorMap[err.ErrorCode()] = err
}
