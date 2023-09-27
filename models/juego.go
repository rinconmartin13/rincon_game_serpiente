package models

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2" // este es el tcell para usar la consola o bueno terminal cmd
)

const (
	Up = iota
	Left
	Right
	Down
)

type Apple Position

func (apple *Apple) Position() Position {
	return newPosition(apple.x, apple.y)
}

// Agregar un canal de control para detener el juego
type control struct {
	stop chan struct{}
}

// Agregar un campo de control al juego
type Game struct {
	mu        sync.Mutex
	direction int
	speed     time.Duration
	isStart   bool
	isOver    bool
	score     int
	Apple     *Apple
	Board     *Board
	Snake     *Snake
	screen    tcell.Screen
	sound     Sound
	ctrl      control
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Crea una nueva manzana.
func newApple(x int, y int) *Apple {
	return &Apple{x, y}
}

// Cambiar el tamaño de la pantalla si se cambia el terminal.
func (g *Game) resizeScreen() {
	g.mu.Lock()
	g.screen.Sync()
	g.mu.Unlock()
}

// Salir del juego.
func (g *Game) exit() {
	g.mu.Lock()
	g.screen.Fini()
	g.mu.Unlock()
	os.Exit(0)
}

// Establece la posición de la nueva manzana en el tablero.
func (g *Game) setNewApplePosition() {
	var availablePositions []Position

	for _, position := range g.Board.area {
		if !g.Snake.contains(position) {
			availablePositions = append(availablePositions, position)
		}
	}
	applePosition := availablePositions[rand.Intn(len(availablePositions))]
	g.Apple = newApple(applePosition.x, applePosition.y)
}

// Muestra la pantalla de carga.
func (g *Game) drawLoading() {
	if !g.hasStarted() {
		message := "PULSE <ENTER> PARA EMPEZAR CON EL JUEGO"
		textLength := len(message)
		x1 := (g.Board.width - textLength) / 2
		x2 := x1 + textLength
		y := g.Board.height / 2
		g.drawText(x1, y, x2, y, message)
	}
}

// Mostrar texto en la terminal.
func (g *Game) drawText(x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorCornflowerBlue)
	for _, r := range text {
		g.screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

// Muestra la pantalla final.
func (g *Game) drawEnding() {
	if g.hasEnded() {
		message := "Perdiste :("
		messageLength := len(message)
		x1 := g.Board.width/2 - messageLength/2
		// x2 := g.Board.width/2 + messageLength/2
		y := g.Board.height / 2
		style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
		// Dibuja el mensaje en la pantalla con el estilo personalizado
		for col, char := range message {
			g.screen.SetContent(x1+col, y, char, nil, style)
		}
	}
}

// Determina si el juego debe continuar.
func (g *Game) shouldContinue() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return !g.isOver && g.isStart
}

// Actualiza el estado del juego a terminado.
func (g *Game) over() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.isOver = true
	g.sound.GameOver()
}

// Comprueba si el juego ha terminado.
func (g *Game) hasEnded() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.isOver
}

// Muestra la manzana en el tablero.
func (g *Game) drawApple() {
	style := tcell.StyleDefault.Background(tcell.ColorGold).Foreground(tcell.ColorGold)
	g.screen.SetContent(g.Apple.x, g.Apple.y, '🍎', nil, style)
}

// Muestra el tablero de juego.
func (g *Game) drawBoard() {
	width, height := g.Board.width, g.Board.height
	boardStyle := tcell.StyleDefault.Background(tcell.ColorAqua).Foreground(tcell.ColorOrangeRed)
	g.screen.SetContent(0, 0, tcell.RuneULCorner, nil, boardStyle)
	for i := 1; i < width; i++ {
		g.screen.SetContent(i, 0, tcell.RuneHLine, nil, boardStyle)
	}
	g.screen.SetContent(width, 0, tcell.RuneURCorner, nil, boardStyle)

	for i := 1; i < height; i++ {
		g.screen.SetContent(0, i, tcell.RuneVLine, nil, boardStyle)
	}

	g.screen.SetContent(0, height, tcell.RuneLLCorner, nil, boardStyle)

	for i := 1; i < height; i++ {
		g.screen.SetContent(width, i, tcell.RuneVLine, nil, boardStyle)
	}

	g.screen.SetContent(width, height, tcell.RuneLRCorner, nil, boardStyle)

	for i := 1; i < width; i++ {
		g.screen.SetContent(i, height, tcell.RuneHLine, nil, boardStyle)
	}

	g.drawText(1, height+1, width, height+10, fmt.Sprintf("Puntuación:%d", g.score))
	g.drawText(1, height+3, width, height+10, "Pulsa ESC o Ctrl+C para salir del juego")
	g.drawText(1, height+4, width, height+10, "Pulsa las flechas para controlar la dirección de la serpiente ⬅⬆⬇➡")
}

// Muestra la serpiente.
func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorSkyblue)
	current := g.Snake.head

	for current != nil {
		g.screen.SetContent(current.x, current.y, tcell.RuneCkBoard, nil, snakeStyle)
		current = current.next
	}
}

// Actualiza el estado de los elementos del juego (serpiente y manzana).
func (g *Game) updateItemState() {
	if g.Snake.canMove(g.Board, g.direction) {
		g.Snake.move(g.direction)

		if g.Snake.CanEat(g.Apple) {
			g.sound.Hiss()
			g.Snake.Eat(g.Apple)
			g.score += 5 //incremen
			g.setNewApplePosition()
		}
	} else {
		g.over()
	}
}

// Comprueba si necesitamos cambiar la dirección según la nueva dirección.
func (g *Game) shouldUpdateDirection(direction int) bool {
	if g.direction == direction {
		return false
	}
	if g.direction == Left && direction != Right {
		return true
	}
	if g.direction == Up && direction != Down {
		return true
	}
	if g.direction == Down && direction != Up {
		return true
	}
	if g.direction == Right && direction != Left {
		return true
	}

	return false
}

// Actualiza la pantalla del juego.
func (g *Game) updateScreen() {
	g.screen.Clear()
	g.drawLoading()
	g.drawApple()
	g.drawBoard()
	g.drawSnake()
	g.drawEnding()
	g.screen.Show()
}

// Actualiza el estado del juego para comenzar.
func (g *Game) start() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.isStart = true
}

// Comprueba si el juego ha comenzado.
func (g *Game) hasStarted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.isStart
}

// Ejecuta el juego usando gorutien para controlar el movimiento de la serpiente.
func (g *Game) run(directionChan chan int) {
	ticker := time.NewTicker(g.speed)
	defer ticker.Stop()

	for {
		select {
		case newDirection := <-directionChan:
			if g.shouldUpdateDirection(newDirection) {
				g.mu.Lock()
				g.direction = newDirection
				g.mu.Unlock()
			}
		case <-ticker.C:
			if g.shouldContinue() {
				g.updateItemState()
			}
			g.updateScreen()
		}
	}
}

// Actualiza el estado del juego según los eventos del teclado.
func (g *Game) handleKeyBoardEvents(directionChan chan int) {
	defer close(directionChan)

	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventResize:
			g.resizeScreen()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				g.exit()
			}
			if !g.hasStarted() && event.Key() == tcell.KeyEnter {
				g.start()
			}

			if !g.hasEnded() {
				if event.Key() == tcell.KeyLeft {
					directionChan <- Left
				}
				if event.Key() == tcell.KeyRight {
					directionChan <- Right
				}
				if event.Key() == tcell.KeyDown {
					directionChan <- Down
				}
				if event.Key() == tcell.KeyUp {
					directionChan <- Up
				}
			}
		}
	}
}

// Crea un nuevo juego.
func newGame(board *Board, silent bool) *Game {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorLightGoldenrodYellow).Foreground(tcell.ColorWhiteSmoke)
	screen.SetStyle(defStyle)
	sound, err := NewSound(silent)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	game := &Game{
		Board:     board,
		Snake:     NewSnake(),
		direction: Up,
		speed:     time.Millisecond * 100,
		screen:    screen,
		sound:     sound,
	}

	game.setNewApplePosition()
	game.updateScreen()

	return game
}

// Agregar un método para esperar a que el juego termine
func (g *Game) waitForGameEnd() {
	<-g.ctrl.stop
}

// En el método StartGame
func StartGame(silent bool) {
	directionChan := make(chan int, 10)
	game := newGame(newBoard(70, 20), silent)

	// Crear una goroutine para la lógica del juego
	go game.run(directionChan)
	// Crear una goroutine para la detección de eventos del teclado
	go game.handleKeyBoardEvents(directionChan)
	// Crear una goroutine para la actualización de pantalla para que la pantalla se refresque de manera continua
	go game.updateScreen()
	// Esperar a que el juego termine
	game.waitForGameEnd()
}
