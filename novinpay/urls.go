package novinpay

const (
	paymentGatewayURL                = "https://pna.shaparak.ir/_ipgw_/payment/"
	baseURI                          = "https://pna.shaparak.ir/ref-payment2/RestServices/mts"
	loginURI                         = baseURI + "/merchantLogin/"
	generateSignedDataTokenURL       = baseURI + "/generateSignedDataToken/"
	generateTransactionDataToSignURL = baseURI + "/generateTransactionDataToSign/"
	verifyTransactionURL             = baseURI + "/verifyMerchantTrans/"
)
