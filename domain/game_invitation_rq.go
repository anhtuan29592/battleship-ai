package domain

type GameInvitationRQ struct {
	SessionId string    `json:"sessionId"`
	GameRule  GameRule `json:"gameRule"`
}
