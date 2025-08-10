package internal

import "regexp"

const (
	ErrInvalidAmount         = "invalid amount/balance %s"
	ErrInvalidAccount        = "invalid internal id %s"
	ErrInvalidJson           = "invalid input %s"
	ErrInvalidTransferAmount = "invalid transfer amount %s, err %s"
	ErrAccountTaken          = "internal ID %d already taken, please choose another"

	ErrTransactionFailed = "transaction failed %s"
)

func isStringValidNumber(s string) bool {
	re := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return re.MatchString(s)
}
