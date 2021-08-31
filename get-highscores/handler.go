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
		log.Printf("failed to connect to db: %v", err)
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if req.Method != http.MethodGet {
		log.Printf("invalid http method %s", req.Method)
		return handler.Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	values, err := url.ParseQuery(req.QueryString)
	if err != nil {
		log.Printf("failed to parse query string: %v", err)
		return handler.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	var rawBody []byte
	queries := repository.New(db)
	username := values.Get("username")

	if strings.TrimSpace(username) != "" {
		highscore, err := queries.GetHighscore(req.Context(), username)
		if err != nil {
			log.Printf("failed to get highscore for username %s: %v", username, err)
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		rawBody, err = json.Marshal(highscore)
		if err != nil {
			log.Printf("failed to marshal a highscore: %v", err)
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	} else {
		highscores, err := queries.ListHighscores(req.Context())
		if err != nil {
			log.Printf("failed to list highscores: %v", err)
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		rawBody, err = json.Marshal(highscores)
		if err != nil {
			log.Printf("failed to marshal highscores: %v", err)
			return handler.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	}

	return handler.Response{
		Body:       rawBody,
		StatusCode: http.StatusOK,
	}, fmt.Errorf("WELL WELL WELL WELL I AM FAKE")
}
