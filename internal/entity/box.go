package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kathleenfrench/common/gui"
)

// Box represents a machine
type Box struct {
	Name        string `boldholdKey:"Name"`
	IP          string
	Completed   bool `boltholdIndex:"Completed"` // when root + user flags captured
	Active      bool `boltholdIndex:"Active"`
	Hostname    string
	OS          string `boltholdIndex:"OS"`
	Difficulty  string
	Notes       string
	Description string
	Up          bool
	Flags       Flags
	Created     time.Time
	LastUpdated time.Time
}

// Flags represents the capture flags on a box
type Flags struct {
	Root string
	User string
}

// Validate validates whether a new box has all of the required fields
func (b Box) Validate() error {
	if b.Name == "" {
		return errors.New("setting a name for the box is required")
	}

	if !govalidator.IsIP(b.IP) {
		gui.Warn("invalid IP address", b.IP)
		b.IP = gui.InputPromptWithResponse("what is its IP?", "", true)
		if !govalidator.IsIP(b.IP) {
			return errors.New(b.IP + " is not a valid IP address")
		}
	}

	return nil
}

// BucketName returns the expected bucket name for a box
func (b *Box) BucketName() string {
	return "Box"
}
