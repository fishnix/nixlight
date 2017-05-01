package light

// Controller is an interace that has a method to Start() the
// controller and to set the Brightness(int) of a light
type Controller interface {
	Start() error
	Brightness(int) error
}
