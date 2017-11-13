package domain

import "github.com/anhtuan29592/paladin/lib"

type ShipPosition struct {
	Type      string      `json:"type"`
	Positions []lib.Point `json:"positions"`
}
