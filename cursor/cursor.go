package cursor

type CursorType int

const (
	Vertical CursorType = iota
	Horizontal
	Underline
	Arrow
	Full
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Cursor struct {
	Type CursorType

	// Positioning
	X        int
	Y        int
	Pointing Direction

	Blinking           bool
	BlinkTimeMs        int
	currentBlinkTimeMs int
}

func New() Cursor {
	cursor := Cursor{}
	cursor.Reset()
	return cursor
}

func (c *Cursor) Reset() {
	c.Type = Vertical
	c.X = 0
	c.Y = 0
	c.Blinking = false
	c.BlinkTimeMs = 500
	c.currentBlinkTimeMs = 0
}

func (c *Cursor) Up() {
	if c.Y > 0 {
		c.Y--
	}
}

func (c *Cursor) Down() {
	// TODO: Max out based on max view height
	c.Y++
}

func (c *Cursor) Left() {
	if c.X > 0 {
		c.X--
	}
}

func (c *Cursor) Right() {
	// TODO: See Down TODO
	c.X++
}
