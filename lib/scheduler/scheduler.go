package scheduler

import (
	"sort"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

// Scheduler determines when the next scheduled actions are for a given set of rules in a given location
type Scheduler struct {
	rules          []Rule
	locationConfig LocationConfig
}

// New creates a new Scheduler
func New(locationConfig LocationConfig, rules []Rule) *Scheduler {
	return &Scheduler{
		locationConfig: locationConfig,
		rules:          rules,
	}
}

// GetNextScheduledRules returns an ordered listing of all rules, ordered by soonest scheduled time
func (s *Scheduler) GetNextScheduledRules(currentTime time.Time) []ScheduledRule {
	currentTimeInLocalTZ := currentTime.In(s.locationConfig.Timezone.Location)
	tomorrowInLocalTZ := currentTimeInLocalTZ.AddDate(0, 0, 1)

	nextSunrise, nextSunset := s.GetNextSunriseAndSunset(currentTime)

	scheduledRules := make([]ScheduledRule, len(s.rules))

	for i, rule := range s.rules {
		// Get the next time that this rule is supposed to happen
		var scheduledFor time.Time
		if rule.TimeTrigger.LocalTime != nil {
			localTimeToday := rule.TimeTrigger.LocalTime.ForTime(currentTimeInLocalTZ)
			if localTimeToday.After(currentTime) {
				scheduledFor = localTimeToday
			} else {
				scheduledFor = rule.TimeTrigger.LocalTime.ForTime(tomorrowInLocalTZ)
			}
		} else if rule.TimeTrigger.SolarEvent != nil {
			switch rule.TimeTrigger.SolarEvent.SolarEvent {
			case Sunrise:
				scheduledFor = nextSunrise
			case Sunset:
				scheduledFor = nextSunset
			}
		}

		scheduledRules[i] = ScheduledRule{
			Rule:         rule,
			ScheduledFor: scheduledFor.In(s.locationConfig.Timezone.Location),
		}
	}

	sort.SliceStable(
		scheduledRules,
		func(i, j int) bool {
			return scheduledRules[i].ScheduledFor.Before(scheduledRules[j].ScheduledFor)
		},
	)

	return scheduledRules
}

// GetNextSunriseAndSunset returns the next sunrise and sunset after a given time
func (s *Scheduler) GetNextSunriseAndSunset(currentTime time.Time) (time.Time, time.Time) {
	var nextSunrise, nextSunset *time.Time

	todaySunrise, todaySunset := sunrise.SunriseSunset(
		s.locationConfig.Latitude, s.locationConfig.Longitude,
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
	)

	if todaySunrise.After(currentTime) {
		nextSunrise = &todaySunrise
	}
	if todaySunset.After(currentTime) {
		nextSunset = &todaySunset
	}

	if nextSunrise != nil && nextSunset != nil {
		return *nextSunrise, *nextSunset
	}

	tomorrow := currentTime.AddDate(0, 0, 1)
	tomorrowSunrise, tomorrowSunset := sunrise.SunriseSunset(
		s.locationConfig.Latitude, s.locationConfig.Longitude,
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
	)
	if nextSunrise == nil {
		nextSunrise = &tomorrowSunrise
	}
	if nextSunset == nil {
		nextSunset = &tomorrowSunset
	}
	return *nextSunrise, *nextSunset
}
