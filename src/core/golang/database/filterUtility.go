package database

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"strings"
	"time"
)

func buildNumericFilter(ref string, f core.NumericFilterData) string {
	var template string = ""

	switch f.Op {
	case core.Disabled:
		return fmt.Sprintf("%s IS NOT NULL", ref)
	case core.Equal:
		template = "%s = %d"
	case core.NotEqual:
		template = "%s != %d"
	case core.Greater:
		template = "%s > %d"
	case core.GreaterEqual:
		template = "%s >= %d"
	case core.Less:
		template = "%s < %d"
	case core.LessEqual:
		template = "%s <= %d"
	}

	return fmt.Sprintf(template, ref, f.Target)
}

func buildStringFilter(ref string, f core.StringFilterData) string {
	var template string = ""

	switch f.Op {
	case core.SOpDisabled:
		return fmt.Sprintf("%s IS NOT NULL", ref)
	case core.SOpEqual:
		template = "%s = '%s'"
	case core.SOpNotEqual:
		template = "%s != '%s'"
	case core.SOpContains:
		template = "%s LIKE  '%%%s%%'"
	case core.SOpExclude:
		template = "%s NOT LIKE '%%%s%%'"
	}

	return fmt.Sprintf(template, ref, f.Target)
}

func buildTimeRangeFilter(ref string, f core.TimeRangeFilterData) string {
	if !f.Enabled || (!f.Start.NullTime.Valid && !f.End.NullTime.Valid) {
		return fmt.Sprintf("(%s IS NOT NULL OR %s IS NULL)", ref, ref)
	}

	builder := strings.Builder{}
	useAnd := f.Start.NullTime.Valid && f.End.NullTime.Valid

	if f.Start.NullTime.Valid {
		builder.WriteString(fmt.Sprintf("%s >= '%s'", ref, f.Start.NullTime.Time.Format(time.RFC3339)))
	}

	if useAnd {
		builder.WriteString(" AND ")
	}

	if f.End.NullTime.Valid {
		builder.WriteString(fmt.Sprintf("%s <= '%s'", ref, f.End.NullTime.Time.Format(time.RFC3339)))
	}

	return builder.String()
}

func buildDocRequestStatusFilter(completionTime string, feedbackTime string, progressTime string, dueDate string, f core.DocRequestStatusFilterData) string {
	if len(f.ValidStatuses) == 0 {
		return "TRUE = TRUE"
	}

	builder := strings.Builder{}

	for idx, status := range f.ValidStatuses {
		builder.WriteString("(")
		switch status {
		case core.DocRequestOpen:
			builder.WriteString(fmt.Sprintf("%s IS NULL", completionTime))
			builder.WriteString(" AND ")
			builder.WriteString(fmt.Sprintf("%s IS NULL", progressTime))
			builder.WriteString(" AND ")
			{
				builder.WriteString("(")
				builder.WriteString(fmt.Sprintf("%s IS NULL", dueDate))
				builder.WriteString(" OR ")
				builder.WriteString(fmt.Sprintf("NOW() <= %s", dueDate))
				builder.WriteString(")")
			}
		case core.DocRequestInProgress:
			builder.WriteString(fmt.Sprintf("%s IS NULL", completionTime))
			builder.WriteString(" AND ")
			builder.WriteString(fmt.Sprintf("%s IS NOT NULL", progressTime))
			builder.WriteString(" AND ")
			{
				builder.WriteString("(")
				builder.WriteString(fmt.Sprintf("%s IS NULL", dueDate))
				builder.WriteString(" OR ")
				builder.WriteString(fmt.Sprintf("NOW() <= %s", dueDate))
				builder.WriteString(")")
			}
		case core.DocRequestFeedback:
			builder.WriteString(fmt.Sprintf("%s IS NOT NULL", completionTime))
			builder.WriteString(" AND ")
			builder.WriteString(fmt.Sprintf("%s IS NOT NULL", feedbackTime))
			builder.WriteString(" AND ")
			builder.WriteString(fmt.Sprintf("%s < %s", completionTime, feedbackTime))
			builder.WriteString(" AND ")
			{
				builder.WriteString("(")
				builder.WriteString(fmt.Sprintf("%s IS NULL", dueDate))
				builder.WriteString(" OR ")
				builder.WriteString(fmt.Sprintf("NOW() <= %s", dueDate))
				builder.WriteString(")")
			}
		case core.DocRequestComplete:
			builder.WriteString(fmt.Sprintf("%s IS NOT NULL", completionTime))
			builder.WriteString(" AND ")
			{
				builder.WriteString("(")
				builder.WriteString(fmt.Sprintf("%s IS NULL", feedbackTime))
				builder.WriteString(" OR ")
				builder.WriteString(fmt.Sprintf("%s > %s", completionTime, feedbackTime))
				builder.WriteString(")")
			}
		case core.DocRequestOverdue:
			builder.WriteString(fmt.Sprintf("%s IS NOT NULL", dueDate))
			builder.WriteString(" AND ")
			builder.WriteString(fmt.Sprintf("NOW() > %s", dueDate))
			builder.WriteString(" AND ")
			{
				builder.WriteString("(")
				builder.WriteString(fmt.Sprintf("%s IS NULL", completionTime))
				builder.WriteString(" OR ")
				builder.WriteString(fmt.Sprintf("%s < %s", completionTime, feedbackTime))
				builder.WriteString(")")
			}
		}
		builder.WriteString(")")

		if idx != len(f.ValidStatuses)-1 {
			builder.WriteString(" OR ")
		}
	}

	return builder.String()
}

func buildUserFilter(ref string, f core.UserFilterData) string {
	if len(f.UserIds) == 0 {
		return "TRUE = TRUE"
	}

	hasNull := false
	hasUsers := false

	for _, u := range f.UserIds {
		if u.NullInt64.Valid {
			hasUsers = true
		} else {
			hasNull = true
		}
	}

	builder := strings.Builder{}
	builder.WriteString("(")

	if hasUsers {
		builder.WriteString(ref)
		builder.WriteString(" IN (")

		for idx, u := range f.UserIds {
			if u.NullInt64.Valid {
				builder.WriteString(fmt.Sprintf("%d", u.NullInt64.Int64))
				if idx != len(f.UserIds)-1 {
					builder.WriteString(",")
				}
			}
		}
		builder.WriteString(")")
	}

	if hasNull {
		if hasUsers {
			builder.WriteString(" OR ")
		}
		builder.WriteString(fmt.Sprintf("%s IS NULL", ref))
	}

	builder.WriteString(")")

	return builder.String()
}

func buildDocRequestFilter(docReqTable string, f core.DocRequestFilterData) string {
	builder := strings.Builder{}

	builder.WriteString(buildTimeRangeFilter(fmt.Sprintf("%s.request_time", docReqTable), f.RequestTimeFilter))
	builder.WriteString(" AND ")
	builder.WriteString(buildTimeRangeFilter(fmt.Sprintf("%s.due_date", docReqTable), f.DueDateFilter))
	builder.WriteString(" AND ")
	builder.WriteString(buildDocRequestStatusFilter(
		fmt.Sprintf("%s.completion_time", docReqTable),
		fmt.Sprintf("%s.feedback_time", docReqTable),
		fmt.Sprintf("%s.progress_time", docReqTable),
		fmt.Sprintf("%s.due_date", docReqTable),
		f.StatusFilter,
	))
	builder.WriteString(" AND ")
	builder.WriteString(buildUserFilter(fmt.Sprintf("%s.requested_user_id", docReqTable), f.RequesterFilter))
	builder.WriteString(" AND ")
	builder.WriteString(buildUserFilter(fmt.Sprintf("%s.assignee", docReqTable), f.AssigneeFilter))
	return builder.String()
}
