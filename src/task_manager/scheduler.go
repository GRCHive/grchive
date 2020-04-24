package main

import (
	"errors"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"sync"
	"time"
)

type Scheduler struct {
	Clock core.Clock
	Log   bool

	// Jobs to run on the next tick.
	jobs     map[int64]*Job
	jobMutex sync.RWMutex
}

func CreateScheduler(c core.Clock) *Scheduler {
	return &Scheduler{
		Clock: c,
		jobs:  map[int64]*Job{},
	}
}

func (s *Scheduler) SyncRun() {
	for {
		s.RunImmediate(false)
		time.Sleep(1 * time.Second)
	}
}

func (s *Scheduler) RunImmediate(force bool) {
	jobsToRemove := []int64{}

	s.jobMutex.RLock()
	if s.Log {
		strT, _ := s.Clock.Now().MarshalText()
		core.Info(fmt.Sprintf("[%s] Tick %d Jobs", string(strT), len(s.jobs)))
	}

	for id, j := range s.jobs {
		hasNext, err := j.Tick(s.Clock, force)
		if err != nil {
			// Keep running the job as maybe it's just an infra blip that'll recover.
			// If it's a faulty job then it's probably our fault anyway so we'll just
			// have to fix it for the clients.
			core.Warning(fmt.Sprintf("Failed to Run Job [%d]: %s", id, err.Error()))
		}

		if !hasNext {
			jobsToRemove = append(jobsToRemove, id)
		}
	}
	s.jobMutex.RUnlock()

	s.jobMutex.Lock()
	for _, j := range jobsToRemove {
		delete(s.jobs, j)
	}
	s.jobMutex.Unlock()
}

func (s *Scheduler) AddJob(j *Job) error {
	if j == nil {
		return nil
	}

	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	_, ok := s.jobs[j.Id()]
	if ok {
		return errors.New("Job already exists.")
	}

	s.jobs[j.Id()] = j
	return nil
}

func (s *Scheduler) RemoveJob(jobId int64) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()
	delete(s.jobs, jobId)
}