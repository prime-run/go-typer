package ui

type CursorType int // CursorType represents the type of cursor to be used in the UI.

const (
	BlockCursor     CursorType = iota // Block cursor
	UnderlineCursor                   // Underline cursor
)

var DefaultCursorType CursorType = BlockCursor // Default cursor type

type Cursor struct {
	style CursorType
}

// NewCursor creates a new cursor with the specified style.
func NewCursor(style CursorType) *Cursor {
	return &Cursor{
		style: style,
	}
}

// render returns the string representation of the cursor based on its style.
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
