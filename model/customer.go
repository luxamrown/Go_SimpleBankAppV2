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
	Id        string `db:"id,omitempty"`
	HistoryId string `db:"history_id,omitempty"`
	CusomerId string `db:"customer_id,omitempty"`
	Amount    int    `db:"amount,omitempty"`
	Message   string `db:"message,omitempty"`
}

type HistoryMerchant struct {
	Id                 string  `db:"id,omitempty"`
	ReceiverMerchantId *string `db:"receiver_merchant_id,omitempty"`
	ReceiverCustomerId *string `db:"receiver_customer_id,omitempty"`
	SuccesAt           string  `db:"success_at"`
}

type TransactionDetail struct {
	Id                    string  `json:"id"`
	ReceiverAccountNumber *string `json:"receiver_account_number"`
	Amount                int     `json:"amount"`
	Message               string  `json:"message"`
	Date                  string  `json:"date"`
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

func NewTransactionDetail(id string, receiverAccNum *string, message, date string, amount int) TransactionDetail {
	return TransactionDetail{
		Id:                    id,
		ReceiverAccountNumber: receiverAccNum,
		Amount:                amount,
		Message:               message,
		Date:                  date,
	}
}

func NewMultipleTransactionDetail(idx int, id string, receiverAccNum *string, message, date string, amount int) []TransactionDetail {
	var multipleTransaction []TransactionDetail
	for i := 0; i < idx; i++ {
		newSingleTransaction := NewTransactionDetail(id, receiverAccNum, message, date, amount)
		multipleTransaction = append(multipleTransaction, newSingleTransaction)
	}
	return multipleTransaction
}
