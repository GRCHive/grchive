package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"sort"
	"time"
)

func createPbcNotificationJob(c core.Clock) *Job {
	// This task should be run twice every day at
	// 9am ET and 5pm ET (kind of arbitrary).
	utcLoc, _ := time.LoadLocation("UTC")

	dailyTask := core.ScheduledDailyTaskRawInput{
		Times: []time.Time{
			time.Date(1969, 12, 0, 13, 0, 0, 0, utcLoc),
			time.Date(1969, 12, 0, 1, 0, 0, 0, utcLoc),
		},
	}

	recurring, err := dailyTask.GenerateRecurringTasks(utcLoc)
	if err != nil {
		core.Error("Failed to create recurring task: " + err.Error())
	}

	schedule, err := CreateRecurringJobSchedule(recurring, c)
	if err != nil {
		core.Error("Failed to create job schedule: " + err.Error())
	}

	job := Job{
		id:           PbcNotificationJobId,
		schedule:     schedule,
		handler:      &PbcNotificationHandler{},
		typ:          core.KInternalTask,
		backoffCount: 0,
		backoffTime:  time.Duration(0),
		lastTick:     time.Now(),
	}
	return &job
}

type PbcNotificationHandler struct {
}

func dateOnly(t time.Time) time.Time {
	tm := t.In(time.UTC)
	return time.Date(
		tm.Year(),
		tm.Month(),
		tm.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func createOrgUrl(orgName string) string {
	return core.EnvConfig.SelfUri + "/dashboard/org/" + orgName
}

func sendPbcNotificationToUsers(orgName string, pbc *core.DocumentRequest, daysBefore int32, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}

	users, err := database.FindMultipleUsersFromIds(userIds)
	if err != nil {
		return err
	}

	subject := ""
	message := ""

	if daysBefore == 0 {
		// Day of.
		subject = fmt.Sprintf("GRCHive PBC [%s] Due Today", pbc.Name)
		message = fmt.Sprintf(`Hello,

This is a reminder that the PBC <a href="%s/requests/doc/%d">%s</a> is due <b>today</b>.

Thank you!`, createOrgUrl(orgName), pbc.Id, pbc.Name)
	} else if daysBefore > 0 {
		// Before due date.
		subject = fmt.Sprintf("GRCHive PBC [%s] Due in %d Days", pbc.Name, daysBefore)
		message = fmt.Sprintf(`Hello,

This is a reminder that the PBC <a href="%s/requests/doc/%d">%s</a> is due in %d days.

Thank you!`, createOrgUrl(orgName), pbc.Id, pbc.Name, daysBefore)
	} else {
		// After due date.
		subject = fmt.Sprintf("GRCHive PBC [%s] is Overdue (%d Days)", pbc.Name, -daysBefore)
		message = fmt.Sprintf(`Hello,

This is a reminder that the PBC <a href="%s/requests/doc/%d">%s</a> is overdue by %d days.

Thank you!`, createOrgUrl(orgName), pbc.Id, pbc.Name, -daysBefore)
	}

	for _, u := range users {
		mailPayload := mail.MailPayload{
			From: mail.Email{
				Name:  "GRCHive Support",
				Email: "support@grchive.com",
			},
			To: mail.Email{
				Name:  u.FullName(),
				Email: u.Email,
			},
			Subject: subject,
			Message: message,
		}

		// Oh well if an email notification doesn't get sent.
		err = mail.MailProvider.SendMail(mailPayload)
		if err != nil {
			core.Warning("Failed to send email: " + err.Error())
		}
	}

	return nil
}

func (h *PbcNotificationHandler) Tick(c core.Clock) error {
	tm := dateOnly(c.Now())

	// Get all orgs.
	allOrgs, err := database.FindOrganizations()
	if err != nil {
		return err
	}

	// Iterate through orgs - get PBC notification setting.
	for _, org := range allOrgs {
		settings, err := database.GetOrgPbcNotificationCadenceSettings(org.Id)
		if err != nil {
			return err
		}

		// Iterate through all PBCs and figure out which one needs a notification event generated.
		pbcs, err := database.GetAllDocumentRequestsForOrganization(
			org.Id,
			core.ValidDueDateDocRequestFilter,
			core.ServerRole,
		)

		if err != nil {
			return err
		}

		notificationRecord, err := database.GetPbcNotificationRecord(org.Id)
		if err != nil {
			return err
		}

		// Sort PBCs by due date so that we can efficiently process which PBCs match up which settings.
		// 	N: number of settings
		// 	M : number of PBCs.
		// If we didn't sort, we would have to do a O(NM) task to determine which PBCs need to fire off which notification.
		// If we do sort, we can do it in O(M log M) + O(N) time instead.
		//
		// For each pbc, we can define the "time to due date" (TTD) as TTD = DUE_DATE - NOW().
		// If we sort the PBC's by due date in ascending order, we are guaranteed that the TTD to be in ascending order as well.
		// A smaller TTD indicates that the PBC is overdue while a larger PBC indicates that there's more time until the PBC is due.
		sort.Slice(pbcs, func(i int, j int) bool {
			return pbcs[i].DueDate.NullTime.Time.Before(pbcs[j].DueDate.NullTime.Time)
		})

		// Need to go backwards since the settings are ordered
		// from earliest (pre-due date) notification to latest (post due-date).
		// Since we ordered PBCs with increasing TTD, the latest settings are what
		// we need to handle first.
		settingIdx := len(settings) - 1
		for _, pbc := range pbcs {
			// This shouldn't ever happen but putting it in here to be safe.
			if !pbc.DueDate.NullTime.Valid {
				continue
			}

			dueDate := dateOnly(pbc.DueDate.NullTime.Time)
			ttd := dueDate.Sub(tm).Hours() / 24.0

			for settingIdx >= 0 {
				if ttd >= float64(settings[settingIdx].DaysBeforeDue) {
					settingIdx = settingIdx - 1
				} else {
					break
				}
			}

			if settingIdx < 0 {
				break
			}

			usableSetting := settings[settingIdx]

			// Need to make sure we didn't send a notification for this particular
			// setting already.
			recordKey := core.PbcNotificationRecordKey{
				CadenceId: usableSetting.Id,
				OrgId:     org.Id,
				RequestId: pbc.Id,
			}
			_, alreadySent := notificationRecord[recordKey]

			if alreadySent {
				continue
			}

			// Send notification
			userIdsToSendTo := make([]int64, 0)
			userIdsToSendTo = append(userIdsToSendTo, usableSetting.AdditionalUsers...)
			if usableSetting.SendToRequester {
				userIdsToSendTo = append(userIdsToSendTo, pbc.RequestedUserId)
			}

			if usableSetting.SendToAssignee && pbc.AssigneeUserId.NullInt64.Valid {
				userIdsToSendTo = append(userIdsToSendTo, pbc.AssigneeUserId.NullInt64.Int64)
			}

			err = sendPbcNotificationToUsers(org.OktaGroupName, pbc, usableSetting.DaysBeforeDue, userIdsToSendTo)
			if err != nil {
				return err
			}

			// Mark notification as sent.
			tx := database.CreateTx()
			err = database.WrapTx(tx, func() error {
				return database.MarkPbcNotificationRecordWithTx(tx, recordKey)
			})
			// Maybe not needed?
			notificationRecord[recordKey] = true

			if err != nil {
				return err
			}
		}
	}

	return nil
}
