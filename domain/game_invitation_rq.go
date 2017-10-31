package domain

type GameInviteRQ struct {
	SessionId string `json:"sessionId"`
	GameRule *GameRule `json:"gameRule"`
}
