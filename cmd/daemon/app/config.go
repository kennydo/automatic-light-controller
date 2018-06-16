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
	Days        []string    `toml:"days"`
	TimeTrigger TimeTrigger `toml:"time_trigger"`
	LightState  LightState  `toml:"light_state"`
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
