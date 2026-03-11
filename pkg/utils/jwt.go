package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrEmptyID      = errors.New("empty id in token claims")
)

// GetJWTClaimsFromRequest extracts JWT claims from the Gin context.
// It assumes that the claims have been set in the context by prior middleware.
func GetJWTClaimsFromRequest(ctx *gin.Context) (jwt.MapClaims, error) {
	claimsValue, _ := ctx.Get("claims")
	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

var PrivateKeyCompromisedError = newrelic.Error{
	Message: "JWT is invalid format - Private key might be compromised",
	Class:   "JWTError",
}

// GetUserIDFromJWTClaims retrieves the user ID from JWT claims in the Gin context.
// It looks for the "sub" claim which is expected to contain the user ID.
func GetUserIDFromJWTClaims(ctx *gin.Context) (string, error) {
	claims, err := GetJWTClaimsFromRequest(ctx)
	if err != nil {
		return "", err
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		newrelic.FromContext(ctx).NoticeError(PrivateKeyCompromisedError)
		return "", ErrEmptyID
	}

	return userID, nil
}
