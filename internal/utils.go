package internal

import "regexp"

const (
	ErrInvalidAmount         = "invalid amount/balance %s"
	ErrInvalidTransferAmount = "invalid transfer amount %s, err %s"
)

func isStringValidNumber(s string) bool {
	re := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return re.MatchString(s)
}
