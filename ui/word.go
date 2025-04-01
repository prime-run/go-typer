package ui

import (
	"strings"
	"time"
)

type WordState int

const (
	Untyped WordState = iota
	Perfect
	Imperfect
	Error
)

type Word struct {
	target []rune
	typed  []rune
	state  WordState
	active bool
	cursor *Cursor
}

func NewWord(target []rune) *Word {
	targetCopy := make([]rune, len(target))
	copy(targetCopy, target)

	return &Word{
		target: targetCopy,
		typed:  make([]rune, 0, len(target)),
		state:  Untyped,
		active: false,
		cursor: NewCursor(DefaultCursorType),
	}
}

func (w *Word) Type(r rune) {
	if w.IsSpace() {
		if r == ' ' {
			w.typed = []rune{' '}
			w.state = Perfect
		} else {
			w.typed = []rune{r}
			w.state = Error
		}
		return
	}

	if len(w.typed) < len(w.target) {
		w.typed = append(w.typed, r)
	} else if len(w.typed) == len(w.target) {
		w.typed[len(w.typed)-1] = r
	}

	w.updateState()
}

func (w *Word) Skip() {
	if len(w.typed) == 0 {
		for i := 0; i < len(w.target); i++ {
			w.typed = append(w.typed, '\x00')
		}
	} else {
		currentLen := len(w.typed)
		for i := currentLen; i < len(w.target); i++ {
			w.typed = append(w.typed, '\x00')
		}
	}
	w.state = Error
}

func (w *Word) Backspace() bool {
	if len(w.typed) == 0 {
		return false
	}
	w.typed = w.typed[:len(w.typed)-1]
	w.updateState()
	return true
}

func (w *Word) updateState() {
	if len(w.typed) == 0 {
		w.state = Untyped
		return
	}

	if w.IsSpace() {
		if len(w.typed) == 1 && w.typed[0] == ' ' {
			w.state = Perfect
		} else {
			w.state = Error
		}
		return
	}

	for _, r := range w.typed {
		if r == '\x00' {
			w.state = Error
			return
		}
	}

	minLen := min(len(w.typed), len(w.target))
	perfect := true
	for i := 0; i < minLen; i++ {
		if w.typed[i] != w.target[i] {
			perfect = false
			break
		}
	}

	if perfect && len(w.typed) == len(w.target) {
		w.state = Perfect
	} else if perfect && len(w.typed) < len(w.target) {
		w.state = Imperfect
	} else {
		w.state = Error
	}
}

func (w *Word) IsComplete() bool {
	return len(w.typed) >= len(w.target)
}

func (w *Word) HasStarted() bool {
	return len(w.typed) > 0
}

func (w *Word) IsSpace() bool {
	return len(w.target) == 1 && w.target[0] == ' '
}

func (w *Word) SetActive(active bool) {
	w.active = active
}

func (w *Word) SetCursorType(cursorType CursorType) {
	w.cursor = NewCursor(cursorType)
}

func (w *Word) Render(showCursor bool) string {
	startTime := time.Now()

	var result strings.Builder

	if w.IsSpace() {
		if len(w.typed) == 0 {
			if showCursor && w.active {
				return w.cursor.Render(' ')
			}
			return DimStyle.Render(" ")
		} else if len(w.typed) == 1 && w.typed[0] == ' ' {
			return InputStyle.Render(" ")
		} else {
			return ErrorStyle.Render(string(w.typed[0]))
		}
	}

	targetLen := len(w.target)
	typedLen := len(w.typed)

	for i := 0; i < max(targetLen, typedLen); i++ {
		if showCursor && w.active && i == typedLen {
			if i < targetLen {
				result.WriteString(w.cursor.Render(w.target[i]))
			} else {
				result.WriteString(w.cursor.Render(' '))
			}
			continue
		}

		if i >= typedLen {
			result.WriteString(DimStyle.Render(string(w.target[i])))
			continue
		}

		if i >= targetLen {
			result.WriteString(ErrorStyle.Render(string(w.typed[i])))
			continue
		}

		if w.typed[i] == '\x00' {
			result.WriteString(DimStyle.Render(string(w.target[i])))
			continue
		}

		if w.typed[i] == w.target[i] {
			if w.state == Error {
				result.WriteString(PartialErrorStyle.Render(string(w.target[i])))
			} else {
				result.WriteString(InputStyle.Render(string(w.target[i])))
			}
		} else {
			result.WriteString(ErrorStyle.Render(string(w.typed[i])))
		}
	}

	rendered := result.String()

	// Only log active word rendering for performance analysis
	if w.active {
		renderTime := time.Since(startTime)
		DebugLog("Word: Active word render completed in %s, length: %d", renderTime, len(rendered))
	}

	return rendered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
