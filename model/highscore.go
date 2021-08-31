package model

type Highscore struct {
	Username string `json:"username"`
	Score    int64  `json:"score"`
}
