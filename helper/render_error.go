package helper

type ReturnErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ErrorMap = map[string]ReturnErr{
	"not_found": {
		Code:    404,
		Message: "record not found",
	},
}
