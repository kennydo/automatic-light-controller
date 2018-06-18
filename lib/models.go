package lib

import (
	"fmt"
	"strconv"
)

type LightState struct {
	Brightness Brightness `toml:"brightness"`
}

type Brightness struct {
	Percent int
}

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
