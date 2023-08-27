package domain

import "errors"

var (
	ErrorBadOrderNumber         error = errors.New("bad order number format")
	ErrorOrderExistsSameUser    error = errors.New("order exists")
	ErrorOrderExistsAnotherUser error = errors.New("order exists but different user")
	ErrorLoginExists            error = errors.New("user with same login exists")
)
