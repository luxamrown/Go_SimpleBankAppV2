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
	Logout(id string) error
}

type customerRepoImpl struct {
	customerDB *sqlx.DB
}

func (c *customerRepoImpl) RegisterNewAccount(customer model.Customer) error {
	fmt.Println(customer.UserPassword)
	hashedPass, err := util.HashPassword(customer.UserPassword)
	if err != nil {
		return err
	}
	_, err = c.customerDB.Exec("INSERT INTO customers(id, account_number, user_name, user_password, balance) VALUES($1, $2, $3, $4, $5)", customer.Id, customer.AccountNumber, customer.UserName, hashedPass, customer.UserBalance)
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
	match := util.CheckPasswordHash(password, selectedAuth.SelectedHashPassword)
	if !match {
		fmt.Println(selectedAuth.SelectedHashPassword)
		return "", fmt.Errorf("wrong credential")
	}
	return selectedAuth.Id, nil

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
