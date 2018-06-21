package app

import (
	"fmt"
	"time"

	"github.com/kennydo/automatic-light-controller/lib/huebridge"
	"github.com/kennydo/automatic-light-controller/lib/scheduler"
)

type App struct {
	config    *Config
	scheduler *scheduler.Scheduler
	huebridge *huebridge.HueBridge
}

func New(config *Config) (*App, error) {
	mHueBridge, err := huebridge.New(config.HueBridge.IPAddress.String(), config.HueBridge.Username)
	if err != nil {
		return nil, err
	}

	mScheduler := scheduler.New(config.Location, config.Rules)

	return &App{
		config:    config,
		scheduler: mScheduler,
		huebridge: mHueBridge,
	}, nil
}

func (a *App) Run() error {
	var currentTime time.Time
	var err error

	for {
		currentTime = time.Now()
		scheduledActions := a.scheduler.GetNextScheduledActions(currentTime)

		for _, scheduledAction := range scheduledActions {
			currentTime = time.Now()
			durationToSleep := scheduledAction.ScheduledFor.Sub(currentTime)
			fmt.Printf("Sleeping %v minutes for scheduled action %+v\n", durationToSleep.Minutes(), scheduledAction)
			time.Sleep(durationToSleep)

			for _, groupName := range scheduledAction.Rule.LightGroups {
				conditionsAreSatisfied, err := a.conditionsAreSatisfied(groupName, scheduledAction.Rule.Conditions)
				if err != nil {
					break
				}

				if !conditionsAreSatisfied {
					fmt.Printf(
						"Not executing rule %+v on group %v because conditions are not satisfied\n",
						scheduledAction.Rule,
						groupName,
					)
					continue
				}

				err = a.huebridge.SetGroupLightState(groupName, scheduledAction.Rule.LightState)
				if err != nil {
					break
				}

			}

			if err != nil {
				break
			}
		}

		if err != nil {
			break
		}
	}

	return err
}

func (a *App) conditionsAreSatisfied(lightGroup string, conditions []scheduler.Condition) (bool, error) {
	satisfied := true

	for _, condition := range conditions {
		switch condition.Type.ConditionType {
		case scheduler.LightsAreOff:
			currentLightState, err := a.huebridge.GetGroupLightState(lightGroup)
			if err != nil {
				return false, err
			}
			satisfied = currentLightState.Brightness.Percent == 0
		case scheduler.LightsAreOn:
			currentLightState, err := a.huebridge.GetGroupLightState(lightGroup)
			if err != nil {
				return false, err
			}
			satisfied = currentLightState.Brightness.Percent > 0
		}

		if !satisfied {
			break
		}
	}

	return satisfied, nil
}
