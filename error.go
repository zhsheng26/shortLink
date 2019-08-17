package main

import "net/http"

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

func NewNotFindErr(err error) ReqError {
	return ReqError{
		Code: http.StatusNotFound,
		Err:  err,
	}
}

func NewBadReqErr(err error) ReqError {
	return ReqError{
		Code: http.StatusBadRequest,
		Err:  err,
	}
}
