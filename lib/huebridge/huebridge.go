package huebridge

import (
	"fmt"
	"math"

	"github.com/heatxsink/go-hue/groups"
	"github.com/heatxsink/go-hue/lights"
	"github.com/kennydo/automatic-light-controller/lib"
	"go.uber.org/zap"
)

type HueBridge struct {
	Hostname        string
	Username        string
	groupIDByName   map[string]int
	groupController *groups.Groups
	logger          *zap.Logger
}

func New(logger *zap.Logger, hostname string, username string) (*HueBridge, error) {
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
	logger.Info("Fetched Hue groups", zap.Any("groups", groupIDByName))

	return &HueBridge{
		Hostname:        hostname,
		Username:        username,
		groupController: groupController,
		groupIDByName:   groupIDByName,
		logger:          logger,
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

	b.logger.Info("Setting group to desired state", zap.Int("groupID", groupID), zap.Any("desiredState", desiredState))

	response, err := b.groupController.SetGroupState(groupID, desiredState)
	if err != nil {
		return err
	}

	b.logger.Info("Got response from Hue", zap.Any("response", response))

	return nil
}

func (b *HueBridge) GetGroupLightState(groupName string) (*lib.LightState, error) {
	var err error

	groupID, ok := b.groupIDByName[groupName]
	if !ok {
		return nil, fmt.Errorf("Did not recognize Hue group name: %v", groupName)
	}

	group, err := b.groupController.GetGroup(groupID)
	if err != nil {
		return nil, err
	}

	var percentage int

	if group.Action.On {
		percentage = int(math.Floor((float64(group.Action.Bri) / 254) * 100))
	} else {
		percentage = 0
	}

	return &lib.LightState{
		Brightness: lib.Brightness{
			Percent: percentage,
		},
	}, nil
}
