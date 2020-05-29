export enum ShellTypes {
    Bash = 1,
    PowerShell = 2,
}

export let ShellTypeToCodeMirror = new Map<ShellTypes, string>()
ShellTypeToCodeMirror.set(ShellTypes.Bash, "text/x-sh")
ShellTypeToCodeMirror.set(ShellTypes.PowerShell, "application/x-powershell")

export interface ShellScript {
    Id : number
    OrgId : number
    TypeId : number
    Name : string
    Description : string
    BucketId : string
    StorageId: string
}

export interface ShellScriptVersion {
    Id: number
    ShellId: number
    OrgId: number
    UploadTime: Date
    UploadUserId : number
    GcsGeneration : number
}

export function cleanShellScriptVersionFromJson(v : ShellScriptVersion) {
    v.UploadTime = new Date(v.UploadTime)
}
