package game

import "errors"

var (
	ErrGameExit = errors.New("game ended by player")
)