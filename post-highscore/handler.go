package function

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/mrwormhole/highscore-api/middleware"
	"github.com/mrwormhole/highscore-api/model"
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
	if req.Method != http.MethodPost {
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

	var highscore model.Highscore
	err = json.Unmarshal(req.Body, &highscore)
	if err != nil {
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("failed to unmarshal highscore")
	}

	queries := repository.New(db)
	existingHighscore, err := queries.GetHighscore(req.Context(), highscore.Username)
	if err != nil && err != sql.ErrNoRows {
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("failed to get a highscore: %v", err)
	}

	if existingHighscore.ID == 0 && existingHighscore.Score == 0 {
		params := repository.CreateHighscoreParams{Username: highscore.Username, Score: highscore.Score}
		createdHighscore, err := queries.CreateHighscore(req.Context(), params)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to create a highscore: %v", err)
		}

		raw, err := json.Marshal(createdHighscore)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to marshal created highscore")
		}

		return handler.Response{
			Body:       []byte(raw),
			StatusCode: http.StatusOK,
		}, nil
	}

	if highscore.Score > existingHighscore.Score {
		params := repository.UpdateHighscoreParams{ID: existingHighscore.ID, Score: highscore.Score}
		updatedHighscore, err := queries.UpdateHighscore(req.Context(), params)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to update a highscore: %v", err)
		}

		raw, err := json.Marshal(updatedHighscore)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to marshal updated highscore")
		}

		return handler.Response{
			Body:       []byte(raw),
			StatusCode: http.StatusOK,
		}, nil
	}

	return handler.Response{
		StatusCode: http.StatusOK,
	}, nil
}
