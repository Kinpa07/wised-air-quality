package auth

import (
	"context"
	"errors"
)

type ctxKeyTokenType uint

const (
	ctxKeyToken ctxKeyTokenType = iota + 1
	ctxKeyTokenRaw
)

// GetClaimsFromContext retrieves custom claims from the context
func GetClaimsFromContext[T any](ctx context.Context) (T, error) {
	var zero T
	maybeToken, ok := ctx.Value(ctxKeyToken).(T)
	if !ok {
		return zero, errors.New("invalid token provided")
	}
	return maybeToken, nil
}

// NewContextWithClaims adds custom claims to the context
func NewContextWithClaims[T any](ctx context.Context, claims T) context.Context {
	return context.WithValue(ctx, ctxKeyToken, claims)
}

// GetTokenFromContext retrieves the default Token from the context
// This function is kept for backward compatibility
func GetTokenFromContext(ctx context.Context) (*Token, error) {
	maybeToken, ok := ctx.Value(ctxKeyToken).(*Token)
	if !ok {
		return nil, errors.New("invalid token provided")
	}
	return maybeToken, nil
}

// NewContextWithToken adds the default Token to the context
// This function is kept for backward compatibility
func NewContextWithToken(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, ctxKeyToken, token)
}

func GetTokenRawFromContext(ctx context.Context) (string, error) {
	maybeToken, ok := ctx.Value(ctxKeyTokenRaw).(string)
	if !ok {
		return "", errors.New("invalid token raw provided")
	}
	return maybeToken, nil
}

func NewContextWithRawToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxKeyTokenRaw, token)
}
