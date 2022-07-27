package model

type Customer struct {
	Id            string `db:"id"`
	AccountNumber string `db:"account_number"`
	UserName      string `db:"user_name"`
	UserPin       string `db:"user_pin"`
	UserPassword  string `db:"user_password"`
	UserBalance   int    `db:"balance"`
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
