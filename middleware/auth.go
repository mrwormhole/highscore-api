package middleware

import (
	"fmt"
	"os"
	"strings"

	handler "github.com/openfaas/templates-sdk/go-http"
)

func Authorization(req handler.Request) error {
	authHeader := req.Header.Get("Authorization")
	authHeaderValues := strings.Split(authHeader, " ")
	if len(authHeaderValues) != 2 || authHeaderValues[0] != "Bearer" {
		msg := "authorization header is in the wrong format"
		return fmt.Errorf(msg)
	}
	if authHeaderValues[1] != os.Getenv("BEARER_TOKEN") {
		msg := "bearer token is not valid"
		return fmt.Errorf(msg)
	}

	return nil
}
