package errs

import "errors"

var ErrLoginUniqueViolation = errors.New("login should be unique")
var ErrLoginOrPasswordNotFound = errors.New("login or password not found")
var ErrOrderNumberFormat = errors.New("wrong order number format")
var ErrOrderDuplicate = errors.New("this order number is already added")
var ErrOrderOtherUserDuplicate = errors.New("this order number is already added by other user")
var ErrOrderNotFound = errors.New("order not found")
var ErrUserNotFound = errors.New("user not found")
var ErrBalanceNotEnoughFunds = errors.New("not enough funds")
