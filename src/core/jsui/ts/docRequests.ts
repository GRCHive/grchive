export interface DocumentRequest {
    Id:              number
    Name:            string
    Description:     string
    OrgId:           number
    RequestedUserId: number
    AssigneeUserId:  number | null
    DueDate         : Date | null
    CompletionTime:  Date | null
    RequestTime:     Date
}

export enum RequestLinkageMode {
    None = 0,
    DocCat = 1,
    Controls = 2
}

export let requestLinkageItems : any[] = [
    {
        text: "None",
        value: RequestLinkageMode.None
    },
    {
        text: "Document Category",
        value: RequestLinkageMode.DocCat
    },
    {
        text: "Control",
        value: RequestLinkageMode.Controls
    },
]

export function cleanJsonDocumentRequest(f : DocumentRequest) {
    if (!!f.CompletionTime) {
        f.CompletionTime = new Date(f.CompletionTime)
    }

    if (!!f.DueDate) {
        f.DueDate = new Date(f.DueDate)
    }
    f.RequestTime = new Date(f.RequestTime)
}
