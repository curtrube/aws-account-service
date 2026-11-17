package account

import "errors"

var (
	ErrNoAccount = errors.New("account: no matching account found")
)
