package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/model"
	"mohamadelabror.me/simplebankappv2/util"
)

type CustomerRepo interface {
	RegisterNewAccount(customer model.Customer) *response.ErrorResp
	Login(username, password, token string) (string, *response.ErrorResp)
	Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber string, amount int, isMerchant bool) *response.ErrorResp
	AddLogToHistory(senderAccountNumber, receiverAccountNumber, time, id string, isMerchant bool) *response.ErrorResp
	AddTransactionDetail(id, historyId, customerId, message string, amount int) *response.ErrorResp
	GetTransactionDetail(idDetail, idHistory string, isMerchant bool) (model.TransactionDetail, *response.ErrorResp)
	GetAllTransactionDetail(idCustomer string) ([]model.TransactionDetail, *response.ErrorResp)
	GetBalanceUser(idCustomer, pin string) (*int, *response.ErrorResp)
	Logout(id string) *response.ErrorResp
}

type customerRepoImpl struct {
	customerDB *sqlx.DB
}

func (c *customerRepoImpl) RegisterNewAccount(customer model.Customer) *response.ErrorResp {
	var isUserNameExist int
	// type credToHash struct {
	// 	password string
	// 	pin      string
	// }
	err := c.customerDB.Get(&isUserNameExist, "SELECT COUNT(account_number) FROM customers WHERE user_name = $1", customer.UserName)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	if isUserNameExist == 1 {
		return response.NewError(util.ERROR_CODE_USERNAMETAKEN, fmt.Errorf(util.ERROR_MSG_USERNAMETAKEN))
	}
	hashedPass, err := util.Hashing(customer.UserPassword)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	hashedPin, err := util.Hashing(customer.UserPin)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	_, err = c.customerDB.Exec("INSERT INTO customers(id, account_number, user_name, user_pin, user_password, balance) VALUES($1, $2, $3, $4, $5, $6)", customer.Id, customer.AccountNumber, customer.UserName, hashedPin, hashedPass, customer.UserBalance)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	return nil
}

func (c *customerRepoImpl) Login(username, password, token string) (string, *response.ErrorResp) {
	type AuthGet struct {
		SelectedHashPassword string `db:"user_password"`
		Id                   string `db:"id"`
	}
	var isUserExist int
	selectedAuth := AuthGet{}
	err := c.customerDB.Get(&isUserExist, "SELECT COUNT(account_number) FROM customers WHERE user_name = $1", username)
	if err != nil {
		return "", response.NewError(util.ERROR_CODE_DB, err)
	}
	if isUserExist == 0 {
		return "", response.NewError(util.ERROR_CODE_CREDENTIALERROR, fmt.Errorf(util.ERROR_MSG_CREDENTIALERROR))
	}
	err = c.customerDB.Get(&selectedAuth, "SELECT user_password, id FROM customers WHERE user_name = $1", username)
	if err != nil {
		return "", response.NewError(util.ERROR_CODE_DB, err)
	}
	match := util.CheckHash(password, selectedAuth.SelectedHashPassword)
	if !match {
		return "", response.NewError(util.ERROR_CODE_CREDENTIALERROR, fmt.Errorf(util.ERROR_MSG_CREDENTIALERROR))
	}
	errR := c.SaveToken(selectedAuth.Id, token)
	if errR != nil {
		return "", errR
	}
	return selectedAuth.Id, nil
}

func (c *customerRepoImpl) Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber string, amount int, isMerchant bool) *response.ErrorResp {
	var isAuth int
	err := c.customerDB.Get(&isAuth, "SELECT COUNT(id) FROM customers WHERE id = $1 AND account_number = $2", sender_id, senderAccountNumber)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	if isAuth == 0 {
		return response.NewError(util.ERROR_CODE_CREDENTIALERROR, fmt.Errorf(util.ERROR_MSG_CREDENTIALERROR))
	}
	err = c.PinChecker(sender_id, sender_pin)
	if err != nil {
		return response.NewError(util.ERROR_CODE_WRONGPIN, fmt.Errorf(util.ERROR_MSG_WRONGPIN))
	}
	err = c.ReceiverExistsChecker(receiverAccountNumber, isMerchant)
	if err != nil {
		return response.NewError(util.ERROR_CODE_RECEIVERNOTFOUND, fmt.Errorf(util.ERROR_MSG_RECEIVERNOTFOUND))
	}
	err = c.BalanceValidator(senderAccountNumber, amount)
	if err != nil {
		return response.NewError(util.ERROR_CODE_BALANCE, fmt.Errorf(util.ERROR_MSG_BALANCE))
	}
	_, err = c.customerDB.Exec("UPDATE customers SET balance = balance - $1 WHERE account_number = $2", amount, senderAccountNumber)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	if isMerchant {
		_, err = c.customerDB.Exec("UPDATE merchants SET balance = balance + $1 WHERE id = $2", amount, receiverAccountNumber)
		if err != nil {
			return response.NewError(util.ERROR_CODE_DB, err)
		}
		return nil
	}
	_, err = c.customerDB.Exec("UPDATE customers SET balance = balance + $1 WHERE account_number = $2", amount, receiverAccountNumber)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	return nil
}

func (c *customerRepoImpl) AddLogToHistory(senderAccountNumber, receiverAccountNumber, time, id string, isMerchant bool) *response.ErrorResp {
	var senderId string
	var receiverId string
	err := c.customerDB.Get(&senderId, "SELECT id FROM customers WHERE account_number = $1", senderAccountNumber)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	if isMerchant {
		_, err := c.customerDB.Exec("INSERT INTO history(id, sender_id, receiver_merchant_id, success_at) VALUES ($1, $2, $3, $4)", id, senderId, receiverAccountNumber, time)
		if err != nil {
			return response.NewError(util.ERROR_CODE_DB, err)
		}
		return nil
	}
	err = c.customerDB.Get(&receiverId, "SELECT id FROM customers WHERE account_number = $1", receiverAccountNumber)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	_, err = c.customerDB.Exec("INSERT INTO history(id, sender_id, receiver_customer_id, success_at) VALUES ($1, $2, $3, $4)", id, senderId, receiverId, time)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	return nil
}

func (c *customerRepoImpl) AddTransactionDetail(id, historyId, customerId, message string, amount int) *response.ErrorResp {
	_, err := c.customerDB.Exec("INSERT INTO transaction_detail(id, history_id, customer_id, amount, message) VALUES ($1, $2, $3, $4, $5)", id, historyId, customerId, amount, message)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	return nil
}

func (c *customerRepoImpl) GetTransactionDetail(idDetail, idHistory string, isMerchant bool) (model.TransactionDetail, *response.ErrorResp) {
	var getTransactionDep model.TransactionDetailT
	var NewTransactionDetail model.TransactionDetail
	var getHistoryDep model.HistoryMerchant
	var receiverAccNumber *string
	err := c.customerDB.Get(&getTransactionDep, "SELECT amount, message FROM transaction_detail WHERE id = $1", idDetail)
	if err != nil {
		return model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
	}

	if isMerchant {
		err = c.customerDB.Get(&getHistoryDep, "SELECT receiver_merchant_id, success_at FROM history WHERE id = $1", idHistory)
		if err != nil {
			return model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
		}
		NewTransactionDetail = model.NewTransactionDetail(idDetail, getHistoryDep.ReceiverMerchantId, getTransactionDep.Message, getHistoryDep.SuccesAt, getTransactionDep.Amount)
		return NewTransactionDetail, nil
	}

	err = c.customerDB.Get(&getHistoryDep, "SELECT receiver_customer_id, success_at FROM history WHERE id = $1", idHistory)
	if err != nil {
		return model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
	}
	err = c.customerDB.Get(&receiverAccNumber, "SELECT account_number FROM customers WHERE id = $1", getHistoryDep.ReceiverCustomerId)
	if err != nil {
		return model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
	}
	NewTransactionDetail = model.NewTransactionDetail(idDetail, receiverAccNumber, getTransactionDep.Message, getHistoryDep.SuccesAt, getTransactionDep.Amount)
	return NewTransactionDetail, nil
}

func (c *customerRepoImpl) GetAllTransactionDetail(idCustomer string) ([]model.TransactionDetail, *response.ErrorResp) {
	var ListTransactionMerchant []model.TransactionDetail
	var ListTransactionCustomer []model.TransactionDetail
	var ListTransaction []model.TransactionDetail
	var getTransactionDep []model.TransactionDetailT
	var getHistoryDep []model.HistoryMerchant
	var accountNumbers []*string
	err := c.customerDB.Select(&getHistoryDep, "SELECT id, receiver_customer_id, receiver_merchant_id, success_at FROM history WHERE sender_id = $1", idCustomer)
	if err != nil {
		return []model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
	}
	for _, elemH := range getHistoryDep {
		err := c.customerDB.Select(&getTransactionDep, "SELECT id, amount, message FROM transaction_detail WHERE history_id = $1", elemH.Id)
		if err != nil {
			return []model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
		}
		for _, elemT := range getTransactionDep {
			err := c.customerDB.Select(&accountNumbers, "SELECT account_number FROM customers WHERE id = $1", elemH.ReceiverCustomerId)
			if err != nil {
				return []model.TransactionDetail{}, response.NewError(util.ERROR_CODE_DB, err)
			}
			if elemH.ReceiverMerchantId == nil {
				for _, elemNum := range accountNumbers {
					newListTransactionCustomer := model.NewMultipleTransactionDetail(len(getTransactionDep), elemT.Id, elemNum, elemT.Message, elemH.SuccesAt, elemT.Amount)
					ListTransactionCustomer = append(ListTransactionCustomer, newListTransactionCustomer...)
				}
			} else {
				newListTransactionMerchant := model.NewMultipleTransactionDetail(len(getTransactionDep), elemT.Id, elemH.ReceiverMerchantId, elemT.Message, elemH.SuccesAt, elemT.Amount)
				ListTransactionMerchant = append(ListTransactionMerchant, newListTransactionMerchant...)
			}
		}
	}
	ListTransaction = append(ListTransaction, ListTransactionCustomer...)
	ListTransaction = append(ListTransaction, ListTransactionMerchant...)

	return ListTransaction, nil
}

func (c *customerRepoImpl) GetBalanceUser(idCustomer, pin string) (*int, *response.ErrorResp) {
	var selectedBalance *int
	err := c.PinChecker(idCustomer, pin)
	if err != nil {
		return nil, response.NewError(util.ERROR_CODE_WRONGPIN, fmt.Errorf(util.ERROR_MSG_WRONGPIN))
	}
	err = c.customerDB.Get(&selectedBalance, "SELECT balance FROM customers WHERE id = $1", idCustomer)
	if err != nil {
		return nil, response.NewError(util.ERROR_CODE_DB, err)
	}
	return selectedBalance, nil
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

func (c *customerRepoImpl) PinChecker(idCustomer, pin string) error {
	var selectedPin string
	err := c.customerDB.Get(&selectedPin, "SELECT user_pin FROM customers WHERE id = $1", idCustomer)
	if err != nil {
		return err
	}
	matchPin := util.CheckHash(pin, selectedPin)
	if !matchPin {
		return fmt.Errorf("wrong pin")
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

func (c *customerRepoImpl) SaveToken(id, token string) *response.ErrorResp {
	_, err := c.customerDB.Exec("UPDATE customers SET token = $1 WHERE id = $2", token, id)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	return nil
}

func (c *customerRepoImpl) Logout(id string) *response.ErrorResp {
	_, err := c.customerDB.Exec("UPDATE customers SET token = '' WHERE id = $1", id)
	if err != nil {
		return response.NewError(util.ERROR_CODE_DB, err)
	}
	return nil
}

func NewCustomerRepo(custDb *sqlx.DB) CustomerRepo {
	return &customerRepoImpl{
		customerDB: custDb,
	}
}
