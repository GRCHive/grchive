export interface DocumentRequest {
    Id:              number
    Name:            string
    Description:     string
    CatId:           number
    OrgId:           number
    RequestedUserId: number
    CompletionTime:  Date | null
    RequestTime:     Date
}
