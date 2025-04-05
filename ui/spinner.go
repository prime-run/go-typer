package ui

type Spinner struct {
	frames []string
	index  int
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{
			"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷",
		},
		index: 0,
	}
}

func (s *Spinner) Update() {
	s.index = (s.index + 1) % len(s.frames)
}

func (s *Spinner) View() string {
	return s.frames[s.index]
}
