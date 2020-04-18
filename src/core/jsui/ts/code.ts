export interface ManagedCode {
    Id : number
    OrgId : number
    GitHash : string
    ActionTime : Date
    GitPath : string
    GiteaFileSha: string
    UserId : number
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

export interface ScriptRun {
    Id: number
    LinkId: number
    StartTime: Date
    RequiresBuild: boolean
    BuildStartTime: Date | null
    BuildFinishTime: Date | null
    BuildSuccess: boolean
    RunStartTime: Date | null
    RunFinishTime: Date | null
    RunSuccess: boolean
    UserId : number
}

export interface DroneCiStatus {
	CodeId     : number
	OrgId      : number
	CommitHash : string
	TimeStart  : Date
	TimeEnd    : Date | null
	Success    : boolean
	Logs       : string
	Jar        : string
}

export function cleanDroneCiStatusFromJson(r : DroneCiStatus) {
    r.TimeStart = new Date(r.TimeStart)
    if (!!r.TimeEnd) {
        r.TimeEnd = new Date(r.TimeEnd)
    }
}

export function cleanScriptRunFromJson(r : ScriptRun) {
    r.StartTime = new Date(r.StartTime)
    if (!!r.BuildStartTime) {
        r.BuildStartTime = new Date(r.BuildStartTime)
    }

    if (!!r.BuildFinishTime) {
        r.BuildFinishTime = new Date(r.BuildFinishTime)
    }

    if (!!r.RunStartTime) {
        r.RunStartTime = new Date(r.RunStartTime)
    }

    if (!!r.RunFinishTime) {
        r.RunFinishTime = new Date(r.RunFinishTime)
    }
}
