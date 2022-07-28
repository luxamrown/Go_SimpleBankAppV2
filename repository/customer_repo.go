package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"mohamadelabror.me/simplebankappv2/model"
	"mohamadelabror.me/simplebankappv2/util"
)

type CustomerRepo interface {
	RegisterNewAccount(customer model.Customer) error
	Login(username, password string) (string, error)
	SaveToken(accountNumber, token string) error
	Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber string, amount int, isMerchant bool) error
	AddLogToHistory(senderAccountNumber, receiverAccountNumber, time, id string, isMerchant bool) error
	AddTransactionDetail(id, historyId, message string, amount int) error
	GetTransactionDetail(idDetail, idHistory string, isMerchant bool) (model.TransactionDetail, error)
	Logout(id string) error
}

type customerRepoImpl struct {
	customerDB *sqlx.DB
}

func (c *customerRepoImpl) RegisterNewAccount(customer model.Customer) error {
	var isUserNameExist int
	// type credToHash struct {
	// 	password string
	// 	pin      string
	// }
	err := c.customerDB.Get(&isUserNameExist, "SELECT COUNT(account_number) FROM customers WHERE user_name = $1", customer.UserName)
	if err != nil {
		return err
	}
	if isUserNameExist == 1 {
		return fmt.Errorf("username is taken")
	}
	hashedPass, err := util.Hashing(customer.UserPassword)
	if err != nil {
		return err
	}
	hashedPin, err := util.Hashing(customer.UserPin)
	if err != nil {
		return err
	}
	_, err = c.customerDB.Exec("INSERT INTO customers(id, account_number, user_name, user_pin, user_password, balance) VALUES($1, $2, $3, $4, $5, $6)", customer.Id, customer.AccountNumber, customer.UserName, hashedPin, hashedPass, customer.UserBalance)
	if err != nil {
		return err
	}
	return nil
}

func (c *customerRepoImpl) Login(username, password string) (string, error) {
	type AuthGet struct {
		SelectedHashPassword string `db:"user_password"`
		Id                   string `db:"id"`
	}
	var isUserExist int
	selectedAuth := AuthGet{}
	err := c.customerDB.Get(&isUserExist, "SELECT COUNT(account_number) FROM customers WHERE user_name = $1", username)
	if err != nil {
		return "", err
	}
	if isUserExist == 0 {
		return "", fmt.Errorf("user not found")
	}
	err = c.customerDB.Get(&selectedAuth, "SELECT user_password, id FROM customers WHERE user_name = $1", username)
	if err != nil {
		return "", err
	}
	match := util.CheckHash(password, selectedAuth.SelectedHashPassword)
	if !match {
		return "", fmt.Errorf("wrong credential")
	}
	return selectedAuth.Id, nil
}

func (c *customerRepoImpl) Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber string, amount int, isMerchant bool) error {
	var isAuth int
	var selectedPin string
	err := c.customerDB.Get(&isAuth, "SELECT COUNT(id) FROM customers WHERE id = $1 AND account_number = $2", sender_id, senderAccountNumber)
	if err != nil {
		return err
	}
	if isAuth == 0 {
		return fmt.Errorf("unauthorized user")
	}
	err = c.customerDB.Get(&selectedPin, "SELECT user_pin FROM customers WHERE id = $1", sender_id)
	if err != nil {
		return err
	}
	matchPin := util.CheckHash(sender_pin, selectedPin)
	if !matchPin {
		return fmt.Errorf("wrong pin")
	}
	err = c.ReceiverExistsChecker(receiverAccountNumber, isMerchant)
	if err != nil {
		return err
	}
	err = c.BalanceValidator(senderAccountNumber, amount)
	if err != nil {
		return err
	}
	_, err = c.customerDB.Exec("UPDATE customers SET balance = balance - $1 WHERE account_number = $2", amount, senderAccountNumber)
	if err != nil {
		return err
	}
	if isMerchant {
		_, err = c.customerDB.Exec("UPDATE merchants SET balance = balance + $1 WHERE id = $2", amount, receiverAccountNumber)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = c.customerDB.Exec("UPDATE customers SET balance = balance + $1 WHERE account_number = $2", amount, receiverAccountNumber)
	if err != nil {
		return err
	}
	return nil
}

func (c *customerRepoImpl) AddLogToHistory(senderAccountNumber, receiverAccountNumber, time, id string, isMerchant bool) error {
	var senderId string
	var receiverId string
	err := c.customerDB.Get(&senderId, "SELECT id FROM customers WHERE account_number = $1", senderAccountNumber)
	if err != nil {
		return err
	}
	if isMerchant {
		_, err := c.customerDB.Exec("INSERT INTO history(id, sender_id, receiver_merchant_id, success_at) VALUES ($1, $2, $3, $4)", id, senderId, receiverAccountNumber, time)
		if err != nil {
			return err
		}
		return nil
	}
	err = c.customerDB.Get(&receiverId, "SELECT id FROM customers WHERE account_number = $1", receiverAccountNumber)
	if err != nil {
		return err
	}
	_, err = c.customerDB.Exec("INSERT INTO history(id, sender_id, receiver_customer_id, success_at) VALUES ($1, $2, $3, $4)", id, senderId, receiverId, time)
	if err != nil {
		return err
	}
	return nil
}

func (c *customerRepoImpl) AddTransactionDetail(id, historyId, message string, amount int) error {
	_, err := c.customerDB.Exec("INSERT INTO transaction_detail(id, history_id, amount, message) VALUES ($1, $2, $3, $4)", id, historyId, amount, message)
	if err != nil {
		return err
	}
	return nil
}

func (c *customerRepoImpl) GetTransactionDetail(idDetail, idHistory string, isMerchant bool) (model.TransactionDetail, error) {
	var getTransactionDep model.TransactionDetailT
	var getHistoryDep model.HistoryCustomer
	var getHistoryDepM model.HistoryMerchant
	var NewTransactionDetail model.TransactionDetail
	var receiverAccNumber string
	err := c.customerDB.Get(&getTransactionDep, "SELECT amount, message FROM transaction_detail WHERE id = $1", idDetail)
	if err != nil {
		return model.TransactionDetail{}, err
	}

	if isMerchant {
		err = c.customerDB.Get(&getHistoryDepM, "SELECT receiver_merchant_id, success_at FROM history WHERE id = $1", idHistory)
		if err != nil {
			return model.TransactionDetail{}, err
		}
		NewTransactionDetail = model.NewTransactionDetail(idDetail, getHistoryDepM.ReceiverMerchantId, getTransactionDep.Message, getHistoryDepM.SuccesAt, getTransactionDep.Amount)
		return NewTransactionDetail, nil
	}

	err = c.customerDB.Get(&getHistoryDep, "SELECT receiver_customer_id, success_at FROM history WHERE id = $1", idHistory)
	if err != nil {
		return model.TransactionDetail{}, err
	}
	err = c.customerDB.Get(&receiverAccNumber, "SELECT account_number FROM customers WHERE id = $1", getHistoryDep.ReceiverId)
	if err != nil {
		return model.TransactionDetail{}, err
	}
	NewTransactionDetail = model.NewTransactionDetail(idDetail, receiverAccNumber, getTransactionDep.Message, getHistoryDep.SuccesAt, getTransactionDep.Amount)
	return NewTransactionDetail, nil
}

func (c *customerRepoImpl) ReceiverExistsChecker(accountNumber string, isMerchant bool) error {
	var isReceiverExist int
	if isMerchant {
		err := c.customerDB.Get(&isReceiverExist, "SELECT COUNT(id) from merchants WHERE id = $1", accountNumber)
		if err != nil {
			return err
		}
		if isReceiverExist == 0 {
			return fmt.Errorf("merchant not found")
		}
		return nil
	}
	err := c.customerDB.Get(&isReceiverExist, "SELECT COUNT(id) FROM customers WHERE account_number = $1", accountNumber)
	if err != nil {
		return err
	}
	if isReceiverExist == 0 {
		return fmt.Errorf("receiver not found")
	}
	return nil
}

func (c *customerRepoImpl) BalanceValidator(accountNumber string, ammount int) error {
	var balance int
	err := c.customerDB.Get(&balance, "SELECT balance FROM customers WHERE account_number = $1", accountNumber)
	if err != nil {
		return err
	}
	if balance < ammount {
		return fmt.Errorf("balance not sufficent")
	}
	return nil
}

func (c *customerRepoImpl) SaveToken(id, token string) error {
	_, err := c.customerDB.Exec("UPDATE customers SET token = $1 WHERE id = $2", token, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *customerRepoImpl) Logout(id string) error {
	_, err := c.customerDB.Exec("UPDATE customers SET token = '' WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewCustomerRepo(custDb *sqlx.DB) CustomerRepo {
	return &customerRepoImpl{
		customerDB: custDb,
	}
}
