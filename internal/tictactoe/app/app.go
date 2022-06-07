package app

import (
	"github.com/GarmaTs/golevel2/internal/tictactoe/delivery"
	"github.com/GarmaTs/golevel2/internal/tictactoe/usecase"
)

type App struct {
}

func (app App) NewApp() {
	var Game usecase.Game
	var Player usecase.CliPlayer
	var Input delivery.CliUserInteract
	Game.Run(&Player, Input)
}
