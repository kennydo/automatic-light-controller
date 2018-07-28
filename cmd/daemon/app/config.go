package app

import (
	"net"

	"github.com/kennydo/automatic-light-controller/lib/scheduler"
)

// Config holds the config for the app
type Config struct {
	Location  scheduler.LocationConfig `toml:"location"`
	HueBridge HueBridgeConfig          `toml:"hue_bridge"`
	Rules     []scheduler.Rule         `toml:"rules"`
}

// HueBridgeConfig holds information about what the bridge's IP is and how to authenticate to it
type HueBridgeConfig struct {
	IPAddress net.IP `toml:"ip_address"`
	Username  string `tom:"username"`
}
