package huebridge

import (
	"fmt"
	"log"
	"math"

	"github.com/heatxsink/go-hue/groups"
	"github.com/heatxsink/go-hue/lights"
	"github.com/kennydo/automatic-light-controller/lib"
)

type HueBridge struct {
	Hostname        string
	Username        string
	groupIDByName   map[string]int
	groupController *groups.Groups
}

func New(hostname string, username string) (*HueBridge, error) {
	groupController := groups.New(hostname, username)

	// Get the mapping of group ID by name
	allGroups, err := groupController.GetAllGroups()
	if err != nil {
		return nil, err
	}
	groupIDByName := make(map[string]int)
	for _, group := range allGroups {
		groupIDByName[group.Name] = group.ID
	}
	log.Printf("Fetched these Hue groups: %v", groupIDByName)

	return &HueBridge{
		Hostname:        hostname,
		Username:        username,
		groupController: groupController,
		groupIDByName:   groupIDByName,
	}, nil
}

func (b *HueBridge) SetGroupLightState(groupName string, lightState lib.LightState) error {
	var err error

	groupID, ok := b.groupIDByName[groupName]
	if !ok {
		return fmt.Errorf("Did not recognize Hue group name: %v", groupName)
	}

	desiredState := lights.State{
		On:  lightState.Brightness.Percent > 0,
		Bri: uint8(math.Ceil(254 * (float64(lightState.Brightness.Percent) / 100.0))),
	}

	log.Printf("Going to set group %v to %+v", groupID, desiredState)

	response, err := b.groupController.SetGroupState(groupID, desiredState)
	if err != nil {
		return err
	}

	log.Printf("Response from Hue: %v", response)

	return nil
}
