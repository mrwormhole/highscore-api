package function

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mrwormhole/highscore-api/middleware"
	"github.com/mrwormhole/highscore-api/repository"
	handler "github.com/openfaas/templates-sdk/go-http"
)

func Handle(req handler.Request) (handler.Response, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB")))
	defer func() {
		err = db.Close()
		if err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()
	if err != nil {
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("failed to connect to db: %v", err)
	}
	if req.Method != http.MethodDelete {
		return handler.Response{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("invalid http method %s", req.Method)
	}

	err = middleware.Authorization(req)
	if err != nil {
		return handler.Response{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("%v", err)
	}

	values, err := url.ParseQuery(req.QueryString)
	if err != nil {
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("failed to parse query string: %v", err)
	}

	queries := repository.New(db)
	username := values.Get("username")

	if strings.TrimSpace(username) != "" {
		err = queries.DeleteHighscore(req.Context(), username)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to delete a highscore for username %s: %v", username, err)
		}
	}

	return handler.Response{
		StatusCode: http.StatusOK,
	}, nil
}
