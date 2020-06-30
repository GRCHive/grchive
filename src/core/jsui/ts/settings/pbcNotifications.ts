export interface PbcNotificationCadenceSettings {
    Id: number
    OrgId : number
    DaysBeforeDue : number
    SendToAssignee: boolean
    SendToRequester: boolean
    AdditionalUsers: number[]
}
