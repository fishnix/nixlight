package light

import (
	"log"
)

// Dummy is a dummy light controller
type Dummy struct {
	Name string
}

// Start initializes the dummy light controller
func (d *Dummy) Start() error {
	log.Printf("Starting dummy light %s controller\n", d.Name)
	return nil
}

// SetBrightness sets the brightness level of the dummy light
func (d *Dummy) SetBrightness(b int) error {
	log.Printf("Setting dummy light %s's brightness to %d\n", d.Name, b)
	return nil
}
