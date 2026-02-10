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
	ErrResponseInvalidUserOrPass = errors.New("invalid username or password")
	ErrResponseInvalidSourceIp   = errors.New("invalid ip address")
	ErrResponseInvalidData       = errors.New("invalid data")

	ErrMismatchVerificationRefnum = errors.New("mismatch verification refnum")
	ErrMismatchVerificationAmount = errors.New("mismatch verification amount")
)

func GetResponseError(response string) error {
	switch response {
	case ResponseInvalidUserOrPass:
		return ErrResponseInvalidUserOrPass
	case ResponseInvalidSourceIp:
		return ErrResponseInvalidSourceIp
	case ResponseInvalidData:
		return ErrResponseInvalidData
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
