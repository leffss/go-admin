package codemsg

import "net/http"

const (
	Success             = http.StatusOK
	Error               = http.StatusInternalServerError

	ErrorNoToken 		= 1000
	ErrorToken 			= 1001
	ErrorTokenCreate  	= 1002
	ErrorTokenExpire 	= 1003
	ErrorTokenFormat 	= 1004

	ErrorInvalidParams  = 2000

	ErrorUserNotExist   = 3000
	ErrorUserWrongPass  = 3001
	ErrorUserStatus     = 3002
)
