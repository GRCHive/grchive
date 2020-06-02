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

export interface ShellScriptRun {
    Id: number
    ScriptVersionId: number
    RunUserId: number
    CreateTime: Date
    RunTime: Date | null
    EndTime: Date | null
}

export function cleanShellScriptRunFromJson(r : ShellScriptRun) {
    r.CreateTime = new Date(r.CreateTime)
    if (!!r.RunTime) {
        r.RunTime = new Date(r.RunTime)
    }

    if (!!r.EndTime) {
        r.EndTime = new Date(r.EndTime)
    }
}

export interface ShellScriptRunPerServer {
    RunId: number
    OrgId: number
    ServerId: number
    EncryptedLog: string | null
    RunTime: Date | null
    EndTime: Date | null
    Success: boolean
}

export function cleanShellScriptRunPerServerFromJson(r : ShellScriptRunPerServer) {
    if (!!r.RunTime) {
        r.RunTime = new Date(r.RunTime)
    }

    if (!!r.EndTime) {
        r.EndTime = new Date(r.EndTime)
    }
}
