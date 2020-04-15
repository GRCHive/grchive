export interface ManagedCode {
    Id : number
    OrgId : number
    GitHash : string
    ActionTime : Date
    GitPath : string
    GiteaFileSha: string
}

export function cleanManagedCodeFromJson(c : ManagedCode) {
    c.ActionTime = new Date(c.ActionTime)
}

export interface SupportedParamType {
    Id : number
    Name : string
}

export interface CodeParamType {
    Name: string
    ParamId: number
}
