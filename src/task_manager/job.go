package main

import (
	"errors"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"math"
	"time"
)

type JobTaskHandler interface {
	Tick(c core.Clock) error
}

type Job struct {
	id       int64
	schedule *Schedule
	typ      core.TaskType

	handler  JobTaskHandler
	lastTick time.Time

	backoffCount int
	backoffTime  time.Duration
}

func (j Job) Id() int64 {
	return j.id
}

// Returns whether or not the job is currently "backed off".
func (j *Job) TickBackoff(c core.Clock) bool {
	if j.backoffCount == 0 {
		return false
	}

	elapsed := c.Now().Sub(j.lastTick)
	j.backoffTime = j.backoffTime - elapsed
	if j.backoffTime > 0 {
		return true
	}

	return false
}

func (j *Job) Backoff() {
	j.backoffCount = j.backoffCount + 1

	// Each backoff should take longer than the last.
	// Don't want to backoff for more than 5 minutes.
	// Use the function y = e^{0.5x}. This reaches max backoff in ~11 tries.
	backoffSeconds := int64(math.Min(
		math.Max(
			math.Exp(0.5*float64(j.backoffCount)),
			0,
		),
		300,
	))

	j.backoffTime = time.Duration(backoffSeconds) * time.Second
}

// Returns a boolean that indicates whether this job will want to run again.
func (j *Job) Tick(c core.Clock, force bool) (bool, error) {
	defer func() {
		j.lastTick = c.Now()
	}()

	if j.TickBackoff(c) {
		return true, nil
	}

	if !j.schedule.ShouldRun(c) && !force {
		return true, nil
	}

	core.Info(fmt.Sprintf("\tRunning Job %d", j.Id()))
	err := j.handler.Tick(c)
	if err != nil {
		// If the job fails we don't want to immediately retry.
		j.Backoff()
		return true, err
	} else {
		j.backoffCount = 0
		j.backoffTime = time.Duration(0)
	}

	j.schedule.MarkRun(c)
	return j.schedule.HasNextRun(c), nil
}

func CreateJobFromTaskMetadata(task core.ScheduledTaskMetadata, schedule *Schedule) (*Job, error) {
	job := Job{
		id:           task.Id,
		schedule:     schedule,
		typ:          task.TaskType,
		backoffCount: 0,
		backoffTime:  time.Duration(0),
		lastTick:     time.Now(),
	}

	switch task.TaskType {
	case core.KGrchiveApiTask:
		job.handler = &GrchiveApiJobHandler{
			taskId: task.Id,
			data:   task.TaskData,
			userId: task.UserId,
		}
	default:
		return nil, errors.New("Unsupported task type.")
	}

	return &job, nil
}

func createJob(t core.ScheduledTaskMetadata, oneTime *core.ScheduledTaskOneTime, recurring *core.ScheduledTaskRecurrence, c core.Clock) (*Job, error) {
	var schedule *Schedule
	var err error
	if oneTime != nil {
		schedule, err = CreateOneTimeJobSchedule(oneTime, c)
	} else if recurring != nil {
		schedule, err = CreateRecurringJobSchedule(recurring, c)
	} else {
		return nil, errors.New("No job schedule found.")
	}

	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, nil
	}

	job, err := CreateJobFromTaskMetadata(t, schedule)
	if err != nil {
		return nil, err
	}
	return job, nil
}
