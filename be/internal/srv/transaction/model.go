package transaction

type ReqGetTransactions struct {
	UserId string `json:"userId"`
}

type ReqPostTransaction struct {
	UserId                 string  `json:"userId"`
	SenderAccountNumber    string  `json:"senderAccountNumber"`
	RecipientAccountNumber string  `json:"recipientAccountNumber"`
	Amount                 float64 `json:"amount"`
	Description            string  `json:"description"`
	OTP                    string  `json:"otp"`
}
