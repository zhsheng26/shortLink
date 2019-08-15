package main

type MiError interface {
	error
	Status() int
}

type ReqError struct {
	Code int
	Err  error
}

func (re ReqError) Error() string {
	return re.Err.Error()
}

func (re ReqError) Status() int {
	return re.Code
}
