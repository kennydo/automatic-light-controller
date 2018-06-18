package app

import (
	"net"

	"github.com/kennydo/automatic-light-controller/lib/scheduler"
)

type Config struct {
	Location  scheduler.LocationConfig `toml:"location"`
	HueBridge HueBridgeConfig          `toml:"hue_bridge"`
	Rules     []scheduler.Rule         `toml:"rules"`
}

type HueBridgeConfig struct {
	IPAddress net.IP `toml:"ip_address"`
	Username  string `tom:"username"`
}
