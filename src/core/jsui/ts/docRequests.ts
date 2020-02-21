export interface DocumentRequest {
    Id:              number
    Name:            string
    Description:     string
    OrgId:           number
    RequestedUserId: number
    CompletionTime:  Date | null
    RequestTime:     Date
}

export function cleanJsonDocumentRequest(f : DocumentRequest) {
    if (!!f.CompletionTime) {
        f.CompletionTime = new Date(f.CompletionTime)
    }
    f.RequestTime = new Date(f.RequestTime)
}
