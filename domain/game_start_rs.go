package domain

type GameStartRS struct {
	Ships []*ShipPosition `json:"ships"`
}
