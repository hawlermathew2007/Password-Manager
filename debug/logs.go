package debug

import (
	"fmt"
)

type LogType uint8 // File, app, ui yada yada
 
const (
	MasterLogin LogType = iota
	MasterCreateAcc
	MasterDeleteAcc
	MasterUpdateAcc
	PwdScanned
	PwdCopied
	PwdShown
	TreeNodeSelected
	TreeChildSelected
)

type LogContext struct {
	Username 		string
	Domain 			string
	ScanResult	string
}

func GetLogFormat(logType LogType, context LogContext) string {
	var msg string
	switch logType {
		case MasterLogin:
			msg = fmt.Sprintf("[%s] Master logged in successfully.", logType)
		case MasterCreateAcc:
			msg = fmt.Sprintf("[%s] Account created: %s in %s", logType, context.Username, context.Domain)
		case MasterDeleteAcc:
			msg = fmt.Sprintf("[%s] Account deleted: %s in %s", logType, context.Username, context.Domain)
		case MasterUpdateAcc:
			msg = fmt.Sprintf("[%s] Account updated: %s in %s", logType, context.Username, context.Domain)
		case PwdCopied:
			msg = fmt.Sprintf("[%s] Password copied: %s in %s", logType, context.Username, context.Domain)
		case PwdShown:
			msg = fmt.Sprintf("[%s] Password revealed: %s in %s", logType, context.Username, context.Domain)
		case PwdScanned:
			msg = fmt.Sprintf("[%s] Password scanned: %s in %s — Result: %s", logType, context.Username, context.Domain, context.ScanResult)
		case TreeNodeSelected:
			msg = fmt.Sprintf("[%s] Accounts in Domain:%s is viewed", logType, context.Domain)
		case TreeChildSelected:
			msg = fmt.Sprintf("[%s] Account:%s in Domain:%s is viewed", logType, context.Username, context.Domain)
		default:
			msg = "[Unknown] The log context is invalid for the provided log type."
	}
	return msg
}
