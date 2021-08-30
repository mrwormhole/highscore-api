package function

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mrwormhole/highscore-api/model"
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

	if req.Method != http.MethodPost {
		return handler.Response{
			Body:       []byte("invalid http method."),
			StatusCode: 400,
		}, nil
	}

	var highscore model.Highscore
	json.Unmarshal(req.Body, &highscore)
	highscore.Username = strings.ToLower(highscore.Username)

	queries := repository.New(db)
	existingHighscore, err := queries.GetHighscore(req.Context(), highscore.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			params := repository.CreateHighscoreParams{Username: highscore.Username, Score: int32(highscore.Score)}
			_, err = queries.CreateHighscore(req.Context(), params)
			if err != nil {
				errMsg := fmt.Sprintf("failed to create a highscore: %v", err)
				return handler.Response{
					Body:       []byte(errMsg),
					StatusCode: 500,
				}, fmt.Errorf(errMsg)
			}
			return handler.Response{
				Body:       []byte(fmt.Sprintf("created a highscore for username %v", highscore.Username)),
				StatusCode: 200,
			}, nil
		}
		errMsg := fmt.Sprintf("failed to get a highscore: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	params := repository.UpdateHighscoreParams{ID: existingHighscore.ID, Score: existingHighscore.Score}
	_, err = queries.UpdateHighscore(req.Context(), params)
	if err != nil {
		errMsg := fmt.Sprintf("failed to update a highscore: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	err = db.Close()
	if err != nil {
		errMsg := fmt.Sprintf("failed to close db: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}
	return handler.Response{
		Body:       []byte(fmt.Sprintf("update a highscore for username %v", existingHighscore.Username)),
		StatusCode: http.StatusOK,
	}, nil
}
