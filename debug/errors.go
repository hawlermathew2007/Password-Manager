package debug

import (
)

type UsrErrType uint8
type SysErrType uint8 // File, app, ui yada yada

const (
	WrongMasterPwd UsrErrType = iota
	BadPwdPolicy
)

func GetErrDet(usrType UsrErrType) string  {
	var msg string
	switch usrType {
		case WrongMasterPwd:
			msg = "Incorrect password. Please try again"
		default:
			msg = "Unknown Error..."
	}
	return msg
}
