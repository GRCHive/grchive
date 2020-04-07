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
