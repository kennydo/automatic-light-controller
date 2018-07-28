package app

import (
	"time"

	"github.com/kennydo/automatic-light-controller/lib/huebridge"
	"github.com/kennydo/automatic-light-controller/lib/scheduler"
	"go.uber.org/zap"
)

// App is application that schedules and executes Hue light changes
type App struct {
	config    *Config
	scheduler *scheduler.Scheduler
	huebridge *huebridge.HueBridge
	logger    *zap.Logger
}

// New creates a new app
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

// Run is the main loop
func (a *App) Run() error {
	var currentTime time.Time

	for {
		currentTime = time.Now()
		scheduledRules := a.scheduler.GetNextScheduledRules(currentTime)

		for _, scheduledRule := range scheduledRules {
			currentTime = time.Now()
			durationToSleep := scheduledRule.ScheduledFor.Sub(currentTime)
			a.logger.Info("Sleeping until next scheduled rule", zap.Float64("minutes", durationToSleep.Minutes()), zap.Any("rule", scheduledRule))
			time.Sleep(durationToSleep)

			for _, groupName := range scheduledRule.Rule.LightGroups {
				conditionsAreSatisfied, err := a.conditionsAreSatisfied(groupName, scheduledRule.Rule.Conditions)
				if err != nil {
					return err
				}

				if !conditionsAreSatisfied {
					a.logger.Info("Not executing rule because conditions are not satisfied", zap.Any("rule", scheduledRule.Rule), zap.String("group", groupName))
					continue
				}

				err = a.huebridge.SetGroupLightState(groupName, scheduledRule.Rule.LightState)
				if err != nil {
					return err
				}
			}
		}
	}
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
