import { 
    StringFilterData, NullStringFilterData, StringComparisonOperators,
    TimeRangeFilterData, NullTimeRangeFilterDate, cleanTimeRangeFilterDataFromJson
} from './filters'

export interface AuditEventEntry {
    Id                : number
    OrgId             : number
    ResourceType      : string
    ResourceId        : number
    ResourceExtraData : object
    Action            : string
    PerformedAt       : Date
    UserId            : number | null // If null, it means it was a system generated event or a deleted user.
}

export function cleanAuditEventEntryFromJson(e : AuditEventEntry) {
    e.PerformedAt = new Date(e.PerformedAt) 
}

export interface AuditTrailFilterData {
    ResourceTypeFilter : StringFilterData
    ActionFilter: StringFilterData
    UserFilter: StringFilterData
    TimeRangeFilter: TimeRangeFilterData
}

export let NullAuditTrailFilterData : AuditTrailFilterData = {
    ResourceTypeFilter : JSON.parse(JSON.stringify(NullStringFilterData)),
    ActionFilter: JSON.parse(JSON.stringify(NullStringFilterData)),
    UserFilter: JSON.parse(JSON.stringify(NullStringFilterData)),
    TimeRangeFilter: cleanTimeRangeFilterDataFromJson(JSON.parse(JSON.stringify(NullTimeRangeFilterDate))),
}
