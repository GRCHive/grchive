export interface GenericRequest {
    Id : number
    OrgId : number
    UploadTime : Date
    UploadUserId : number
    Name : string
    Assignee : number | null
    DueDate : Date | null
    Description : string
}

export interface GenericApproval {
    Id : number
    RequestId : number
    ResponseTime: Date
    ResponderUserId : number
    Response : boolean
    Reason : string
}

export function cleanGenericRequestFromJson(g : GenericRequest) {
    g.UploadTime = new Date(g.UploadTime)
    if (!!g.DueDate) {
        g.DueDate = new Date(g.DueDate)
    }
}

export function cleanGenericApprovalFromJson(g : GenericApproval) {
    g.ResponseTime = new Date(g.ResponseTime)
}
