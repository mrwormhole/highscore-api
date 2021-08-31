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
		log.Printf("failed to connect to db: %v", err)
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if req.Method != http.MethodPost {
		log.Printf("invalid http method %s", req.Method)
		return handler.Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	err = middleware.Authorization(req)
	if err != nil {
		log.Printf("%v", err)
		return handler.Response{
			StatusCode: http.StatusBadRequest,
		}, err
	}

	var highscore model.Highscore
	err = json.Unmarshal(req.Body, &highscore)
	if err != nil {
		log.Printf("failed to unmarshal highscore")
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	queries := repository.New(db)
	existingHighscore, err := queries.GetHighscore(req.Context(), highscore.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("failed to get a highscore: %v", err)
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if existingHighscore.ID == 0 && existingHighscore.Score == 0 {
		params := repository.CreateHighscoreParams{Username: highscore.Username, Score: int32(highscore.Score)}
		createdHighscore, err := queries.CreateHighscore(req.Context(), params)
		if err != nil {
			log.Printf("failed to create a highscore: %v", err)
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		raw, err := json.Marshal(createdHighscore)
		if err != nil {
			log.Printf("failed to marshal created highscore")
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return handler.Response{
			Body:       []byte(raw),
			StatusCode: http.StatusOK,
		}, nil
	}

	if int32(highscore.Score) > existingHighscore.Score {
		params := repository.UpdateHighscoreParams{ID: existingHighscore.ID, Score: int32(highscore.Score)}
		updatedHighscore, err := queries.UpdateHighscore(req.Context(), params)
		if err != nil {
			errMsg := fmt.Sprintf("failed to update a highscore: %v", err)
			return handler.Response{
				Body:       []byte(errMsg),
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf(errMsg)
		}

		raw, err := json.Marshal(updatedHighscore)
		if err != nil {
			log.Printf("failed to marshal updated highscore")
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
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
