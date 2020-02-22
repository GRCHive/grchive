export interface DocumentRequest {
    Id:              number
    Name:            string
    Description:     string
    OrgId:           number
    RequestedUserId: number
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
    f.RequestTime = new Date(f.RequestTime)
}
