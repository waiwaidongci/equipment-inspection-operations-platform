package auth

import "context"

type contextKey string

const claimsKey contextKey = "auth-claims"

func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFrom(ctx context.Context) (Claims, bool) {
	value := ctx.Value(claimsKey)
	claims, ok := value.(Claims)
	return claims, ok
}
