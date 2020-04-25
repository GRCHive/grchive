package core

import (
	"github.com/teambition/rrule-go"
	"time"
)

type TaskType int32

const (
	KRabbitMQTask   TaskType = 1
	KGrchiveApiTask          = 2
)

type ScheduledTaskMetadata struct {
	Id            int64       `db:"id"`
	Name          string      `db:"name"`
	Description   string      `db:"description"`
	OrgId         int32       `db:"org_id"`
	UserId        int64       `db:"user_id"`
	TaskType      TaskType    `db:"task_type"`
	TaskData      interface{} `db:"task_data"`
	ScheduledTime time.Time   `db:"scheduled_time"`
}

type ScheduledTaskOneTime struct {
	EventId       int64     `db:"event_id"`
	EventDateTime time.Time `db:"event_date_time"`
}

type ScheduledTaskRecurrence struct {
	EventId       int64     `db:"event_id"`
	StartDateTime time.Time `db:"start_date_time"`
	RRule         rrule.Set `db:"rrule"`
	Timezone      string    `db:"timezone"`
}

type FullScheduledTask struct {
	Metadata  ScheduledTaskMetadata
	OneTime   *ScheduledTaskOneTime
	Recurring *ScheduledTaskRecurrence
}

// Raw inputs are what we expect the clients to give us in a sort of expanded/raw
// form that we convert into iCal RRules to store in the database.

type CronFrequency int32

const (
	KCronFreqDaily   CronFrequency = 0
	KCronFreqWeekly  CronFrequency = 1
	KCronFreqMonthly CronFrequency = 2
)

type CronWeekdayHash int32

const (
	KCronHashFirst  CronWeekdayHash = 0
	KCronHashSecond CronWeekdayHash = 1
	KCronHashThird  CronWeekdayHash = 2
	KCronHashFourth CronWeekdayHash = 3
	KCronHashLast   CronWeekdayHash = 4
)

type ScheduledDailyTaskRawInput struct {
	Times []time.Time
}

func (t ScheduledDailyTaskRawInput) GenerateRecurringTasks(tz *time.Location) (*ScheduledTaskRecurrence, error) {
	ret := ScheduledTaskRecurrence{
		StartDateTime: time.Now().In(tz),
		RRule:         rrule.Set{},
		Timezone:      tz.String(),
	}

	for _, tm := range t.Times {
		dt := CombineDateWithTime(ret.StartDateTime, tm.In(tz))

		rule, err := rrule.NewRRule(rrule.ROption{
			Freq:     rrule.DAILY,
			Dtstart:  dt,
			Interval: 1,
		})

		if err != nil {
			return nil, err
		}

		ret.RRule.RRule(rule)
	}

	return &ret, nil
}

type ScheduledWeeklyTaskRawInput struct {
	Days []Days
	Time time.Time
}

func (t ScheduledWeeklyTaskRawInput) GenerateRecurringTasks(tz *time.Location) (*ScheduledTaskRecurrence, error) {
	ret := ScheduledTaskRecurrence{
		StartDateTime: time.Now().In(tz),
		RRule:         rrule.Set{},
		Timezone:      tz.String(),
	}

	opt := rrule.ROption{
		Freq:      rrule.WEEKLY,
		Dtstart:   CombineDateWithTime(ret.StartDateTime, t.Time.In(tz)),
		Interval:  1,
		Wkst:      rrule.SU,
		Byweekday: make([]rrule.Weekday, 0),
	}

	for _, d := range t.Days {
		opt.Byweekday = append(opt.Byweekday, DaysToRRule(d))
	}

	rule, err := rrule.NewRRule(opt)
	if err != nil {
		return nil, err
	}

	ret.RRule.RRule(rule)
	return &ret, nil
}

type ScheduledMonthlyTaskRawInput struct {
	UseDate bool
	Dates   []int32
	Nth     CronWeekdayHash
	Day     Days
	Time    time.Time
}

func (t ScheduledMonthlyTaskRawInput) GenerateRecurringTasks(tz *time.Location) (*ScheduledTaskRecurrence, error) {
	ret := ScheduledTaskRecurrence{
		StartDateTime: time.Now().In(tz),
		RRule:         rrule.Set{},
		Timezone:      tz.String(),
	}

	opt := rrule.ROption{
		Freq:     rrule.MONTHLY,
		Dtstart:  CombineDateWithTime(ret.StartDateTime, t.Time.In(tz)),
		Interval: 1,
	}

	if t.UseDate {
		opt.Bymonthday = make([]int, 0)
		for _, dt := range t.Dates {
			opt.Bymonthday = append(opt.Bymonthday, int(dt))
		}
	} else {
		day := DaysToRRule(t.Day)

		switch t.Nth {
		case KCronHashFirst:
			day = day.Nth(1)
		case KCronHashSecond:
			day = day.Nth(2)
		case KCronHashThird:
			day = day.Nth(3)
		case KCronHashFourth:
			day = day.Nth(4)
		case KCronHashLast:
			day = day.Nth(-1)
		}
		opt.Byweekday = []rrule.Weekday{day}
	}

	rule, err := rrule.NewRRule(opt)
	if err != nil {
		return nil, err
	}

	ret.RRule.RRule(rule)
	return &ret, nil
}

type ScheduledTaskRawInput struct {
	Repeat      bool
	OneTimeDate NullTime
	Frequency   CronFrequency
	Name        string
	Description string
	Daily       *ScheduledDailyTaskRawInput
	Weekly      *ScheduledWeeklyTaskRawInput
	Monthly     *ScheduledMonthlyTaskRawInput
	Timezone    string
}

func (t ScheduledTaskRawInput) GenerateTaskMetadata(userId int64, orgId int32, taskType TaskType, data interface{}) *ScheduledTaskMetadata {
	return &ScheduledTaskMetadata{
		Name:        t.Name,
		Description: t.Description,
		OrgId:       orgId,
		UserId:      userId,
		TaskType:    taskType,
		TaskData:    data,
	}
}

func (t ScheduledTaskRawInput) GenerateOneTimeTask(tz *time.Location) *ScheduledTaskOneTime {
	return &ScheduledTaskOneTime{
		EventDateTime: t.OneTimeDate.NullTime.Time.In(tz),
	}
}

// Different data needed for various tasks.
type GrchiveApiTaskData struct {
	Endpoint string
	Method   string
	Payload  interface{}
}
