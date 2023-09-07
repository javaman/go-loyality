package domain

import "errors"

var (
	ErrorOrderNotFound          error = errors.New("order not found")
	ErrorOrderExists            error = errors.New("order with same number exists")
	ErrorBadOrderNumber         error = errors.New("bad order number format")
	ErrorOrderExistsSameUser    error = errors.New("order exists")
	ErrorOrderExistsAnotherUser error = errors.New("order exists but different user")
	ErrorLoginExists            error = errors.New("user with same login exists")
	ErrorPayMoney               error = errors.New("not enoght money")
	ErrorTooFast                error = errors.New("not enoght money")
)
