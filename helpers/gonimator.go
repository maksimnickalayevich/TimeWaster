package helpers

import "time"

// Gonimator wipes out terminal screen correctly,
// without scrolling it. Handles simple animation for terminal
type Gonimator struct {
	tick time.Time
}

// ClearTerminal removes everything from terminal.
// doesn't use scrolling, actually wipes out terminal window
func (gon *Gonimator) ClearTerminal() {

}
