package function

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mrwormhole/highscore-api/repository"
	handler "github.com/openfaas/templates-sdk/go-http"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB")))
	if err != nil {
		errMsg := fmt.Sprintf("failed to connect to db: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	if req.Method != http.MethodDelete {
		return handler.Response{
			Body:       []byte("invalid http method."),
			StatusCode: 400,
		}, nil
	}

	authorizationHeader := req.Header.Get("Authorization")
	authorizationHeaderValues := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderValues) != 2 || authorizationHeaderValues[0] != "Bearer" {
		errMsg := "authorization heaeder is in the wrong format"
		return handler.Response{
			Body: []byte(errMsg),
		}, fmt.Errorf(errMsg)
	}
	if authorizationHeaderValues[1] != os.Getenv("AUTH_HEADER_TOKEN") {
		errMsg := "authorization heaeder token is not valid"
		return handler.Response{
			Body: []byte(errMsg),
		}, fmt.Errorf(errMsg)
	}

	values, err := url.ParseQuery(req.QueryString)
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse query string: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	queries := repository.New(db)
	username := values.Get("username")

	if strings.TrimSpace(username) != "" {
		err = queries.DeleteHighscore(req.Context(), username)
		if err != nil {
			errMsg := fmt.Sprintf("failed to delete a highscore for username %s: %v", username, err)
			return handler.Response{
				Body:       []byte(errMsg),
				StatusCode: 500,
			}, fmt.Errorf(errMsg)
		}
	}

	return handler.Response{
		StatusCode: http.StatusOK,
	}, err
}
