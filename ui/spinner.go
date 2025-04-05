package ui

// Spinner provides a simple text-based spinner animation
type Spinner struct {
	frames []string
	index  int
}

// NewSpinner creates a new spinner with the default animation
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{
			"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷",
		},
		index: 0,
	}
}

// Update advances the spinner animation
func (s *Spinner) Update() {
	s.index = (s.index + 1) % len(s.frames)
}

// View returns the current frame of the spinner
func (s *Spinner) View() string {
	return s.frames[s.index]
}
