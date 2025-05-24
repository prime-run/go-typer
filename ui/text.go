package ui

import (
	"strings"
	"time"

	devlog "github.com/prime-run/go-typer/log"
)

type Text struct {
	words      []*Word
	cursorPos  int
	showCursor bool
	cursorType CursorType
	sourceText string
}

func NewText(text string) *Text {

	estimatedWordCount := len(text)/6 + 1
	words := make([]*Word, 0, estimatedWordCount)
	var currentWord []rune

	for _, r := range text {
		if r == ' ' {
			if len(currentWord) > 0 {
				words = append(words, NewWord(currentWord))
				currentWord = make([]rune, 0, 8)
			}
			words = append(words, NewWord([]rune{' '}))
		} else {
			currentWord = append(currentWord, r)
		}
	}

	if len(currentWord) > 0 {
		words = append(words, NewWord(currentWord))
	}

	t := &Text{
		words:      words,
		cursorPos:  0,
		showCursor: true,
		cursorType: UnderlineCursor,
		sourceText: text,
	}

	if len(t.words) > 0 {
		t.words[0].SetActive(true)
	}

	return t
}

func (t *Text) CurrentWord() *Word {
	if t.cursorPos >= len(t.words) {
		return nil
	}
	return t.words[t.cursorPos]
}

func (t *Text) Type(r rune) {
	if t.cursorPos >= len(t.words) {
		return
	}

	currentWord := t.words[t.cursorPos]

	if currentWord.IsSpace() {
		if r == ' ' {
			currentWord.Type(r)
			if t.cursorPos < len(t.words)-1 {
				currentWord.SetActive(false)
				t.cursorPos++
				t.words[t.cursorPos].SetActive(true)
			}
		} else {
			currentWord.Type(r)
		}
		return
	}

	if r == ' ' {
		if !currentWord.HasStarted() {
			return
		}

		if !currentWord.IsComplete() {
			currentWord.Skip()
		}

		if t.cursorPos < len(t.words)-2 {
			currentWord.SetActive(false)
			t.cursorPos += 2
			t.words[t.cursorPos].SetActive(true)
		} else if t.cursorPos < len(t.words)-1 {
			currentWord.SetActive(false)
			t.cursorPos++
			t.words[t.cursorPos].SetActive(true)
		} else {
			currentWord.SetActive(false)
			if !currentWord.IsComplete() {
				currentWord.Skip()
			}
		}
	} else {
		currentWord.Type(r)

		if currentWord.IsComplete() {
			if t.cursorPos < len(t.words)-1 {
				nextWord := t.words[t.cursorPos+1]
				currentWord.SetActive(false)
				t.cursorPos++
				nextWord.SetActive(true)
			} else {
				currentWord.SetActive(false)
				if !currentWord.IsComplete() {
					currentWord.Skip()
				}
			}
		}
	}
}

func (t *Text) Backspace() {
	if t.cursorPos >= len(t.words) {
		return
	}

	currentWord := t.words[t.cursorPos]
	if !currentWord.Backspace() && t.cursorPos > 0 {
		currentWord.SetActive(false)
		t.cursorPos--
		currentWord = t.words[t.cursorPos]
		currentWord.SetActive(true)

		if currentWord.IsSpace() {
			currentWord.Backspace()
			if t.cursorPos > 0 {
				currentWord.SetActive(false)
				t.cursorPos--
				t.words[t.cursorPos].SetActive(true)
			}
		}
	}
}

func (t *Text) Render() string {
	startTime := time.Now()
	devlog.Log("Text: Render started")

	estimatedSize := 0
	for _, word := range t.words {
		estimatedSize += len(word.target) * 3
	}

	var result strings.Builder
	result.Grow(estimatedSize)

	showCursor := t.showCursor
	if t.cursorType == UnderlineCursor {
		showCursor = true
	}

	for _, word := range t.words {
		result.WriteString(word.Render(showCursor))
	}

	rendered := TextContainerStyle.Render(result.String())

	renderTime := time.Since(startTime)
	devlog.Log("Text: Render completed in %s, length: %d", renderTime, len(rendered))

	return rendered
}

func (t *Text) Update() {
	t.showCursor = true
}

func (t *Text) SetCursorType(cursorType CursorType) {
	t.cursorType = cursorType
	for _, word := range t.words {
		word.SetCursorType(cursorType)
	}
}

func (t *Text) IsComplete() bool {
	for _, word := range t.words {
		if !word.IsComplete() {
			return false
		}
	}
	return true
}

func (t *Text) GetCursorPos() int {
	return t.cursorPos
}

func (t *Text) Stats() (total, correct, errors int) {
	for _, word := range t.words {
		if word.IsSpace() {
			continue
		}

		if word.state == Perfect {
			correct++
		} else if word.state == Error {
			errors++
		}
		total++
	}
	return
}

func (t *Text) GetText() string {
	if t.sourceText != "" {
		return t.sourceText
	}

	var builder strings.Builder
	builder.Grow(len(t.words) * 8)

	for _, word := range t.words {
		if !word.IsSpace() {
			builder.WriteString(string(word.target))
		} else {
			builder.WriteRune(' ')
		}
	}
	return builder.String()
}
