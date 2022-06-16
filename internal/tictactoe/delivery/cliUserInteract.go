package delivery

import "fmt"

type CliUserInteract struct {
}

func (cli CliUserInteract) GetUserMove(isFirstPlayer bool, isVersusAI bool) (int, int) {
	var rowIdx, colIdx int
	playerNum := "Первый"
	if !isFirstPlayer {
		playerNum = "Второй"
	}
	for {
		if isVersusAI {
			fmt.Printf("Игрок, введите номер строки и номер столбца (начиная с нуля): ")
		} else {
			fmt.Printf("%s игрок, введите номер строки "+
				"и номер столбца (начиная с нуля): ", playerNum)
		}
		if _, err := fmt.Scan(&rowIdx, &colIdx); err != nil {
			fmt.Println("Ошибка: Должны быть целые числа через пробел")
		} else {
			break
		}
	}

	return rowIdx, colIdx
}

func (cli CliUserInteract) SetFieldSize() int {
	size := 3
	fmt.Print("Введите размер поля (по умолчанию три): ")
	fmt.Scan(&size)
	if size < 2 {
		size = 3
	}
	return size
}

func (cli CliUserInteract) IsVersusAI() bool {
	var text string
	fmt.Print("Против компьютера? (y/n): ")

	_, _ = fmt.Scan(&text)

	if text == "y" {
		return true
	}
	return false
}

func (cli CliUserInteract) ShowError(err error) {
	fmt.Println(err)
}

func (cli CliUserInteract) ShowField(size int, field [][]string) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if len(field[i][j]) > 0 {
				fmt.Print(field[i][j] + " ")
			} else {
				fmt.Print("_ ")
			}
		}
		fmt.Println()
	}
}

func (cli CliUserInteract) ShowGameResult(res string) {
	fmt.Println(res)
}

func (cli CliUserInteract) AskForRepeat() bool {
	var text string
	fmt.Print("Начать заново? (y/n): ")

	_, _ = fmt.Scan(&text)

	if text == "y" {
		return true
	}
	return false
}

func (cli CliUserInteract) ShowTotalScore(msg string) {
	fmt.Println(msg)
}

func (cli CliUserInteract) ShowMessage(msg string) {
	fmt.Println(msg)
}
