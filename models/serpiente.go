package models

type node struct {
	Position
	next *node
}

type Snake struct {
	head *node
}

type Position struct {
	x, y int
}

// Crea una nueva posición.
func newPosition(x, y int) Position {
	return Position{x, y}
}

// La siguiente posición de la cabeza de serpiente en el tablero.
func (s *Snake) nextHeadPosition(direction int) Position {
	switch direction {
	case Up:
		return newPosition(s.head.x, s.head.y-1)
	case Left:
		return newPosition(s.head.x-1, s.head.y)
	case Right:
		return newPosition(s.head.x+1, s.head.y)
	case Down:
		return newPosition(s.head.x, s.head.y+1)
	default:
		panic("error: Dirección inválida")
	}
}

// Comprueba si la serpiente puede moverse en la dirección indicada.
func (s *Snake) canMove(board *Board, direction int) bool {
	position := s.nextHeadPosition(direction)
	// Si el cuerpo actual contiene la siguiente posición del encabezado, devuelve falso.
	if s.contains(position) {
		return false
	}

	switch direction {
	case Up:
		return position.y > 0
	case Left:
		return position.x > 0
	case Right:
		return position.x < board.width
	case Down:
		return position.y < board.height
	}

	return true
}

// Comprueba si la serpiente ya ha tenido la posición.
func (s *Snake) contains(position Position) bool {
	current := s.head

	for current != nil {
		if current.x == position.x && current.y == position.y {
			return true
		}
		current = current.next
	}

	return false
}

// Comprueba si la serpiente puede comerse la manzana.
func (s *Snake) CanEat(apple *Apple) bool {
	return s.head.x == apple.x && s.head.y == apple.y
}

// La serpiente Come la manzana y agrega la posición de la manzana al corte del cuerpo de la serpiente.
func (s *Snake) Eat(apple *Apple) {
	s.add(apple.Position())
}

// Agrega una nueva posición a la cabeza.
func (s *Snake) add(position Position) *Snake {
	s.head = &node{
		Position: position,
		next:     s.head,
	}
	return s
}

// Mueve la serpiente según la dirección.
func (s *Snake) move(direction int) {
	position := s.nextHeadPosition(direction)
	s.add(position)

	if s.head.next.next != nil {
		current := s.head
		for ; current.next.next != nil; current = current.next {
		}
		current.next = nil
	}
}

// Crea una nueva serpiente con posición y longitud predeterminadas.
func NewSnake() *Snake {
	snake := &Snake{}

	return snake.add(newPosition(9, 10)).
		add(newPosition(9, 10)).
		add(newPosition(9, 9)).
		add(newPosition(9, 8)).
		add(newPosition(9, 7))
}
