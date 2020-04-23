package main

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
)

type JobTaskHandler interface {
	Tick(c core.Clock) error
}

type Job struct {
	id       int64
	schedule *Schedule
	typ      core.TaskType

	handler JobTaskHandler
}

func (j Job) Id() int64 {
	return j.id
}

// Returns a boolean that indicates whether this job will want to run again.
func (j *Job) Tick(c core.Clock, force bool) (bool, error) {
	if !j.schedule.ShouldRun(c) && !force {
		return true, nil
	}

	err := j.handler.Tick(c)
	if err != nil {
		return true, err
	}

	j.schedule.MarkRun(c)
	return j.schedule.HasNextRun(c), nil
}

func CreateJobFromTaskMetadata(task core.ScheduledTaskMetadata, schedule *Schedule) (*Job, error) {
	job := Job{
		id:       task.Id,
		schedule: schedule,
		typ:      task.TaskType,
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
