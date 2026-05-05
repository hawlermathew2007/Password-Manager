package debug

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Defer function with Recover() need for storage maybe
// Show errors
// Log and Errors should have master password to view

type DisplayErrFunc func(title, message string)
type ErrorType uint8
 
const (
	TypeUnknown  ErrorType = iota
	TypeAuth          // wrong master password, locked vault
	TypeCrypto        // encryption / decryption failures
	TypeStorage       // file I/O, database
	TypeInput         // invalid user input / validation
	TypeInternal      // bugs, unexpected nil pointers, etc.
	TypeUI            // tview widget / rendering issues
)
 
func (t ErrorType) String() string {
	switch t {
	case TypeAuth:
		return "AuthError"
	case TypeCrypto:
		return "CryptoError"
	case TypeStorage:
		return "StorageError"
	case TypeInput:
		return "InputError"
	case TypeInternal:
		return "InternalError"
	case TypeUI:
		return "UIError"
	default:
		return "UnknownError"
	}
}

type Error struct {
	Type    ErrorType
	Message string
	Err     error
	Stack   string
	Time    time.Time
}

type Handler struct {
	Logger   *Logger
	DisplayFunc DisplayErrFunc
}

func GetUserError(usrType UsrErrType) error {
	userErrMsg := GetErrDet(usrType)
	return NewError(usrType, userErrMsg, nil)
}

func (e *Error) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func NewHandler(logger *Logger, displayFunc DisplayErrFunc) *Handler {
	return &Handler{Logger: logger, DisplayFunc: displayFunc}
}

func NewError(tp ErrorType, msg string, err error) *Error {
	return &Error{
		Type: tp,
		Message: msg,
		Err: err,
		Stack: captureStack(0),
		Time: time.Now(),
	}
}

func (h *Handler) Handle(err error) {
	if err == nil {
		return
	}
	h.Logger.LogErr(err)
 
	if h.DisplayFunc == nil {
		return
	}
 
	ae, ok := AsError(err)
	if !ok {
		h.DisplayFunc("Error", err.Error())
		return
	}
 
	title := ae.Type.String()
	h.DisplayFunc(title, ae.Message)
}

func Recover(h *Handler) {
	r := recover()
	if r == nil {
		return
	}

	var err error

	switch v := r.(type) {
	case error:
		err = NewError(
			TypeInternal,
			"An unexpected error occurred. Please restart.",
			v,
		)
	default:
		err = NewError(
			TypeInternal,
			"An unexpected error occurred. Please restart.",
			fmt.Errorf("%v", v),
		)
	}

	if h != nil {
		h.Handle(err)
	}
}

func AsError(err error) (*Error, bool) {
	var ae *Error
	return ae, errors.As(err, &ae)
}

func captureStack(skip int) string {
	pcs := make([]uintptr, 16)
	n := runtime.Callers(skip+2, pcs) // May need to change here based on how many functions we need to skip // Temp: Remove "skip" by skip = 0
	frames := runtime.CallersFrames(pcs[:n])
 
	var b strings.Builder
	for {
		f, more := frames.Next()
		// Skip runtime internals
		if strings.HasPrefix(f.Function, "runtime.") {
			if !more {
				break
			}
			continue
		}
		fmt.Fprintf(&b, "  %s\n    %s:%d\n", f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
	return b.String()
}
