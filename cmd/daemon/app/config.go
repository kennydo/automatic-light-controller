package app

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Config struct {
	Timezone  Timezone        `toml:"timezone"`
	HueBridge HueBridgeConfig `toml:"hue_bridge"`
	Rules     []Rule          `toml:"rules"`
}

type HueBridgeConfig struct {
	IPAddress net.IP `toml:"ip_address"`
	Username  string `tom:"username"`
}

type Rule struct {
	Days        []Weekday   `toml:"days"`
	TimeTrigger TimeTrigger `toml:"time_trigger"`
	LightState  LightState  `toml:"light_state"`
	Conditions  []Condition `toml:"conditions"`
}

type TimeTrigger struct {
	SolarEvent string `toml:"solar_event"`
	LocalTime  string `toml:"local_time"`
}

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

type Timezone struct {
	time.Location
}

func (t *Timezone) UnmarshalText(text []byte) error {
	loc, err := time.LoadLocation(string(text))
	if err != nil {
		return err
	}

	t.Location = *loc
	return nil
}

type Weekday struct {
	time.Weekday
}

func (w *Weekday) UnmarshalText(text []byte) error {
	stringText := string(text)
	var weekday time.Weekday

	switch stringText {
	case "MO":
		weekday = time.Monday
	case "TU":
		weekday = time.Tuesday
	case "WE":
		weekday = time.Wednesday
	case "TH":
		weekday = time.Thursday
	case "FR":
		weekday = time.Friday
	case "SA":
		weekday = time.Saturday
	case "SU":
		weekday = time.Sunday
	default:
		return fmt.Errorf("Unrecognized weekday: %v", stringText)
	}

	w.Weekday = weekday

	return nil
}

type Condition struct {
	Type ConditionType `toml:"type"`
}

type ConditionType int

const (
	LightsAreOn ConditionType = iota
	LightsAreOff
)

var conditions = [...]string{
	"LightsAreOn",
	"LightsAreOff",
}

// String returns the English name of the day ("Sunday", "Monday", ...).
func (c ConditionType) String() string { return conditions[c] }

func (c *ConditionType) UnmarshalText(text []byte) error {
	stringText := string(text)
	var conditionType ConditionType

	switch stringText {
	case "lights_are_on":
		conditionType = LightsAreOn
	case "lights_are_off":
		conditionType = LightsAreOff
	default:
		return fmt.Errorf("Unrecognized condition type: %v", stringText)
	}

	c = &conditionType
	return nil
}
