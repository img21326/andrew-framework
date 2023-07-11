package helper

type ErrorInterface interface {
	Code() int
	HttpCode() int
	Message() string
}

var ErrorMap = map[int]ErrorInterface{}

func AddError(err ErrorInterface) {
	ErrorMap[err.Code()] = err
}
