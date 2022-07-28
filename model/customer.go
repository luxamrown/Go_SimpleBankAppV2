package model

type Customer struct {
	Id            string `db:"id"`
	AccountNumber string `db:"account_number"`
	UserName      string `db:"user_name"`
	UserPin       string `db:"user_pin"`
	UserPassword  string `db:"user_password"`
	UserBalance   int    `db:"balance"`
}

type TransactionDetailT struct {
	Id      string `db:"id"`
	Amount  int    `db:"amount"`
	Message string `db:"message"`
}

type HistoryMerchant struct {
	ReceiverMerchantId string `db:"receiver_merchant_id"`
	SuccesAt           string `db:"success_at"`
}

type HistoryCustomer struct {
	ReceiverId string `db:"receiver_customer_id"`
	SuccesAt   string `db:"success_at"`
}

type TransactionDetail struct {
	Id                    string
	ReceiverAccountNumber string
	Amount                int
	Message               string
	Date                  string
}

func NewCustomer(id, accountNumber, userName, userPin, userPassword string, balance int) Customer {
	return Customer{
		Id:            id,
		AccountNumber: accountNumber,
		UserName:      userName,
		UserPin:       userPin,
		UserPassword:  userPassword,
		UserBalance:   balance,
	}
}

func NewTransactionDetail(id, receiverAccNum, message, date string, amount int) TransactionDetail {
	return TransactionDetail{
		Id:                    id,
		ReceiverAccountNumber: receiverAccNum,
		Amount:                amount,
		Message:               message,
		Date:                  date,
	}
}
