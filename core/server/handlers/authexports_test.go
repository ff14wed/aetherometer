package handlers

import (
	"context"
)

func ContextWithToken(token string) context.Context {
	return context.WithValue(context.Background(), authCtxKey, token)
}
