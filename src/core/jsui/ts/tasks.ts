export interface ScheduledTaskMetadata {
    Id : number
    Name : string
    Description : string
    OrgId : number
    UserId : number
    ScheduledTime : Date
}

export function cleanScheduledTaskMetadataFromJson(m : ScheduledTaskMetadata) {
    m.ScheduledTime = new Date(m.ScheduledTime)
}
