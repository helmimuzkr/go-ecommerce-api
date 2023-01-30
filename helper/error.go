package helper

type CustomErr struct {
	Code int
	Msg  string
}

func NewCustErr(c int, m string) error {
	return CustomErr{
		Code: c,
		Msg:  m,
	}
}

func (ce CustomErr) Error() string {
	return ce.Msg
}

func ErrResponse(err error) (int, map[string]string) {
	errcust := err.(CustomErr)
	return errcust.Code, map[string]string{"message": errcust.Msg}
}
