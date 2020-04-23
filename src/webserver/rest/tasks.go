package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type AllScheduledTasksInput struct {
	OrgId    int32          `webcore:"orgId"`
	Range    core.TimeRange `webcore:"range"`
	Timezone string         `webcore:"timezone"`
}

func allScheduledTasks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllScheduledTasksInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tasks, err := database.GetAllScheduledTasksForOrgId(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get tasks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tz, err := time.LoadLocation(inputs.Timezone)
	if err != nil {
		core.Warning("Failed to load timezone: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type RetData struct {
		Tasks []*core.ScheduledTaskMetadata
		Times map[int64][]core.TimeRange
	}

	ret := RetData{
		Tasks: make([]*core.ScheduledTaskMetadata, 0),
		Times: map[int64][]core.TimeRange{},
	}
	for _, t := range tasks {
		ret.Tasks = append(ret.Tasks, &t.Metadata)
		ret.Times[t.Metadata.Id] = make([]core.TimeRange, 0)

		timesToAdd := make([]time.Time, 0)

		// Determine if this tasks exists within the given time range.
		if t.OneTime != nil && inputs.Range.InRange(t.OneTime.EventDateTime) {
			timesToAdd = append(timesToAdd, t.OneTime.EventDateTime.In(tz))
		} else if t.Recurring != nil {
			timesToAdd = append(
				timesToAdd,
				t.Recurring.RRule.Between(
					inputs.Range.Start.In(tz),
					inputs.Range.End.In(tz),
					true,
				)...)
		}

		// Assume all tasks last 30 minutes for now just to get a calendar event up.
		// In the future we want to more smartly predict how long the task will take.
		for _, add := range timesToAdd {
			ret.Times[t.Metadata.Id] = append(
				ret.Times[t.Metadata.Id],
				core.TimeRange{
					Start: add,
					End:   add.Add(30 * time.Minute),
				},
			)
		}
	}

	jsonWriter.Encode(ret)
}
