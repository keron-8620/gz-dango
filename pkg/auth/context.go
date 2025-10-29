package auth

import (
	"context"
)

type contextKey string

const ctxUserClaimsKey contextKey = "user_claims"

func GetUserClaims(ctx context.Context) (*UserClaims, error) {
	uc, ok := ctx.Value(ctxUserClaimsKey).(*UserClaims)
	if !ok {
		return nil, ErrGetUserClaims
	}
	if uc == nil {
		return nil, ErrUserClaimsMissing
	}
	return uc, nil
}

func SetUserClaims(ctx context.Context, uc *UserClaims) context.Context {
	return context.WithValue(ctx, ctxUserClaimsKey, uc)
}
