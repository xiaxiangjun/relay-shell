package common

type Error struct {
	code int
	msg  string
}

func (self *Error) Code() int {
	return self.code
}

func (self *Error) Msg() string {
	return self.msg
}

func ErrInternalServerError(msg string) *Error {
	return &Error{code: 500, msg: msg}
}

// 定义错误码, 采用http标准错误码
var ErrOK = &Error{code: 200, msg: "OK"}
var ErrForbidden = &Error{code: 403, msg: "Forbidden"}
var ErrNotFound = &Error{code: 404, msg: "Not Found"}
var ErrMethodNotAlowed = &Error{code: 405, msg: "Method Not Allowed"}
var ErrNotAcceptable = &Error{code: 406, msg: "Not Acceptable"}
