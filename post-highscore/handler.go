package function

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

	queries := repository.New(db)
	_, err = queries.GetHighscore(req.Context(), highscore.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("123failed to get a highscore: %v", err)
			return handler.Response{
				Body:       []byte(errMsg),
				StatusCode: 500,
			}, fmt.Errorf(errMsg)
		}
		errMsg := fmt.Sprintf("failed to get a highscore: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}

	//queries.CreateHighscore(req.Context(), repository.CreateHighscoreParams{Username: "JACK", Score: 555})
	//queries.UpdateHighscore(req.Context(), repository.UpdateHighscoreParams{ID: 4, Score: 555})

	err = db.Close()
	if err != nil {
		errMsg := fmt.Sprintf("failed to close db: %v", err)
		return handler.Response{
			Body:       []byte(errMsg),
			StatusCode: 500,
		}, fmt.Errorf(errMsg)
	}
	return handler.Response{
		StatusCode: http.StatusOK,
	}, nil
}
