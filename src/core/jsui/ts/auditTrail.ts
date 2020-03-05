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
