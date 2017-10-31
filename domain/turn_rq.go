package domain

type TurnRQ struct {
	SessionId  string `json:"sessionId"`
	TurnNumber int    `json:"turnNumber"`
}
