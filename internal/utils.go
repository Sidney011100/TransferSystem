package internal

import (
	"regexp"
	"transferSystem/model"
)

const (
	ErrInvalidAmount         = "invalid amount/balance %s"
	ErrInvalidTransferAmount = "invalid transfer amount %s, err %s"
)

var cacheAccountLRU = make(map[int64]*model.Account)

func getAccountFromCache(id int64) (*model.Account, bool) {
	cache, ok := cacheAccountLRU[id]
	if !ok {
		return nil, false
	}
	return cache, true
}

func updateCachedAccount(id int64, newAcc *model.Account) {
	cacheAccountLRU[id] = newAcc
}

func removeCachedAccount(id int64) {
	delete(cacheAccountLRU, id)
}

func isStringValidNumber(s string) bool {
	re := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return re.MatchString(s)
}
