package domain

type GameRule struct {
	BoardWidth  int            `json:"boardWidth"`
	BoardHeight int            `json:"boardHeight"`
	Ships       []ShipQuantity `json:"ships"`
}

type ShipQuantity struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}
