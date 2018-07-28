package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kennydo/automatic-light-controller/lib"
)

// LocationConfig describes the coordinates of a geographical location and the timezone in that location
type LocationConfig struct {
	Timezone  Timezone `toml:"timezone"`
	Latitude  float64  `toml:"latitude"`
	Longitude float64  `toml:"longitude"`
}

// Timezone describes a timezone
type Timezone struct {
	*time.Location
}

// UnmarshalText converts the text form of a location to the Timezone object
func (t *Timezone) UnmarshalText(text []byte) error {
	loc, err := time.LoadLocation(string(text))
	if err != nil {
		return err
	}

	t.Location = loc
	return nil
}

// Rule describes under which conditions which lights should be set into which state
type Rule struct {
	Days        []Weekday      `toml:"days"`
	LightGroups []string       `toml:"light_groups"`
	TimeTrigger TimeTrigger    `toml:"time_trigger"`
	LightState  lib.LightState `toml:"light_state"`
	Conditions  []Condition    `toml:"conditions"`
}

// TimeTrigger describes either a solar event or a hard-coded local time for something to happen
type TimeTrigger struct {
	SolarEvent *SolarEventWrapper `toml:"solar_event"`
	LocalTime  *TimeInDay         `toml:"local_time"`
}

// TimeInDay describes a certain time in a day
type TimeInDay struct {
	Hour   int
	Minute int
}

// UnmarshalText converts text in "HH:MM" format into a TimeInDay
func (t *TimeInDay) UnmarshalText(text []byte) error {
	var err error
	stringText := string(text)

	if strings.Count(stringText, ":") != 1 {
		return fmt.Errorf("Time in day must contain exactly one \":\": %v", stringText)
	}

	elements := strings.Split(stringText, ":")
	t.Hour, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}

	t.Minute, err = strconv.Atoi(elements[1])
	return err
}

// String returns a string in the form "HH:MM" for a time in day
func (t *TimeInDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

// ForTime returns a timestamp where the year, month, day, and timezone are the same as the given object, but with the time of the argument
func (t *TimeInDay) ForTime(ti time.Time) time.Time {
	return time.Date(
		ti.Year(),
		ti.Month(),
		ti.Day(),
		t.Hour,
		t.Minute,
		0,
		0,
		ti.Location(),
	)
}

// SolarEventWrapper wraps a SolarEvent
type SolarEventWrapper struct {
	SolarEvent
}

// SolarEvent is an enum for solar events
type SolarEvent int

const (
	// Sunrise refers to the instant the sun rises
	Sunrise SolarEvent = iota
	// Sunset refers to the instant the sun sets
	Sunset
)

var solarEvents = [...]string{
	"Sunrise",
	"Sunset",
}

// String returns the English name of the condition type
func (s SolarEvent) String() string { return solarEvents[s] }

// UnmarshalText converts text to the appropriate SolarEvent enum
func (s *SolarEventWrapper) UnmarshalText(text []byte) error {
	stringText := string(text)
	var solarEvent SolarEvent

	switch stringText {
	case "sunrise":
		solarEvent = Sunrise
	case "sunset":
		solarEvent = Sunset
	default:
		return fmt.Errorf("Unrecognized solar event: %v", stringText)
	}

	s.SolarEvent = solarEvent
	return nil
}

// Weekday refers to day of week
type Weekday struct {
	time.Weekday
}

// UnmarshalText converts text into Weekday objects
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

// Condition describes things that must be true for a rule to take effect
type Condition struct {
	Type *ConditionTypeWrapper `toml:"type"`
}

// ConditionTypeWrapper wraps ConditionType
type ConditionTypeWrapper struct {
	ConditionType
}

// ConditionType is an enum representing types of conditions
type ConditionType int

const (
	// LightsAreOn is a condition type requiring the lights to be on
	LightsAreOn ConditionType = iota
	// LightsAreOff is a condition type requiring the lights to be off
	LightsAreOff
)

var conditions = [...]string{
	"LightsAreOn",
	"LightsAreOff",
}

// String returns the English name of the condition type
func (c ConditionType) String() string { return conditions[c] }

// UnmarshalText converts text to condition type enum
func (c *ConditionTypeWrapper) UnmarshalText(text []byte) error {
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

	c.ConditionType = conditionType
	return nil
}

// ScheduledRule represents a specific instance of a rule execution that should be executed at a certain time
type ScheduledRule struct {
	Rule         Rule
	ScheduledFor time.Time
}
