package huebridge

import (
	"log"

	"github.com/heatxsink/go-hue/groups"
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
