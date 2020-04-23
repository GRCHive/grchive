package main

import (
	"github.com/teambition/rrule-go"
	"gitlab.com/grchive/grchive/core"
	"time"
)

type Schedule struct {
	lastRunTime *time.Time
	nextRunTime *time.Time

	recurrence *rrule.Set
	timezone   string
}

func (s Schedule) ShouldRun(c core.Clock) bool {
	if s.nextRunTime == nil {
		return false
	}
	return c.Now().UTC().After(*s.nextRunTime)
}

func (s *Schedule) MarkRun(c core.Clock) {
	now := c.Now().UTC()
	s.lastRunTime = &now
	s.nextRunTime = nil

	if s.recurrence != nil {
		next := s.recurrence.After(now, false)
		s.nextRunTime = &next
	}
}

func (s Schedule) HasNextRun(c core.Clock) bool {
	return (s.nextRunTime != nil)
}

func CreateOneTimeJobSchedule(s *core.ScheduledTaskOneTime, c core.Clock) (*Schedule, error) {
	if c.Now().After(s.EventDateTime) {
		return nil, nil
	}

	return &Schedule{
		lastRunTime: nil,
		nextRunTime: &s.EventDateTime,
		recurrence:  nil,
	}, nil
}

func CreateRecurringJobSchedule(s *core.ScheduledTaskRecurrence, c core.Clock) (*Schedule, error) {
	// Recreate the rrule.Set since it won't have the
	// right timezone unless we do that since we can't be sure
	// whether this rule set came from a place where we can trust
	// that it's been re-created properly.
	loc, err := time.LoadLocation(s.Timezone)
	if err != nil {
		return nil, err
	}

	recurrence, err := core.RebuildRRuleSet(s.RRule, func(opt *rrule.ROption) {
		opt.Dtstart = opt.Dtstart.In(loc)
		opt.RFC = false
	})

	if err != nil {
		return nil, err
	}

	next := recurrence.After(c.Now(), true)
	return &Schedule{
		lastRunTime: nil,
		nextRunTime: &next,
		recurrence:  recurrence,
		timezone:    s.Timezone,
	}, nil
}
