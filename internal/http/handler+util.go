package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Serves as a global bucket for storing errors that occur in handlers.
// Each error is associated to an UUID.
var handlerErrorsBucket = map[string]error{}

// Self redirects and appends an error uuid to the url.
func RedirectWithError(ectx echo.Context, err error) error {
	u := ectx.Request().URL
	q := u.Query()

	uuid := StoreHandlerError(err)
	q.Set(redirectErrorQueryParam, uuid)
	u.RawQuery = q.Encode()

	return ectx.Redirect(301, u.RequestURI())
}

// Tries to extract an error that may have originated from a self redirect.
func FromRedirectWithError(ectx echo.Context) (error, bool) {
	uuid := ectx.Request().URL.Query().Get(redirectErrorQueryParam)
	err, ok := handlerErrorsBucket[uuid]
	if ok {
		delete(handlerErrorsBucket, uuid)
	}

	return err, ok
}

func StoreHandlerError(err error) string {
	uuid := uuid.NewString()
	handlerErrorsBucket[uuid] = err

	return uuid
}