package base

import "errors"

var (
	ErrInvalidCartNumber               = errors.New("invalid card number")
	ErrInsufficientBalance             = errors.New("insufficient balance")
	ErrIncorrectPassword               = errors.New("incorrect password")
	ErrMaxPasswordAttemptsExceeded     = errors.New("maximum password attempts exceeded")
	ErrInvalidCard                     = errors.New("invalid card")
	ErrMaxWithdrawalAttemptsExceeded   = errors.New("maximum withdrawal attempts exceeded")
	ErrUserCancelledTransaction        = errors.New("user canceled the transaction")
	ErrCardExpired                     = errors.New("card has expired")
	ErrWithdrawalAmountExceedsLimit    = errors.New("withdrawal amount exceeds limit")
	ErrUserKVCNotPresent               = errors.New("user KVC not present")
	ErrInvalidCardIssuer               = errors.New("invalid card issuer")
	ErrCardIssuerSwitchError           = errors.New("card issuer switch error")
	ErrNoResponseFromCardIssuer        = errors.New("no response from card issuer")
	ErrCardholderNotAuthorized         = errors.New("cardholder not authorized for this transaction")
	ErrInvalidMerchant                 = errors.New("invalid merchant")
	ErrSecurityError                   = errors.New("security error occurred")
	ErrInvalidMerchantCredentials      = errors.New("invalid merchant credentials")
	ErrInvalidAmount                   = errors.New("invalid amount")
	ErrInvalidResponse                 = errors.New("invalid response")
	ErrPreviousRequestInProgress       = errors.New("previous request is in progress")
	ErrInvalidInputFormat              = errors.New("invalid input format")
	ErrInvalidAccount                  = errors.New("invalid account")
	ErrSystemError                     = errors.New("system error")
	ErrInvalidDate                     = errors.New("invalid date")
	ErrDuplicateRequestNumber          = errors.New("duplicate request number")
	ErrSaleTransactionNotFound         = errors.New("sale transaction not found")
	ErrVerifyRequestAlreadySubmitted   = errors.New("verify request already submitted")
	ErrVerifyRequestNotFound           = errors.New("verify request not found")
	ErrTransactionAlreadySettled       = errors.New("transaction already settled")
	ErrTransactionNotSettled           = errors.New("transaction not settled")
	ErrSettleTransactionNotFound       = errors.New("settle transaction not found")
	ErrTransactionReversed             = errors.New("transaction reversed")
	ErrRefundTransactionNotFound       = errors.New("refund transaction not found")
	ErrInvalidBillID                   = errors.New("invalid bill ID")
	ErrInvalidPaymentID                = errors.New("invalid payment ID")
	ErrInvalidBillIssuerOrganization   = errors.New("invalid bill issuer organization")
	ErrTransactionTimeLimitExpired     = errors.New("transaction time limit expired")
	ErrDataRegistrationError           = errors.New("error in data registration")
	ErrInvalidPayerID                  = errors.New("invalid payer ID")
	ErrCustomerDataDefinitionError     = errors.New("error in customer data definition")
	ErrMaxDataEntryAttemptsExceeded    = errors.New("maximum data entry attempts exceeded")
	ErrInvalidIP                       = errors.New("invalid IP")
	ErrDuplicateTransaction            = errors.New("duplicate transaction")
	ErrReferenceTransactionNotFound    = errors.New("reference transaction not found")
	ErrInvalidTransaction              = errors.New("invalid transaction")
	ErrDepositError                    = errors.New("deposit error")
	ErrReturnURLNotInRegisteredDomain  = errors.New("return URL not in merchant's registered domain")
	ErrStaticPasswordUsageLimitReached = errors.New("static password usage limit reached")
	ErrCardOwnershipVerificationFailed = errors.New("card ownership verification failed")
	ErrInvalidResponseFormat           = errors.New("invalid response format")

	ErrInvalidResponseStatusCode = errors.New("invalid response status code")
)

func ConvertError(code string) error {
	switch code {
	case "0":
		return nil
	case "11":
		return ErrInvalidCartNumber
	case "12":
		return ErrInsufficientBalance
	case "13":
		return ErrIncorrectPassword
	case "14":
		return ErrMaxPasswordAttemptsExceeded
	case "15":
		return ErrInvalidCard
	case "16":
		return ErrMaxWithdrawalAttemptsExceeded
	case "17":
		return ErrUserCancelledTransaction
	case "18":
		return ErrCardExpired
	case "19":
		return ErrWithdrawalAmountExceedsLimit
	case "20":
		return ErrUserKVCNotPresent
	case "21":
		return ErrInvalidCardIssuer
	case "23":
		return ErrSecurityError
	case "24":
		return ErrCardIssuerSwitchError
	case "25":
		return ErrInvalidAmount
	case "30":
		return ErrPreviousRequestInProgress
	case "31":
		return ErrInvalidResponse
	case "32":
		return ErrInvalidInputFormat
	case "33":
		return ErrInvalidAccount
	case "34":
		return ErrSystemError
	case "35":
		return ErrInvalidDate
	case "41":
		return ErrDuplicateRequestNumber
	case "42":
		return ErrSaleTransactionNotFound
	case "43":
		return ErrVerifyRequestAlreadySubmitted
	case "44":
		return ErrVerifyRequestNotFound
	case "45":
		return ErrTransactionAlreadySettled
	case "46":
		return ErrTransactionNotSettled
	case "47":
		return ErrSettleTransactionNotFound
	case "48":
		return ErrTransactionReversed
	case "51":
		return ErrDuplicateTransaction
	case "54":
		return ErrReferenceTransactionNotFound
	case "55":
		return ErrInvalidTransaction
	case "56":
		return ErrCardOwnershipVerificationFailed
	case "57":
		return ErrCardOwnershipVerificationFailed
	case "61":
		return ErrDepositError
	case "62":
		return ErrReturnURLNotInRegisteredDomain
	case "98":
		return ErrStaticPasswordUsageLimitReached
	case "111":
		return ErrInvalidBillIssuerOrganization
	case "112":
		return ErrCardIssuerSwitchError
	case "113":
		return ErrNoResponseFromCardIssuer
	case "114":
		return ErrCardholderNotAuthorized
	case "412":
		return ErrInvalidBillID
	case "413":
		return ErrInvalidPaymentID
	case "414":
		return ErrInvalidBillIssuerOrganization
	case "415":
		return ErrTransactionTimeLimitExpired
	case "416":
		return ErrDataRegistrationError
	case "417":
		return ErrInvalidPayerID
	case "418":
		return ErrCustomerDataDefinitionError
	case "419":
		return ErrMaxDataEntryAttemptsExceeded
	case "421":
		return ErrInvalidIP
	case "995":
		return ErrCardOwnershipVerificationFailed
	case "997":
		return ErrInvalidMerchant
	default:
		return ErrInvalidResponseFormat
	}
}
