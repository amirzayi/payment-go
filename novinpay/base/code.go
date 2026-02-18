package base

import (
	"errors"
	"strings"
)

type Code uint8

const (
	Success Code = iota + 1
	InvalidUserOrPass
	InvalidSourceIp
	InvalidData
)

const (
	PaidStatus      = "OK"
	CancelledStatus = "Canceled By User"
)
const (
	ResponseSuccess           = "erSucceed"
	ResponseInvalidUserOrPass = "erAAS_InvalidUseridOrPass"
	ResponseInvalidSourceIp   = "erAAS_InvalidSourceIp"
	ResponseInvalidData       = "erAAS_InvalidData"
)

var (
	ErrInvalidUserOrPass = errors.New("invalid username or password")
	ErrInvalidSourceIp   = errors.New("invalid ip address")
	ErrInvalidData       = errors.New("invalid data")

	ErrMismatchVerificationRefnum = errors.New("mismatch verification refnum")
	ErrMismatchVerificationAmount = errors.New("mismatch verification amount")

	ErrInvalidResponseStatusCode = errors.New("invalid response status code")
)

func GetResponseError(response string) error {
	switch response {
	case ResponseInvalidUserOrPass:
		return ErrInvalidUserOrPass
	case ResponseInvalidSourceIp:
		return ErrInvalidSourceIp
	case ResponseInvalidData:
		return ErrInvalidData
	default:
		return errors.New("unknown")
	}
}

func GetPayCheck(status string) string {
	if strings.EqualFold(status, PaidStatus) {
		return PaidStatus
	}
	return CancelledStatus
}
