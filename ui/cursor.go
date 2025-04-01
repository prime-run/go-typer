package ui

type CursorType int

const (
	BlockCursor CursorType = iota
	UnderlineCursor
)

// Global default cursor type that can be set via command line
var DefaultCursorType CursorType = BlockCursor

type Cursor struct {
	style CursorType
}

func NewCursor(style CursorType) *Cursor {
	return &Cursor{
		style: style,
	}
}

func (c *Cursor) Render(char rune) string {
	switch c.style {
	case BlockCursor:
		return BlockCursorStyle.Render(string(char))
	case UnderlineCursor:
		return UnderlineCursorStyle.Render(string(char))
	default:
		return BlockCursorStyle.Render(string(char))
	}
}
