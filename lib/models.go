package lib

import (
	"fmt"
	"strconv"
)

// LightState describes the state of lights
type LightState struct {
	Brightness Brightness `toml:"brightness"`
}

// Brightness describes how bright something should be
type Brightness struct {
	Percent int
}

// UnmarshalText unmarshals text about brightness (in a format like "100%" or "0%")
func (b *Brightness) UnmarshalText(text []byte) error {
	var err error
	if text[len(text)-1] != '%' {
		err = fmt.Errorf("Brightness did not end in a percentage")
		return err
	}

	var percent int
	percent, err = strconv.Atoi(string(text[:len(text)-1]))
	if err != nil {
		return err
	}
	b.Percent = percent

	return nil
}
