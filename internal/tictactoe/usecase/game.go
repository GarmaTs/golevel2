package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/GarmaTs/golevel2/internal/tictactoe/common"
)

const (
	FIRST_PLAYER_VAL  = "x"
	SECOND_PLAYER_VAL = "o"
)

type Game struct {
	size               int
	gameField          [][]string
	isFirstPlayerMoves bool
	resultTable        map[string]int
	isVersusAI         bool
	aiPlayer           AIPlayer
}

func (game *Game) initResultTable() {
	game.resultTable = make(map[string]int)
}

func (game *Game) init(size int, isVSAI bool) {
	game.isFirstPlayerMoves = true
	game.size = size
	game.isVersusAI = isVSAI

	game.gameField = make([][]string, game.size)
	for i := range game.gameField {
		game.gameField[i] = make([]string, game.size)
	}
}

func (game Game) Run(pl common.IPlayer, deliv common.IUserInteract) {
	isVSAI := deliv.IsVersusAI()
	size := deliv.SetFieldSize()
	game.initResultTable()
	game.init(size, isVSAI)

	for {
		var i, j int
		if !game.isVersusAI || game.isFirstPlayerMoves {
			i, j = deliv.GetUserMove(game.isFirstPlayerMoves, game.isVersusAI)
			if err := game.isCorrectMove(i, j); err != nil {
				deliv.ShowError(err)
				continue
			}
		} else {
			i, j = game.aiPlayer.MakeMove(game.gameField, game.size, SECOND_PLAYER_VAL)
			deliv.ShowMessage("Ход комьютера:")
		}

		if game.isFirstPlayerMoves {
			game.setCell(i, j, FIRST_PLAYER_VAL)
		} else {
			game.setCell(i, j, SECOND_PLAYER_VAL)
		}

		res, msg := game.isGameFinished()
		if !res {
			deliv.ShowField(game.size, game.gameField)
		} else {
			if game.isFirstPlayerMoves {
				game.setResultTable(FIRST_PLAYER_VAL)
			} else {
				game.setResultTable(SECOND_PLAYER_VAL)
			}

			deliv.ShowField(game.size, game.gameField)
			deliv.ShowGameResult(msg)
			deliv.ShowTotalScore(game.totalScoreToStr())

			// Запрос на повтор
			startAgain := deliv.AskForRepeat()
			if !startAgain {
				break
			} else {
				isVSAI := deliv.IsVersusAI()
				size := deliv.SetFieldSize()
				game.init(size, isVSAI)
				continue
			}
		}
		game.isFirstPlayerMoves = !game.isFirstPlayerMoves
	}
}

func (game Game) isCorrectMove(i, j int) error {
	if i < 0 {
		return errors.New("Номер строки не может быть меньше нуля")
	}
	if i > game.size-1 {
		return errors.New(fmt.Sprintf("Номер строки не может быть больше %s",
			strconv.Itoa(game.size-1)))
	}

	if j < 0 {
		return errors.New("Номер столбца не может быть меньше нуля")
	}
	if j > game.size-1 {
		return errors.New(fmt.Sprintf("Номер столбца не может быть больше %s",
			strconv.Itoa(game.size-1)))
	}

	if len(game.gameField[i][j]) > 0 {
		return errors.New("Эта ячейка занята")
	}

	return nil
}

func (game *Game) setCell(i, j int, value string) {
	game.gameField[i][j] = value
}

func (game Game) isGameFinished() (bool, string) {
	var gameIsFinished bool
	var message string

	// Проверка на победителя
	gameIsFinished = game.haveWinner()
	if gameIsFinished {
		playerNum := FIRST_PLAYER_VAL
		if !game.isFirstPlayerMoves {
			playerNum = SECOND_PLAYER_VAL
		}
		message = fmt.Sprintf("Конец игры. Игрок \"%s\" победил.", playerNum)
	}

	// Проверка на ничью
	if !gameIsFinished {
		gameIsFinished = game.haveDraw()
		if gameIsFinished {
			message = "Ничья. Ходов не осталось."
		}
	}

	return gameIsFinished, message
}

func (game Game) haveDraw() bool {
	res := true
	for i := 0; i < game.size; i++ {
		for j := 0; j < game.size; j++ {
			if len(game.gameField[i][j]) == 0 {
				res = false
				break
			}
		}
	}
	return res
}

func (game Game) haveWinner() bool {
	tgtVal := FIRST_PLAYER_VAL
	if !game.isFirstPlayerMoves {
		tgtVal = SECOND_PLAYER_VAL
	}

	isWinnerFound := false
	// Проверка по диагоналям
	res := game.diagonalsHasWinner(tgtVal)
	if res {
		isWinnerFound = true
	}

	// Проверка по строкам
	for i := 0; i < game.size; i++ {
		res := game.rowHasWinner(i, tgtVal)
		if res {
			isWinnerFound = true
		}
	}

	// Проверка по столбцам
	for j := 0; j < game.size; j++ {
		res := game.colHasWinner(j, tgtVal)
		if res {
			isWinnerFound = true
		}
	}
	return isWinnerFound
}

func (game Game) rowHasWinner(rowIdx int, tgtVal string) bool {
	res := true
	for j := 0; j < game.size; j++ {
		if game.gameField[rowIdx][j] != tgtVal {
			res = false
			break
		}
	}
	return res
}

func (game Game) colHasWinner(colIdx int, tgtVal string) bool {
	res := true
	for i := 0; i < game.size; i++ {
		if game.gameField[i][colIdx] != tgtVal {
			res = false
			break
		}
	}
	return res
}

func (game Game) diagonalsHasWinner(tgtVal string) bool {
	res := true
	for i := 0; i < game.size; i++ {
		if game.gameField[i][i] != tgtVal {
			res = false
			break
		}
	}
	if res {
		return res
	}

	res = true
	j := game.size
	for i := 0; i < game.size; i++ {
		j--
		if game.gameField[i][j] != tgtVal {
			res = false
			break
		}
	}
	return res
}

func (game *Game) setResultTable(playerValue string) {
	prevResult := game.resultTable[playerValue]
	game.resultTable[playerValue] = prevResult + 1
}

func (game Game) totalScoreToStr() string {
	message := fmt.Sprintf("Итого: игрок \"%s\" победил %d раз, ",
		FIRST_PLAYER_VAL, game.resultTable[FIRST_PLAYER_VAL])
	message += fmt.Sprintf("игрок \"%s\" победил %d раз.",
		SECOND_PLAYER_VAL, game.resultTable[SECOND_PLAYER_VAL])
	return message
}
