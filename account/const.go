package account

const (
	ErrAccountNotFound = "account %d not found"
	ConversionFailed   = "parse float err %s"

	ErrAccountHasInsufficientFunds = "account %d has insufficient funds"
	ErrInitialBalanceNotPositive   = "initial balance not positive"
	ErrFailedToCreateTransaction   = "failed to create transaction %s"
)
