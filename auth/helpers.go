package auth

import (
	"context"

	"github.com/spf13/cast"
)

type contextKey string

func (c contextKey) String() string {
	return "ck/" + string(c)
}

var enterpriseContextKey = contextKey("enterprise-id")

func SetEnterpriseID(ctx context.Context, entID uint64) context.Context {
	return context.WithValue(ctx, enterpriseContextKey, entID)
}

func GetEnterpriseID(ctx context.Context) uint64 {
	return cast.ToUint64(ctx.Value(enterpriseContextKey))
}

var userIDContextKey = contextKey("user-id")

func SetUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetUserID(ctx context.Context) uint64 {
	return cast.ToUint64(ctx.Value(userIDContextKey))
}
