export interface DocumentRequest {
    Id:              number
    Name:            string
    Description:     string
    OrgId:           number
    RequestedUserId: number
    AssigneeUserId:  number | null
    DueDate         : Date | null
    CompletionTime:  Date | null
    FeedbackTime: Date | null
    ProgressTime: Date | null
    RequestTime:     Date
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

    f.RequestTime = new Date(f.RequestTime)
}

export enum DocRequestStatus {
    Open,
    InProgress,
    Feedback,
    Complete,
    Overdue,
}

export function getDocumentRequestStatus(r : DocumentRequest) : DocRequestStatus {
    let currentTime = new Date()
    if (!r.CompletionTime) {
        // No completion time: open, in progress, overdue.
        if (!!r.DueDate) {
            if (currentTime <= r.DueDate) {
                return DocRequestStatus.Open
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
