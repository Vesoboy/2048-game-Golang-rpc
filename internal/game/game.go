package game

import (
	"math/rand"
	"time"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Game struct {
	Grid      [][]int `json:"grid"`
	Score     int     `json:"score"`
	BestScore int     `json:"bestScore"`
	GameOver  bool    `json:"gameOver"`
	size      int
}

func NewGame() *Game {
	g := &Game{
		size: 4,
		Grid: make([][]int, 4),
	}
	for i := range g.Grid {
		g.Grid[i] = make([]int, 4)
	}
	g.addRandomTile()
	g.addRandomTile()
	return g
}

func (g *Game) Move(dir Direction) bool {
	if g.GameOver {
		return false
	}

	// Сохраняем старое состояние для сравнения
	oldGrid := make([][]int, g.size)
	for i := range oldGrid {
		oldGrid[i] = make([]int, g.size)
		copy(oldGrid[i], g.Grid[i])
	}

	// Выполняем ход
	switch dir {
	case Up:
		g.moveUp()
	case Right:
		g.moveRight()
	case Down:
		g.moveDown()
	case Left:
		g.moveLeft()
	}

	// Проверяем, было ли движение
	moved := false
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if oldGrid[i][j] != g.Grid[i][j] {
				moved = true
				break
			}
		}
		if moved {
			break
		}
	}

	// Если был ход, добавляем новую плитку
	if moved {
		g.addRandomTile()
		if g.Score > g.BestScore {
			g.BestScore = g.Score
		}
	}

	// Проверяем окончание игры
	g.GameOver = !g.canMove()

	return moved
}

// processLine обрабатывает одну линию (строку или столбец)
func (g *Game) processLine(line []int) ([]int, int) {
	// Удаляем нули и сжимаем числа
	var numbers []int
	for _, num := range line {
		if num != 0 {
			numbers = append(numbers, num)
		}
	}

	// Добавляем нули в конец
	for len(numbers) < len(line) {
		numbers = append(numbers, 0)
	}

	// Объединяем одинаковые числа
	scoreIncrease := 0
	for i := 0; i < len(numbers)-1; i++ {
		if numbers[i] != 0 && numbers[i] == numbers[i+1] {
			numbers[i] *= 2
			scoreIncrease += numbers[i]
			// Сдвигаем все числа влево
			copy(numbers[i+1:], numbers[i+2:])
			numbers[len(numbers)-1] = 0
		}
	}

	return numbers, scoreIncrease
}

func (g *Game) moveLeft() {
	for i := 0; i < g.size; i++ {
		row := make([]int, g.size)
		copy(row, g.Grid[i])
		newRow, score := g.processLine(row)
		copy(g.Grid[i], newRow)
		g.Score += score
	}
}

func (g *Game) moveRight() {
	for i := 0; i < g.size; i++ {
		row := make([]int, g.size)
		// Разворачиваем строку
		for j := 0; j < g.size; j++ {
			row[j] = g.Grid[i][g.size-1-j]
		}
		newRow, score := g.processLine(row)
		// Разворачиваем обратно
		for j := 0; j < g.size; j++ {
			g.Grid[i][j] = newRow[g.size-1-j]
		}
		g.Score += score
	}
}

func (g *Game) moveUp() {
	for j := 0; j < g.size; j++ {
		column := make([]int, g.size)
		for i := 0; i < g.size; i++ {
			column[i] = g.Grid[i][j]
		}
		newColumn, score := g.processLine(column)
		for i := 0; i < g.size; i++ {
			g.Grid[i][j] = newColumn[i]
		}
		g.Score += score
	}
}

func (g *Game) moveDown() {
	for j := 0; j < g.size; j++ {
		column := make([]int, g.size)
		for i := 0; i < g.size; i++ {
			column[i] = g.Grid[g.size-1-i][j]
		}
		newColumn, score := g.processLine(column)
		for i := 0; i < g.size; i++ {
			g.Grid[i][j] = newColumn[g.size-1-i]
		}
		g.Score += score
	}
}

func (g *Game) addRandomTile() {
	// Находим все пустые клетки
	emptyCells := make([][2]int, 0)
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.Grid[i][j] == 0 {
				emptyCells = append(emptyCells, [2]int{i, j})
			}
		}
	}

	if len(emptyCells) > 0 {
		// Выбираем случайную пустую клетку
		rand.Seed(time.Now().UnixNano())
		cell := emptyCells[rand.Intn(len(emptyCells))]
		// С вероятностью 0.9 добавляем 2, иначе 4
		if rand.Float64() < 0.9 {
			g.Grid[cell[0]][cell[1]] = 2
		} else {
			g.Grid[cell[0]][cell[1]] = 4
		}
	}
}

func (g *Game) canMove() bool {
	// Проверяем наличие пустых клеток
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.Grid[i][j] == 0 {
				return true
			}
		}
	}

	// Проверяем возможность объединения по горизонтали
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size-1; j++ {
			if g.Grid[i][j] == g.Grid[i][j+1] {
				return true
			}
		}
	}

	// Проверяем возможность объединения по вертикали
	for i := 0; i < g.size-1; i++ {
		for j := 0; j < g.size; j++ {
			if g.Grid[i][j] == g.Grid[i+1][j] {
				return true
			}
		}
	}

	return false
}
