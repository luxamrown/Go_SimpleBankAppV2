package appreq

type RegisterReq struct {
	UserName     string `json:"user_name"`
	UserPin      string `json:"user_pin"`
	UserPassword string `json:"user_password"`
	Balance      int    `json:"balance"`
}

type LogoutReq struct {
	Id string `json:"account_id"`
}

type TransactionReq struct {
	SenderAccNumber       string `json:"sender_account_number"`
	SenderId              string `json:"sender_id"`
	SenderPin             string `json:"sender_pin"`
	Message               string `json:"transfer_message"`
	ReceiverAccountNumber string `json:"receiver_account_number"`
	Amount                int    `json:"amount_transfer"`
	IsMerchant            bool   `json:"is_merchant"`
}

// type RequestTest struct {
// 	Id          string `json:"id"`
// 	PublishDate string `json:"publish_date"`
// 	Lain        int    `json:"lain"`
// }
