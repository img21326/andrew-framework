package helper

type ErrorInterface interface {
	Code() int
	HttpCode() int
	Message() string
}

var ErrorMap = map[int]ErrorInterface{}

func AddError(err ErrorInterface) {
	if _, ok := ErrorMap[err.Code()]; ok {
		panic("Error code already exists")
	}
	ErrorMap[err.Code()] = err
}
