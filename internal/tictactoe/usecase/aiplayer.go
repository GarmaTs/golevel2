package usecase

import (
	"sync"
)

type Cell struct {
	i int
	j int
}

type mapStruct struct {
	sync.Mutex
	dict map[Cell]int
	size int
}

func NewMapStruct() *mapStruct {
	m := new(mapStruct)
	m.dict = make(map[Cell]int)
	return m
}

func (m *mapStruct) Get(key Cell) int {
	m.Lock()
	defer m.Unlock()
	if v, has := m.dict[key]; has {
		return v
	}
	return 0
}

func (m *mapStruct) Set(key Cell, val int) {
	m.Lock()
	defer m.Unlock()
	m.dict[key] = val
}

func (m *mapStruct) GetCellWithMaxPrior() Cell {
	cell := Cell{0, 0}
	maxPrior := m.Get(Cell{0, 0})
	for _, val := range m.dict {
		if val > maxPrior {
			maxPrior = val
		}
	}
	for key, val := range m.dict {
		if val == maxPrior {
			cell.i = key.i
			cell.j = key.j
			break
		}
	}
	return cell
}

type AIPlayer struct {
	marker   string
	priorMap *mapStruct // Для каждой ячейки выставляется приоритет
}

var wg sync.WaitGroup

func (ai AIPlayer) MakeMove(gameField [][]string, size int, marker string) (int, int) {
	ai.priorMap = NewMapStruct()
	ai.marker = marker

	wg.Add(size * size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			go ai.setPriority(ai.priorMap, gameField, i, j)
		}
	}
	wg.Wait()

	cell := ai.priorMap.GetCellWithMaxPrior()

	return cell.i, cell.j
}

func (ai AIPlayer) setPriority(m *mapStruct, gameField [][]string, rowIdx, colIdx int) {
	defer wg.Done()

	m.Set(Cell{rowIdx, colIdx}, 0)
	size := len(gameField)
	if len(gameField[rowIdx][colIdx]) > 0 {
		m.Set(Cell{rowIdx, colIdx}, -1) // Занятую ячейку не смотрим
		return
	}

	maxPrior := size * size * 10
	// Если ai-игроку остался последний ход - то эта ячейка получает maxPrior
	// Если противнику остался последний ход - то эта ячейка получает maxPrior-1
	// Проверить кол-во потенциальных побед
	// Проверить кол-во ходов до победы

	var countPossibleWins int
	// Проверяем строку
	{
		cEmptyCells := 0
		enemyMarkersCount := 0
		ownMarkersCount := 0
		for j := 0; j < size; j++ {
			if len(gameField[rowIdx][j]) == 0 {
				cEmptyCells++
			} else {
				if gameField[rowIdx][j] != ai.marker {
					enemyMarkersCount++
				} else {
					ownMarkersCount++
				}
			}
		}

		resultForLine := makeDecisionForLine(maxPrior, size, cEmptyCells, enemyMarkersCount, ownMarkersCount)
		if resultForLine == maxPrior {
			m.Set(Cell{rowIdx, colIdx}, maxPrior)
			return
		}
		if resultForLine == maxPrior-1 {
			m.Set(Cell{rowIdx, colIdx}, maxPrior-1)
			return
		}
		countPossibleWins += resultForLine
	}

	// Проверяем столбец
	{
		cEmptyCells := 0
		enemyMarkersCount := 0
		ownMarkersCount := 0
		for i := 0; i < size; i++ {
			if len(gameField[i][colIdx]) == 0 {
				cEmptyCells++
			} else {
				if gameField[i][colIdx] != ai.marker {
					enemyMarkersCount++
				} else {
					ownMarkersCount++
				}
			}
		}

		resultForLine := makeDecisionForLine(maxPrior, size, cEmptyCells, enemyMarkersCount, ownMarkersCount)
		if resultForLine == maxPrior {
			m.Set(Cell{rowIdx, colIdx}, maxPrior)
			return
		}
		if resultForLine == maxPrior-1 {
			m.Set(Cell{rowIdx, colIdx}, maxPrior-1)
			return
		}
		countPossibleWins += resultForLine
	}

	// Проверяем по первой диагонали
	{
		isOnFirstFirstDiagonal := false
		for i := 0; i < size; i++ {
			if rowIdx == i && colIdx == i {
				isOnFirstFirstDiagonal = true
				break
			}
		}

		if isOnFirstFirstDiagonal {
			cEmptyCells := 0
			enemyMarkersCount := 0
			ownMarkersCount := 0

			for i := 0; i < size; i++ {
				if len(gameField[i][i]) == 0 {
					cEmptyCells++
				} else {
					if gameField[i][i] != ai.marker {
						enemyMarkersCount++
					} else {
						ownMarkersCount++
					}
				}
			}

			resultForLine := makeDecisionForLine(maxPrior, size, cEmptyCells, enemyMarkersCount, ownMarkersCount)
			if resultForLine == maxPrior {
				m.Set(Cell{rowIdx, colIdx}, maxPrior)
				return
			}
			if resultForLine == maxPrior-1 {
				m.Set(Cell{rowIdx, colIdx}, maxPrior-1)
				return
			}
			countPossibleWins += resultForLine
		}
	}

	// Проверяем по второй диагонали
	{
		isOnSecondDiagonal := false
		j := size
		for i := 0; i < size; i++ {
			j--
			if i == rowIdx && j == colIdx {
				isOnSecondDiagonal = true
				break
			}
		}

		if isOnSecondDiagonal {
			cEmptyCells := 0
			enemyMarkersCount := 0
			ownMarkersCount := 0

			j := size
			for i := 0; i < size; i++ {
				j--
				if len(gameField[i][j]) == 0 {
					cEmptyCells++
				} else {
					if gameField[i][j] != ai.marker {
						enemyMarkersCount++
					} else {
						ownMarkersCount++
					}
				}
			}

			resultForLine := makeDecisionForLine(maxPrior, size, cEmptyCells, enemyMarkersCount, ownMarkersCount)
			if resultForLine == maxPrior {
				m.Set(Cell{rowIdx, colIdx}, maxPrior)
				return
			}
			if resultForLine == maxPrior-1 {
				m.Set(Cell{rowIdx, colIdx}, maxPrior-1)
				return
			}
			countPossibleWins += resultForLine
		}
	}

	m.Set(Cell{rowIdx, colIdx}, countPossibleWins)
}

func makeDecisionForLine(maxPrior, size, cEmptyCells, enemyMarkersCount, ownMarkersCount int) int {
	countPossibleWins := 0

	if cEmptyCells == size {
		countPossibleWins++
	}
	if size-ownMarkersCount == 1 {
		return maxPrior
	}
	if size-enemyMarkersCount == 1 {
		return maxPrior - 1
	}
	if enemyMarkersCount == 0 && ownMarkersCount > 0 {
		// Если в текущей линии меньше ходов до победы, то увеличиваем приоритет
		countPossibleWins += ownMarkersCount
	}

	return countPossibleWins
}
