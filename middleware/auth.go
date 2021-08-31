package middleware

import (
	"errors"
	"os"
	"strings"

	handler "github.com/openfaas/templates-sdk/go-http"
)

func Authorization(req handler.Request) error {
	authHeader := req.Header.Get("Authorization")
	authHeaderValues := strings.Split(authHeader, " ")
	if len(authHeaderValues) != 2 || authHeaderValues[0] != "Bearer" {
		return errors.New("authorization header is in the wrong format")
	}
	if authHeaderValues[1] != os.Getenv("BEARER_TOKEN") {
		return errors.New("bearer token is not valid")
	}

	return nil
}
