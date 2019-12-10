export interface DocumentRequest {
    Name:            string
    Description:     string
    CatId:           number
    OrgId:           number
    RequestedUserId: number
    CompletionTime:  Date | null
}
