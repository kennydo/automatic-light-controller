package app

import (
	"time"

	"github.com/kennydo/automatic-light-controller/lib/huebridge"
	"github.com/kennydo/automatic-light-controller/lib/scheduler"
	"go.uber.org/zap"
)

type App struct {
	config    *Config
	scheduler *scheduler.Scheduler
	huebridge *huebridge.HueBridge
	logger    *zap.Logger
}

func New(config *Config) (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	mHueBridge, err := huebridge.New(logger, config.HueBridge.IPAddress.String(), config.HueBridge.Username)
	if err != nil {
		return nil, err
	}

	mScheduler := scheduler.New(config.Location, config.Rules)

	return &App{
		config:    config,
		scheduler: mScheduler,
		huebridge: mHueBridge,
		logger:    logger,
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
			a.logger.Info("Sleeping until next scheduled action", zap.Float64("minutes", durationToSleep.Minutes()), zap.Any("action", scheduledAction))
			time.Sleep(durationToSleep)

			for _, groupName := range scheduledAction.Rule.LightGroups {
				conditionsAreSatisfied, err := a.conditionsAreSatisfied(groupName, scheduledAction.Rule.Conditions)
				if err != nil {
					break
				}

				if !conditionsAreSatisfied {
					a.logger.Info("Not executing rule because conditions are not satisfied", zap.Any("rule", scheduledAction.Rule), zap.String("group", groupName))
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
