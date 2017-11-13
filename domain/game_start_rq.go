package domain

type GameStartRQ struct {
	SessionId string `json:"sessionId"`
	Player1   Player `json:"player1"`
	Player2   Player `json:"player2"`
}
