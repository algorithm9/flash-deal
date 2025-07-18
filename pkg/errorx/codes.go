package errorx

import (
	"fmt"
	"net/http"
)

type Code int

const (
	CodeOK Code = 0

	CodeBadRequest   Code = 400
	CodeUnauthorized Code = 401
	CodeInternal     Code = 500

	CodeLimitRequest Code = 1000
)

func NewCode(code int) Code {
	return Code(code)
}

func (c Code) Int() int {
	return int(c)
}

func (c Code) String() string {
	code := c.Int()
	str := http.StatusText(code)
	if str == "" {
		str = "Unknown Error"
	}
	return fmt.Sprintf("%d: %s", code, str)
}
