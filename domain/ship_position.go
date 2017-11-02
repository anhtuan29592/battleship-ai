package domain

import "github.com/anhtuan29592/battleship-ai/lib"

type ShipPosition struct {
	Type      string       `json:"type"`
	Positions []lib.Point `json:"positions"`
}
