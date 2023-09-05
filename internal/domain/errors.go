package domain

import "errors"

var (
	ErrorOrderExists             error = errors.New("order with same number exists")
	ErrorBadOrderNumber         error = errors.New("bad order number format")
	ErrorOrderExistsSameUser    error = errors.New("order exists")
	ErrorOrderExistsAnotherUser error = errors.New("order exists but different user")
	ErrorLoginExists            error = errors.New("user with same login exists")
)
