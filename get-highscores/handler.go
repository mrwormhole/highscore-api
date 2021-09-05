package function

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
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
	if req.Method != http.MethodGet {
		return handler.Response{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("invalid http method %s", req.Method)
	}

	values, err := url.ParseQuery(req.QueryString)
	if err != nil {
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("failed to parse query string: %v", err)
	}

	var rawBody []byte
	queries := repository.New(db)
	username := values.Get("username")

	if strings.TrimSpace(username) != "" {
		highscore, err := queries.GetHighscore(req.Context(), username)
		if err != nil {
			if err == sql.ErrNoRows {
				return handler.Response{
					StatusCode: http.StatusNotFound,
				}, nil
			}
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to get highscore for username %s: %v", username, err)
		}

		rawBody, err = json.Marshal(highscore)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to marshal a highscore: %v", err)
		}
	} else {
		highscores, err := queries.ListHighscores(req.Context())
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to list highscores: %v", err)
		}

		rawBody, err = json.Marshal(highscores)
		if err != nil {
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("failed to marshal highscores: %v", err)
		}
	}

	return handler.Response{
		Body:       rawBody,
		StatusCode: http.StatusOK,
	}, nil
}
