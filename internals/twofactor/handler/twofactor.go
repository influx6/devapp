package handler

import (
	sessionsapi "github.com/influx6/devapp/internals/sessions/handler"
)

// TwoFactorAPI implements http methods for facing with two factor authentication.
type TwoFactorAPI struct {
	Sessions sessionsapi.SessionAPI
}
