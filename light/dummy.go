package light

import (
	"log"
)

type Dummy struct {
	Name string
}

func (d *Dummy) Start() error {
	log.Printf("Starting dummy light %s controller\n", d.Name)

	if l, ok := interface{}(d).(Controller); ok {
		log.Printf("It looks like %+v satisfies the Controller interface!", l)
	}

	return nil
}

func (d *Dummy) Brightness(b int) error {
	log.Printf("Setting dummy light %s's brightness to %d\n", d.Name, b)
	return nil
}
