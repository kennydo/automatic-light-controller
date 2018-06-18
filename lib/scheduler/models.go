package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kennydo/automatic-light-controller/lib"
)

type LocationConfig struct {
	Timezone  Timezone `toml:"timezone"`
	Latitude  float64  `toml:"latitude"`
	Longitude float64  `toml:"longitude"`
}

type Timezone struct {
	*time.Location
}

func (t *Timezone) UnmarshalText(text []byte) error {
	loc, err := time.LoadLocation(string(text))
	if err != nil {
		return err
	}

	t.Location = loc
	return nil
}

type Rule struct {
	Days        []Weekday      `toml:"days"`
	LightGroups []string       `toml:"light_groups"`
	TimeTrigger TimeTrigger    `toml:"time_trigger"`
	LightState  lib.LightState `toml:"light_state"`
	Conditions  []Condition    `toml:"conditions"`
}

type TimeTrigger struct {
	SolarEvent *SolarEventWrapper `toml:"solar_event"`
	LocalTime  *TimeInDay         `toml:"local_time"`
}

type TimeInDay struct {
	Hour   int
	Minute int
}

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
	if err != nil {
		return err
	}
	return nil
}

func (t *TimeInDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

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

type SolarEventWrapper struct {
	SolarEvent
}

type SolarEvent int

const (
	Sunrise SolarEvent = iota
	Sunset
)

var solarEvents = [...]string{
	"Sunrise",
	"Sunset",
}

// String returns the English name of the condition type
func (s SolarEvent) String() string { return solarEvents[s] }

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
	Type *ConditionTypeWrapper `toml:"type"`
}

type ConditionTypeWrapper struct {
	ConditionType
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

// String returns the English name of the condition type
func (c ConditionType) String() string { return conditions[c] }

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

type ScheduledAction struct {
	Rule         Rule
	ScheduledFor time.Time
}
