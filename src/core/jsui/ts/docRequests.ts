export interface DocumentRequest {
    Id:              number
    Name:            string
    Description:     string
    OrgId:           number
    RequestedUserId: number
    AssigneeUserId:  number | null
    DueDate:         Date | null
    CompletionTime:  Date | null
    FeedbackTime:    Date | null
    ProgressTime:    Date | null
    RequestTime:     Date
    ApproveTime:     Date | null
    ApproveUserId:   number | null
}

export function cleanJsonDocumentRequest(f : DocumentRequest) {
    if (!!f.CompletionTime) {
        f.CompletionTime = new Date(f.CompletionTime)
    }

    if (!!f.DueDate) {
        f.DueDate = new Date(f.DueDate)
    }

    if (!!f.FeedbackTime) {
        f.FeedbackTime = new Date(f.FeedbackTime)
    }

    if (!!f.ProgressTime) {
        f.ProgressTime = new Date(f.ProgressTime)
    }

    if (!!f.ApproveTime) {
        f.ApproveTime = new Date(f.ApproveTime)
    }

    f.RequestTime = new Date(f.RequestTime)
}

export enum DocRequestStatus {
    Open,
    InProgress,
    Feedback,
    Complete,
    Overdue,
    Approved
}

export function getDocRequestStatusString(s : DocRequestStatus) : string {
    switch (s) {
        case DocRequestStatus.Open:
            return 'Open'
        case DocRequestStatus.InProgress:
            return 'In Progress'
        case DocRequestStatus.Feedback:
            return 'Feedback'
        case DocRequestStatus.Complete:
            return 'Pending Approval'
        case DocRequestStatus.Overdue:
            return 'Overdue'
        case DocRequestStatus.Approved:
            return 'Approved'
    }
}

export const allDocRequestStatus : DocRequestStatus[] = [
    DocRequestStatus.Open,
    DocRequestStatus.InProgress,
    DocRequestStatus.Feedback,
    DocRequestStatus.Complete,
    DocRequestStatus.Overdue,
    DocRequestStatus.Approved,
]

export function getDocumentRequestStatus(r : DocumentRequest) : DocRequestStatus {
    let currentTime = new Date()
    if (!!r.ApproveTime && !!r.ApproveUserId) {
        return DocRequestStatus.Approved
    } else if (!r.CompletionTime) {
        // No completion time: open, in progress, overdue.
        if (!!r.DueDate) {
            if (currentTime <= r.DueDate) {
                if (!!r.ProgressTime) {
                    return DocRequestStatus.InProgress
                } else {
                    return DocRequestStatus.Open
                }
            } else {
                return DocRequestStatus.Overdue
            }
        } else if (!!r.ProgressTime) {
            return DocRequestStatus.InProgress
        } else {
            return DocRequestStatus.Open
        }
    } else if (!!r.FeedbackTime) {
        // Has feedback time: complete, feedback or overdue.
        // Completion time should be filled out already!
        if (!!r.DueDate && currentTime > r.DueDate) {
            return DocRequestStatus.Overdue
        } else if (r.FeedbackTime <= r.CompletionTime) {
            return DocRequestStatus.Complete
        } else {
            return DocRequestStatus.Feedback
        }
    } else {
        // No feedback but has completion: complete
        return DocRequestStatus.Complete
    }
}

export interface DocRequestStatusFilterData {
    ValidStatuses: DocRequestStatus[]
}

export let NullDocRequestStatusFilterData : DocRequestStatusFilterData = {
    ValidStatuses: []
}

export function copyDocRequestStatusFilterData(c : DocRequestStatusFilterData) : DocRequestStatusFilterData {
    let ret = JSON.parse(JSON.stringify(c))
    return ret
}

import {
    TimeRangeFilterData, NullTimeRangeFilterDate, copyTimeRangeFilterData, cleanTimeRangeFilterDataFromJson,
    UserFilterData, NullUserFilterData, copyUserFilterData,
} from './filters'

export interface DocRequestFilterData {
    RequestTimeFilter: TimeRangeFilterData
    DueDateFilter: TimeRangeFilterData
    StatusFilter: DocRequestStatusFilterData
    RequesterFilter: UserFilterData
    AssigneeFilter: UserFilterData
}

export let NullDocRequestFilterData : DocRequestFilterData = {
    RequestTimeFilter : copyTimeRangeFilterData(NullTimeRangeFilterDate),
    DueDateFilter: copyTimeRangeFilterData(NullTimeRangeFilterDate),
    StatusFilter: copyDocRequestStatusFilterData(NullDocRequestStatusFilterData),
    RequesterFilter: copyUserFilterData(NullUserFilterData),
    AssigneeFilter: copyUserFilterData(NullUserFilterData),
}

export function copyDocRequestFilterData(c : DocRequestFilterData) : DocRequestFilterData {
    let ret = JSON.parse(JSON.stringify(c))
    cleanDocRequestFilterDataFromJson(ret)
    return ret
}

export function cleanDocRequestFilterDataFromJson(c : DocRequestFilterData) {
    cleanTimeRangeFilterDataFromJson(c.DueDateFilter)    
    cleanTimeRangeFilterDataFromJson(c.RequestTimeFilter)    
}
