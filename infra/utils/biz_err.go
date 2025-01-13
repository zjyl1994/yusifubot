package utils

import "errors"

type BizErr interface {
	GetBizMsg() string
}

type bizErr struct {
	message string
	base    error
}

func NewBizErr(msg string, base error) bizErr {
	return bizErr{
		message: msg,
		base:    base,
	}
}

func (e bizErr) Error() string {
	return e.message
}

func (e bizErr) Is(err error) bool {
	return errors.Is(e.base, err)
}

func (e bizErr) GetBizMsg() string {
	return e.message
}
