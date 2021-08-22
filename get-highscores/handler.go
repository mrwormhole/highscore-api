package function

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
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
		}, fmt.Errorf("failed to connect to db: %v", err)
	}

	queries := repository.New(db)

	highscores, err := queries.ListHighscores(req.Context())
	if err != nil {
		errMsg := fmt.Sprintf("failed to list highscores: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	rawBody, err := json.Marshal(highscores)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal highscores: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	return handler.Response{
		Body:       rawBody,
		StatusCode: http.StatusOK,
	}, err
}
