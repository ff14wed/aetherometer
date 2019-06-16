package handlers

import (
	"context"
	"sync/atomic"
)

func ContextWithToken(token string) context.Context {
	return context.WithValue(context.Background(), authCtxKey, token)
}

func (a *Auth) ResetOTPUsed() {
	atomic.StoreInt32(&a.otpUsed, 0)
}
