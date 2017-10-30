package ship

import "github.com/anhtuan29592/battleship-ai/lib"

type Ship interface {
	get_positions() ([]*lib.Point)
	conflict_with(other *Ship) (bool)
	is_valid(board_w int, board_h int) (bool)
}
