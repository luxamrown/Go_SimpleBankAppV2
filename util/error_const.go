package util

const (
	ERROR_CODE_DB               = "X01"
	ERROR_CODE_WRONGPIN         = "X02"
	ERROR_CODE_BALANCE          = "X03"
	ERROR_CODE_RECEIVERNOTFOUND = "X04"
	ERROR_CODE_CREDENTIALERROR  = "X05"
	ERROR_CODE_USERNAMETAKEN    = "X06"
)

const (
	ERROR_MSG_WRONGPIN         = "wrong pin"
	ERROR_MSG_BALANCE          = "balance not sufficent"
	ERROR_MSG_RECEIVERNOTFOUND = "receiver/merchant is not found"
	ERROR_MSG_CREDENTIALERROR  = "unauthorized/credential error"
	ERROR_MSG_USERNAMETAKEN    = "username is taken"
)
