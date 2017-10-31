package domain

import "github.com/anhtuan29592/battleship-ai/lib"

type NotifyRQ struct {
	SessionId  string      `json:"sessionId"`
	ShotResult *ShotResult `json:"shotResult"`
}

type ShotResult struct {
	PlayerId            string        `json:"playerId"`
	Position            *lib.Point    `json:"position"`
	Status              string        `json:"status"`
	RecognizedWholeShip *ShipPosition `json:"recognizedWholeShip"`
}
