package appreq

type RegisterReq struct {
	AccountNumber string `json:"account_number"`
	UserName      string `json:"user_name"`
	UserPassword  string `json:"user_password"`
	Balance       int    `json:"balance"`
}

type LogoutReq struct {
	Id string `json:"account_id"`
}
